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
	assert.Equal(t, len(model.Controllers), 2)
	assert.Equal(t, len(model.Nodes), 2)
	assert.Equal(t, model.Controllers["controller1"].Port, 36421)
	assert.Equal(t, model.Controllers["controller2"].Port, 36421)
	assert.Equal(t, model.ServiceModels["kpm"].Version, "1.0.0")
	assert.Equal(t, model.ServiceModels["ni"].ID, 2)

}
