// SPDX-FileCopyrightText: 2022-present Intel Corporation
// SPDX-FileCopyrightText: 2020-present Open Networking Foundation <info@opennetworking.org>
//
// SPDX-License-Identifier: Apache-2.0

package mobility

import (
	"context"
	"fmt"
	"math"
	"math/rand"
	"sync"
	"time"

	e2sm_mho "github.com/onosproject/onos-e2-sm/servicemodels/e2sm_mho_go/v2/e2sm-mho-go"

	"github.com/onosproject/onos-api/go/onos/ransim/types"
	"github.com/onosproject/onos-lib-go/pkg/logging"
	"github.com/onosproject/ran-simulator/pkg/handover"
	"github.com/onosproject/ran-simulator/pkg/measurement"
	"github.com/onosproject/ran-simulator/pkg/model"
	"github.com/onosproject/ran-simulator/pkg/store/cells"
	"github.com/onosproject/ran-simulator/pkg/store/routes"
	"github.com/onosproject/ran-simulator/pkg/store/ues"
	"github.com/onosproject/ran-simulator/pkg/utils"
	"github.com/onosproject/rrm-son-lib/pkg/model/id"
)

var log = logging.GetLogger()

// Driver is an abstraction of an entity driving the UE mobility
type Driver interface {
	// Start starts the driving engine
	Start(ctx context.Context)

	// Stop stops the driving engine
	Stop()

	// GenerateRoutes generates routes for all UEs that currently do not have a route; remove routes with no UEs
	GenerateRoutes(ctx context.Context, minSpeed uint32, maxSpeed uint32, speedStdDev uint32, routeEndPoints []model.RouteEndPoint, directRoute bool)

	// GetMeasCtrl returns the Measurement Controller
	GetMeasCtrl() measurement.MeasController

	// GetHoCtrl
	GetHoCtrl() handover.HOController

	// GetRrcCtrl returns the Rrc Controller
	GetRrcCtrl() RrcCtrl

	// Handover
	Handover(ctx context.Context, imsi types.IMSI, tCell *model.UECell)

	//GetHoLogic
	GetHoLogic() string

	//SetHoLogic
	SetHoLogic(hoLogic string)

	// AddRrcChan
	AddRrcChan(ch chan model.UE)
}

type driver struct {
	cellStore               cells.Store
	routeStore              routes.Store
	ueStore                 ues.Store
	apiKey                  string
	ticker                  *time.Ticker
	done                    chan bool
	stopLocalHO             chan bool
	min                     *model.Coordinate
	max                     *model.Coordinate
	measCtrl                measurement.MeasController
	hoCtrl                  handover.HOController
	hoLogic                 string
	rrcCtrl                 RrcCtrl
	ueLock                  map[types.IMSI]*sync.Mutex
	rrcStateChangesDisabled bool
	wayPointRoute           bool
}

// NewMobilityDriver returns a driving engine capable of "driving" UEs along pre-specified routes
func NewMobilityDriver(cellStore cells.Store, routeStore routes.Store, ueStore ues.Store, apiKey string, hoLogic string, ueCountPerCell uint, rrcStateChangesDisabled bool, wayPointRoute bool) Driver {
	return &driver{
		cellStore:               cellStore,
		routeStore:              routeStore,
		ueStore:                 ueStore,
		hoLogic:                 hoLogic,
		rrcCtrl:                 NewRrcCtrl(ueCountPerCell),
		rrcStateChangesDisabled: rrcStateChangesDisabled,
		wayPointRoute:           wayPointRoute,
	}
}

var tickUnit = time.Second

const tickFrequency = 1

const measType = "EventA3" // ToDo: should be programmable

const hoType = "A3" // ToDo: should be programmable

