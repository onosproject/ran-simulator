// SPDX-FileCopyrightText: 2020-present Open Networking Foundation <info@opennetworking.org>
//
// SPDX-License-Identifier: LicenseRef-ONF-Member-1.0
//

package manager

import (
	"fmt"
	"math"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/onosproject/config-models/modelplugin/e2node-1.0.0/e2node_1_0_0"
	"github.com/onosproject/onos-topo/api/topo"
	"github.com/onosproject/ran-simulator/pkg/southbound/kubernetes"

	"github.com/onosproject/ran-simulator/api/trafficsim"
	"github.com/onosproject/ran-simulator/api/types"
	"github.com/onosproject/ran-simulator/pkg/dispatcher"
	"github.com/onosproject/ran-simulator/pkg/utils"
)

const (
	// DefaultTxPower - all cells start with this power level
	DefaultTxPower = 10

	// PowerFactor - relate power to distance in decimal degrees
	PowerFactor = 0.001

	// PowerBase - baseline for power to distance in decimal degrees
	PowerBase = 0.013

	// MaxNeighbours to find - useful limit for hex layouts
	MaxNeighbours = 6
)

const defaultColor = "#000000"

const (
	maxPowerdB = 30.0
	minPowerdB = -15.0
)

// MaxCrnti - Maximum value of CRNTI
const MaxCrnti = 65523

// InvalidCrnti ...
const InvalidCrnti = "0000"

// CellIf :
type CellIf interface {
	GetPosition() types.Point
}

// MaxNumUesPerCell is the variable to configure each Cell's maximum load - the number of maximum UEs
var MaxNumUesPerCell uint32

// CellCreator - wrap the NewCell function
func CellCreator(object *topo.Object) error {
	cell, err := NewCell(object)
	if err != nil {
		return nil
	}

	mgr := GetManager()

	mgr.CellsLock.Lock()
	defer mgr.CellsLock.Unlock()

	if _, ok := mgr.Cells[*cell.Ecgi]; ok {
		return nil
	}

	log.Infof("Cell created %s", cell.Ecgi)

	mgr.Cells[*cell.Ecgi] = cell
	mgr.CellConfigs[*cell.Ecgi] = &e2node_1_0_0.Device{
		E2Node: &e2node_1_0_0.E2Node_E2Node{
			Intervals: &e2node_1_0_0.E2Node_E2Node_Intervals{
				PdcpMeasReportPerUe:    nil,
				RadioMeasReportPerCell: nil,
				RadioMeasReportPerUe:   nil,
				SchedMeasReportPerCell: nil,
				SchedMeasReportPerUe:   nil,
			},
		}}
	mgr.cellCreateTimer.Reset(time.Second)
	// If no new cells have been created in the last second, then
	// afterCellCreation below gets run. This allows bulk changes to be handled
	// together

	return nil
}

// CellDeleter - wrap the cell deletion function
func CellDeleter(object *topo.Object) error {
	ecgi, err := ecgiFromtopo(object)
	if err != nil {
		return err
	}

	mgr := GetManager()
	// TODO stop the gRPC server for this cell

	mgr.CellsLock.Lock()
	defer mgr.CellsLock.Unlock()

	delete(mgr.CellConfigs, *ecgi)
	delete(mgr.Cells, *ecgi)

	log.Infof("Cell deleted %s", ecgi)
	return nil
}

