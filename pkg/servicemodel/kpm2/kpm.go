// SPDX-FileCopyrightText: 2020-present Open Networking Foundation <info@opennetworking.org>
//
// SPDX-License-Identifier: Apache-2.0

package kpm2

import (
	"context"
	"encoding/binary"
	"encoding/csv"
	"math"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/onosproject/onos-lib-go/api/asn1/v1/asn1"

	"github.com/onosproject/ran-simulator/pkg/utils"

	"github.com/onosproject/ran-simulator/pkg/utils/e2sm/kpm2/id/cellglobalid"

	"github.com/onosproject/ran-simulator/pkg/utils/e2sm/kpm2/measobjectitem"

	"github.com/onosproject/ran-simulator/pkg/utils/e2sm/kpm2/reportstyle"

	"github.com/onosproject/ran-simulator/pkg/utils/e2sm/kpm2/ranfuncdescription"

	"github.com/onosproject/ran-simulator/pkg/utils/e2sm/kpm2/nodeitem"

	"github.com/onosproject/ran-simulator/pkg/utils/e2sm/kpm2/measurments"

	kpm2gNBID "github.com/onosproject/ran-simulator/pkg/utils/e2sm/kpm2/id/gnbid"
	kpm2IndicationHeader "github.com/onosproject/ran-simulator/pkg/utils/e2sm/kpm2/indication"
	kpm2MessageFormat1 "github.com/onosproject/ran-simulator/pkg/utils/e2sm/kpm2/indication/messageformat1"

	ransimtypes "github.com/onosproject/onos-api/go/onos/ransim/types"
	"github.com/onosproject/onos-e2-sm/servicemodels/e2sm_kpm_v2_go/pdubuilder"
	e2smkpmv2sm "github.com/onosproject/onos-e2-sm/servicemodels/e2sm_kpm_v2_go/servicemodel"
	e2smkpmv2 "github.com/onosproject/onos-e2-sm/servicemodels/e2sm_kpm_v2_go/v2/e2sm-kpm-v2-go"
	e2apies "github.com/onosproject/onos-e2t/api/e2ap/v2/e2ap-ies"
	e2appducontents "github.com/onosproject/onos-e2t/api/e2ap/v2/e2ap-pdu-contents"
	e2aptypes "github.com/onosproject/onos-e2t/pkg/southbound/e2ap/types"
	"github.com/onosproject/onos-lib-go/pkg/errors"
	"github.com/onosproject/onos-lib-go/pkg/logging"
	"github.com/onosproject/ran-simulator/pkg/model"
	"github.com/onosproject/ran-simulator/pkg/servicemodel"
	"github.com/onosproject/ran-simulator/pkg/servicemodel/registry"
	"github.com/onosproject/ran-simulator/pkg/store/nodes"
	"github.com/onosproject/ran-simulator/pkg/store/subscriptions"
	"github.com/onosproject/ran-simulator/pkg/store/ues"
	e2apIndicationUtils "github.com/onosproject/ran-simulator/pkg/utils/e2ap/indication"
	subutils "github.com/onosproject/ran-simulator/pkg/utils/e2ap/subscription"
	subdeleteutils "github.com/onosproject/ran-simulator/pkg/utils/e2ap/subscriptiondelete"
	"google.golang.org/protobuf/proto"
)

var _ servicemodel.Client = &Client{}

var log = logging.GetLogger()

const (
	modelVersion           = "v2"
	ricStyleType           = 1
	ricStyleName           = "Periodic Report"
	ricFormatType          = 1
	ricIndMsgFormat        = 1
	ricIndHdrFormat        = 1
	ranFunctionDescription = "KPM 2.0 Monitor"
	ranFunctionShortName   = "ORAN-E2SM-KPM"
	ranFunctionE2SmOid     = "1.3.6.1.4.1.53148.1.2.2.2"
	ranFunctionInstance    = 1
)

const (
	fileFormatVersion1 string = "version1"
	senderName         string = "RAN Simulator"
	senderType         string = ""
	vendorName         string = "ONF"
)

var (
	nsDataFile   *os.File
	nsDataReader *csv.Reader
)

// Client kpm service model client
type Client struct {
	ServiceModel *registry.ServiceModel
}

