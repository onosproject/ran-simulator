// SPDX-FileCopyrightText: 2021-present Open Networking Foundation <info@opennetworking.org>
//
// SPDX-License-Identifier: LicenseRef-ONF-Member-1.0
//

package mobility

import (
	"context"
	"crypto/tls"
	"fmt"
	"github.com/onosproject/onos-api/go/onos/ransim/types"
	"github.com/onosproject/ran-simulator/pkg/model"
	"github.com/onosproject/ran-simulator/pkg/utils"
	"googlemaps.github.io/maps"
	"math"
	"math/rand"
	"net/http"
)

const googleAPIKeyMinLen = 38
const stepsPerDecimalDegree = 500

func (d *driver) GenerateRoutes(ctx context.Context, minSpeed uint32, maxSpeed uint32, speedStdDev uint32) {
	// TODO: Determine the area for choosing random end-point locations
	for _, ue := range d.ueStore.ListAllUEs(ctx) {
		_, err := d.routeStore.Get(ctx, ue.IMSI)
		if err != nil {
			_ = d.generateRoute(ctx, ue.IMSI, uint32(rand.Intn(int(maxSpeed-minSpeed))), speedStdDev)
		}
	}
}

func (d *driver) generateRoute(ctx context.Context, imsi types.IMSI, speedAvg uint32, speedStdDev uint32) error {
	var err error
	start := d.randomCoordinate()
	end := d.randomCoordinate()

	var points []*model.Coordinate
	if len(d.apiKey) >= googleAPIKeyMinLen {
		points, err = googleRoute(start, end, d.apiKey)
		log.Infof("Generated new Route %d with %d points using Google Directions", imsi, len(points))
	} else {
		points, err = randomRoute(start, end)
		log.Infof("Generated new Route %d with %d points using Random Directions", imsi, len(points))
	}
	if err != nil {
		return err
	}

	route := &model.Route{
		IMSI:        imsi,
		Points:      points,
		SpeedAvg:    speedAvg,
		SpeedStdDev: speedStdDev,
		Color:       utils.RandomColor(),
	}
	return d.routeStore.Add(ctx, route)
}

func (d *driver) randomCoordinate() *model.Coordinate {
	return &model.Coordinate{
		Lat: 0,
		Lng: 0,
	}
}

func googleRoute(startLoc *model.Coordinate, endLoc *model.Coordinate, apiKey string) ([]*model.Coordinate, error) {
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
		Origin:      fmt.Sprintf("%f,%f", startLoc.Lat, startLoc.Lng),
		Destination: fmt.Sprintf("%f,%f", endLoc.Lat, endLoc.Lng),
	}

	googleRoute, _, err := googleMapsClient.Directions(context.Background(), dirReq)
	if err != nil {
		return nil, err
	}
	points := make([]*model.Coordinate, 0)
	for _, groute := range googleRoute {
		latLngs, err := groute.OverviewPolyline.Decode()
		if err != nil {
			return nil, err
		}
		for _, ll := range latLngs {
			point := model.Coordinate{
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
func randomRoute(startLoc *model.Coordinate, endLoc *model.Coordinate) ([]*model.Coordinate, error) {
	routeWidth := endLoc.Lng - startLoc.Lng
	routeHeight := endLoc.Lat - startLoc.Lat

	directLength := math.Hypot(routeWidth, routeHeight)
	// Try to have a step evey 1/stepsPerDecimalDegree of a decimal degree
	points := make([]*model.Coordinate, int(math.Floor(directLength*stepsPerDecimalDegree)))

	for i := range points {
		randFactor := (rand.Float64() - 0.5) / stepsPerDecimalDegree
		if i == 0 {
			randFactor = 0.0
		}
		deltaX := routeWidth*float64(i)/float64(len(points)) + randFactor
		deltaY := routeHeight*float64(i)/float64(len(points)) + randFactor

		points[i] = &model.Coordinate{
			Lng: startLoc.Lng + deltaX,
			Lat: startLoc.Lat + deltaY,
		}
	}
	points = append(points, endLoc)

	return points, nil
}
