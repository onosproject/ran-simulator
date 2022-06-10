// SPDX-FileCopyrightText: 2020-present Open Networking Foundation <info@opennetworking.org>
//
// SPDX-License-Identifier: Apache-2.0

package mho

import (
	"context"
	ransimtypes "github.com/onosproject/onos-api/go/onos/ransim/types"
	e2sm_mho "github.com/onosproject/onos-e2-sm/servicemodels/e2sm_mho_go/v2/e2sm-mho-go"
	e2sm_v2_ies "github.com/onosproject/onos-e2-sm/servicemodels/e2sm_mho_go/v2/e2sm-v2-ies"
	"github.com/onosproject/onos-lib-go/api/asn1/v1/asn1"
	"github.com/onosproject/ran-simulator/pkg/model"
	"github.com/onosproject/ran-simulator/pkg/store/subscriptions"
	"github.com/onosproject/ran-simulator/pkg/utils"
	e2apIndicationUtils "github.com/onosproject/ran-simulator/pkg/utils/e2ap/indication"
	subutils "github.com/onosproject/ran-simulator/pkg/utils/e2ap/subscription"
	indHdr "github.com/onosproject/ran-simulator/pkg/utils/e2sm/mho/indication/header"
	indMsgFmt1 "github.com/onosproject/ran-simulator/pkg/utils/e2sm/mho/indication/message_format1"
	indMsgFmt2 "github.com/onosproject/ran-simulator/pkg/utils/e2sm/mho/indication/message_format2"
)

func (m *Mho) sendRicIndication(ctx context.Context, subscription *subutils.Subscription) error {
	node := m.ServiceModel.Node
	// Creates and sends an indication message for each cell in the node
	for _, ncgi := range node.Cells {
		log.Debugf("Send MHO indications for cell ncgi:%d", ncgi)
		for _, ue := range m.ServiceModel.UEs.ListUEs(ctx, ncgi) {
			// Ignore idle UEs
			if ue.RrcState == e2sm_mho.Rrcstatus_RRCSTATUS_IDLE {
				continue
			}
			log.Debugf("Send MHO indications for cell ncgi:%d, IMSI:%d", ncgi, ue.IMSI)
			err := m.sendRicIndicationFormat1(ctx, ncgi, ue, subscription)
			if err != nil {
				log.Warn(err)
				continue
			}
		}
	}
	return nil
}

func (m *Mho) sendRicIndicationFormat1(ctx context.Context, ncgi ransimtypes.NCGI, ue *model.UE, subscription *subutils.Subscription) error {
	subID := subscriptions.NewID(subscription.GetRicInstanceID(), subscription.GetReqID(), subscription.GetRanFuncID())
	sub, err := m.ServiceModel.Subscriptions.Get(subID)
	if err != nil {
		return err
	}

	indicationHeaderBytes, err := m.createIndicationHeaderBytes(ctx, ncgi)
	if err != nil {
		return err
	}

	indicationMessageBytes, err := m.createIndicationMsgFormat1(ue)
	if err != nil {
		return err
	}
	if indicationMessageBytes == nil {
		return nil
	}

	indication := e2apIndicationUtils.NewIndication(
		e2apIndicationUtils.WithRicInstanceID(subscription.GetRicInstanceID()),
		e2apIndicationUtils.WithRanFuncID(subscription.GetRanFuncID()),
		e2apIndicationUtils.WithRequestID(subscription.GetReqID()),
		e2apIndicationUtils.WithIndicationHeader(indicationHeaderBytes),
		e2apIndicationUtils.WithIndicationMessage(indicationMessageBytes))

	ricIndication, err := indication.Build()
	if err != nil {
		return err
	}

	err = sub.E2Channel.RICIndication(ctx, ricIndication)
	if err != nil {
		return err
	}

	return nil
}

func (m *Mho) sendRicIndicationFormat2(ctx context.Context, ncgi ransimtypes.NCGI, ue *model.UE, subscription *subutils.Subscription) error {
	subID := subscriptions.NewID(subscription.GetRicInstanceID(), subscription.GetReqID(), subscription.GetRanFuncID())
	sub, err := m.ServiceModel.Subscriptions.Get(subID)
	if err != nil {
		return err
	}

	indicationHeaderBytes, err := m.createIndicationHeaderBytes(ctx, ncgi)
	if err != nil {
		return err
	}

	indicationMessageBytes, err := m.createIndicationMsgFormat2(ue)
	if err != nil {
		return err
	}
	if indicationMessageBytes == nil {
		return nil
	}

	indication := e2apIndicationUtils.NewIndication(
		e2apIndicationUtils.WithRicInstanceID(subscription.GetRicInstanceID()),
		e2apIndicationUtils.WithRanFuncID(subscription.GetRanFuncID()),
		e2apIndicationUtils.WithRequestID(subscription.GetReqID()),
		e2apIndicationUtils.WithIndicationHeader(indicationHeaderBytes),
		e2apIndicationUtils.WithIndicationMessage(indicationMessageBytes))

	ricIndication, err := indication.Build()
	if err != nil {
		return err
	}

	err = sub.E2Channel.RICIndication(ctx, ricIndication)
	if err != nil {
		return err
	}

	return nil
}