// NewCell - create a new cell from a topo
func NewCell(object *topo.Object) (*types.Cell, error) {
	ecgi, err := ecgiFromtopo(object)
	if err != nil {
		return nil, err
	}

	var latitude float64
	if latitudeStr, ok := object.GetAttributes()[types.LatitudeKey]; ok {
		if latitude, err = strconv.ParseFloat(latitudeStr, 64); err != nil {
			return nil, err
		}
	} else {
		return nil, fmt.Errorf("device %s does not have 'latitude' in attributes", object.ID)
	}

	var longitude float64
	if longitudeStr, ok := object.GetAttributes()[types.LongitudeKey]; ok {
		if longitude, err = strconv.ParseFloat(longitudeStr, 64); err != nil {
			return nil, err
		}
	} else {
		return nil, fmt.Errorf("device %s does not have 'longitude' in attributes", object.ID)
	}
	cellLoc := types.Point{
		Lat: latitude,
		Lng: longitude,
	}

	var azimuth int64
	if azimuthStr, ok := object.GetAttributes()[types.AzimuthKey]; ok {
		if azimuth, err = strconv.ParseInt(azimuthStr, 10, 32); err != nil {
			return nil, err
		}
	} else {
		return nil, fmt.Errorf("device %s does not have 'azimuth' in attributes", object.ID)
	}

	var arc int64
	if arcStr, ok := object.GetAttributes()[types.ArcKey]; ok {
		if arc, err = strconv.ParseInt(arcStr, 10, 32); err != nil {
			return nil, err
		}
	} else {
		return nil, fmt.Errorf("device %s does not have 'arc' in attributes", object.ID)
	}

	var grpcPort int64
	if grpcPortStr, ok := object.GetAttributes()[types.GrpcPortKey]; ok {
		if grpcPort, err = strconv.ParseInt(grpcPortStr, 10, 32); err != nil {
			return nil, err
		}
	} else {
		// Try to parse it from the address
		if addressStr, ok := object.GetAttributes()[types.AddressKey]; ok {
			parts := strings.Split(addressStr, ":")
			if len(parts) != 2 {
				return nil, fmt.Errorf("cannot parse address to get port number %s", addressStr)
			}
			if grpcPort, err = strconv.ParseInt(parts[1], 10, 32); err != nil {
				return nil, err
			}
		}
	}

	cell := &types.Cell{
		Location:   &cellLoc,
		Color:      utils.RandomColor(),
		Ecgi:       ecgi,
		TxPowerdB:  DefaultTxPower,
		Port:       uint32(grpcPort),
		CrntiMap:   make(map[types.Crnti]types.Imsi),
		CrntiIndex: 0,
		MaxUEs:     MaxNumUesPerCell,
		Sector: &types.Sector{
			Azimuth: int32(azimuth),
			Arc:     int32(arc),
		},
	}
	cell.Sector.Centroid = centroidPosition(cell)

	return cell, nil
}

func (m *Manager) afterCellCreation(mapLayoutParams types.MapLayout,
	routesParams RoutesParams, serverParams utils.ServerParams,
	newServerHandler NewServerHandler) {

	for {
		<-m.cellCreateTimer.C // 1 second after cell(s) have been loaded
		if len(m.Cells) == 0 {
			continue
		}
		m.MapLayout.Center, m.Locations = NewLocations(m.Cells,
			int(mapLayoutParams.MaxUes), mapLayoutParams.LocationsScale)
		log.Infof("Cell creation post action. %d cells. Centre %f, %f",
			len(GetManager().Cells), m.MapLayout.GetCenter().GetLat(), m.MapLayout.GetCenter().GetLng())

		for _, cell := range m.Cells {
			cell := cell //pin
			go func() {
				// Blocks here when server running
				err := newServerHandler(*cell.Ecgi, uint16(cell.Port), serverParams)
				if err != nil {
					log.Error("Unable to start server ", err)
				}
			}()
		}

		if serverParams.AddK8sSvcPorts {
			allGrpcPorts := make([]uint16, 0)
			for _, cell := range m.Cells {
				cell := cell //pin
				allGrpcPorts = append(allGrpcPorts, uint16(cell.Port))
			}
			err := kubernetes.AddK8SServicePorts(allGrpcPorts)
			if err != nil {
				log.Fatal(err.Error())
			}
		}

		for _, cell := range m.Cells {
			//cell := cell //pin
			cell.Neighbors = makeNeighbors(cell, m.Cells)
		}

		// Only start after 3 cells and if not started before
		if len(m.Cells) > 3 && len(m.Routes) == 0 {
			var err error
			m.Routes, err = m.NewRoutes(mapLayoutParams, routesParams)
			if err != nil {
				log.Fatalf("Error calculating routes %s", err.Error())
			}
			m.UserEquipments, err = m.NewUserEquipments(mapLayoutParams, routesParams)
			if err != nil {
				log.Fatalf("Error creating new UEs %s", err.Error())
			}
			go m.startMoving(routesParams)
		} else if len(m.Cells) <= 3 {
			log.Warnf("Not creating UEs - must have 3 or more cells. Currently %d", len(m.Cells))
		}
	}
}

