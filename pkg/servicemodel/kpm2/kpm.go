// SPDX-FileCopyrightText: 2020-present Open Networking Foundation <info@opennetworking.org>
//
// SPDX-License-Identifier: LicenseRef-ONF-Member-1.0

package kpm2

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/onosproject/ran-simulator/pkg/utils/e2sm/kpm2/measurments"

	"github.com/onosproject/ran-simulator/pkg/utils/e2sm/kpm2/labelinfo"

	e2smtypes "github.com/onosproject/onos-api/go/onos/e2t/e2sm"
	kpm2IndicationHeader "github.com/onosproject/ran-simulator/pkg/utils/e2sm/kpm2/indication"
	kpm2MessageFormat1 "github.com/onosproject/ran-simulator/pkg/utils/e2sm/kpm2/indication/messageformat1"
	kpm2NodeIDutils "github.com/onosproject/ran-simulator/pkg/utils/e2sm/kpm2/nodeid"

	ransimtypes "github.com/onosproject/onos-api/go/onos/ransim/types"
	"github.com/onosproject/onos-e2-sm/servicemodels/e2sm_kpm_v2/pdubuilder"
	e2smkpmv2 "github.com/onosproject/onos-e2-sm/servicemodels/e2sm_kpm_v2/v2/e2sm-kpm-v2"
	e2apies "github.com/onosproject/onos-e2t/api/e2ap/v1beta2/e2ap-ies"
	e2appducontents "github.com/onosproject/onos-e2t/api/e2ap/v1beta2/e2ap-pdu-contents"
	e2aptypes "github.com/onosproject/onos-e2t/pkg/southbound/e2ap101/types"
	"github.com/onosproject/onos-lib-go/pkg/errors"
	"github.com/onosproject/onos-lib-go/pkg/logging"
	"github.com/onosproject/ran-simulator/pkg/model"
	"github.com/onosproject/ran-simulator/pkg/modelplugins"
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

var log = logging.GetLogger("sm", "kpm2")

const (
	modelVersion           = "v2"
	ricStyleType           = 1
	ricStyleName           = "Periodic Report"
	ricFormatType          = 5
	ricIndMsgFormat        = 1
	ricIndHdrFormat        = 1
	ranFunctionDescription = "KPM 2.0 Monitor"
	ranFunctionShortName   = "ORAN-E2SM-KPM"
	ranFunctionE2SmOid     = "1.3.6.1.4.1.53148.1.2.2.2"
	ranFunctionInstance    = 1
)

// TODO hard coded values for indication messages and should be replaced by
// real values
const (
	fileFormatVersion string = "txt"
	senderName        string = "ONF"
	senderType        string = "test-type"
	vendorName        string = "ONF"
)

// Client kpm service model client
type Client struct {
	ServiceModel *registry.ServiceModel
}

