// SPDX-FileCopyrightText: 2020-present Open Networking Foundation <info@opennetworking.org>
//
// SPDX-License-Identifier: Apache-2.0
//

package config

import (
	"fmt"
	"github.com/onosproject/onos-config/pkg/config/load"
	cfggnmi "github.com/onosproject/onos-config/pkg/northbound/gnmi"
	topodevice "github.com/onosproject/onos-topo/api/device"
	"github.com/onosproject/onos-topo/pkg/bulk"
	"github.com/onosproject/ran-simulator/api/types"
	"github.com/onosproject/ran-simulator/pkg/utils"
	"github.com/openconfig/gnmi/proto/gnmi"
	"github.com/pmcxs/hexgrid"
	"math"
)

// HoneycombTopoGenerator - used by the cli tool "honeycomb"
func HoneycombTopoGenerator(numTowers uint, sectorsPerTower uint, latitude float64,
	longitude float64, plmnid types.PlmnID, ecidStart uint16, portstart uint16, pitch float32) (*bulk.DeviceConfig, error) {

	mapCentre := types.Point{
		Lat: latitude,
		Lng: longitude,
	}

	aspectRatio := utils.AspectRatio(&mapCentre)
	newTopoConfig := bulk.DeviceConfig{TopoDevices: make([]topodevice.Device, 0)}

	points := hexMesh(float64(pitch), numTowers)
	var t, s uint
	for t = 0; t < numTowers; t++ {
		var azOffset uint = 0
		if sectorsPerTower == 6 {
			azOffset = uint(math.Mod(float64(t), 2) * 30)
		}
		for s = 0; s < sectorsPerTower; s++ {
			ecID := types.EcID(fmt.Sprintf("%07x", ecidStart+uint16(t*sectorsPerTower)+uint16(s)))
			grpcPort := portstart + uint16(t*sectorsPerTower) + uint16(s)
			topoDevice := topodevice.Device{
				ID:      topodevice.ID(fmt.Sprintf("%s-%s", string(plmnid), string(ecID))),
				Address: fmt.Sprintf("ran-simulator:%d", grpcPort),
				Version: types.E2NodeVersion100,
				TLS: topodevice.TlsConfig{
					Insecure: true,
				},
				Type:        types.E2NodeType,
				Displayname: fmt.Sprintf("Tower %d Cell %d", t, s),
				Attributes:  make(map[string]string),
			}
			azimuth := azOffset
			if s > 0 {
				azimuth = 360.0*s/sectorsPerTower + azOffset
			}
			topoDevice.Attributes[types.AzimuthKey] = fmt.Sprintf("%d", azimuth)
			topoDevice.Attributes[types.ArcKey] = fmt.Sprintf("%d", 360.0/uint16(sectorsPerTower))
			topoDevice.Attributes[types.PlmnIDKey] = string(plmnid)
			topoDevice.Attributes[types.EcidKey] = string(ecID)
			topoDevice.Attributes[types.LatitudeKey] = fmt.Sprintf("%f", latitude+points[t].Lat)
			topoDevice.Attributes[types.LongitudeKey] = fmt.Sprintf("%f", longitude+points[t].Lng/aspectRatio)
			topoDevice.Attributes[types.GrpcPortKey] = fmt.Sprintf("%d", grpcPort)
			newTopoConfig.TopoDevices = append(newTopoConfig.TopoDevices, topoDevice)
		}
	}
	return &newTopoConfig, nil
}

// HoneycombConfigGenerator - used by the cli tool "honeycomb"
func HoneycombConfigGenerator(numTowers uint, sectorsPerTower uint, plmnid types.PlmnID,
	ecidStart uint16) (*load.ConfigGnmiSimple, error) {

	newConfigConfig := load.ConfigGnmiSimple{
		SetRequest: load.SetRequest{
			Prefix: &gnmi.Path{
				Elem: []*gnmi.PathElem{
					{
						Name: "e2node",
					},
				},
			},
			Update:  make([]*load.Update, 0),
			Replace: make([]*load.Update, 0),
			Delete:  make([]*gnmi.Path, 0),
		},
	}
	var t, s uint
	for t = 0; t < numTowers; t++ {
		for s = 0; s < sectorsPerTower; s++ {
			ecID := types.EcID(fmt.Sprintf("%07x", ecidStart+uint16(t*sectorsPerTower)+uint16(s)))

			updateRadioMeasReportPerUe := &load.Update{
				Path: &gnmi.Path{
					Target: fmt.Sprintf("%s-%s", string(plmnid), string(ecID)),
					Elem: []*gnmi.PathElem{
						{Name: "intervals"},
						{Name: "RadioMeasReportPerUe"},
					},
				},
				Val: &load.TypedValue{UIntValue: &gnmi.TypedValue_UintVal{UintVal: 20}},
			}
			newConfigConfig.SetRequest.Update = append(newConfigConfig.SetRequest.Update, updateRadioMeasReportPerUe)

			updateRadioMeasReportPerCell := &load.Update{
				Path: &gnmi.Path{
					Target: fmt.Sprintf("%s-%s", string(plmnid), string(ecID)),
					Elem: []*gnmi.PathElem{
						{Name: "intervals"},
						{Name: "RadioMeasReportPerCell"},
					},
				},
				Val: &load.TypedValue{UIntValue: &gnmi.TypedValue_UintVal{UintVal: 21}},
			}
			newConfigConfig.SetRequest.Update = append(newConfigConfig.SetRequest.Update, updateRadioMeasReportPerCell)
		}
	}

	newConfigConfig.SetRequest.Extension = []*load.Extension{
		{
			ID:    cfggnmi.GnmiExtensionVersion,
			Value: types.E2NodeVersion100,
		},
		{
			ID:    cfggnmi.GnmiExtensionDeviceType,
			Value: types.E2NodeType,
		},
	}
	return &newConfigConfig, nil
}

func hexMesh(pitch float64, numTowers uint) []*types.Point {
	rings, _ := numRings(numTowers)
	points := make([]*types.Point, 0)
	hexArray := hexgrid.HexRange(hexgrid.NewHex(0, 0), int(rings))

	for _, h := range hexArray {
		x, y := hexgrid.Point(hexgrid.HexToPixel(hexgrid.LayoutPointY00(pitch, pitch), h))
		points = append(points, &types.Point{
			Lat: x,
			Lng: y,
		})
	}
	return points
}

// Number of cells in the hexagon layout 3x^2+9x+7
func numRings(numTowers uint) (uint, error) {
	switch n := numTowers; {
	case n <= 7:
		return 1, nil
	case n <= 19:
		return 2, nil
	case n <= 37:
		return 3, nil
	case n <= 61:
		return 4, nil
	case n <= 91:
		return 5, nil
	case n <= 127:
		return 6, nil
	case n <= 169:
		return 7, nil
	case n <= 217:
		return 8, nil
	case n <= 271:
		return 9, nil
	case n <= 331:
		return 10, nil
	case n <= 469:
		return 11, nil
	default:
		return 0, fmt.Errorf(">469 not handled %d", numTowers)
	}

}
