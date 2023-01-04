// SPDX-FileCopyrightText: 2020-present Open Networking Foundation <info@opennetworking.org>
//
// SPDX-License-Identifier: Apache-2.0

package xnap

import (
	"github.com/onosproject/onos-api/go/onos/ransim/types"
	v1 "github.com/onosproject/onos-e2t/api/xnap/v1"
	xnapcommondatatypesv1 "github.com/onosproject/onos-e2t/api/xnap/v1/xnap-commondatatypes"
	xnapiesv1 "github.com/onosproject/onos-e2t/api/xnap/v1/xnap-ies"
	xnappducontentsv1 "github.com/onosproject/onos-e2t/api/xnap/v1/xnap-pdu-contents"
	xnappdudescriptionsv1 "github.com/onosproject/onos-e2t/api/xnap/v1/xnap-pdu-descriptions"
	"github.com/onosproject/onos-e2t/pkg/southbound/xnap/encoder"
	"github.com/onosproject/onos-e2t/pkg/southbound/xnap/pdubuilder"
	"github.com/onosproject/onos-lib-go/api/asn1/v1/asn1"
	"github.com/onosproject/onos-lib-go/pkg/logging"
)

type XnItemInfo struct {
	PlmnIDBytes []byte
}

type XnItemSlice struct {
	Sst []byte
	Sd  []byte
}

type XnItemAMFRegion struct {
	AmfRegionID    []byte
	AmfRegionIDLen uint32
}

type XnItemCellInfo struct {
	NCGIKey                         types.NCGI
	NrCellIDBytes                   []byte
	NrCellIDLen                     uint32
	NrPCI                           int32
	NrArfcn                         int32
	SulFreqBand                     int32
	FreqBand                        int32
	MeasureTimingConfigurationBytes []byte
	RanAC                           int32
}

var log = logging.GetLogger()

