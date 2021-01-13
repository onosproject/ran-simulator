// SPDX-FileCopyrightText: 2020-present Open Networking Foundation <info@opennetworking.org>
//
// SPDX-License-Identifier: LicenseRef-ONF-Member-1.0
//

package adapters

import (
	"context"
	"fmt"
	"github.com/onosproject/ran-simulator/api/types"
	"github.com/onosproject/ran-simulator/pkg/model"
	"io"
	"math"
	"strconv"
	"strings"
	"time"

	"github.com/onosproject/onos-lib-go/pkg/certs"
	"github.com/onosproject/onos-lib-go/pkg/logging"
	"github.com/onosproject/onos-lib-go/pkg/southbound"
	"github.com/onosproject/onos-topo/api/topo"
	"github.com/onosproject/ran-simulator/pkg/utils"
	"google.golang.org/grpc"
)

var log = logging.GetLogger("adapter", "topo")

const simNodeKindID = "SimE2Node"

// TopoAdapter provides means to populate simulated nodes from ONOS topology service
type TopoAdapter struct {
	TopoEndpoint string
	ServerParams utils.ServerParams
	SimNodes     model.SimNodes
}

// Connect initiates connection to the ONOS topology service and uses its events to populate
// the model with simulated nodes
func (a *TopoAdapter) Connect(ctx context.Context) error {
	log.Infof("Connecting to ONOS Topo...%s", a.TopoEndpoint)
	// Attempt to create connection to the Topo
	opts, err := certs.HandleCertPaths(a.ServerParams.CaPath, a.ServerParams.KeyPath, a.ServerParams.CertPath, true)
	if err != nil {
		log.Fatal(err)
	}
	opts = append(opts, grpc.WithStreamInterceptor(southbound.RetryingStreamClientInterceptor(time.Second)))
	conn, err := southbound.Connect(ctx, a.TopoEndpoint, "", "", opts...)
	if err != nil {
		log.Fatal("Failed to connect to %s. Retry. %s", a.TopoEndpoint, err)
	}

	topoClient := topo.NewTopoClient(conn)
	stream, err := topoClient.Watch(context.Background(), &topo.WatchRequest{})
	if err != nil {
		return err
	}
	for {
		resp, err := stream.Recv() // Block here and wait for events from topo
		if err == io.EOF {
			// read done.
			return nil
		}
		if err != nil {
			return err
		}

		// If the event is relevant, process it
		if isRelevant(resp) {
			a.processResponse(resp)
		}
	}
}

func isRelevant(resp *topo.WatchResponse) bool {
	return resp.Event.Object.Type == topo.Object_ENTITY &&
		resp.Event.Object.GetEntity().KindID == simNodeKindID
}

func (a *TopoAdapter) processResponse(resp *topo.WatchResponse) {
	switch resp.Event.Type {
	case topo.EventType_NONE:
		a.addNode(resp.Event.Object)
	case topo.EventType_ADDED:
		a.addNode(resp.Event.Object)
	case topo.EventType_REMOVED:
		a.removeNode(resp.Event.Object)
	default:
	}
}

func (a *TopoAdapter) addNode(obj topo.Object) {
	cell, err := toCell(obj)
	if err == nil {
		a.SimNodes.Add(&model.SimNode{Cell: *cell})
	}
}

func (a *TopoAdapter) removeNode(obj topo.Object) {

	// TODO
}

const (
	// DefaultTxPower - all cells start with this power level
	DefaultTxPower = 10

	// PowerFactor - relate power to distance in decimal degrees
	PowerFactor = 0.001

	// PowerBase - baseline for power to distance in decimal degrees
	PowerBase = 0.013
)

func toCell(object topo.Object) (*types.Cell, error) {
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
		MaxUEs:     10, // MaxNumUesPerCell,
		Sector: &types.Sector{
			Azimuth: int32(azimuth),
			Arc:     int32(arc),
		},
	}
	cell.Sector.Centroid = centroidPosition(cell)

	return cell, nil
}

func ecgiFromtopo(object topo.Object) (*types.ECGI, error) {
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
