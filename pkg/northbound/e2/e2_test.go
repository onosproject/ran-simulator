// SPDX-FileCopyrightText: 2020-present Open Networking Foundation <info@opennetworking.org>
//
// SPDX-License-Identifier: LicenseRef-ONF-Member-1.0
//

package e2

import (
	"github.com/onosproject/onos-topo/pkg/bulk"
	"github.com/onosproject/ran-simulator/api/types"
	"github.com/onosproject/ran-simulator/pkg/manager"
)

func setUpManager() (*manager.Manager, error) {
	routesParams := manager.RoutesParams{
		APIKey:    "",
		StepDelay: 1000,
	}

	topoConfig, err := bulk.GetTopoConfig("berlin-rectangular-4-1-topo.yaml")
	if err != nil {
		return nil, err
	}
	mapLayout := types.MapLayout{
		Zoom:   12,
		MinUes: 3,
	}

	cells := make(map[types.ECGI]*types.Cell)

	for _, td := range topoConfig.TopoEntities {
		td := td //pin
		cell, err := manager.NewCell(bulk.TopoEntityToTopoObject(&td))
		if err != nil {
			return nil, err
		}
		cells[*cell.Ecgi] = cell
	}
	centre, locations := manager.NewLocations(cells, 5, 1)
	mapLayout.Center = centre

	mgr, err := manager.NewManager()
	if err != nil {
		return nil, err
	}
	mgr.MapLayout = mapLayout
	mgr.CellsLock.Lock()
	mgr.Cells = cells
	mgr.CellsLock.Unlock()
	mgr.Locations = locations

	mgr.Routes, err = mgr.NewRoutes(mapLayout, routesParams)
	if err != nil {
		return nil, err
	}
	mgr.UserEquipments, err = mgr.NewUserEquipments(mapLayout, routesParams)
	if err != nil {
		return nil, err
	}
	return mgr, nil
}

func stopManager(m *manager.Manager) {
	close(m.CellsChannel)
	close(m.UeChannel)
	close(m.RouteChannel)
	close(m.LatencyChannel)
	for r := range m.Routes {
		delete(m.Routes, r)
	}
	for l := range m.Locations {
		delete(m.Locations, l)
	}
	m.CellsLock.Lock()
	for tid := range m.Cells {
		delete(m.Cells, tid)
	}
	m.CellsLock.Unlock()
}
