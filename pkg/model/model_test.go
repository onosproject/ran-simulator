// SPDX-FileCopyrightText: 2020-present Open Networking Foundation <info@opennetworking.org>
//
// SPDX-License-Identifier: Apache-2.0

package model

import (
	"github.com/onosproject/onos-api/go/onos/ransim/types"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestModel(t *testing.T) {
	model := &Model{}
	err := LoadConfig(model, "test")
	assert.NoError(t, err)
	t.Log(model)
	assert.Equal(t, 2, len(model.Controllers))
	assert.Equal(t, 2, len(model.Nodes))
	assert.Equal(t, 4, len(model.Cells))
	assert.Equal(t, 36421, model.Controllers["controller1"].Port)
	assert.Equal(t, 36421, model.Controllers["controller2"].Port)
	assert.Equal(t, "1.0.0", model.ServiceModels["kpm"].Version)
	assert.Equal(t, 3, model.ServiceModels["rc"].ID)
	assert.Equal(t, 2, model.ServiceModels["ni"].ID)
	assert.Equal(t, uint(12), model.UECount)
	assert.Equal(t, "314628", model.Plmn)
	assert.Equal(t, types.PlmnID(0x314628), model.PlmnID)

	assert.Equal(t, types.NCGI(84325717761), model.Cells["cell3"].NCGI)
	assert.Equal(t, 2, len(model.Nodes["node1"].Cells))
	assert.Equal(t, 44.0, model.Cells["cell3"].Sector.Center.Lat)

	assert.Equal(t, true, model.MapLayout.FadeMap)
	assert.Equal(t, 45.0, model.MapLayout.Center.Lat)
}