func (d *driver) Start(ctx context.Context) {
	log.Info("Driver starting")

	// Iterate over all routes and position the UEs at the start of their routes
	for _, route := range d.routeStore.List(ctx) {
		d.initializeUEPosition(ctx, route)
	}

	d.ueLock = make(map[types.IMSI]*sync.Mutex)
	for _, ue := range d.ueStore.ListAllUEs(ctx) {
		d.ueLock[ue.IMSI] = &sync.Mutex{}
	}

	d.ticker = time.NewTicker(tickFrequency * tickUnit)
	d.done = make(chan bool)
	d.stopLocalHO = make(chan bool)

	// Add measController
	d.measCtrl = measurement.NewMeasController(measType, d.cellStore, d.ueStore)
	d.measCtrl.Start(ctx)
	d.hoCtrl = handover.NewHOController(hoType, d.cellStore, d.ueStore)
	d.hoCtrl.Start(ctx)
	// link measController with hoController
	go d.linkMeasCtrlHoCtrl()

	// Add hoController
	if d.hoLogic == "local" {
		log.Info("HO logic is running locally")
		// process handover decision
		go d.processHandoverDecision(ctx)
	} else if d.hoLogic == "mho" {
		log.Info("HO logic is running outside - mho")
	} else {
		log.Warn("There is no handover logic - running measurement only")
	}

	go d.drive(ctx)
}

func (d *driver) Stop() {
	log.Info("Driver stopping")
	d.ticker.Stop()
	d.done <- true
}

func (d *driver) GetMeasCtrl() measurement.MeasController {
	return d.measCtrl
}

func (d *driver) GetHoCtrl() handover.HOController {
	return d.hoCtrl
}

func (d *driver) GetRrcCtrl() RrcCtrl {
	return d.rrcCtrl
}

func (d *driver) SetHoLogic(hoLogic string) {
	if d.hoLogic == "local" && hoLogic == "mho" {
		log.Info("Stopping local HO")
		d.stopLocalHO <- true
	} else if d.hoLogic == "mho" && hoLogic == "local" {
		log.Info("Starting local HO")
		go d.linkMeasCtrlHoCtrl()
	}
	d.hoLogic = hoLogic
}

func (d *driver) AddRrcChan(ch chan model.UE) {
	d.addRrcChan(ch)
}

func (d *driver) lockUE(imsi types.IMSI) {
	d.ueLock[imsi].Lock()
}

func (d *driver) unlockUE(imsi types.IMSI) {
	if _, ok := d.ueLock[imsi]; !ok {
		log.Errorf("lock not found for IMSI %d", imsi)
		return
	}
	d.ueLock[imsi].Unlock()
}

func (d *driver) drive(ctx context.Context) {
	for {
		select {
		case <-d.done:
			ctx.Done()
			close(d.done)
			return
		case <-d.ticker.C:
			for _, route := range d.routeStore.List(ctx) {
				go d.processRoute(ctx, route)

			}
		}
	}
}

func (d *driver) processRoute(ctx context.Context, route *model.Route) {
	d.lockUE(route.IMSI)
	defer d.unlockUE(route.IMSI)
	if route.NextPoint == 0 && !route.Reverse {
		d.initializeUEPosition(ctx, route)
	}
	d.updateUEPosition(ctx, route)
	d.updateUESignalStrength(ctx, route.IMSI)
	if !d.rrcStateChangesDisabled {
		d.updateRrc(ctx, route.IMSI)
	}
	d.updateFiveQI(ctx, route.IMSI)
	d.reportMeasurement(ctx, route.IMSI)
}

// Initializes UE positions to the start of its routes.
func (d *driver) initializeUEPosition(ctx context.Context, route *model.Route) {
	bearing := utils.InitialBearing(*route.Points[0], *route.Points[1])
	_ = d.ueStore.MoveToCoordinate(ctx, route.IMSI, *route.Points[0], uint32(math.Round(bearing)))
	_ = d.routeStore.Start(ctx, route.IMSI, route.SpeedAvg, route.SpeedStdDev)
}

