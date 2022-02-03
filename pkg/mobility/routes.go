// SPDX-FileCopyrightText: 2021-present Open Networking Foundation <info@opennetworking.org>
//
// SPDX-License-Identifier: Apache-2.0
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

const latMargin = 0.04 // ~ 4.4km at equator; ~3.1km at 45
const lngMargin = 0.01 // ~ 4.4km

var routeEndPointIndex = 0

func (d *driver) GenerateRoutes(ctx context.Context, minSpeed uint32, maxSpeed uint32, speedStdDev uint32, routeEndPoints []model.RouteEndPoint, directRoute bool) {
	d.establishArea(ctx)
	log.Infof("Generating routes in area min=%v; max=%v\n", d.min, d.max)
	for _, ue := range d.ueStore.ListAllUEs(ctx) {
		_, err := d.routeStore.Get(ctx, ue.IMSI)
		if err != nil {
			err = d.generateRoute(ctx, ue.IMSI, uint32(rand.Intn(int(maxSpeed-minSpeed))), speedStdDev, routeEndPoints, directRoute)
			if err != nil {
				log.Warnf("Unable to generate route for %d, %v", ue.IMSI, err)
			}
		}
	}
}

// Determines the area for choosing random end-point locations
func (d *driver) establishArea(ctx context.Context) {
	cells, err := d.cellStore.List(ctx)
	if err != nil {
		return
	}

	d.min = &model.Coordinate{Lat: 90.0, Lng: 180.0}
	d.max = &model.Coordinate{Lat: -90.0, Lng: -180.0}
	for _, cell := range cells {
		d.min.Lat = math.Min(cell.Sector.Center.Lat, d.min.Lat)
		d.min.Lng = math.Min(cell.Sector.Center.Lng, d.min.Lng)
		d.max.Lat = math.Max(cell.Sector.Center.Lat, d.max.Lat)
		d.max.Lng = math.Max(cell.Sector.Center.Lng, d.max.Lng)
	}

	// Widen the area slightly to allow UEs to move at the edges of the RAN topology
	// No, this does not account for Earth curvature, but should be good enough
	d.min.Lat = d.min.Lat - latMargin
	d.min.Lng = d.min.Lng - lngMargin
	d.max.Lat = d.max.Lat + latMargin
	d.max.Lng = d.max.Lng + lngMargin
}

func (d *driver) generateRoute(ctx context.Context, imsi types.IMSI, speedAvg uint32, speedStdDev uint32, routeEndPoints []model.RouteEndPoint, directRoute bool) error {
	var err error
	var start, end *model.Coordinate

	if len(routeEndPoints) == 0 {
		// chose random end points
		start = d.randomCoordinate()
		end = d.randomCoordinate()
	} else {
		// round-robin through the model's end points
		start = &routeEndPoints[routeEndPointIndex].Start
		end = &routeEndPoints[routeEndPointIndex].End
		routeEndPointIndex = (routeEndPointIndex + 1) % len(routeEndPoints)
	}

	var points []*model.Coordinate
	if len(d.apiKey) >= googleAPIKeyMinLen {
		points, err = googleRoute(start, end, d.apiKey)
		log.Infof("Generated route for UE %d with %d points using Google Directions", imsi, len(points))
	} else {
		points, err = randomRoute(start, end, directRoute)
		log.Infof("Generated route for UE %d with %d points using Random Directions, start:%v, end:%v", imsi, len(points), start, end)
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
		Lat: rand.Float64()*(d.max.Lat-d.min.Lat) + d.min.Lat,
		Lng: rand.Float64()*(d.max.Lng-d.min.Lng) + d.min.Lng,
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

func randomRoute(startLoc *model.Coordinate, endLoc *model.Coordinate, directRoute bool) ([]*model.Coordinate, error) {
	routeWidth := endLoc.Lng - startLoc.Lng
	routeHeight := endLoc.Lat - startLoc.Lat

	directLength := math.Hypot(routeWidth, routeHeight)
	// Try to have a step evey 1/stepsPerDecimalDegree of a decimal degree
	points := make([]*model.Coordinate, int(math.Floor(directLength*stepsPerDecimalDegree)))

	for i := range points {
		randFactor := (rand.Float64() - 0.5) / stepsPerDecimalDegree
		if i == 0 || directRoute {
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