// E2ConnectionUpdate implements connection update procedure
func (sm *Client) E2ConnectionUpdate(ctx context.Context, request *e2appducontents.E2ConnectionUpdate) (response *e2appducontents.E2ConnectionUpdateAcknowledge, failure *e2appducontents.E2ConnectionUpdateFailure, err error) {
	return nil, nil, errors.NewNotSupported("connection update is not supported")
}

// NewServiceModel creates a new service model
func NewServiceModel(node model.Node, model *model.Model,
	subStore *subscriptions.Subscriptions, nodeStore nodes.Store, ueStore ues.Store) (registry.ServiceModel, error) {
	kpmSm := registry.ServiceModel{
		RanFunctionID: registry.Kpm2,
		ModelName:     ranFunctionShortName,
		Revision:      1,
		OID:           ranFunctionE2SmOid,
		Version:       modelVersion,
		Node:          node,
		Model:         model,
		Subscriptions: subStore,
		Nodes:         nodeStore,
		UEs:           ueStore,
	}
	kpmClient := &Client{
		ServiceModel: &kpmSm,
	}

	kpmSm.Client = kpmClient

	plmnID := ransimtypes.NewUint24(uint32(kpmSm.Model.PlmnID))

	cells := node.Cells
	cellMeasObjectItems := make([]*e2smkpmv2.CellMeasurementObjectItem, 0)
	for _, cellNcgi := range cells {
		nci := ransimtypes.GetNCI(cellNcgi)
		ncibs := &asn1.BitString{
			Value: utils.Uint64ToBitString(uint64(nci), 36),
			Len:   36,
		}
		cellGlobalID, err := cellglobalid.
			NewGlobalNRCGIID(cellglobalid.WithPlmnID(plmnID),
				cellglobalid.WithNRCellID(ncibs)).
			Build()
		if err != nil {
			return registry.ServiceModel{}, err
		}

		cellMeasObjItem := measobjectitem.NewCellMeasObjectItem(
			measobjectitem.WithCellObjectID(strconv.FormatUint(uint64(cellNcgi), 16)),
			measobjectitem.WithCellGlobalID(cellGlobalID)).
			Build()

		cellMeasObjectItems = append(cellMeasObjectItems, cellMeasObjItem)
	}

	// Creates an indication header
	gNBID := &asn1.BitString{
		Value: utils.Uint64ToBitString(uint64(node.GnbID), 22),
		Len:   22,
	}

	globalKPMNodeID, err := kpm2gNBID.NewGlobalGNBID(
		kpm2gNBID.WithPlmnID(plmnID.Value()),
		kpm2gNBID.WithGNBIDChoice(gNBID)).Build()
	if err != nil {
		log.Error(err)
		return registry.ServiceModel{}, err
	}

	kpmNodeItem := nodeitem.NewNodeItem(
		nodeitem.WithGlobalKpmNodeID(globalKPMNodeID),
		nodeitem.WithCellMeasurementObjectItems(cellMeasObjectItems)).
		Build()

	reportKpmNodeList := make([]*e2smkpmv2.RicKpmnodeItem, 0)
	reportKpmNodeList = append(reportKpmNodeList, kpmNodeItem)

	ricEventTriggerStyleItem := pdubuilder.CreateRicEventTriggerStyleItem(ricStyleType, ricStyleName, ricFormatType)

	ricEventTriggerStyleList := make([]*e2smkpmv2.RicEventTriggerStyleItem, 0)
	ricEventTriggerStyleList = append(ricEventTriggerStyleList, ricEventTriggerStyleItem)

	measInfoActionList := e2smkpmv2.MeasurementInfoActionList{
		Value: make([]*e2smkpmv2.MeasurementInfoActionItem, 0),
	}

	for _, measType := range measTypes {
		log.Debug("Measurement Name and ID:", measType.measTypeName, measType.measTypeID)
		measInfoActionItem, _ := measurments.NewMeasurementInfoActionItem(
			measurments.WithMeasTypeName(measType.measTypeName.String()),
			measurments.WithMeasTypeID(measType.measTypeID)).Build()

		measInfoActionList.Value = append(measInfoActionList.Value, measInfoActionItem)

	}

	reportStyleItem := reportstyle.NewReportStyleItem(
		reportstyle.WithRICStyleType(ricStyleType),
		reportstyle.WithRICStyleName(ricStyleName),
		reportstyle.WithRICFormatType(ricFormatType),
		reportstyle.WithMeasInfoActionList(&measInfoActionList),
		reportstyle.WithIndicationHdrFormatType(ricIndHdrFormat),
		reportstyle.WithIndicationMsgFormatType(ricIndMsgFormat)).
		Build()

	ricReportStyleList := make([]*e2smkpmv2.RicReportStyleItem, 0)
	ricReportStyleList = append(ricReportStyleList, reportStyleItem)

	ranFuncDescPdu, err := ranfuncdescription.NewRANFunctionDescription(
		ranfuncdescription.WithRANFunctionShortName(ranFunctionShortName),
		ranfuncdescription.WithRANFunctionE2SmOID(ranFunctionE2SmOid),
		ranfuncdescription.WithRANFunctionDescription(ranFunctionDescription),
		ranfuncdescription.WithRANFunctionInstance(ranFunctionInstance),
		ranfuncdescription.WithRICKPMNodeList(reportKpmNodeList),
		ranfuncdescription.WithRICEventTriggerStyleList(ricEventTriggerStyleList),
		ranfuncdescription.WithRICReportStyleList(ricReportStyleList)).
		Build()

	if err != nil {
		log.Error(err)
		return registry.ServiceModel{}, err
	}

	protoBytes, err := proto.Marshal(ranFuncDescPdu)
	if err != nil {
		log.Error(err)
		return registry.ServiceModel{}, err
	}

	var kpm2ServiceModel e2smkpmv2sm.Kpm2ServiceModel
	ranFuncDescBytes, err := kpm2ServiceModel.RanFuncDescriptionProtoToASN1(protoBytes)
	if err != nil {
		log.Error(err)
		return registry.ServiceModel{}, err
	}
	kpmSm.Description = ranFuncDescBytes
	return kpmSm, nil
}

