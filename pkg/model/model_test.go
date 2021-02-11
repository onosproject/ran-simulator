// SPDX-FileCopyrightText: 2020-present Open Networking Foundation <info@opennetworking.org>
//
// SPDX-License-Identifier: LicenseRef-ONF-Member-1.0

package model

import (
	"io/ioutil"
	"testing"

	"github.com/stretchr/testify/assert"
	"gopkg.in/yaml.v2"
)

func TestModel(t *testing.T) {
	model := Model{}
	bytes, err := ioutil.ReadFile("test.yaml")
	assert.NoError(t, err)
	err = yaml.Unmarshal(bytes, &model)
	assert.NoError(t, err)
	t.Log(model)
	assert.Equal(t, 2, len(model.Controllers))
	assert.Equal(t, 2, len(model.Nodes))
	assert.Equal(t, 36421, model.Controllers["controller1"].Port)
	assert.Equal(t, 36421, model.Controllers["controller2"].Port)
	assert.Equal(t, "1.0.0", model.ServiceModels["kpm"].Version)
	assert.Equal(t, 3, model.ServiceModels["rc"].ID)
	assert.Equal(t, 2, model.ServiceModels["ni"].ID)
	assert.Equal(t, uint(12), model.UECount)
	assert.Equal(t, PlmnID(314), model.PlmnID)

	assert.Equal(t, 2, len(model.Nodes["node1"].Cells))
	assert.Equal(t, 44.0, model.Nodes["node2"].Cells["cell1"].Sector.Center.Lat)
}