// Find the strongest power signal cell to any point - return strongest, candidate1 and candidate2
// in order of power. This is derived from
// 1. the distance of the point from cell
// 2. the arc and the azimuth of the cell
// 3. the power setting of the antenna
// Note this does not take any account of who's serving - it's just about power
// values are in dB
func (m *Manager) findStrongestCells(point *types.Point) ([]*types.ECGI, []float64, error) {
	var (
		strongest  *types.ECGI
		candidate1 *types.ECGI
		candidate2 *types.ECGI
	)

	strongestStr := -math.MaxFloat64
	candidate1Str := -math.MaxFloat64
	candidate2Str := -math.MaxFloat64

	m.CellsLock.RLock()
	for _, cell := range m.Cells {
		strength := strengthAtPoint(point, cell)

		if strength > strongestStr {
			candidate2 = candidate1
			candidate2Str = candidate1Str
			candidate1 = strongest
			candidate1Str = strongestStr
			strongest = cell.Ecgi
			strongestStr = strength
		} else if strength > candidate1Str {
			candidate2 = candidate1
			candidate2Str = candidate1Str
			candidate1 = cell.Ecgi
			candidate1Str = strength
		} else if strength > candidate2Str {
			candidate2 = cell.Ecgi
			candidate2Str = strength
		}
	}
	m.CellsLock.RUnlock()

	return []*types.ECGI{strongest, candidate1, candidate2},
		[]float64{strongestStr, candidate1Str, candidate2Str}, nil
}

func strengthAtPoint(point *types.Point, cell *types.Cell) float64 {
	distAtt := distanceAttenuation(point, cell)
	angleAtt := angleAttenuation(point, cell)

	return cell.TxPowerdB + distAtt + angleAtt
}

// distanceAttenuation is the antenna Gain as a function of the dist
// a very rough approximation to take in to account the width of
// the antenna beam. A 120° wide beam with 30° height will span ≅ 2x0.5 = 1 steradians
// A 60° wide beam will be half that and so will have double the gain
// https://en.wikipedia.org/wiki/Sector_antenna
// https://en.wikipedia.org/wiki/Steradian
func distanceAttenuation(point *types.Point, cell *types.Cell) float64 {
	latDist := point.GetLat() - cell.GetLocation().GetLat()
	realLngDist := (point.GetLng() - cell.GetLocation().GetLng()) / utils.AspectRatio(cell.GetLocation())
	r := math.Hypot(latDist, realLngDist)
	gain := 120.0 / float64(cell.GetSector().GetArc())
	return 10 * math.Log10(gain*math.Sqrt(PowerFactor/r))
}

// angleAttenuation is the attenuation of power reaching a UE due to its
// position off the centre of the beam in dB
// It is an approximation of the directivity of the antenna
// https://en.wikipedia.org/wiki/Radiation_pattern
// https://en.wikipedia.org/wiki/Sector_antenna
func angleAttenuation(point *types.Point, cell *types.Cell) float64 {

	azRads := utils.AzimuthToRads(float64(cell.Sector.Azimuth))
	pointRads := math.Atan2(point.Lat-cell.Location.Lat, point.Lng-cell.Location.Lng)
	angularOffset := math.Abs(azRads - pointRads)
	angleScaling := float64(cell.Sector.Arc) / 120.0 // Compensate for narrower beams

	// We just use a simple linear formula 0 => no loss
	// 33° => -3dB for a 120° sector according to [2]
	// assume this is 1:1 rads:attenuation e.g. 0.50 rads = 0.5 = -3dB attenuation
	return 10 * math.Log10(1-(angularOffset/math.Pi/angleScaling))
}

