// SPDX-FileCopyrightText: 2020-present Open Networking Foundation <info@opennetworking.org>
//
// SPDX-License-Identifier: Apache-2.0

package messageformat1

import (
	"encoding/hex"
	"github.com/onosproject/onos-e2-sm/servicemodels/e2sm_mho_go/pdubuilder"
	e2sm_mho_go "github.com/onosproject/onos-e2-sm/servicemodels/e2sm_mho_go/v2/e2sm-mho-go"
	"github.com/onosproject/onos-lib-go/api/asn1/v1/asn1"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCreateIndicationMessage(t *testing.T) {

	cgi, err := pdubuilder.CreateCgiNrCGI([]byte{0xAA, 0xFD, 0xD4}, &asn1.BitString{
		Value: []byte{0x00, 0x00, 0x00, 0x40, 0x00},
		Len:   36,
	})
	assert.NoError(t, err)
	rsrp := &e2sm_mho_go.Rsrp{
		Value: 1234,
	}
	measItem, err := pdubuilder.CreateMeasurementRecordItem(cgi, rsrp)
	assert.NoError(t, err)
	measItem.SetFiveQi(21)
	measReport := make([]*e2sm_mho_go.E2SmMhoMeasurementReportItem, 0)
	measReport = append(measReport, measItem)

	indicationMessage := NewIndicationMessage(WithUeID(1), WithMeasReport(measReport))
	assert.NotNil(t, indicationMessage)
	assert.Equal(t, indicationMessage.ueID, int64(1))
	assert.Equal(t, len(indicationMessage.MeasReport), 1)
	assert.Equal(t, indicationMessage.MeasReport[0].GetFiveQi().GetValue(), int32(21))

	aper, err := indicationMessage.ToAsn1Bytes()
	assert.NoError(t, err)
	t.Logf("E2SM-MHO-IndicationMessage (Format 1) APER bytes are\n%v", hex.Dump(aper))
}