func float_encoder(data float32) int64 {
	var buf [4]byte
	binary.BigEndian.PutUint32(buf[:], math.Float32bits(data))
	int_data := int64(binary.BigEndian.Uint32(buf[:]))
	log.Infof("data : %+v", data)
	log.Infof("int_data : %+v", int_data)
	return int_data
}

func (sm *Client) collect(ctx context.Context,
	actionDefinition *e2smkpmv2.E2SmKpmActionDefinition,
	cellNCGI ransimtypes.NCGI) (*e2smkpmv2.MeasurementDataItem, error) {
	measInfoList := actionDefinition.GetActionDefinitionFormats().GetActionDefinitionFormat1().GetMeasInfoList()
	measRecord := e2smkpmv2.MeasurementRecord{
		Value: make([]*e2smkpmv2.MeasurementRecordItem, 0),
	}

	nsData, err := nsDataReader.Read()
	if err != nil {
		log.Debugf("can't get ns data from file")
	}

	for _, measInfo := range measInfoList.Value {
		for _, measType := range measTypes {
			if measType.measTypeName.String() == measInfo.MeasType.GetMeasName().Value {
				switch measType.measTypeName {
				case PrbUsedDL:
					if len(nsData) != 0 {
						log.Debugf("utilization : %v", nsData[4])
						utilization, err := strconv.ParseFloat(strings.TrimSpace(nsData[4]), 64)
						if err != nil {
							log.Errorf("%v", err)
						}
						log.Debugf("utilization per slice for Cell %v set for value: %+v",
							cellNCGI, utilization)
						measRecordReal := measurments.NewMeasurementRecordItemInteger(
							measurments.WithIntegerValue(float_encoder(float32(utilization))),
						).Build()
						measRecord.Value = append(measRecord.Value, measRecordReal)
					}
				case PdcpPduVolumeDL:
					if len(nsData) != 0 {
						log.Debugf("volume : %v", nsData[5])
						volume, err := strconv.ParseFloat(strings.TrimSpace(nsData[5]), 64)
						if err != nil {
							log.Errorf("%v", err)
						}
						log.Debugf("volume for Cell %v set for value: %+v",
							cellNCGI, volume)
						measRecordReal := measurments.NewMeasurementRecordItemInteger(
							measurments.WithIntegerValue(float_encoder(float32(volume))),
						).Build()
						measRecord.Value = append(measRecord.Value, measRecordReal)
					}
				case PdcpRatePerPRBDL:
					if len(nsData) != 0 {
						log.Debugf("pdcp_rate : %+v", nsData[3])
						pdcp_rate, err := strconv.ParseFloat(strings.TrimSpace(nsData[3]), 64)
						if err != nil {
							log.Errorf("%v", err)
						}
						log.Debugf("pdcp rate for Cell %v set for value: %+v",
							cellNCGI, pdcp_rate)
						measRecordReal := measurments.NewMeasurementRecordItemInteger(
							measurments.WithIntegerValue(float_encoder(float32(pdcp_rate))),
						).Build()
						measRecord.Value = append(measRecord.Value, measRecordReal)
					}
				case RRCConnMax:
					log.Debugf("Max number of UEs for Cell %v set for RRC Con Max: %v",
						cellNCGI, int64(sm.ServiceModel.UEs.MaxUEsPerCell(ctx, uint64(cellNCGI))))
					measRecordInteger := measurments.NewMeasurementRecordItemInteger(
						measurments.WithIntegerValue(int64(sm.ServiceModel.UEs.MaxUEsPerCell(ctx, uint64(cellNCGI))))).
						Build()
					measRecord.Value = append(measRecord.Value, measRecordInteger)
				case RRCConnAvg:
					log.Debugf("Avg number of UEs for Cell %v set for RRC Con Max: %v",
						cellNCGI, int64(sm.ServiceModel.UEs.LenPerCell(ctx, uint64(cellNCGI))))
					measRecordInteger := measurments.NewMeasurementRecordItemInteger(
						measurments.WithIntegerValue(int64(sm.ServiceModel.UEs.LenPerCell(ctx, uint64(cellNCGI))))).
						Build()
					measRecord.Value = append(measRecord.Value, measRecordInteger)
				default:
					measRecordNoValue := measurments.NewMeasurementRecordItemNoValue()
					measRecord.Value = append(measRecord.Value, measRecordNoValue)

				}

			}
		}

	}
	measDataItem, err := measurments.NewMeasurementDataItem(
		measurments.WithMeasurementRecord(&measRecord),
		measurments.WithIncompleteFlag(e2smkpmv2.IncompleteFlag_INCOMPLETE_FLAG_TRUE)).
		Build()
	return measDataItem, err
}

