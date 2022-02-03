// SPDX-FileCopyrightText: 2020-present Open Networking Foundation <info@opennetworking.org>
//
// SPDX-License-Identifier: Apache-2.0

package message

import (
	"testing"

	ransimtypes "github.com/onosproject/onos-api/go/onos/ransim/types"

	e2smrcpreies "github.com/onosproject/onos-e2-sm/servicemodels/e2sm_rc_pre_go/v2/e2sm-rc-pre-v2-go"
	"github.com/onosproject/ran-simulator/pkg/utils/e2sm/rc/nrt"
	"github.com/stretchr/testify/assert"
)

func TestCreateIndicationMessage(t *testing.T) {

	plmnID := ransimtypes.NewUint24(12345)
	nrt1, err := nrt.NewNeighbour(
		nrt.WithPci(10),
		nrt.WithNrcellIdentity(15),
		nrt.WithEarfcn(40),
		nrt.WithCellSize(e2smrcpreies.CellSize_CELL_SIZE_MACRO),
		nrt.WithPlmnID(plmnID.Value())).Build()
	assert.NoError(t, err)
	nrt2, err := nrt.NewNeighbour(
		nrt.WithPci(20),
		nrt.WithNrcellIdentity(25),
		nrt.WithEarfcn(50),
		nrt.WithCellSize(e2smrcpreies.CellSize_CELL_SIZE_FEMTO),
		nrt.WithPlmnID(plmnID.Value())).Build()
	assert.NoError(t, err)

	indicationMessage, err := NewIndicationMessage(WithPlmnID(plmnID.Value()),
		WithCellSize(e2smrcpreies.CellSize_CELL_SIZE_MACRO),
		WithEarfcn(20),
		WithPci(10),
		WithNeighbours([]*e2smrcpreies.Nrt{nrt1, nrt2})).Build()
	assert.NoError(t, err)

	assert.Equal(t, indicationMessage.GetIndicationMessageFormat1().Neighbors[0].Pci.GetValue(), int32(10))

}
