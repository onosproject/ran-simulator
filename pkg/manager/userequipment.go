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
	"strings"
	"time"

	"github.com/onosproject/ran-simulator/api/trafficsim"
	"github.com/onosproject/ran-simulator/api/types"
	"github.com/onosproject/ran-simulator/pkg/dispatcher"
)

// NewUserEquipments - create a new set of UEs (phone, car etc)
func (m *Manager) NewUserEquipments(mapLayoutParams types.MapLayout, params RoutesParams) map[types.UEName]*types.Ue {
	ues := make(map[types.UEName]*types.Ue)

	// There is already a route per UE
	var u uint32
	for u = 0; u < mapLayoutParams.MinUes; u++ {
		ue := m.newUe(int(u))
		ues[ue.Name] = ue
	}
	return ues
}

func (m *Manager) newUe(ueIdx int) *types.Ue {
	name := types.UEName(fmt.Sprintf("Ue-%04X", ueIdx+1))
	routeName := types.RouteID(fmt.Sprintf("Route-%d", ueIdx))
	route := m.Routes[routeName]
	towers, distances := m.findClosestTowers(route.Waypoints[0])
	m.TowersLock.RLock()
	servingTowerDist := distanceToTower(m.Towers[towers[0]], route.Waypoints[0])
	m.TowersLock.RUnlock()
	ue := &types.Ue{
		Name:             name,
		Type:             "Car",
		Position:         route.Waypoints[0],
		Rotation:         0,
		Route:            routeName,
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
		},
	}

	crnti, _ := m.NewCrnti(ue.ServingTower, ue.Name)
	ue.Crnti = crnti

	// Now would be a good time to update the Route colour
	for _, t := range m.Towers {
		if t.EcID == towers[0] {
			m.Routes[routeName].Color = t.Color
			break
		}
	}

	return ue
}

// SetNumberUes - change the number of active UEs
func (m *Manager) SetNumberUes(numUes int) error {
	mgr.UserEquipmentsLock.Lock()
	defer mgr.UserEquipmentsLock.Unlock()

	currentNum := len(m.UserEquipments)
	if numUes < int(m.MapLayout.MinUes) {
		return fmt.Errorf("number of UEs requested %d is below minimum %d", numUes, m.MapLayout.MinUes)
	} else if numUes > int(m.MapLayout.MaxUes) {
		return fmt.Errorf("number of UEs requested %d is above maximum %d", numUes, m.MapLayout.MaxUes)
	} else if numUes < currentNum {
		log.Infof("Decreasing number of UEs from %d to %d", currentNum, numUes)
		for ueidx := currentNum - 1; ueidx >= numUes; ueidx-- {
			ueName := ueName(ueidx)
			routeName := routeName(ueidx)
			log.Infof("Removing Route %s, UE %s", routeName, ueName)
			m.removeRoute(routeName)
			ue, ok := m.UserEquipments[ueName]
			if !ok {
				return fmt.Errorf("error removing UE %s (%d)", ueName, ueidx)
			}
			delete(m.UserEquipments, ueName)
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
			newRoute, err := m.newRoute(startLoc, ueidx, m.googleAPIKey, defaultColor)
			if err != nil {
				return err
			}
			m.Routes[newRoute.GetRouteID()] = newRoute
			m.RouteChannel <- dispatcher.Event{
				Type:   trafficsim.Type_ADDED,
				Object: newRoute,
			}
			ue := m.newUe(ueidx)
			m.UserEquipments[ue.GetName()] = ue
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
func (m *Manager) GetUe(name types.UEName) (*types.Ue, error) {
	m.UserEquipmentsLock.RLock()
	defer m.UserEquipmentsLock.RUnlock()
	ue, ok := m.UserEquipments[name]
	if !ok {
		return nil, fmt.Errorf("ue %s not found", name)
	}
	return ue, nil
}

// UeHandover perform the handover on simulated UE
func (m *Manager) UeHandover(name types.UEName, tower types.EcID) {
	ue, err := m.GetUe(name)
	if err != nil {
		log.Error(err)
		return
	}
	names, _ := m.findClosestTowers(ue.Position)
	err = m.DelCrnti(ue.ServingTower, ue.Crnti)
	if err != nil {
		log.Errorf(err.Error())
		return
	}
	m.UserEquipmentsLock.Lock()
	ue.Crnti = InvalidCrnti
	ue.ServingTower = tower
	newCrnti, err := m.NewCrnti(tower, ue.Name)
	if err != nil {
		m.UserEquipmentsLock.Unlock()
		log.Errorf(err.Error())
		return
	}
	ue.Crnti = newCrnti
	ue.Tower1 = names[0]
	ue.Tower2 = names[1]
	ue.Tower3 = names[2]
	m.UserEquipmentsLock.Unlock()
	m.UeChannel <- dispatcher.Event{
		Type:       trafficsim.Type_UPDATED,
		UpdateType: trafficsim.UpdateType_HANDOVER,
		Object:     ue,
	}
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
		for ueidx := 0; ueidx < len(m.Routes); ueidx++ {
			ueName := ueName(ueidx)
			routeName := routeName(ueidx)
			ue, err := m.GetUe(ueName)
			if err != nil {
				log.Errorf(err.Error())
				continue
			}
			err = m.moveUe(ue, m.Routes[routeName])
			if err != nil && strings.HasPrefix(err.Error(), "end of route") {
				oldRouteFinish := m.Routes[routeName].GetWaypoints()[len(m.Routes[routeName].GetWaypoints())-1]
				log.Errorf("Need to do a new route for %s Start %v %v", ueName, oldRouteFinish, err)
				newRoute, err := m.newRoute(&Location{
					Name:     "noname",
					Position: *oldRouteFinish,
				}, ueidx, params.APIKey, m.getColorForUe(ueName))
				if err != nil {
					log.Fatalf("Error %s", err.Error())
					breakout = true
				}
				m.Routes[routeName] = newRoute
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
		time.Sleep(params.StepDelay)
		if breakout {
			break
		}
	}
	log.Warnf("Stopped driving")
}

func (m *Manager) getColorForUe(ueName types.UEName) string {
	ue, ok := m.UserEquipments[ueName]
	if !ok {
		return ""
	}
	for _, t := range m.Towers {
		if t.EcID == ue.ServingTower {
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
				return fmt.Errorf("end of route %s %d", route.GetRouteID(), idx)
			}
			ue.Position = route.Waypoints[idx+1]
			ue.Rotation = uint32(getRotationDegrees(route.Waypoints[idx], route.Waypoints[idx+1]) + 180)
			names, distances := m.findClosestTowers(ue.Position)
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
	return fmt.Errorf("unexpectedly hit end of route %s %v %v", route.GetRouteID(), ue.Position, route.GetWaypoints()[0])
}

func ueName(idx int) types.UEName {
	return types.UEName(fmt.Sprintf("Ue-%04X", idx+1))
}

// UeDeepCopy ...
func UeDeepCopy(original *types.Ue) *types.Ue {
	return &types.Ue{
		Name: original.GetName(),
		Type: original.GetType(),
		Position: &types.Point{
			Lat: original.GetPosition().GetLat(),
			Lng: original.GetPosition().GetLng(),
		},
		Rotation:         original.GetRotation(),
		Route:            original.GetRoute(),
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
