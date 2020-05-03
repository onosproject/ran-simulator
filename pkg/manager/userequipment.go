// Copyright 2020-present Open Networking Foundation.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package manager

import (
	"fmt"
	"github.com/onosproject/ran-simulator/pkg/utils"
	"strings"
	"time"

	"github.com/onosproject/ran-simulator/api/trafficsim"
	"github.com/onosproject/ran-simulator/api/types"
	"github.com/onosproject/ran-simulator/pkg/dispatcher"
)

// NewUserEquipments - create a new set of UEs (phone, car etc)
func (m *Manager) NewUserEquipments(mapLayoutParams types.MapLayout, params RoutesParams) (map[types.Imsi]*types.Ue, error) {
	ues := make(map[types.Imsi]*types.Ue)

	// There is already a route per UE
	var u uint32
	for u = 0; u < mapLayoutParams.MinUes; u++ {
		ue, err := m.newUe(int(u))
		if err != nil {
			return nil, err
		}
		ues[ue.Imsi] = ue
	}
	return ues, nil
}

func (m *Manager) newUe(ueIdx int) (*types.Ue, error) {
	imsi := utils.ImsiGenerator(ueIdx)
	route := m.Routes[imsi]
	towers, distances, err := m.findClosestCells(route.Waypoints[0])
	if err != nil {
		return nil, err
	}
	m.CellsLock.RLock()
	servingTowerDist, err := distanceToCellCentroid(m.Cells[*towers[0]], route.Waypoints[0])
	if err != nil {
		m.CellsLock.RUnlock()
		return nil, err
	}
	m.CellsLock.RUnlock()
	ue := &types.Ue{
		Imsi:             imsi,
		Type:             "Car",
		Position:         route.Waypoints[0],
		Rotation:         0,
		ServingTower:     towers[0],
		ServingTowerDist: servingTowerDist,
		Tower1:           towers[0],
		Tower1Dist:       distances[0],
		Tower2:           towers[1],
		Tower2Dist:       distances[1],
		Tower3:           towers[2],
		Tower3Dist:       distances[2],
		Admitted:         false,
		Crnti:            InvalidCrnti,
		Metrics: &types.UeMetrics{
			HoLatency:         0,
			HoReportTimestamp: 0,
			IsFirst:           true,
		},
	}

	crnti, _ := m.NewCrnti(ue.ServingTower, ue.Imsi)
	ue.Crnti = crnti

	// Now would be a good time to update the Route colour
	for _, t := range m.Cells {
		if t.Ecgi == towers[0] {
			m.Routes[imsi].Color = t.Color
			break
		}
	}

	return ue, nil
}

// SetNumberUes - change the number of active UEs
func (m *Manager) SetNumberUes(numUes int) error {
	m.UserEquipmentsMapLock.Lock()
	defer m.UserEquipmentsMapLock.Unlock()
	currentNum := len(m.UserEquipments)
	if numUes < currentNum {
		log.Infof("Decreasing number of UEs from %d to %d", currentNum, numUes)
		for ueidx := currentNum - 1; ueidx >= numUes; ueidx-- {
			imsi := utils.ImsiGenerator(ueidx)
			log.Infof("Removing Route and UE %d", imsi)
			mgr.UserEquipmentsLock.Lock()
			m.removeRoute(imsi)
			ue, ok := m.UserEquipments[imsi]
			if !ok {
				mgr.UserEquipmentsLock.Unlock()
				return fmt.Errorf("error removing UE %d (%d)", imsi, ueidx)
			}
			delete(m.UserEquipments, imsi)
			mgr.UserEquipmentsLock.Unlock()
			m.UeChannel <- dispatcher.Event{
				Type:   trafficsim.Type_REMOVED,
				Object: ue,
			}
		}
	} else {
		log.Infof("Increasing number of UEs from %d to %d", currentNum, numUes)
		for ueidx := currentNum; ueidx < numUes; ueidx++ {
			startLoc, err := m.getRandomLocation("")
			if err != nil {
				return err
			}
			imsi := utils.ImsiGenerator(ueidx)
			newRoute, err := m.newRoute(startLoc, imsi, m.googleAPIKey, defaultColor)
			if err != nil {
				return err
			}
			m.UserEquipmentsLock.Lock()
			m.Routes[imsi] = newRoute
			m.UserEquipmentsLock.Unlock()
			m.RouteChannel <- dispatcher.Event{
				Type:   trafficsim.Type_ADDED,
				Object: newRoute,
			}
			ue, err := m.newUe(ueidx)
			if err != nil {
				return err
			}
			mgr.UserEquipmentsLock.Lock()
			m.UserEquipments[ue.GetImsi()] = ue
			mgr.UserEquipmentsLock.Unlock()
			m.UeChannel <- dispatcher.Event{
				Type:   trafficsim.Type_ADDED,
				Object: ue,
			}
		}
	}
	m.MapLayout.CurrentRoutes = uint32(numUes)

	return nil
}

// GetUe returns Ue based on its name
func (m *Manager) GetUe(imsi types.Imsi) (*types.Ue, error) {
	m.UserEquipmentsLock.RLock()
	defer m.UserEquipmentsLock.RUnlock()
	ue, ok := m.UserEquipments[imsi]
	if !ok {
		return nil, fmt.Errorf("ue %d not found", imsi)
	}
	return ue, nil
}

