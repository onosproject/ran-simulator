// SPDX-FileCopyrightText: 2020-present Open Networking Foundation <info@opennetworking.org>
//
// SPDX-License-Identifier: Apache-2.0

package f1ap

import (
	f1apv1 "github.com/onosproject/onos-e2t/api/f1ap/v1"
	f1apcommondatatypesv1 "github.com/onosproject/onos-e2t/api/f1ap/v1/f1ap_commondatatypes"
	f1apiesv1 "github.com/onosproject/onos-e2t/api/f1ap/v1/f1ap_ies"
	f1appducontentsv1 "github.com/onosproject/onos-e2t/api/f1ap/v1/f1ap_pdu_contents"
	f1appdudescriptionsv1 "github.com/onosproject/onos-e2t/api/f1ap/v1/f1ap_pdu_descriptions"
	"github.com/onosproject/onos-e2t/pkg/southbound/f1ap/encoder"
	"github.com/onosproject/onos-e2t/pkg/southbound/f1ap/pdubuilder"
	"github.com/onosproject/onos-lib-go/api/asn1/v1/asn1"
	"github.com/onosproject/onos-lib-go/pkg/logging"
)

var log = logging.GetLogger()

type SCellItemInfo struct {
	PlmnIDBytes                     []byte
	NrCellIDBytes                   []byte
	NrCellIDLen                     uint32
	NrPCI                           int32
	SulFreqBandIndicationNr         int32
	FreqBandIndicatorNr             int32
	NrArfcn                         int32
	MeasureTimingConfigurationBytes []byte
}