func (sm *Client) createIndicationMsgFormat1(ctx context.Context,
	cellNCGI ransimtypes.NCGI, actionDefinition *e2smkpmv2.E2SmKpmActionDefinition, interval int64) ([]byte, error) {
	log.Debug("Create Indication message format 1 based on action defs for cell:", cellNCGI)
	format1 := actionDefinition.GetActionDefinitionFormats().GetActionDefinitionFormat1()
	measInfoList := format1.GetMeasInfoList()
	measData := &e2smkpmv2.MeasurementData{
		Value: make([]*e2smkpmv2.MeasurementDataItem, 0),
	}
	granularity := actionDefinition.GetActionDefinitionFormats().GetActionDefinitionFormat1().GetGranulPeriod().Value
	numDataItems := int(interval / granularity)

	for i := 0; i < numDataItems; i++ {
		measDataItem, err := sm.collect(ctx, actionDefinition, cellNCGI)
		if err != nil {
			log.Warn(err)
			return nil, err
		}

		measData.Value = append(measData.Value, measDataItem)
	}
	subID := format1.SubscriptId.GetValue()

	// Creating an indication message format 1
	indicationMessage := kpm2MessageFormat1.NewIndicationMessage(
		kpm2MessageFormat1.WithCellObjID(strconv.FormatUint(uint64(cellNCGI), 16)),
		kpm2MessageFormat1.WithGranularity(uint32(granularity)), // TODO: check if this is a sensible conversion
		kpm2MessageFormat1.WithSubscriptionID(subID),
		kpm2MessageFormat1.WithMeasData(measData),
		kpm2MessageFormat1.WithMeasInfoList(measInfoList))

	indicationMessageBytes, err := indicationMessage.ToAsn1Bytes()
	if err != nil {
		log.Warn(err)
		return nil, err
	}

	return indicationMessageBytes, nil
}