// NewServiceModel creates a new service model
func NewServiceModel(node model.Node, model *model.Model, modelPluginRegistry modelplugins.ModelRegistry,
	subStore *subscriptions.Subscriptions, nodeStore nodes.Store, ueStore ues.Store) (registry.ServiceModel, error) {
	kpmSm := registry.ServiceModel{
		RanFunctionID:       registry.Kpm2,
		ModelName:           ranFunctionShortName,
		Revision:            1,
		OID:                 ranFunctionE2SmOid,
		Version:             modelVersion,
		ModelPluginRegistry: modelPluginRegistry,
		Node:                node,
		Model:               model,
		Subscriptions:       subStore,
		Nodes:               nodeStore,
		UEs:                 ueStore,
	}
	kpmClient := &Client{
		ServiceModel: &kpmSm,
	}

	kpmSm.Client = kpmClient

	plmnID := ransimtypes.NewUint24(uint32(kpmSm.Model.PlmnID)).ToBytes()
	bs := e2smkpmv2.BitString{
		Value: 0x9bcd4,
		Len:   22,
	}
	// TODO - Fix hardcoded cellID
	cellGlobalID, err := pdubuilder.CreateCellGlobalIDNRCGI(plmnID, 0xabcdef012<<28) // 36 bit
	if err != nil {
		log.Error(err)
		return registry.ServiceModel{}, err
	}

	// TODO - Fix hardcoded cellObjID
	var cellObjID = "ONF"
	cellMeasObjItem := pdubuilder.CreateCellMeasurementObjectItem(cellObjID, cellGlobalID)

	// TODO - Fix hardcoded IDs
	var gnbCuUpID int64 = 12345
	var gnbDuID int64 = 6789
	globalKpmnodeID, err := pdubuilder.CreateGlobalKpmnodeIDgNBID(&bs, plmnID, gnbCuUpID, gnbDuID)
	if err != nil {
		log.Error(err)
		return registry.ServiceModel{}, err
	}

	cmol := make([]*e2smkpmv2.CellMeasurementObjectItem, 0)
	cmol = append(cmol, cellMeasObjItem)

	kpmNodeItem := pdubuilder.CreateRicKpmnodeItem(globalKpmnodeID, cmol)

	rknl := make([]*e2smkpmv2.RicKpmnodeItem, 0)
	rknl = append(rknl, kpmNodeItem)

	retsi := pdubuilder.CreateRicEventTriggerStyleItem(ricStyleType, ricStyleName, ricFormatType)

	retsl := make([]*e2smkpmv2.RicEventTriggerStyleItem, 0)
	retsl = append(retsl, retsi)

	measInfoActionList := e2smkpmv2.MeasurementInfoActionList{
		Value: make([]*e2smkpmv2.MeasurementInfoActionItem, 0),
	}

	var measTypeName1 = "RRC.ConnEstabAtt.Tot"
	var measTypeID1 int32 = 1
	measInfoActionItem1 := pdubuilder.CreateMeasurementInfoActionItem(measTypeName1, measTypeID1)
	measInfoActionList.Value = append(measInfoActionList.Value, measInfoActionItem1)

	var measTypeName2 = "RRC.ConnEstabSucc.Tot"
	var measTypeID2 int32 = 2
	measInfoActionItem2 := pdubuilder.CreateMeasurementInfoActionItem(measTypeName2, measTypeID2)
	measInfoActionList.Value = append(measInfoActionList.Value, measInfoActionItem2)

	var measTypeName3 = "RRC.ConnReEstabAtt.Tot"
	var measTypeID3 int32 = 3
	measInfoActionItem3 := pdubuilder.CreateMeasurementInfoActionItem(measTypeName3, measTypeID3)
	measInfoActionList.Value = append(measInfoActionList.Value, measInfoActionItem3)

	var measTypeName4 = "RRC.ConnReEstabAtt.reconfigFail"
	var measTypeID4 int32 = 4
	measInfoActionItem4 := pdubuilder.CreateMeasurementInfoActionItem(measTypeName4, measTypeID4)
	measInfoActionList.Value = append(measInfoActionList.Value, measInfoActionItem4)

	var measTypeName5 = "RRC.ConnReEstabAtt.HOFail"
	var measTypeID5 int32 = 5
	measInfoActionItem5 := pdubuilder.CreateMeasurementInfoActionItem(measTypeName5, measTypeID5)
	measInfoActionList.Value = append(measInfoActionList.Value, measInfoActionItem5)

	var measTypeName6 = "RRC.ConnReEstabAtt.Other"
	var measTypeID6 int32 = 6
	measInfoActionItem6 := pdubuilder.CreateMeasurementInfoActionItem(measTypeName6, measTypeID6)
	measInfoActionList.Value = append(measInfoActionList.Value, measInfoActionItem6)

	var measTypeName7 = "RRC.Conn.Avg"
	var measTypeID7 int32 = 7
	measInfoActionItem7 := pdubuilder.CreateMeasurementInfoActionItem(measTypeName7, measTypeID7)
	measInfoActionList.Value = append(measInfoActionList.Value, measInfoActionItem7)

	var measTypeName8 = "RRC.Conn.Max"
	var measTypeID8 int32 = 8
	measInfoActionItem8 := pdubuilder.CreateMeasurementInfoActionItem(measTypeName8, measTypeID8)
	measInfoActionList.Value = append(measInfoActionList.Value, measInfoActionItem8)

	rrsi := pdubuilder.CreateRicReportStyleItem(ricStyleType, ricStyleName, ricFormatType, &measInfoActionList, ricIndHdrFormat, ricIndMsgFormat)

	rrsl := make([]*e2smkpmv2.RicReportStyleItem, 0)
	rrsl = append(rrsl, rrsi)

	ranFuncDescPdu, err := pdubuilder.CreateE2SmKpmRanfunctionDescription(
		ranFunctionShortName,
		ranFunctionE2SmOid,
		ranFunctionDescription,
		ranFunctionInstance,
		rknl,
		retsl,
		rrsl)
	if err != nil {
		log.Error(err)
		return registry.ServiceModel{}, err
	}

	protoBytes, err := proto.Marshal(ranFuncDescPdu)
	if err != nil {
		log.Error(err)
		return registry.ServiceModel{}, err
	}
	kpmModelPlugin, _ := modelPluginRegistry.GetPlugin(ranFunctionE2SmOid)
	if kpmModelPlugin == nil {
		return registry.ServiceModel{}, errors.New(errors.Invalid, "model plugin is nil")
	}
	ranFuncDescBytes, err := kpmModelPlugin.RanFuncDescriptionProtoToASN1(protoBytes)
	if err != nil {
		log.Error(err)
		return registry.ServiceModel{}, err
	}
	kpmSm.Description = ranFuncDescBytes
	return kpmSm, nil
}