func (d *driver) updateUEPosition(ctx context.Context, route *model.Route) {
	// Get the UE
	ue, err := d.ueStore.Get(ctx, route.IMSI)
	if err != nil {
		log.Warn("Unable to find UE %d", route.IMSI)
		return
	}

	// Determine speed and heading
	speed := float64(route.SpeedAvg) + rand.NormFloat64()*float64(route.SpeedStdDev)
	distanceDriven := (tickFrequency * speed) / 3600.0

	// Determine bearing and distance to the next point
	bearing := utils.InitialBearing(ue.Location, *route.Points[route.NextPoint])
	remainingDistance := utils.Distance(ue.Location, *route.Points[route.NextPoint])

	// If distance is less than to the next waypoint, determine the coordinate along that vector
	// Otherwise just use the next waypoint
	newPoint := *route.Points[route.NextPoint]
	reachedWaypoint := remainingDistance <= distanceDriven
	if d.wayPointRoute {
		reachedWaypoint = true
	}
	if !reachedWaypoint {
		newPoint = utils.TargetPoint(ue.Location, bearing, distanceDriven)
	}

	// Move the UE to the determined coordinate; update heading if necessary
	err = d.ueStore.MoveToCoordinate(ctx, route.IMSI, newPoint, uint32(math.Round(bearing)))
	if err != nil {
		log.Warn("Unable to update UE %d coordinates", route.IMSI)
	}

	// Update the route if necessary
	if reachedWaypoint {
		_ = d.routeStore.Advance(ctx, route.IMSI)
	}
}

func (d *driver) reportMeasurement(ctx context.Context, imsi types.IMSI) {
	ue, err := d.ueStore.Get(ctx, imsi)
	if err != nil {
		log.Warn("Unable to find UE %d", imsi)
		return
	}

	// Skip reporting measurement for IDLE UE
	if ue.RrcState == e2sm_mho.Rrcstatus_RRCSTATUS_IDLE {
		return
	}

	d.measCtrl.GetInputChan() <- ue
}

func (d *driver) linkMeasCtrlHoCtrl() {
	log.Info("Connecting measurement and handover controllers")
	for report := range d.measCtrl.GetOutputChan() {
		d.hoCtrl.GetInputChan() <- report
	}
}

func (d *driver) processHandoverDecision(ctx context.Context) {
	log.Info("Handover decision process starting")
	for {
		select {
		case hoDecision := <-d.hoCtrl.GetOutputChan():
			log.Debugf("Received HO Decision: %v", hoDecision)
			imsi := hoDecision.UE.GetID().GetID().(id.UEID).IMSI
			tCellcgi := hoDecision.TargetCell.GetID().GetID().(id.ECGI)
			tCell := &model.UECell{
				ID:   types.GnbID(tCellcgi),
				NCGI: types.NCGI(tCellcgi),
			}
			d.Handover(ctx, types.IMSI(imsi), tCell)
		case <-d.stopLocalHO:
			log.Info("local HO stopped")
			return
		}
	}
}

// Handover handovers ue to target cell
func (d *driver) Handover(ctx context.Context, imsi types.IMSI, tCell *model.UECell) {
	log.Infof("Handover() imsi:%v, tCell:%v", imsi, tCell)
	d.lockUE(imsi)
	defer d.unlockUE(imsi)

	// Update RRC state on handover
	ue, err := d.ueStore.Get(ctx, imsi)
	if err != nil {
		log.Warn("Unable to find UE %d", imsi)
		return
	}

	if ue.Cell.NCGI == tCell.NCGI {
		log.Infof("Duplicate HO skipped imsi%d, cgi:%v", imsi, tCell.NCGI)
		return
	}

	if ue.RrcState != e2sm_mho.Rrcstatus_RRCSTATUS_CONNECTED {
		//d.cellStore.DecrementRrcIdleCount(ctx, ue.Cell.NCGI)
		//d.cellStore.IncrementRrcIdleCount(ctx, tCell.NCGI)
		log.Warnf("HO skipped for not connected UE %d", imsi)
		return
	}

	d.cellStore.DecrementRrcConnectedCount(ctx, ue.Cell.NCGI)
	d.cellStore.IncrementRrcConnectedCount(ctx, tCell.NCGI)

	err = d.ueStore.UpdateCell(ctx, imsi, tCell)
	if err != nil {
		log.Warn("Unable to update UE %d cell info", imsi)
	}

	// after changing serving cell, calculate channel quality/signal strength again
	d.updateUESignalStrength(ctx, imsi)

	// update the maximum number of UEs
	d.ueStore.UpdateMaxUEsPerCell(ctx)

	log.Infof("HO is done successfully: %v to %v", imsi, tCell)
}

