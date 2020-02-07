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
	"github.com/onosproject/ran-simulator/api/types"
	"googlemaps.github.io/maps"
	log "k8s.io/klog"
	"math"
	"math/rand"
	"net/http"
	"time"
)

const googleAPIKeyMinLen = 38
const stepsPerDecimalDegree = 500

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

// If a googleApiKey is given, them call the Google Directions API to get steps that
// follow known streets and traffic rules
func (m *Manager) newRoute(startLoc *Location, rID int, apiKey string, color string) (*types.Route, error) {
	endLoc, err := m.getRandomLocation(startLoc.Name)
	if err != nil {
		return nil, err
	}

	var points []*types.Point
	if len(apiKey) >= googleAPIKeyMinLen {
		points, err = googleRoute(startLoc, endLoc, apiKey)
		log.Infof("Generated new Route-%d with %d points using Google Directions", rID, len(points))
	} else {
		points, err = randomRoute(startLoc, endLoc)
		log.Infof("Generated new Route-%d with %d points using Random Directions", rID, len(points))
	}
	if err != nil {
		return nil, err
	}

	name := fmt.Sprintf("Route-%d", rID)
	route := types.Route{
		Name:      name,
		Waypoints: points,
		Color:     color,
	}

	return &route, nil
}

func googleRoute(startLoc *Location, endLoc *Location, apiKey string) ([]*types.Point, error) {
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
		for _, ll := range latLngs {
			point := types.Point{
				Lat: float32(ll.Lat),
				Lng: float32(ll.Lng),
			}
			points = append(points, &point)
		}
	}
	return points, nil
}

// If the no google API Key is given, we cannot access the Directions API and so
// generate a route randomly - it will not follow the streets
// Warning - this is a simple calculator that expect points to be on the same hemisphere
func randomRoute(startLoc *Location, endLoc *Location) ([]*types.Point, error) {
	routeWidth := float64(endLoc.Position.GetLng() - startLoc.Position.GetLng())
	routeHeight := float64(endLoc.Position.GetLat() - startLoc.Position.GetLat())

	directLength := math.Hypot(routeWidth, routeHeight)
	// Try to have a step evey 1/stepsPerDecimalDegree of a decimal degree
	points := make([]*types.Point, int(math.Floor(directLength*stepsPerDecimalDegree)))

	for i := range points {
		randFactor := (rand.Float64() - 0.5) / stepsPerDecimalDegree
		if i == 0 {
			randFactor = 0.0
		}
		deltaX := float32(routeWidth*float64(i)/float64(len(points)) + randFactor)
		deltaY := float32(routeHeight*float64(i)/float64(len(points)) + randFactor)

		points[i] = &types.Point{
			Lng: startLoc.Position.GetLng() + deltaX,
			Lat: startLoc.Position.GetLat() + deltaY,
		}
	}
	points = append(points, &endLoc.Position)

	return points, nil
}
