// SPDX-FileCopyrightText: 2020-present Open Networking Foundation <info@opennetworking.org>
//
// SPDX-License-Identifier: Apache-2.0
//

package manager

import (
	"context"
	"crypto/tls"
	"fmt"
	"github.com/onosproject/ran-simulator/api/trafficsim"
	"github.com/onosproject/ran-simulator/api/types"
	"github.com/onosproject/ran-simulator/pkg/dispatcher"
	"github.com/onosproject/ran-simulator/pkg/utils"
	"googlemaps.github.io/maps"
	"math"
	"math/rand"
	"net/http"
	"time"
)

const googleAPIKeyMinLen = 38
const stepsPerDecimalDegree = 500

// RoutesParams :
type RoutesParams struct {
	APIKey    string
	StepDelay time.Duration
}

// NewRoutes Create new routes, by taking two random locations and asking Google for
// directions to get from one to the other
func (m *Manager) NewRoutes(mapLayoutParams types.MapLayout, params RoutesParams) (map[types.Imsi]*types.Route, error) {
	routes := make(map[types.Imsi]*types.Route)

	for r := 0; r < int(mapLayoutParams.MinUes); r++ {
		startLoc, err := m.getRandomLocation("")
		if err != nil {
			return nil, err
		}
		// Colour is dependent on UE tower and is not known at this stage
		route, err := m.newRoute(startLoc, utils.ImsiGenerator(r), params.APIKey, defaultColor)
		if err != nil {
			return nil, err
		}

		routes[route.RouteID] = route
	}

	return routes, nil
}

func (m *Manager) removeRoute(routeName types.Imsi) {
	r, ok := m.Routes[routeName]
	if ok {
		delete(m.Routes, routeName)
		m.RouteChannel <- dispatcher.Event{
			Type:   trafficsim.Type_REMOVED,
			Object: r,
		}
	}
}

// If a googleAPIKey is given, them call the Google Directions API to get steps that
// follow known streets and traffic rules
func (m *Manager) newRoute(startLoc *Location, rID types.Imsi, apiKey string, color string) (*types.Route, error) {
	endLoc, err := m.getRandomLocation(startLoc.Name)
	if err != nil {
		return nil, err
	}

	var points []*types.Point
	if len(apiKey) >= googleAPIKeyMinLen {
		points, err = googleRoute(startLoc, endLoc, apiKey)
		log.Infof("Generated new Route %d with %d points using Google Directions", rID, len(points))
	} else {
		points, err = randomRoute(startLoc, endLoc)
		log.Infof("Generated new Route %d with %d points using Random Directions", rID, len(points))
	}
	if err != nil {
		return nil, err
	}

	route := types.Route{
		RouteID:   rID,
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
				Lat: ll.Lat,
				Lng: ll.Lng,
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
	routeWidth := endLoc.Position.GetLng() - startLoc.Position.GetLng()
	routeHeight := endLoc.Position.GetLat() - startLoc.Position.GetLat()

	directLength := math.Hypot(routeWidth, routeHeight)
	// Try to have a step evey 1/stepsPerDecimalDegree of a decimal degree
	points := make([]*types.Point, int(math.Floor(directLength*stepsPerDecimalDegree)))

	for i := range points {
		randFactor := (rand.Float64() - 0.5) / stepsPerDecimalDegree
		if i == 0 {
			randFactor = 0.0
		}
		deltaX := routeWidth*float64(i)/float64(len(points)) + randFactor
		deltaY := routeHeight*float64(i)/float64(len(points)) + randFactor

		points[i] = &types.Point{
			Lng: startLoc.Position.GetLng() + deltaX,
			Lat: startLoc.Position.GetLat() + deltaY,
		}
	}
	points = append(points, &endLoc.Position)

	return points, nil
}
