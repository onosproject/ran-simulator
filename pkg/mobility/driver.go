// SPDX-FileCopyrightText: 2020-present Open Networking Foundation <info@opennetworking.org>
//
// SPDX-License-Identifier: LicenseRef-ONF-Member-1.0

package mobility

import (
	"context"
	"github.com/onosproject/onos-lib-go/pkg/logging"
	"github.com/onosproject/ran-simulator/pkg/model"
	"github.com/onosproject/ran-simulator/pkg/store/routes"
	"github.com/onosproject/ran-simulator/pkg/store/ues"
	"math"
	"math/rand"
	"time"
)

var log = logging.GetLogger("mobility", "driver")

// Driver is an abstraction of an entity driving the UE mobility
type Driver interface {
	// Start starts the driving engine
	Start()

	// Stop stops the driving engine
	Stop()
}

type driver struct {
	routeStore routes.Store
	ueStore    ues.Store
	ticker     *time.Ticker
	done       chan bool
}

// NewMobilityDriver returns a driving engine capable of "driving" UEs along pre-specified routes
func NewMobilityDriver(routeStore routes.Store, ueStore ues.Store) Driver {
	return &driver{
		routeStore: routeStore,
		ueStore:    ueStore,
	}
}

var tickUnit = time.Second

const tickFrequency = 1

func (d *driver) Start() {
	log.Info("Driver starting")
	ctx := context.Background()

	// Iterate over all routes and position the UEs at the start of their routes
	for _, route := range d.routeStore.List(ctx) {
		d.initializeUEPosition(ctx, route)
	}

	d.ticker = time.NewTicker(tickFrequency * tickUnit)
	d.done = make(chan bool)

	go d.drive(ctx)
}

func (d *driver) Stop() {
	log.Info("Driver stopping")
	d.ticker.Stop()
	d.done <- true
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
			}
		}
	}
}

// Initializes UE positions to the start of its routes.
func (d *driver) initializeUEPosition(ctx context.Context, route *model.Route) {
	_ = d.ueStore.MoveToCoordinate(ctx, route.IMSI, *route.Points[0], uint32(math.Round(initialBearing(*route.Points[0], *route.Points[1]))))
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
	bearing := initialBearing(ue.Location, *route.Points[route.NextPoint])
	remainingDistance := distance(ue.Location, *route.Points[route.NextPoint])

	// If distance is less than to the next waypoint, determine the coordinate along that vector
	// Otherwise just use the next waypoint
	newPoint := *route.Points[route.NextPoint]
	reachedWaypoint := remainingDistance <= distanceDriven
	if !reachedWaypoint {
		newPoint = targetPoint(ue.Location, bearing, distanceDriven)
	}

	// Move the UE to the determined coordinate; update heading if necessary
	_ = d.ueStore.MoveToCoordinate(ctx, route.IMSI, newPoint, uint32(math.Round(bearing)))

	// Update the route if necessary
	if reachedWaypoint {
		_ = d.routeStore.Advance(ctx, route.IMSI)
	}
}

// Earth radius in meters
const earthRadius = 6378100

// http://en.wikipedia.org/wiki/Haversine_formula
func distance(c1 model.Coordinate, c2 model.Coordinate) float64 {
	var la1, lo1, la2, lo2 float64
	la1 = c1.Lat * math.Pi / 180
	lo1 = c1.Lng * math.Pi / 180
	la2 = c2.Lat * math.Pi / 180
	lo2 = c2.Lng * math.Pi / 180

	h := hsin(la2-la1) + math.Cos(la1)*math.Cos(la2)*hsin(lo2-lo1)

	return 2 * earthRadius * math.Asin(math.Sqrt(h))
}

// Returns initial bearing from c1 to c2
func initialBearing(c1 model.Coordinate, c2 model.Coordinate) float64 {
	y := math.Sin(c2.Lng-c1.Lng) * math.Cos(c2.Lat)
	x := math.Cos(c1.Lat)*math.Sin(c2.Lat) - math.Sin(c1.Lat)*math.Cos(c2.Lat)*math.Cos(c2.Lng-c1.Lng)
	theta := math.Atan2(y, x)
	return math.Mod(theta*180/math.Pi+360, 360.0) // in degrees
}

// Returns destination point given starting point and distance along heading.
func targetPoint(c model.Coordinate, bearing float64, dist float64) model.Coordinate {
	var la1, lo1, la2, lo2, azimuth, d float64
	la1 = c.Lat * math.Pi / 180
	lo1 = c.Lng * math.Pi / 180
	azimuth = bearing * math.Pi / 180
	d = dist / earthRadius

	la2 = math.Asin(math.Sin(la1)*math.Cos(d) + math.Cos(la1)*math.Sin(d)*math.Cos(azimuth))
	lo2 = lo1 + math.Atan2(math.Sin(azimuth)*math.Sin(d)*math.Cos(la1), math.Cos(d)-math.Sin(la1)*math.Sin(la2))

	return model.Coordinate{Lat: la2 * 180 / math.Pi, Lng: lo2 * 180 / math.Pi}
}

func hsin(theta float64) float64 {
	return math.Pow(math.Sin(theta/2), 2)
}