func (sm *Client) createIndicationHeaderBytes(fileFormatVersion string) ([]byte, error) {
	// Creates an indication header
	plmnID := ransimtypes.NewUint24(uint32(sm.ServiceModel.Model.PlmnID))
	gNBID := &asn1.BitString{
		Value: utils.Uint64ToBitString(uint64(sm.ServiceModel.Node.GnbID), 22),
		Len:   22,
	}

	kpmNodeID, err := kpm2gNBID.NewGlobalGNBID(
		kpm2gNBID.WithPlmnID(plmnID.Value()),
		kpm2gNBID.WithGNBIDChoice(gNBID)).Build()

	if err != nil {
		log.Warn(err)
		return nil, err
	}
	timestamp := make([]byte, 4)
	binary.BigEndian.PutUint32(timestamp, uint32(time.Now().Unix()))
	header := kpm2IndicationHeader.NewIndicationHeader(
		kpm2IndicationHeader.WithGlobalKpmNodeID(kpmNodeID),
		kpm2IndicationHeader.WithFileFormatVersion(fileFormatVersion),
		kpm2IndicationHeader.WithSenderName(senderName),
		kpm2IndicationHeader.WithSenderType(senderType),
		kpm2IndicationHeader.WithVendorName(vendorName),
		kpm2IndicationHeader.WithTimeStamp(timestamp))

	indicationHeaderAsn1Bytes, err := header.ToAsn1Bytes()
	if err != nil {
		log.Warn(err)
		return nil, err
	}

	return indicationHeaderAsn1Bytes, nil

}

func (sm *Client) sendRicIndicationFormat1(ctx context.Context, ncgi ransimtypes.NCGI,
	subscription *subutils.Subscription,
	actionDefinitions []*e2smkpmv2.E2SmKpmActionDefinition,
	interval int64) error {
	// Creates and sends indication message format 1
	subID := subscriptions.NewID(subscription.GetRicInstanceID(), subscription.GetReqID(), subscription.GetRanFuncID())
	sub, err := sm.ServiceModel.Subscriptions.Get(subID)
	if err != nil {
		return err
	}

	indicationHeaderBytes, err := sm.createIndicationHeaderBytes(fileFormatVersion1)
	if err != nil {
		log.Warn(err)
		return err
	}

	for _, actionDefinition := range actionDefinitions {
		format1 := actionDefinition.GetActionDefinitionFormats().GetActionDefinitionFormat1()
		if format1 != nil {
			cellObjectID := format1.GetCellObjId().Value
			if cellObjectID == strconv.FormatUint(uint64(ncgi), 16) {
				log.Debug("Sending indication message for Cell with ID:", cellObjectID)
				indicationMessageBytes, err := sm.createIndicationMsgFormat1(ctx, ncgi, actionDefinition, interval)
				if err != nil {
					return err
				}

				indication := e2apIndicationUtils.NewIndication(
					e2apIndicationUtils.WithRicInstanceID(subscription.GetRicInstanceID()),
					e2apIndicationUtils.WithRanFuncID(subscription.GetRanFuncID()),
					e2apIndicationUtils.WithRequestID(subscription.GetReqID()),
					e2apIndicationUtils.WithIndicationHeader(indicationHeaderBytes),
					e2apIndicationUtils.WithIndicationMessage(indicationMessageBytes))

				ricIndication, err := indication.Build()
				if err != nil {
					log.Error("creating indication message is failed for Cell with ID", ncgi, err)
					return err
				}

				err = sub.E2Channel.RICIndication(ctx, ricIndication)
				if err != nil {
					return err
				}
			}
		}
	}

	return nil
}

//func (sm *Client) sendRicIndication(ctx context.Context,
//	subscription *subutils.Subscription, actionDefinitions []*e2smkpmv2.E2SmKpmActionDefinition, interval int64) error {
//	node := sm.ServiceModel.Node
//	// Creates and sends an indication message for each cell in the node that are also specified in Action Definition
//	for _, ncgi := range node.Cells {
//		err := sm.sendRicIndicationFormat1(ctx, ncgi, subscription, actionDefinitions, interval)
//		if err != nil {
//			log.Error(err)
//			return err
//		}
//	}
//	return nil
//}

func (sm *Client) reportIndication(ctx context.Context, interval int64, subscription *subutils.Subscription, actionDefinitions []*e2smkpmv2.E2SmKpmActionDefinition) error {
	subID := subscriptions.NewID(subscription.GetRicInstanceID(), subscription.GetReqID(), subscription.GetRanFuncID())

	intervalDuration := time.Duration(interval)
	sub, err := sm.ServiceModel.Subscriptions.Get(subID)
	if err != nil {
		log.Warn(err)
		return err
	}
	sub.Ticker = time.NewTicker(intervalDuration * time.Millisecond)

	nsDataFile, err = os.Open("/usr/local/datasets/cell.csv")
	if err != nil {
		log.Error("can't open the file")
	}

	nsDataReader = csv.NewReader(nsDataFile)

	defer nsDataFile.Close()

	node_cell := sm.ServiceModel.Node.Cells
	var index int = 0
	var node_cell_length int = len(node_cell)

	for {
		select {
		case <-sub.Ticker.C:
			log.Debug("Sending Indication Report for subscription:", sub.ID)
			// err = sm.sendRicIndication(ctx, subscription, actionDefinitions, interval)
			// if err != nil {
			// 	log.Error("creating indication message is failed", err)
			// 	return err
			// }
			err = sm.sendRicIndicationFormat1(ctx, node_cell[index], subscription, actionDefinitions, interval)
			if err != nil {
				log.Error(err)
				return err
			}
			index++
			if index >= node_cell_length {
				index = 0
			}

		case <-sub.E2Channel.Context().Done():
			log.Debug("E2 channel context is done")
			sub.Ticker.Stop()
			return nil

		}
	}
}

