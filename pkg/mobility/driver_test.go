// SPDX-FileCopyrightText: 2020-present Open Networking Foundation <info@opennetworking.org>
//
// SPDX-License-Identifier: LicenseRef-ONF-Member-1.0

package mobility

import (
	"context"
	"github.com/onosproject/ran-simulator/pkg/model"
	"github.com/onosproject/ran-simulator/pkg/store/cells"
	"github.com/onosproject/ran-simulator/pkg/store/nodes"
	"github.com/onosproject/ran-simulator/pkg/store/routes"
	"github.com/onosproject/ran-simulator/pkg/store/ues"
	"testing"
)

func TestDriver(t *testing.T) {
	m := &model.Model{}
	err := model.LoadConfig(m, "test")

	ns := nodes.NewNodeRegistry(m.Nodes)
	cs := cells.NewCellRegistry(m.Cells, ns)
	us := ues.NewUERegistry(1, cs)
	rs := routes.NewRouteRegistry()

	ctx := context.Background()


}