// GetCell returns tower based on its name
func (m *Manager) GetCell(name types.ECGI) *types.Cell {
	m.CellsLock.RLock()
	defer m.CellsLock.RUnlock()
	return m.Cells[name]
}

// UpdateCell Update a cell's properties - usually power level
func (m *Manager) UpdateCell(cell types.ECGI, powerAdjust float32) error {
	// Only the power can be updated at present
	m.CellsLock.Lock()
	c, ok := m.Cells[cell]
	if !ok {
		m.CellsLock.Unlock()
		return fmt.Errorf("unknown cell %s", cell)
	}
	currentPower := c.TxPowerdB
	if currentPower+float64(powerAdjust) < minPowerdB {
		c.TxPowerdB = minPowerdB
	} else if currentPower+float64(powerAdjust) > maxPowerdB {
		c.TxPowerdB = maxPowerdB
	} else {
		c.TxPowerdB += float64(powerAdjust)
	}
	c.GetSector().Centroid = centroidPosition(c)
	m.CellsLock.Unlock()
	m.CellsChannel <- dispatcher.Event{
		Type:   trafficsim.Type_UPDATED,
		Object: c,
	}
	return nil
}

// NewCrnti allocs a new crnti
func (m *Manager) NewCrnti(servingTower *types.ECGI, imsi types.Imsi) (types.Crnti, error) {
	m.CellsLock.Lock()
	defer m.CellsLock.Unlock()
	tower, ok := m.Cells[*servingTower]
	if !ok {
		return "", fmt.Errorf("unknown tower %s", servingTower)
	}
	tower.CrntiIndex++
	crnti := types.Crnti(fmt.Sprintf("%04X", tower.CrntiIndex%MaxCrnti))
	tower.CrntiMap[crnti] = imsi
	return crnti, nil
}

// DelCrnti deletes a crnti
func (m *Manager) DelCrnti(servingTower *types.ECGI, crnti types.Crnti) error {
	m.CellsLock.Lock()
	defer m.CellsLock.Unlock()
	tower, ok := m.Cells[*servingTower]
	if !ok {
		return fmt.Errorf("unknown tower %s", servingTower)
	}
	crntiMap := tower.CrntiMap
	delete(crntiMap, crnti)
	return nil
}

// CrntiToName ...
func (m *Manager) CrntiToName(crnti types.Crnti, ecid types.ECGI) (types.Imsi, error) {
	m.CellsLock.RLock()
	defer m.CellsLock.RUnlock()
	tower, ok := m.Cells[ecid]
	if !ok {
		return 0, fmt.Errorf("tower %s not found", ecid)
	}
	imsi, ok := tower.CrntiMap[crnti]
	if !ok {
		return 0, fmt.Errorf("crnti %s/%s not found", ecid, crnti)
	}
	return imsi, nil
}

// Measure the distance between a point and a cell centroid and return an answer in decimal degrees
// Centroid is used **only** for the display of the beam on the GUI and for
// calculating Neighbours once at startup
// Simple arithmetic is used, do not use for lat or long diff >= 100 degrees
func centroidPosition(cell *types.Cell) *types.Point {
	if cell.Sector.Arc == 360 || cell.Sector.Arc == 0 {
		return cell.Location
	}
	// Work out the location of the centroid of the cell - ref https://en.wikipedia.org/wiki/List_of_centroids
	alpha := 2 * math.Pi * float64(cell.Sector.Arc) / 360 / 2
	dist := 2 * PowerToDist(cell.TxPowerdB) * math.Sin(alpha) / alpha / 3
	var azRads float64 = 0
	if cell.Sector.Azimuth != 90 {
		azRads = math.Pi * 2 * float64(90-cell.Sector.Azimuth) / 360
	}
	aspectRatio := utils.AspectRatio(cell.Location)
	return &types.Point{
		Lat: math.Sin(azRads)*dist + cell.Location.GetLat(),
		Lng: math.Cos(azRads)*dist/aspectRatio + cell.Location.GetLng(),
	}
}

// PowerToDist - convert power in dB to distance in decimal degrees
// Like centroid this is now used only for calculating centroid, which is
// only for the GUI and the neighbours
func PowerToDist(power float64) float64 {
	return math.Sqrt(math.Pow(10, power/10))*PowerFactor + PowerBase
}

