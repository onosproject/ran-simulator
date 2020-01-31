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
	"context"
	"crypto/tls"
	"fmt"
	"github.com/OpenNetworkingFoundation/gmap-ran/api/types"
	"googlemaps.github.io/maps"
	log "k8s.io/klog"
	"net/http"
	"time"
)

// RoutesParams :
type RoutesParams struct {
	NumRoutes int
	APIKey    string
	StepDelay time.Duration
}

// Create new routes, by taking two random locations and asking Google for
// directions to get from one to the other
func (m *Manager) newRoutes(params RoutesParams) (map[string]*types.Route, error) {
	routes := make(map[string]*types.Route)

	for r := 0; r < params.NumRoutes; r++ {
		startLoc, err := m.getRandomLocation("")
		if err != nil {
			return nil, err
		}
		// Colour is dependent on UE tower and is not known at this stage
		route, err := m.newRoute(startLoc, r, params.APIKey, "#000000")
		if err != nil {
			return nil, err
		}

		routes[route.Name] = route
	}

	return routes, nil
}

func (m *Manager) newRoute(startLoc *Location, r int, apiKey string, color string) (*types.Route, error) {
	endLoc, err := m.getRandomLocation(startLoc.Name)
	if err != nil {
		return nil, err
	}

	cfg := &tls.Config{
		InsecureSkipVerify: true,
	}
	transport := &http.Transport{
		TLSClientConfig: cfg,
	}
	client := &http.Client{Transport: transport}
	googleMapsClient, err := maps.NewClient(maps.WithAPIKey(apiKey), maps.WithHTTPClient(client))
	if err != nil {
		return nil, err
	}

	dirReq := &maps.DirectionsRequest{
		Origin:      fmt.Sprintf("%f,%f", startLoc.Position.Lat, startLoc.Position.Lng),
		Destination: fmt.Sprintf("%f,%f", endLoc.Position.Lat, endLoc.Position.Lng),
	}

	googleRoute, _, err := googleMapsClient.Directions(context.Background(), dirReq)
	if err != nil {
		return nil, err
	}
	points := make([]*types.Point, 0)
	for _, groute := range googleRoute {
		latLngs, err := groute.OverviewPolyline.Decode()
		if err != nil {
			return nil, err
		}
		log.Infof("Route-%d is Google route %s #Points %d", r, groute.Summary, len(latLngs))
		for _, ll := range latLngs {
			point := types.Point{
				Lat: float32(ll.Lat),
				Lng: float32(ll.Lng),
			}
			points = append(points, &point)
		}
	}

	name := fmt.Sprintf("Route-%d", r)
	route := types.Route{
		Name:      name,
		Waypoints: points,
		Color:     color,
	}

	return &route, nil
}