func CreateF1SetupRequest(gnbDUID int64, rrcVerBytes []byte, rrcVerLen uint32, f1apSCellItemInfo []SCellItemInfo) ([]byte, error) {
	list := make([]*f1appducontentsv1.F1SetupRequestIes, 0)

	// transaction ID
	trID, err := pdubuilder.CreateTransactionID(1)
	if err != nil {
		return nil, err
	}
	ie1Value, err := pdubuilder.CreateF1SetupRequestIesValueTransactionID(trID)
	if err != nil {
		return nil, err
	}
	ie1, err := pdubuilder.CreateF1SetupRequestIes(int32(f1apv1.ProtocolIeIDTransactionID),
		f1apcommondatatypesv1.Criticality_CRITICALITY_REJECT, ie1Value)
	if err != nil {
		return nil, err
	}
	list = append(list, ie1)

	// GnbDuID
	gnbDuID, err := pdubuilder.CreateGnbDUID(gnbDUID)
	if err != nil {
		return nil, err
	}
	ie2Value, err := pdubuilder.CreateF1SetupRequestIesValueGnbDuID(gnbDuID)
	if err != nil {
		return nil, err
	}
	ie2, err := pdubuilder.CreateF1SetupRequestIes(int32(f1apv1.ProtocolIeIDgNBDUID),
		f1apcommondatatypesv1.Criticality_CRITICALITY_REJECT, ie2Value)
	if err != nil {
		return nil, err
	}
	list = append(list, ie2)

	// GnbDuRRC version
	rrcVersion, err := pdubuilder.CreateRrcVersion(&asn1.BitString{
		Value: rrcVerBytes,
		Len:   rrcVerLen,
	})
	if err != nil {
		return nil, err
	}
	ie3Value, err := pdubuilder.CreateF1SetupRequestIesValueRrcVersion(rrcVersion)
	if err != nil {
		return nil, err
	}
	ie3, err := pdubuilder.CreateF1SetupRequestIes(int32(f1apv1.ProtocolIeIDGNBDURRCVersion),
		f1apcommondatatypesv1.Criticality_CRITICALITY_REJECT, ie3Value)
	if err != nil {
		return nil, err
	}
	list = append(list, ie3)

	// Served Cells List
	gnbDuServedCellsList := make([]*f1appducontentsv1.GnbDUServedCellsItemIes, 0)
	for _, sCell := range f1apSCellItemInfo {
		plmnID, err := pdubuilder.CreatePlmnIdentity(sCell.PlmnIDBytes)
		if err != nil {
			log.Warnf("%+v plmn id is not valid, err: %+v", sCell.PlmnIDBytes, err)
			continue
		}

		nrCellID, err := pdubuilder.CreateNrcellIdentity(&asn1.BitString{
			Value: sCell.NrCellIDBytes,
			Len:   sCell.NrCellIDLen,
		})
		if err != nil {
			return nil, err
		}
		nrcgi, err := pdubuilder.CreateNrcgi(plmnID, nrCellID)
		if err != nil {
			log.Warnf("%+v plmnID and %+v nrcgi is not valid to create ncgi, err: %+v", plmnID, nrCellID, err)
			continue
		}
		nrpci, err := pdubuilder.CreateNrpci(sCell.NrPCI)
		if err != nil {
			log.Warnf("%+v nrpci is not valid, err: %+v", nrpci, err)
			continue
		}

		plmnList := make([]*f1apiesv1.ServedPlmnsItem, 0)
		plmnItem, err := pdubuilder.CreateServedPlmnsItem(plmnID)
		if err != nil {
			log.Warnf("%+v plmnID is not valid, err: %+v", plmnID, err)
			continue
		}
		plmnList = append(plmnList, plmnItem)
		servedPlmns, err := pdubuilder.CreateServedPlmnsList(plmnList)
		if err != nil {
			log.Warnf("%+v plmnIDList is not valid, err: %+v", plmnList, err)
			continue
		}

		sulList := make([]*f1apiesv1.SupportedSulfreqBandItem, 0)
		sulItem, err := pdubuilder.CreateSupportedSulfreqBandItem(sCell.SulFreqBandIndicationNr)
		if err != nil {
			log.Warnf("%+v SulFreqBandIndicationNr is not valid, err: %+v", sCell.SulFreqBandIndicationNr, err)
			continue
		}
		sulList = append(sulList, sulItem)
		fbnrlist := make([]*f1apiesv1.FreqBandNrItem, 0)
		fbnritem, err := pdubuilder.CreateFreqBandNrItem(sCell.FreqBandIndicatorNr, sulList)
		if err != nil {
			log.Warnf("%+v FreqBandIndicatorNr is not valid, err: %+v", sCell.FreqBandIndicatorNr, err)
			continue
		}
		fbnrlist = append(fbnrlist, fbnritem)
		nrFreqInfo := &f1apiesv1.NrfreqInfo{
			NRarfcn:        sCell.NrArfcn,
			FreqBandListNr: fbnrlist,
		}
		transmissionBW, err := pdubuilder.CreateTransmissionBandwIDth(pdubuilder.CreateNrscsScs120(),
			pdubuilder.CreateNrnrbNrb11())
		if err != nil {
			log.Warnf("%+v nrSCS and/or %+v nrNRB is not valid, err: %+v", pdubuilder.CreateNrscsScs120(),
				pdubuilder.CreateNrnrbNrb11(), err)
			continue
		}
		fddInfo, err := pdubuilder.CreateFddInfo(nrFreqInfo, nrFreqInfo, transmissionBW, transmissionBW)
		if err != nil {
			log.Warnf("%+v nrFreqInfo and/or %+v transmissionBW is not valid, err: %+v", nrFreqInfo, transmissionBW, err)
			continue
		}
		nrmode, err := pdubuilder.CreateNrModeInfoFDd(fddInfo)
		if err != nil {
			log.Warnf("%+v fddInfo is not valid, err: %+v", fddInfo, err)
			continue
		}
		servedCellItem := &f1apiesv1.GnbDUServedCellsItem{
			ServedCellInformation: &f1apiesv1.ServedCellInformation{
				NRcgi:                          nrcgi,
				NRpci:                          nrpci,
				ServedPlmns:                    servedPlmns,
				NRModeInfo:                     nrmode,
				MeasurementTimingConfiguration: sCell.MeasureTimingConfigurationBytes,
			},
		}
		servedCellItemValue, err := pdubuilder.CreateGnbDUServedCellsItemIesValueGnbDUServedCellsItem(servedCellItem)
		if err != nil {
			log.Warnf("%+v servedCellItem is not valid, err: %+v", servedCellItem, err)
			continue
		}
		gnbDuServedCellItem, err := pdubuilder.CreateGnbDUServedCellsItemIesValue(int32(f1apv1.ProtocolIeIDGNBDUServedCellsItem),
			f1apcommondatatypesv1.Criticality_CRITICALITY_REJECT, servedCellItemValue)
		if err != nil {
			log.Warnf("%+v gnbDuServedCellItem is not valid, err: %+v", gnbDuServedCellItem, err)
			continue
		}
		gnbDuServedCellsList = append(gnbDuServedCellsList, gnbDuServedCellItem)
	}
	ie4Value, err := pdubuilder.CreateF1SetupRequestIesValueGnbDuServedCellsList(&f1appducontentsv1.GnbDUServedCellsList{
		Value: gnbDuServedCellsList,
	})
	if err != nil {
		return nil, err
	}
	ie4, err := pdubuilder.CreateF1SetupRequestIes(int32(f1apv1.ProtocolIeIDgNBDUServedCellsList),
		f1apcommondatatypesv1.Criticality_CRITICALITY_REJECT, ie4Value)
	if err != nil {
		return nil, err
	}
	list = append(list, ie4)

	f1SetupRequest, err := pdubuilder.CreateF1SetupRequest(list)
	if err != nil {
		return nil, err
	}
	newF1apPdu := &f1appdudescriptionsv1.F1ApPDu{
		F1ApPdu: &f1appdudescriptionsv1.F1ApPDu_InitiatingMessage{
			InitiatingMessage: &f1appdudescriptionsv1.InitiatingMessage{
				ProcedureCode: int32(f1apv1.ProcedureCodeIDF1Setup),
				Criticality:   f1apcommondatatypesv1.Criticality_CRITICALITY_REJECT,
				Value: &f1appdudescriptionsv1.InitiatingMessageF1ApElementaryProcedures{
					ImValues: &f1appdudescriptionsv1.InitiatingMessageF1ApElementaryProcedures_F1SetupRequest{
						F1SetupRequest: f1SetupRequest,
					},
				},
			},
		},
	}
	return encoder.PerEncodeF1ApPdu(newF1apPdu)
}