// find the neighbours of a cell - not distance from towers, but from centroids
func makeNeighbors(self *types.Cell, allCells map[types.ECGI]*types.Cell) []*types.ECGI {
	type distance struct {
		id   *types.ECGI
		dist float64
	}
	distances := make([]distance, 0)

	selfCentroid := self.Sector.Centroid
	for _, otherCell := range allCells {
		if otherCell.Ecgi == self.Ecgi {
			continue
		}
		dist := math.Hypot(
			selfCentroid.Lng-otherCell.Sector.Centroid.Lng,
			selfCentroid.Lat-otherCell.Sector.Centroid.Lat,
		)
		distances = append(distances, distance{id: otherCell.Ecgi, dist: dist})
	}
	sort.Slice(distances, func(i, j int) bool {
		return distances[i].dist < distances[j].dist
	})

	limit := len(allCells) - 1
	if limit > MaxNeighbours {
		limit = MaxNeighbours
	}
	neighbours := make([]*types.ECGI, limit)
	for i := 0; i < limit; i++ {
		neighbours[i] = distances[i].id
	}

	return neighbours
}

func ecgiFromtopo(object *topo.Object) (*types.ECGI, error) {
	var ecid types.EcID
	if ecidStr, ok := object.GetAttributes()[types.EcidKey]; ok {
		ecid = types.EcID(ecidStr)
	}
	var plmnid types.PlmnID
	if plmnidStr, ok := object.GetAttributes()[types.PlmnIDKey]; ok {
		plmnid = types.PlmnID(plmnidStr)
	}
	var ecgi types.ECGI
	var err error
	if ecid == "" || plmnid == "" { // If not found in attrs above use ID
		if ecgi, err = ecgiFromTopoID(object.ID); err != nil {
			return nil, err
		}
	} else {
		ecgi = types.ECGI{PlmnID: plmnid, EcID: ecid}
	}
	return &ecgi, nil
}

func newEcgi(id types.EcID, plmnID types.PlmnID) types.ECGI {
	return types.ECGI{EcID: id, PlmnID: plmnID}
}

// ecgiFromTopoID topo device is formatted like "315010-0001786" PlmnId-Ecid
func ecgiFromTopoID(id topo.ID) (types.ECGI, error) {
	if !strings.Contains(string(id), "-") {
		return types.ECGI{}, fmt.Errorf("unexpected format for E2Node ID %s", id)
	}
	parts := strings.Split(string(id), "-")
	if len(parts) != 2 {
		return types.ECGI{}, fmt.Errorf("unexpected format for E2Node ID %s", id)
	}
	return types.ECGI{EcID: types.EcID(parts[1]), PlmnID: types.PlmnID(parts[0])}, nil
}


// CellDeepCopyNomap - used when sending it up to GUI - no need to copy the map of Crnti or neighbours
// had been causing concurrent write errors inside gRPC processing
func CellDeepCopyNomap(original *types.Cell) *types.Cell {
	return &types.Cell{
		Ecgi: &types.ECGI{
			EcID:   original.GetEcgi().GetEcID(),
			PlmnID: original.GetEcgi().GetPlmnID(),
		},
		Location: &types.Point{
			Lat: original.GetLocation().GetLat(),
			Lng: original.GetLocation().GetLng(),
		},
		Sector: &types.Sector{
			Azimuth: original.GetSector().GetAzimuth(),
			Arc:     original.GetSector().GetArc(),
			Centroid: &types.Point{
				Lat: original.GetSector().GetCentroid().GetLat(),
				Lng: original.GetSector().GetCentroid().GetLng(),
			},
		},
		Color:      original.GetColor(),
		MaxUEs:     original.GetMaxUEs(),
		Neighbors:  nil, // not copied
		TxPowerdB:  original.GetTxPowerdB(),
		CrntiMap:   nil, // not copied
		CrntiIndex: original.GetCrntiIndex(),
		Port:       original.GetPort(),
	}
}
