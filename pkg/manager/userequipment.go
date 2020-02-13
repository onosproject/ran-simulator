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
	log "k8s.io/klog"
)

func (m *Manager) newUserEquipments(params RoutesParams) map[string]*types.Ue {
	ues := make(map[string]*types.Ue)

	// There is already a route per UE
	for u := 0; u < params.NumRoutes; u++ {
		name := fmt.Sprintf("Ue-%04X", u+1)
		routeName := fmt.Sprintf("Route-%d", u)
		route := m.Routes[routeName]
		towers, distances := m.findClosestTowers(route.Waypoints[0])

		ue := types.Ue{
			Name:       name,
			Type:       "Car",
			Position:   route.Waypoints[0],
			Rotation:   0,
			Route:      routeName,
			Tower:      towers[0],
			Tower2:     towers[1],
			Tower2Dist: distances[1],
			Tower3:     towers[2],
			Tower3Dist: distances[2],
		}
		ues[name] = &ue

		// Now would be a good time to update the Route colour
		for _, t := range m.Towers {
			if t.Name == towers[0] {
				m.Routes[routeName].Color = t.Color
				break
			}
		}
	}
	return ues
}

// GetUe returns Ue based on its name
func (m *Manager) GetUe(name string) *types.Ue {
	return m.UserEquipments[name]
}

// UeHandover perform the handover on simulated UE
func (m *Manager) UeHandover(name string, tower string) {
	ue := m.UserEquipments[name]
	ue.Tower = tower
	names, _ := m.findClosestTowers(ue.Position)
	ue.Tower2 = names[1]
	ue.Tower3 = names[2]
	m.UeChannel <- dispatcher.Event{
		Type:       trafficsim.Type_UPDATED,
		UpdateType: trafficsim.UpdateType_HANDOVER,
		Object:     ue,
	}
}

func (m *Manager) startMoving(params RoutesParams) {

	for {
		breakout := false // Needed to breakout of double for loop
		for ueidx := 0; ueidx < params.NumRoutes; ueidx++ {
			ueName := fmt.Sprintf("Ue-%04X", ueidx+1)
			routeName := fmt.Sprintf("Route-%d", ueidx)
			err := m.moveUe(m.UserEquipments[ueName], m.Routes[routeName])
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
				m.UserEquipments[ueName].Position = newRoute.GetWaypoints()[0]
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
	log.Warningf("Stopped driving")
}

func (m *Manager) getColorForUe(ueName string) string {
	ue, ok := m.UserEquipments[ueName]
	if !ok {
		return ""
	}
	for _, t := range m.Towers {
		if t.Name == ue.Tower {
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
				return fmt.Errorf("end of route %s %d", route.GetName(), idx)
			}
			ue.Position = route.Waypoints[idx+1]
			ue.Rotation = uint32(getRotationDegrees(route.Waypoints[idx], route.Waypoints[idx+1]) + 180)
			names, distances := m.findClosestTowers(ue.Position)
			updateType := trafficsim.UpdateType_POSITION
			oldTower2 := ue.Tower2
			oldTower3 := ue.Tower3
			if ue.Tower == names[0] {
				ue.Tower2 = names[1]
				ue.Tower3 = names[2]
				ue.TowerDist = distances[0]
				ue.Tower2Dist = distances[1]
				ue.Tower3Dist = distances[2]
			} else if ue.Tower == names[1] {
				ue.Tower2 = names[0]
				ue.Tower3 = names[2]
				ue.TowerDist = distances[1]
				ue.Tower2Dist = distances[0]
				ue.Tower3Dist = distances[2]
			} else if ue.Tower == names[2] {
				ue.Tower2 = names[0]
				ue.Tower3 = names[1]
				ue.TowerDist = distances[2]
				ue.Tower2Dist = distances[0]
				ue.Tower3Dist = distances[1]
			} else {
				ue.Tower2 = names[0]
				ue.Tower3 = names[1]
				ue.TowerDist = distanceToTower(m.Towers[ue.Tower], ue.Position)
				ue.Tower2Dist = distances[0]
				ue.Tower3Dist = distances[1]
			}
			if ue.Tower2 != oldTower2 || ue.Tower3 != oldTower3 {
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
	return fmt.Errorf("unexpectedly hit end of route %s %v %v", route.GetName(), ue.Position, route.GetWaypoints()[0])
}