// RICControl implements control handler for kpm service model
func (sm *Client) RICControl(ctx context.Context, request *e2appducontents.RiccontrolRequest) (response *e2appducontents.RiccontrolAcknowledge, failure *e2appducontents.RiccontrolFailure, err error) {
	return nil, nil, errors.New(errors.NotSupported, "Control operation is not supported")
}

// RICSubscription implements subscription handler for kpm service model
func (sm *Client) RICSubscription(ctx context.Context, request *e2appducontents.RicsubscriptionRequest) (response *e2appducontents.RicsubscriptionResponse, failure *e2appducontents.RicsubscriptionFailure, err error) {
	log.Infof("RIC Subscription request received for e2 node %d and service model %s:", sm.ServiceModel.Node.GnbID, sm.ServiceModel.ModelName)
	var ricActionsAccepted []*e2aptypes.RicActionID
	ricActionsNotAdmitted := make(map[e2aptypes.RicActionID]*e2apies.Cause)
	actionList := subutils.GetRicActionToBeSetupList(request)
	reqID, err := subutils.GetRequesterID(request)
	if err != nil {
		return nil, nil, err
	}
	ranFuncID, err := subutils.GetRanFunctionID(request)
	if err != nil {
		return nil, nil, err
	}
	ricInstanceID, err := subutils.GetRicInstanceID(request)
	if err != nil {
		return nil, nil, err
	}

	for _, action := range actionList {
		actionID := e2aptypes.RicActionID(action.GetValue().GetRicactionToBeSetupItem().GetRicActionId().GetValue())
		actionType := action.GetValue().GetRicactionToBeSetupItem().GetRicActionType()
		// kpm service model supports report action and should be added to the
		// list of accepted actions
		if actionType == e2apies.RicactionType_RICACTION_TYPE_REPORT {
			ricActionsAccepted = append(ricActionsAccepted, &actionID)
		}
		// kpm service model does not support INSERT and POLICY actions and
		// should be added into the list of not admitted actions
		if actionType == e2apies.RicactionType_RICACTION_TYPE_INSERT ||
			actionType == e2apies.RicactionType_RICACTION_TYPE_POLICY {
			cause := &e2apies.Cause{
				Cause: &e2apies.Cause_RicRequest{
					RicRequest: e2apies.CauseRicrequest_CAUSE_RICREQUEST_ACTION_NOT_SUPPORTED,
				},
			}
			ricActionsNotAdmitted[actionID] = cause
		}
	}

	// At least one required action must be accepted otherwise sends a subscription failure response
	if len(ricActionsAccepted) == 0 {
		log.Warn("no action is accepted")
		cause := &e2apies.Cause{
			Cause: &e2apies.Cause_RicRequest{
				RicRequest: e2apies.CauseRicrequest_CAUSE_RICREQUEST_ACTION_NOT_SUPPORTED,
			},
		}
		subscription := subutils.NewSubscription(
			subutils.WithRequestID(*reqID),
			subutils.WithRanFuncID(*ranFuncID),
			subutils.WithRicInstanceID(*ricInstanceID),
			subutils.WithCause(cause))
		subscriptionFailure, err := subscription.BuildSubscriptionFailure()
		if err != nil {
			return nil, nil, err
		}
		return nil, subscriptionFailure, nil
	}

	reportInterval, err := sm.getReportPeriod(request)
	if err != nil {
		log.Warn(err)
		cause := &e2apies.Cause{
			Cause: &e2apies.Cause_RicRequest{
				RicRequest: e2apies.CauseRicrequest_CAUSE_RICREQUEST_UNSPECIFIED,
			},
		}
		subscription := subutils.NewSubscription(
			subutils.WithRequestID(*reqID),
			subutils.WithRanFuncID(*ranFuncID),
			subutils.WithRicInstanceID(*ricInstanceID),
			subutils.WithCause(cause))
		subscriptionFailure, err := subscription.BuildSubscriptionFailure()
		if err != nil {
			log.Warn(err)
			return nil, subscriptionFailure, nil
		}
		return nil, subscriptionFailure, nil
	}

	actionDefinitions, err := sm.getActionDefinition(actionList, ricActionsAccepted)
	if err != nil {
		log.Warn(err)
		cause := &e2apies.Cause{
			Cause: &e2apies.Cause_RicRequest{
				RicRequest: e2apies.CauseRicrequest_CAUSE_RICREQUEST_UNSPECIFIED,
			},
		}
		subscription := subutils.NewSubscription(
			subutils.WithRequestID(*reqID),
			subutils.WithRanFuncID(*ranFuncID),
			subutils.WithRicInstanceID(*ricInstanceID),
			subutils.WithCause(cause))
		subscriptionFailure, err := subscription.BuildSubscriptionFailure()
		if err != nil {
			log.Warn(err)
			return nil, subscriptionFailure, nil
		}
		return nil, subscriptionFailure, nil
	}

	subscription := subutils.NewSubscription(
		subutils.WithRequestID(*reqID),
		subutils.WithRanFuncID(*ranFuncID),
		subutils.WithRicInstanceID(*ricInstanceID),
		subutils.WithActionsAccepted(ricActionsAccepted),
		subutils.WithActionsNotAdmitted(ricActionsNotAdmitted))
	subscriptionResponse, err := subscription.BuildSubscriptionResponse()
	if err != nil {
		log.Warn(err)
		cause := &e2apies.Cause{
			Cause: &e2apies.Cause_RicRequest{
				RicRequest: e2apies.CauseRicrequest_CAUSE_RICREQUEST_UNSPECIFIED,
			},
		}
		subscription := subutils.NewSubscription(
			subutils.WithRequestID(*reqID),
			subutils.WithRanFuncID(*ranFuncID),
			subutils.WithRicInstanceID(*ricInstanceID),
			subutils.WithCause(cause))
		subscriptionFailure, err := subscription.BuildSubscriptionFailure()
		if err != nil {
			log.Warn(err)
			return nil, subscriptionFailure, nil
		}
		return nil, subscriptionFailure, nil
	}
	go func() {
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()
		err := sm.reportIndication(ctx, reportInterval, subscription, actionDefinitions)
		if err != nil {
			return
		}
	}()
	return subscriptionResponse, nil, nil

}