func (m *Mho) createIndicationHeaderBytes(ctx context.Context, ncgi ransimtypes.NCGI) ([]byte, error) {

	cell, _ := m.ServiceModel.CellStore.Get(ctx, ncgi)
	plmnID := ransimtypes.NewUint24(uint32(m.ServiceModel.Model.PlmnID))
	ncgiTypeNCI := utils.NewNCellIDWithUint64(uint64(ransimtypes.GetNCI(cell.NCGI)))

	header := indHdr.NewIndicationHeader(
		indHdr.WithPlmnID(*plmnID),
		indHdr.WithNrcellIdentity(ncgiTypeNCI.Bytes()))

	indicationHeaderAsn1Bytes, err := header.MhoToAsn1Bytes()
	if err != nil {
		return nil, err
	}

	return indicationHeaderAsn1Bytes, nil
}

func (m *Mho) createIndicationMsgFormat1(ue *model.UE) ([]byte, error) {
	log.Debugf("Create MHO Indication message ueID: %d", ue.IMSI)

	plmnID := ransimtypes.NewUint24(uint32(m.ServiceModel.Model.PlmnID))
	measReport := make([]*e2sm_mho.E2SmMhoMeasurementReportItem, 0)

	if len(ue.Cells) == 0 {
		log.Infof("no neighbor cells found for ueID:%d", ue.IMSI)
		return nil, nil
	}

	nrCellIDTypeNCI := utils.NewNCellIDWithUint64(uint64(ransimtypes.GetNCI(ue.Cell.NCGI)))

	// add serving cell to measReport
	item := &e2sm_mho.E2SmMhoMeasurementReportItem{
		Cgi: &e2sm_v2_ies.Cgi{
			Cgi: &e2sm_v2_ies.Cgi_NRCgi{
				NRCgi: &e2sm_v2_ies.NrCgi{
					PLmnidentity: &e2sm_v2_ies.PlmnIdentity{
						Value: plmnID.ToBytes(),
					},
					NRcellIdentity: &e2sm_v2_ies.NrcellIdentity{
						Value: &asn1.BitString{
							Value: nrCellIDTypeNCI.Bytes(),
							Len:   36,
						},
					},
				},
			},
		},
		Rsrp: &e2sm_mho.Rsrp{
			Value: int32(ue.Cell.Strength),
		},
		FiveQi: &e2sm_v2_ies.FiveQi{
			Value: int32(ue.FiveQi),
		},
	}

	measReport = append(measReport, item)

	for _, cell := range ue.Cells {
		ncgiTypeNCI := utils.NewNCellIDWithUint64(uint64(ransimtypes.GetNCI(cell.NCGI)))

		measReport = append(measReport, &e2sm_mho.E2SmMhoMeasurementReportItem{
			Cgi: &e2sm_v2_ies.Cgi{
				Cgi: &e2sm_v2_ies.Cgi_NRCgi{
					NRCgi: &e2sm_v2_ies.NrCgi{
						PLmnidentity: &e2sm_v2_ies.PlmnIdentity{
							Value: plmnID.ToBytes(),
						},
						NRcellIdentity: &e2sm_v2_ies.NrcellIdentity{
							Value: &asn1.BitString{
								Value: ncgiTypeNCI.Bytes(),
								Len:   36,
							},
						},
					},
				},
			},
			Rsrp: &e2sm_mho.Rsrp{
				Value: int32(cell.Strength),
			},
		})
	}

	ueID := int64(ue.IMSI)

	log.Debugf("MHO measurement report for ueID %s: %v", ueID, measReport)

	indicationMessage := indMsgFmt1.NewIndicationMessage(
		indMsgFmt1.WithUeID(ueID),
		indMsgFmt1.WithMeasReport(measReport))

	log.Debugf("MHO measurement report indication message for ueID %s: %v", ueID, indicationMessage)

	indicationMessageBytes, err := indicationMessage.ToAsn1Bytes()
	if err != nil {
		log.Warn(err)
		return nil, err
	}

	return indicationMessageBytes, nil
}

func (m *Mho) createIndicationMsgFormat2(ue *model.UE) ([]byte, error) {
	log.Debugf("Create MHO RRC indication message ueID: %d", ue.IMSI)

	ueID := int64(ue.AmfUeNgapID)

	indicationMessage := indMsgFmt2.NewIndicationMessage(
		indMsgFmt2.WithUeID(ueID),
		indMsgFmt2.WithRrcStatus(ue.RrcState),
		indMsgFmt2.WithGuami(uint64(m.ServiceModel.Model.PlmnID), m.ServiceModel.Model.Guami.AmfRegionID,
			m.ServiceModel.Model.Guami.AmfSetID, m.ServiceModel.Model.Guami.AmfPointer))

	log.Debugf("MHO RRC state indication message for ueID amf ue ngap id - %v, plmnid - %v, amf region id - %v,"+
		"amf set id - %v, amf pointer - %v: %v", ueID, uint64(m.ServiceModel.Model.PlmnID), m.ServiceModel.Model.Guami.AmfRegionID, m.ServiceModel.Model.Guami.AmfSetID, m.ServiceModel.Model.Guami.AmfPointer,
		indicationMessage)

	indicationMessageBytes, err := indicationMessage.ToAsn1Bytes()
	if err != nil {
		log.Warn(err)
		return nil, err
	}

	return indicationMessageBytes, nil
}
