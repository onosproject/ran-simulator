// SPDX-FileCopyrightText: 2020-present Open Networking Foundation <info@opennetworking.org>
//
// SPDX-License-Identifier: LicenseRef-ONF-Member-1.0

package mobility

import (
	"context"
	"math"
	"math/rand"
	"time"

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

var log = logging.GetLogger("mobility", "driver")

const (
	// Probability [0.0,1.0) of a UE changing its RRC state. The higher this
	// value, the more number of UE RRC state changes will be observed
	probabilityOfRrcStateChange = 0.2
)

// Driver is an abstraction of an entity driving the UE mobility
type Driver interface {
	// Start starts the driving engine
	Start(ctx context.Context)

	// Stop stops the driving engine
	Stop()

	// GenerateRoutes generates routes for all UEs that currently do not have a route; remove routes with no UEs
	GenerateRoutes(ctx context.Context, minSpeed uint32, maxSpeed uint32, speedStdDev uint32)

	// GetMeasCtrl returns the Measurement Controller
	GetMeasCtrl() measurement.MeasController

	// GetRrcCtrl returns the Rrc Controller
	GetRrcCtrl() RrcCtrl
}

type driver struct {
	cellStore  cells.Store
	routeStore routes.Store
	ueStore    ues.Store
	apiKey     string
	ticker     *time.Ticker
	done       chan bool
	min        *model.Coordinate
	max        *model.Coordinate
	measCtrl   measurement.MeasController
	hoCtrl     handover.HOController
	hoLogic    string
	rrcCtrl    RrcCtrl
}

// NewMobilityDriver returns a driving engine capable of "driving" UEs along pre-specified routes
func NewMobilityDriver(cellStore cells.Store, routeStore routes.Store, ueStore ues.Store, apiKey string, hoLogic string) Driver {
	return &driver{
		cellStore:  cellStore,
		routeStore: routeStore,
		ueStore:    ueStore,
		hoLogic:    hoLogic,
	}
}

var tickUnit = time.Second

const tickFrequency = 1

const measType = "EventA3" // ToDo: should be programmable

const hoType = "A3" // ToDo: should be programmable

func (d *driver) Start(ctx context.Context) {
	log.Info("Driver starting")

	rand.Seed(time.Now().UnixNano())

	// Iterate over all routes and position the UEs at the start of their routes
	for _, route := range d.routeStore.List(ctx) {
		d.initializeUEPosition(ctx, route)
	}

	d.ticker = time.NewTicker(tickFrequency * tickUnit)
	d.done = make(chan bool)

	// Add measController
	d.measCtrl = measurement.NewMeasController(measType, d.cellStore, d.ueStore)
	d.measCtrl.Start(ctx)

	d.rrcCtrl.RrcUpdateChan = make(chan model.UE)

	// Add hoController
	if d.hoLogic == "local" {
		log.Info("HO logic is running locally")
		d.hoCtrl = handover.NewHOController(hoType, d.cellStore, d.ueStore)
		d.hoCtrl.Start(ctx)
		// link measController with hoController
		go d.linkMeasCtrlHoCtrl()
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

func (d *driver) GetRrcCtrl() RrcCtrl {
	return d.rrcCtrl
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
				if route.NextPoint == 0 && !route.Reverse {
					d.initializeUEPosition(ctx, route)
				}
				d.updateUEPosition(ctx, route)
				UpdateUESignalStrength(ctx, route.IMSI, d.ueStore, d.cellStore)
				d.reportMeasurement(ctx, route.IMSI)
				d.updateRrc(ctx, route.IMSI, probabilityOfRrcStateChange)
			}
		}
	}
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
	d.measCtrl.GetInputChan() <- ue
}

func (d *driver) linkMeasCtrlHoCtrl() {
	log.Info("Connecting measurement and handover controllers")
	for report := range d.measCtrl.GetOutputChan() {
		d.hoCtrl.GetInputChan() <- report
	}
	log.Info("Measurement and handover controllers disconnected")
}

func (d *driver) processHandoverDecision(ctx context.Context) {
	log.Info("Handover decision process starting")
	for hoDecision := range d.hoCtrl.GetOutputChan() {
		log.Debugf("Received HO Decision: %v", hoDecision)
		imsi := hoDecision.UE.GetID().GetID().(id.UEID).IMSI
		tCellcgi := hoDecision.TargetCell.GetID().GetID().(id.ECGI)
		tCell := &model.UECell{
			ID:   types.GnbID(tCellcgi),
			NCGI: types.NCGI(tCellcgi),
		}
		DoHandover(ctx, types.IMSI(imsi), tCell, d.ueStore, d.cellStore)
	}
	log.Info("HO decision process stopped")
}