func (sm *Client) reportIndication(ctx context.Context, interval int32, subscription *subutils.Subscription) error {
	subID := subscriptions.NewID(subscription.GetRicInstanceID(), subscription.GetReqID(), subscription.GetRanFuncID())
	gNbID, err := strconv.ParseUint(fmt.Sprintf("%d", sm.ServiceModel.Node.EnbID), 10, 64)
	if err != nil {
		log.Warn(err)
		return err
	}
	// Creates an indication header
	plmnID := ransimtypes.NewUint24(uint32(sm.ServiceModel.Model.PlmnID))

	bs := e2smkpmv2.BitString{
		Value: 0x9bcd4,
		Len:   22,
	}

	kpmNodeID, err := kpm2NodeIDutils.NewGlobalGNBID(
		kpm2NodeIDutils.WithPlmnID(plmnID.Value()),
		kpm2NodeIDutils.WithGNBCuUpID(int64(gNbID)),
		kpm2NodeIDutils.WithGNBIDChoice(&bs),
		kpm2NodeIDutils.WithGNBDuID(int64(gNbID))).Build()

	if err != nil {
		log.Warn(err)
		return err
	}

	header := kpm2IndicationHeader.NewIndicationHeader(
		kpm2IndicationHeader.WithGlobalKpmNodeID(kpmNodeID),
		kpm2IndicationHeader.WithFileFormatVersion(fileFormatVersion),
		kpm2IndicationHeader.WithSenderName(senderName),
		kpm2IndicationHeader.WithSenderType(senderType),
		kpm2IndicationHeader.WithVendorName(vendorName),
		kpm2IndicationHeader.WithTimeStamp([]byte{0x21, 0x22, 0x23, 0x24}))

	kpmModelPlugin, _ := sm.ServiceModel.ModelPluginRegistry.GetPlugin(e2smtypes.OID(sm.ServiceModel.OID))
	indicationHeaderAsn1Bytes, err := header.ToAsn1Bytes(kpmModelPlugin)
	if err != nil {
		log.Warn(err)
		return err
	}

	var integerValue int64 = 12345
	var realValue float64 = 6789
	var cellObjID = "onf"
	var granularity int32 = 21
	var fiveQI int32 = 10
	var qfi int32 = 62
	var qci int32 = 15
	var qciMin int32 = 1
	var qciMax int32 = 15
	var arpMax int32 = 15
	var arpMin int32 = 10
	var bitrateRange int32 = 251
	var layerMuMimo int32 = 5
	var distX int32 = 123
	var distY int32 = 456
	var distZ int32 = 789
	startEndIndication := e2smkpmv2.StartEndInd_START_END_IND_START
	sst := []byte{0x01}
	sd := []byte{0x01, 0x02, 0x03}

	// Creates label information
	labelInfo, err := labelinfo.NewLabelInfo(labelinfo.WithFiveQI(fiveQI),
		labelinfo.WithPlmnID(plmnID.Value()),
		labelinfo.WithArpMin(arpMin),
		labelinfo.WithArpMax(arpMax),
		labelinfo.WithBitRateRange(bitrateRange),
		labelinfo.WithDistX(distX),
		labelinfo.WithDistY(distY),
		labelinfo.WithDistZ(distZ),
		labelinfo.WithLayerMuMimo(layerMuMimo),
		labelinfo.WithQCI(qci),
		labelinfo.WithQCIMin(qciMin),
		labelinfo.WithQCIMax(qciMax),
		labelinfo.WithQFI(qfi),
		labelinfo.WithStartEndIndication(startEndIndication),
		labelinfo.WithSD(sd),
		labelinfo.WithSST(sst))

	if err != nil {
		log.Warn(err)
		return err
	}
	labelInfoItem, err := labelInfo.Build()
	if err != nil {
		log.Warn(err)
		return err
	}

	labelInfoList := e2smkpmv2.LabelInfoList{
		Value: make([]*e2smkpmv2.LabelInfoItem, 0),
	}
	labelInfoList.Value = append(labelInfoList.Value, labelInfoItem)

	// Creates measurement type name
	measTypeName, err := measurments.NewMeasurementTypeMeasName(
		measurments.WithMeasurementName("test")).
		Build()
	if err != nil {
		log.Warn(err)
		return err
	}

	if err != nil {
		log.Warn(err)
		return err
	}

	// Creates measurement info item
	measInfoItem, err := measurments.NewMeasurementInfoItem(
		measurments.WithMeasType(measTypeName),
		measurments.WithLabelInfoList(&labelInfoList)).Build()
	if err != nil {
		log.Warn(err)
		return err
	}

	// Creates measurement info list
	measInfoList := e2smkpmv2.MeasurementInfoList{
		Value: make([]*e2smkpmv2.MeasurementInfoItem, 0),
	}
	measInfoList.Value = append(measInfoList.Value, measInfoItem)

	// Creates meas record
	measRecord := e2smkpmv2.MeasurementRecord{
		Value: make([]*e2smkpmv2.MeasurementRecordItem, 0),
	}

	measRecordInteger := measurments.NewMeasurementRecordItemInteger(measurments.WithIntegerValue(integerValue)).Build()
	measRecordReal := measurments.NewMeasurementRecordItemReal(measurments.WithRealValue(realValue)).Build()
	measRecordNoValue := measurments.NewMeasurementRecordItemNoValue()

	measRecord.Value = append(measRecord.Value, measRecordInteger)
	measRecord.Value = append(measRecord.Value, measRecordReal)
	measRecord.Value = append(measRecord.Value, measRecordNoValue)

	measDataItem, err := measurments.NewMeasurementDataItem(
		measurments.WithMeasurementRecord(&measRecord),
		measurments.WithIncompleteFlag(e2smkpmv2.IncompleteFlag_INCOMPLETE_FLAG_TRUE)).Build()
	if err != nil {
		log.Warn(err)
		return err
	}

	measData := e2smkpmv2.MeasurementData{
		Value: make([]*e2smkpmv2.MeasurementDataItem, 0),
	}
	measData.Value = append(measData.Value, measDataItem)

	// Creating an indication message format 1
	indicationMessage := kpm2MessageFormat1.NewIndicationMessage(
		kpm2MessageFormat1.WithCellObjID(cellObjID),
		kpm2MessageFormat1.WithGranularity(granularity),
		kpm2MessageFormat1.WithSubscriptionID(123456),
		kpm2MessageFormat1.WithMeasData(&measData),
		kpm2MessageFormat1.WithMeasInfoList(&measInfoList))

	indicationMessageBytes, err := indicationMessage.ToAsn1Bytes(kpmModelPlugin)
	if err != nil {
		log.Warn(err)
		return err
	}

	intervalDuration := time.Duration(interval)
	sub, err := sm.ServiceModel.Subscriptions.Get(subID)
	if err != nil {
		log.Warn(err)
		return err
	}
	sub.Ticker = time.NewTicker(intervalDuration * time.Millisecond)
	for {
		select {
		case <-sub.Ticker.C:
			log.Debug("Sending Indication Report for subscription:", sub.ID)
			indication := e2apIndicationUtils.NewIndication(
				e2apIndicationUtils.WithRicInstanceID(subscription.GetRicInstanceID()),
				e2apIndicationUtils.WithRanFuncID(subscription.GetRanFuncID()),
				e2apIndicationUtils.WithRequestID(subscription.GetReqID()),
				e2apIndicationUtils.WithIndicationHeader(indicationHeaderAsn1Bytes),
				e2apIndicationUtils.WithIndicationMessage(indicationMessageBytes))

			ricIndication, err := indication.Build()
			if err != nil {
				log.Error("creating indication message is failed", err)
				return err
			}

			err = sub.E2Channel.RICIndication(ctx, ricIndication)
			if err != nil {
				log.Error("Sending indication report is failed:", err)
				return err
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
	log.Infof("RIC Subscription request received for e2 node %d and service model %s:", sm.ServiceModel.Node.EnbID, sm.ServiceModel.ModelName)
	var ricActionsAccepted []*e2aptypes.RicActionID
	ricActionsNotAdmitted := make(map[e2aptypes.RicActionID]*e2apies.Cause)
	actionList := subutils.GetRicActionToBeSetupList(request)
	reqID := subutils.GetRequesterID(request)
	ranFuncID := subutils.GetRanFunctionID(request)
	ricInstanceID := subutils.GetRicInstanceID(request)

	for _, action := range actionList {
		actionID := e2aptypes.RicActionID(action.Value.RicActionId.Value)
		actionType := action.Value.RicActionType
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
					RicRequest: e2apies.CauseRic_CAUSE_RIC_ACTION_NOT_SUPPORTED,
				},
			}
			ricActionsNotAdmitted[actionID] = cause
		}
	}
	subscription := subutils.NewSubscription(
		subutils.WithRequestID(reqID),
		subutils.WithRanFuncID(ranFuncID),
		subutils.WithRicInstanceID(ricInstanceID),
		subutils.WithActionsAccepted(ricActionsAccepted),
		subutils.WithActionsNotAdmitted(ricActionsNotAdmitted))

	// At least one required action must be accepted otherwise sends a subscription failure response
	if len(ricActionsAccepted) == 0 {
		log.Debug("no action is accepted")
		subscriptionFailure, err := subscription.BuildSubscriptionFailure()
		if err != nil {
			return nil, nil, err
		}
		return nil, subscriptionFailure, nil
	}

	reportInterval, err := sm.getReportPeriod(request)
	if err != nil {
		subscriptionFailure, err := subscription.BuildSubscriptionFailure()
		if err != nil {
			return nil, nil, err
		}
		return nil, subscriptionFailure, nil
	}

	subscriptionResponse, err := subscription.BuildSubscriptionResponse()
	if err != nil {
		return nil, nil, err
	}
	go func() {
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()
		err := sm.reportIndication(ctx, reportInterval, subscription)
		if err != nil {
			return
		}
	}()
	return subscriptionResponse, nil, nil

}