// UpdateUESignalStrength updates UE signal strength
func (d *driver) updateUESignalStrength(ctx context.Context, imsi types.IMSI) {
	ue, err := d.ueStore.Get(ctx, imsi)
	if err != nil {
		log.Warn("Unable to find UE %d", imsi)
		return
	}

	// update RSRP from serving cell
	err = d.updateUESignalStrengthServCell(ctx, ue)
	if err != nil {
		log.Warnf("For UE %v: %v", *ue, err)
		return
	}

	// update RSRP from candidate serving cells
	err = d.updateUESignalStrengthCandServCells(ctx, ue)
	if err != nil {
		log.Warnf("For UE %v: %v", *ue, err)
		return
	}
}

// UpdateUESignalStrengthCandServCells updates UE signal strength for serving and candidate cells
func (d *driver) updateUESignalStrengthCandServCells(ctx context.Context, ue *model.UE) error {
	cellList, err := d.cellStore.List(ctx)
	if err != nil {
		return fmt.Errorf("Unable to get all cells")
	}
	var csCellList []*model.UECell
	for _, cell := range cellList {
		rsrp := StrengthAtLocation(ue.Location, *cell)
		if math.IsInf(rsrp, 0) {
			rsrp = 0
		}
		if math.IsNaN(rsrp) {
			continue
		}
		if ue.Cell.NCGI == cell.NCGI {
			continue
		}
		ueCell := &model.UECell{
			ID:       types.GnbID(cell.NCGI),
			NCGI:     cell.NCGI,
			Strength: rsrp,
		}
		csCellList = d.sortUECells(append(csCellList, ueCell), 3) // hardcoded: to be parameterized for the future
	}
	err = d.ueStore.UpdateCells(ctx, ue.IMSI, csCellList)
	if err != nil {
		log.Warn("Unable to update UE %d cells info", ue.IMSI)
	}

	return nil
}

// UpdateUESignalStrengthServCell  updates UE signal strength for serving cell
func (d *driver) updateUESignalStrengthServCell(ctx context.Context, ue *model.UE) error {
	sCell, err := d.cellStore.Get(ctx, ue.Cell.NCGI)
	if err != nil {
		return fmt.Errorf("Unable to find serving cell %d", ue.Cell.NCGI)
	}

	strength := StrengthAtLocation(ue.Location, *sCell)

	if math.IsNaN(strength) {
		strength = -999
	}
	if math.IsInf(strength, 0) {
		strength = 0
	}

	newUECell := &model.UECell{
		ID:       ue.Cell.ID,
		NCGI:     ue.Cell.NCGI,
		Strength: strength,
	}

	err = d.ueStore.UpdateCell(ctx, ue.IMSI, newUECell)
	if err != nil {
		log.Warn("Unable to update UE %d cell info", ue.IMSI)
	}
	return nil
}

// SortUECells sorts ue cells
func (d *driver) sortUECells(ueCells []*model.UECell, numAdjCells int) []*model.UECell {
	// bubble sort
	for i := 0; i < len(ueCells)-1; i++ {
		for j := 0; j < len(ueCells)-i-1; j++ {
			if ueCells[j].Strength < ueCells[j+1].Strength {
				ueCells[j], ueCells[j+1] = ueCells[j+1], ueCells[j]
			}
		}
	}
	if len(ueCells) >= numAdjCells {
		return ueCells[0:numAdjCells]
	}
	return ueCells
}

// GetHoLogic returns the HO Logic ("local" or "mho")
func (d *driver) GetHoLogic() string {
	return d.hoLogic
}