// UeHandover perform the handover on simulated UE
func (m *Manager) UeHandover(imsi types.Imsi, newTowerID *types.ECGI) error {
	ue, err := m.GetUe(imsi)
	if err != nil {
		return err
	}
	err = m.DelCrnti(ue.ServingTower, ue.Crnti)
	if err != nil {
		return err
	}
	m.UserEquipmentsLock.Lock()
	ue.ServingTower = newTowerID
	newCrnti, err := m.NewCrnti(newTowerID, ue.Imsi)
	if err != nil {
		m.UserEquipmentsLock.Unlock()
		return err
	}
	ue.Crnti = newCrnti
	ue.Admitted = false
	m.UserEquipmentsLock.Unlock()
	m.UeChannel <- dispatcher.Event{
		Type:       trafficsim.Type_UPDATED,
		UpdateType: trafficsim.UpdateType_HANDOVER,
		Object:     ue,
	}
	return nil
}

// UeAdmitted - called when the Admission Request for the UE is processed
// This causes the first RadioMeasurementReport to be sent
func (m *Manager) UeAdmitted(ue *types.Ue) {
	// The UEs should not be locked when this is sent or a deadlock will occur
	m.UeChannel <- dispatcher.Event{
		Type:       trafficsim.Type_UPDATED,
		UpdateType: trafficsim.UpdateType_TOWER,
		Object:     ue,
	}
}

func (m *Manager) startMoving(params RoutesParams) {

	for {
		breakout := false // Needed to breakout of double for loop
		m.UserEquipmentsMapLock.Lock()
		for imsi, ue := range m.UserEquipments {
			r, ok := m.Routes[imsi]
			if !ok {
				log.Warnf("Unable to retrieve route for %s", imsi)
				continue
			}
			err := m.moveUe(ue, r)
			if err != nil && strings.HasPrefix(err.Error(), "end of route") {
				oldRouteFinish := r.GetWaypoints()[len(r.GetWaypoints())-1]
				log.Infof("Need to do a new route for %d Start %v %v", imsi, oldRouteFinish, err)
				newRoute, err := m.newRoute(&Location{
					Name:     "noname",
					Position: *oldRouteFinish,
				}, imsi, params.APIKey, m.getColorForUe(imsi))
				if err != nil {
					log.Fatalf("Error %s", err.Error())
					breakout = true
				}
				m.Routes[imsi] = newRoute
				m.RouteChannel <- dispatcher.Event{
					Type:   trafficsim.Type_UPDATED,
					Object: newRoute,
				}
				// Move the UE to this new start point - google might return a
				// start point just a few metres from where we asked
				ue.Position = newRoute.GetWaypoints()[0]
			} else if err != nil {
				log.Errorf("Error %s", err.Error())
				breakout = true
			}
		}
		m.UserEquipmentsMapLock.Unlock()
		time.Sleep(params.StepDelay)
		if breakout {
			break
		}
	}
	log.Warnf("Stopped driving")
}

func (m *Manager) getColorForUe(imsi types.Imsi) string {
	ue, ok := m.UserEquipments[imsi]
	if !ok {
		return ""
	}
	for _, t := range m.Cells {
		if t.Ecgi == ue.ServingTower {
			return t.Color
		}
	}
	return ""
}

// Move the UE to a new position along its route
func (m *Manager) moveUe(ue *types.Ue, route *types.Route) error {
	for idx, wp := range route.GetWaypoints() {
		if ue.Position.GetLng() == wp.GetLng() && ue.Position.GetLat() == wp.GetLat() {
			if idx+1 == len(route.GetWaypoints()) {
				return fmt.Errorf("end of route %d %d", route.GetRouteID(), idx)
			}
			ue.Position = route.Waypoints[idx+1]
			ue.Rotation = uint32(utils.GetRotationDegrees(route.Waypoints[idx], route.Waypoints[idx+1]) + 180)
			names, distances, err := m.findClosestCells(ue.Position)
			if err != nil {
				return err
			}
			updateType := trafficsim.UpdateType_POSITION
			oldTower1 := ue.Tower1
			oldTower2 := ue.Tower2
			oldTower3 := ue.Tower3
			ue.Tower1 = names[0]
			ue.Tower1Dist = distances[0]
			ue.Tower2 = names[1]
			ue.Tower2Dist = distances[1]
			ue.Tower3 = names[2]
			ue.Tower3Dist = distances[2]
			servingTowerDist, _ := distanceToCellCentroid(m.Cells[*ue.ServingTower], ue.Position)
			ue.ServingTowerDist = servingTowerDist

			if ue.Admitted && ue.Tower1 != oldTower1 || ue.Tower2 != oldTower2 || ue.Tower3 != oldTower3 {
				updateType = trafficsim.UpdateType_TOWER
			}
			m.UeChannel <- dispatcher.Event{
				Type:       trafficsim.Type_UPDATED,
				UpdateType: updateType,
				Object:     ue,
			}
			return nil
		}
	}
	return fmt.Errorf("unexpectedly hit end of route %d %v %v", route.GetRouteID(), ue.Position, route.GetWaypoints()[0])
}

// UeDeepCopy ...
func UeDeepCopy(original *types.Ue) *types.Ue {
	return &types.Ue{
		Imsi: original.GetImsi(),
		Type: original.GetType(),
		Position: &types.Point{
			Lat: original.GetPosition().GetLat(),
			Lng: original.GetPosition().GetLng(),
		},
		Rotation:         original.GetRotation(),
		ServingTower:     original.GetServingTower(),
		ServingTowerDist: original.GetServingTowerDist(),
		Tower1:           original.GetTower1(),
		Tower1Dist:       original.GetTower1Dist(),
		Tower2:           original.GetTower2(),
		Tower2Dist:       original.GetTower2Dist(),
		Tower3:           original.GetTower3(),
		Tower3Dist:       original.GetTower3Dist(),
		Crnti:            original.GetCrnti(),
		Admitted:         original.GetAdmitted(),
		Metrics: &types.UeMetrics{
			HoLatency:         original.GetMetrics().GetHoLatency(),
			HoReportTimestamp: original.GetMetrics().GetHoReportTimestamp(),
		},
	}
}