// RICSubscriptionDelete implements subscription delete handler for kpm service model
func (sm *Client) RICSubscriptionDelete(ctx context.Context, request *e2appducontents.RicsubscriptionDeleteRequest) (response *e2appducontents.RicsubscriptionDeleteResponse, failure *e2appducontents.RicsubscriptionDeleteFailure, err error) {
	log.Infof("RIC subscription delete request is received for e2 node %d and  service model %s:", sm.ServiceModel.Node.GnbID, sm.ServiceModel.ModelName)
	reqID, err := subdeleteutils.GetRequesterID(request)
	if err != nil {
		return nil, nil, err
	}
	ranFuncID, err := subdeleteutils.GetRanFunctionID(request)
	if err != nil {
		return nil, nil, err
	}
	ricInstanceID, err := subdeleteutils.GetRicInstanceID(request)
	if err != nil {
		return nil, nil, err
	}
	subID := subscriptions.NewID(*ricInstanceID, *reqID, *ranFuncID)
	sub, err := sm.ServiceModel.Subscriptions.Get(subID)
	if err != nil {
		return nil, nil, err
	}
	subscriptionDelete := subdeleteutils.NewSubscriptionDelete(
		subdeleteutils.WithRequestID(*reqID),
		subdeleteutils.WithRanFuncID(*ranFuncID),
		subdeleteutils.WithRicInstanceID(*ricInstanceID))
	subDeleteResponse, err := subscriptionDelete.BuildSubscriptionDeleteResponse()
	if err != nil {
		return nil, nil, err
	}
	// Stops the goroutine sending the indication messages
	if sub.Ticker != nil {
		sub.Ticker.Stop()
	}
	return subDeleteResponse, nil, nil
}