// RICSubscriptionDelete implements subscription delete handler for kpm service model
func (sm *Client) RICSubscriptionDelete(ctx context.Context, request *e2appducontents.RicsubscriptionDeleteRequest) (response *e2appducontents.RicsubscriptionDeleteResponse, failure *e2appducontents.RicsubscriptionDeleteFailure, err error) {
	log.Infof("RIC subscription delete request is received for e2 node %d and  service model %s:", sm.ServiceModel.Node.EnbID, sm.ServiceModel.ModelName)
	reqID := subdeleteutils.GetRequesterID(request)
	ranFuncID := subdeleteutils.GetRanFunctionID(request)
	ricInstanceID := subdeleteutils.GetRicInstanceID(request)
	subID := subscriptions.NewID(ricInstanceID, reqID, ranFuncID)
	sub, err := sm.ServiceModel.Subscriptions.Get(subID)
	if err != nil {
		return nil, nil, err
	}
	subscriptionDelete := subdeleteutils.NewSubscriptionDelete(
		subdeleteutils.WithRequestID(reqID),
		subdeleteutils.WithRanFuncID(ranFuncID),
		subdeleteutils.WithRicInstanceID(ricInstanceID))
	subDeleteResponse, err := subscriptionDelete.BuildSubscriptionDeleteResponse()
	if err != nil {
		return nil, nil, err
	}
	// Stops the goroutine sending the indication messages
	sub.Ticker.Stop()
	return subDeleteResponse, nil, nil
}
