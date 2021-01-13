// SPDX-FileCopyrightText: 2020-present Open Networking Foundation <info@opennetworking.org>
//
// SPDX-License-Identifier: LicenseRef-ONF-Member-1.0

package model

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestModelBasics(t *testing.T) {
	model := NewModel()
	err := model.Load("test_model.yml")
	assert.NoError(t, err, "unable to load model")
	if err == nil {
		assert.Equal(t, 3, len(model.Nodes.GetAll()), "incorrect number of nodes")
		assert.Equal(t, 2, len(model.Nodes.Get("90125-10002").ServiceModels), "incorrect number of models")

		node := model.Nodes.Remove("90125-10001")
		assert.NotNil(t, node)
		assert.Equal(t, "90125-10001", string(node.ECGI), "incorrect node removed")
		assert.Equal(t, 2, len(model.Nodes.GetAll()), "incorrect number of nodes")
	}
}