// amfregionid 0xdd
func CreateXnSetupRequest(plmnIDByte []byte, gnbIDByte []byte, tacBytes []byte, xnItemSliceList []XnItemSlice, xnAmfRegion XnItemAMFRegion, xnSCellItemInfo []XnItemCellInfo, xnNeighborCells map[types.NCGI][]XnItemCellInfo) ([]byte, error) {
	list := make([]*xnappducontentsv1.XnSetupRequestIEs, 0)

	plmnID, err := pdubuilder.CreatePlmnIdentity(plmnIDByte)
	if err != nil {
		return nil, err
	}

	gnbID, err := pdubuilder.CreateGnbIDChoiceGnbID(&asn1.BitString{
		Value: gnbIDByte,
		Len:   22,
	})
	if err != nil {
		return nil, err
	}

	globalGnBID, err := pdubuilder.CreateGlobalgNbID(plmnID, gnbID)
	if err != nil {
		return nil, err
	}

	ranNodeID, err := pdubuilder.CreateGlobalNgRAnnodeIDGNb(globalGnBID)
	if err != nil {
		return nil, err
	}

	val, err := pdubuilder.CreateXnSetupRequestIEsValueIDGlobalNgRanNodeID(ranNodeID)
	if err != nil {
		return nil, err
	}

	item1 := &xnappducontentsv1.XnSetupRequestIEs{
		Id:          int32(v1.ProtocolIeIDGlobalNGRANnodeID),
		Criticality: xnapcommondatatypesv1.Criticality_CRITICALITY_REJECT,
		Value:       val,
	}
	list = append(list, item1)

	// creating TAIsupportList
	taiSupportList := make([]*xnapiesv1.TaisupportItem, 0)

	tac, err := pdubuilder.CreateTac(tacBytes)
	if err != nil {
		return nil, err
	}

	sliceList := make([]*xnapiesv1.SNSsai, 0)
	for _, sl := range xnItemSliceList {
		tmpSnssai := &xnapiesv1.SNSsai{
			Sst: sl.Sst,
			Sd:  sl.Sd,
		}
		sliceList = append(sliceList, tmpSnssai)
	}
	taiSliceSupportList, err := pdubuilder.CreateSliceSupportList(sliceList)
	if err != nil {
		return nil, err
	}

	broadcastPlmns := make([]*xnapiesv1.BroadcastPlmninTaisupportItem, 0)
	broadcastPlmn, err := pdubuilder.CreatePlmnIdentity(plmnIDByte)
	if err != nil {
		return nil, err
	}
	broadcastPlmnItem, err := pdubuilder.CreateBroadcastPlmninTaisupportItem(broadcastPlmn, taiSliceSupportList)
	if err != nil {
		return nil, err
	}
	broadcastPlmns = append(broadcastPlmns, broadcastPlmnItem)

	taiItem, err := pdubuilder.CreateTaisupportItem(tac, broadcastPlmns)
	if err != nil {
		return nil, err
	}
	taiSupportList = append(taiSupportList, taiItem)

	val2, err := pdubuilder.CreateXnSetupRequestIEsValueIDTaisupportList(&xnapiesv1.TaisupportList{
		Value: taiSupportList,
	})
	if err != nil {
		return nil, err
	}

	item2 := &xnappducontentsv1.XnSetupRequestIEs{
		Id:          int32(v1.ProtocolIeIDTAISupportlist),
		Criticality: xnapcommondatatypesv1.Criticality_CRITICALITY_REJECT,
		Value:       val2,
	}
	list = append(list, item2)

	amfInfoList := make([]*xnapiesv1.GlobalAmfRegionInformation, 0)
	amfInfoItem, err := pdubuilder.CreateGlobalAmfRegionInformation(plmnID, &asn1.BitString{
		Value: xnAmfRegion.AmfRegionID,
		Len:   xnAmfRegion.AmfRegionIDLen,
	})
	if err != nil {
		return nil, err
	}
	amfInfoList = append(amfInfoList, amfInfoItem)

	val3, err := pdubuilder.CreateXnSetupRequestIEsValueIDAmfRegionInformation(&xnapiesv1.AmfRegionInformation{
		Value: amfInfoList,
	})
	if err != nil {
		return nil, err
	}
	item3 := &xnappducontentsv1.XnSetupRequestIEs{
		Id:          int32(v1.ProtocolIeIDAMFRegionInformation),
		Criticality: xnapcommondatatypesv1.Criticality_CRITICALITY_REJECT,
		Value:       val3,
	}
	list = append(list, item3)

	// creating list of served cells NR
	servedCellsNRList := make([]*xnapiesv1.ServedCellsNRItem, 0)
	for _, servCell := range xnSCellItemInfo {
		nrCellidentity, err := pdubuilder.CreateNrCellIdentity(&asn1.BitString{
			Value: servCell.NrCellIDBytes,
			Len:   servCell.NrCellIDLen,
		})
		if err != nil {
			log.Warnf("%+v NrCellIDBytes and/or %+v NrCellIDLen is not valid, err: %+v", servCell.NrCellIDBytes, servCell.NrCellIDLen, err)
			continue
		}
		ncgi, err := pdubuilder.CreateNrCGi(plmnID, nrCellidentity)
		if err != nil {
			return nil, err
		}
		broadcastPlmnsList := make([]*xnapiesv1.PlmnIdentity, 0)
		broadcastPlmnsList = append(broadcastPlmnsList, plmnID)

		freqBandList := make([]*xnapiesv1.NrfrequencyBandItem, 0)
		freqBand, err := pdubuilder.CreateNrfrequencyBand(servCell.FreqBand)
		if err != nil {
			return nil, err
		}
		supportedSulBandList := make([]*xnapiesv1.SupportedSulbandItem, 0)
		sulBandItem := &xnapiesv1.SupportedSulbandItem{
			SulBandItem: &xnapiesv1.SulFrequencyBand{
				Value: servCell.SulFreqBand,
			},
		}
		supportedSulBandList = append(supportedSulBandList, sulBandItem)
		freqBandItem := &xnapiesv1.NrfrequencyBandItem{
			NrFrequencyBand: freqBand,
			SupportedSulBandList: &xnapiesv1.SupportedSulbandList{
				Value: supportedSulBandList,
			},
		}
		freqBandList = append(freqBandList, freqBandItem)
		nrfreqInfo := &xnapiesv1.NrfrequencyInfo{
			NrArfcn: &xnapiesv1.Nrarfcn{
				Value: servCell.NrArfcn,
			},
			FrequencyBandList: &xnapiesv1.NrfrequencyBandList{
				Value: freqBandList,
			},
		}
		transmBW, err := pdubuilder.CreateNrtransmissionBandwIDth(pdubuilder.CreateNrscsScs120(), pdubuilder.CreateNrnrbNrb24())
		if err != nil {
			return nil, err
		}
		nrModeInfo, err := pdubuilder.CreateNrmodeInfoFdd(nrfreqInfo, nrfreqInfo, transmBW, transmBW)
		if err != nil {
			return nil, err
		}
		nrModeInfoch, err := pdubuilder.CreateNrmodeInfoFddChoice(nrModeInfo)
		if err != nil {
			return nil, err
		}
		connSupport, err := pdubuilder.CreateConnectivitySupport(pdubuilder.CreateENdcsupportConnectivitySupportNotSupported())
		if err != nil {
			return nil, err
		}

		servedCellInfoNr := &xnapiesv1.ServedCellInformationNR{
			NrPci: &xnapiesv1.Nrpci{
				Value: servCell.NrPCI,
			},
			CellId: ncgi,
			Tac: &xnapiesv1.Tac{
				Value: tacBytes,
			},
			Ranac: &xnapiesv1.Ranac{
				Value: servCell.RanAC,
			},
			BroadcastPlmn: &xnapiesv1.BroadcastPlmns{
				Value: broadcastPlmnsList,
			},
			NrModeInfo:                     nrModeInfoch,
			MeasurementTimingConfiguration: servCell.MeasureTimingConfigurationBytes,
			ConnectivitySupport:            connSupport,
		}

		neighbourInfoNrList := make([]*xnapiesv1.NeighbourInformationNRItem, 0)
		for _, nCell := range xnNeighborCells[servCell.NCGIKey] {
			neighborNrCellidentity, err := pdubuilder.CreateNrCellIdentity(&asn1.BitString{
				Value: nCell.NrCellIDBytes,
				Len:   nCell.NrCellIDLen,
			})
			if err != nil {
				log.Warnf("%+v NrCellIDBytes and/or %+v NrCellIDLen is not valid, err: %+v", nCell.NrCellIDBytes, nCell.NrCellIDBytes, err)
				continue
			}
			neighborNcgi, err := pdubuilder.CreateNrCGi(plmnID, neighborNrCellidentity)
			if err != nil {
				log.Warnf("failed to create nrcgi: %v", err)
				return nil, err
			}

			freqBandList := make([]*xnapiesv1.NrfrequencyBandItem, 0)
			freqBand, err := pdubuilder.CreateNrfrequencyBand(nCell.FreqBand)
			if err != nil {
				return nil, err
			}
			supportedSulBandList := make([]*xnapiesv1.SupportedSulbandItem, 0)
			sulBandItem := &xnapiesv1.SupportedSulbandItem{
				SulBandItem: &xnapiesv1.SulFrequencyBand{
					Value: nCell.SulFreqBand,
				},
			}
			supportedSulBandList = append(supportedSulBandList, sulBandItem)
			freqBandItem := &xnapiesv1.NrfrequencyBandItem{
				NrFrequencyBand: freqBand,
				SupportedSulBandList: &xnapiesv1.SupportedSulbandList{
					Value: supportedSulBandList,
				},
			}
			freqBandList = append(freqBandList, freqBandItem)
			nrfreqInfo := &xnapiesv1.NrfrequencyInfo{
				NrArfcn: &xnapiesv1.Nrarfcn{
					Value: nCell.NrArfcn,
				},
				FrequencyBandList: &xnapiesv1.NrfrequencyBandList{
					Value: freqBandList,
				},
			}
			connSupport, err := pdubuilder.CreateConnectivitySupport(pdubuilder.CreateENdcsupportConnectivitySupportNotSupported())
			if err != nil {
				return nil, err
			}

			neighbourInfoNrList = append(neighbourInfoNrList, &xnapiesv1.NeighbourInformationNRItem{
				NrPci: &xnapiesv1.Nrpci{
					Value: nCell.NrPCI,
				},
				NrCgi: neighborNcgi,
				Tac: &xnapiesv1.Tac{
					Value: tacBytes,
				},
				Ranac: &xnapiesv1.Ranac{
					Value: nCell.RanAC,
				},
				NrModeInfo: &xnapiesv1.NeighbourInformationNRModeInfo{
					NeighbourInformationNrModeInfo: &xnapiesv1.NeighbourInformationNRModeInfo_FddInfo{
						FddInfo: &xnapiesv1.NeighbourInformationNRModeFddinfo{
							UlNrFreqInfo: nrfreqInfo,
							DlNrFequInfo: nrfreqInfo,
						},
					},
				},
				MeasurementTimingConfiguration: nCell.MeasureTimingConfigurationBytes,
				ConnectivitySupport:            connSupport,
			})
		}
		nrItem := &xnapiesv1.ServedCellsNRItem{
			ServedCellInfoNr: servedCellInfoNr,
			NeighbourInfoNr: &xnapiesv1.NeighbourInformationNR{
				Value: neighbourInfoNrList,
			},
		}
		servedCellsNRList = append(servedCellsNRList, nrItem)
	}

	val4, err := pdubuilder.CreateXnSetupRequestIEsValueIDListOfServedCellsNr(&xnapiesv1.ServedCellsNR{
		Value: servedCellsNRList,
	})
	if err != nil {
		return nil, err
	}
	item4 := &xnappducontentsv1.XnSetupRequestIEs{
		Id:          int32(v1.ProtocolIeIDListofservedcellsNR),
		Criticality: xnapcommondatatypesv1.Criticality_CRITICALITY_REJECT,
		Value:       val4,
	}
	list = append(list, item4)

	xnSetupRequest, err := pdubuilder.CreateXnSetupRequest(list)
	if err != nil {
		return nil, err
	}
	newXnApPdu := &xnappdudescriptionsv1.XnApPDu{
		XnApPdu: &xnappdudescriptionsv1.XnApPDu_InitiatingMessage{
			InitiatingMessage: &xnappdudescriptionsv1.InitiatingMessage{
				ProcedureCode: int32(v1.ProcedureCodeIDxnSetup),
				Criticality:   xnapcommondatatypesv1.Criticality_CRITICALITY_REJECT,
				Value: &xnappdudescriptionsv1.InitiatingMessageXnApElementaryProcedures{
					ImValues: &xnappdudescriptionsv1.InitiatingMessageXnApElementaryProcedures_XnSetupRequest{
						XnSetupRequest: xnSetupRequest,
					},
				},
			},
		},
	}

	return encoder.PerEncodeXnApPdu(newXnApPdu)
}
