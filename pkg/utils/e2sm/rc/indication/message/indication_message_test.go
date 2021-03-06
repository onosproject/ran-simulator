// SPDX-FileCopyrightText: 2020-present Open Networking Foundation <info@opennetworking.org>
//
// SPDX-License-Identifier: LicenseRef-ONF-Member-1.0

package message

import (
	"testing"

	ransimtypes "github.com/onosproject/onos-api/go/onos/ransim/types"

	"github.com/onosproject/ran-simulator/pkg/utils/e2sm/rc/pcirange"

	"github.com/onosproject/ran-simulator/pkg/utils/e2sm/rc/nrt"
	"github.com/stretchr/testify/assert"

	e2smrcpreies "github.com/onosproject/onos-e2-sm/servicemodels/e2sm_rc_pre/v1/e2sm-rc-pre-ies"
)

func TestCreateIndicationMessage(t *testing.T) {

	plmnID := ransimtypes.NewUint24(12345)
	nrt1, err := nrt.NewNeighbour(
		nrt.WithNrIndex(1),
		nrt.WithPci(10),
		nrt.WithEutraCellIdentity(15),
		nrt.WithEarfcn(40),
		nrt.WithCellSize(e2smrcpreies.CellSize_CELL_SIZE_MACRO),
		nrt.WithPlmnID(plmnID.Value())).Build()
	assert.NoError(t, err)
	nrt2, err := nrt.NewNeighbour(
		nrt.WithNrIndex(2),
		nrt.WithPci(20),
		nrt.WithEutraCellIdentity(25),
		nrt.WithEarfcn(50),
		nrt.WithCellSize(e2smrcpreies.CellSize_CELL_SIZE_FEMTO),
		nrt.WithPlmnID(plmnID.Value())).Build()
	assert.NoError(t, err)

	pciRange1, err := pcirange.NewPciRange(pcirange.WithLowerPci(10),
		pcirange.WithUpperPci(30)).Build()
	assert.NoError(t, err)

	pciRange2, err := pcirange.NewPciRange(pcirange.WithLowerPci(20),
		pcirange.WithUpperPci(50)).Build()
	assert.NoError(t, err)

	indicationMessage, err := NewIndicationMessage(WithPlmnID(plmnID.Value()),
		WithCellSize(e2smrcpreies.CellSize_CELL_SIZE_MACRO),
		WithEarfcn(20),
		WithEutraCellIdentity(30),
		WithPci(10),
		WithNeighbours([]*e2smrcpreies.Nrt{nrt1, nrt2}),
		WithPciPool([]*e2smrcpreies.PciRange{pciRange1, pciRange2})).Build()
	assert.NoError(t, err)

	assert.Equal(t, indicationMessage.GetIndicationMessageFormat1().Neighbors[0].NrIndex, int32(1))
	assert.Equal(t, indicationMessage.GetIndicationMessageFormat1().PciPool[0].LowerPci.Value, int32(10))

}
