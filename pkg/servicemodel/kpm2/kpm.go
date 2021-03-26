// SPDX-FileCopyrightText: 2020-present Open Networking Foundation <info@opennetworking.org>
//
// SPDX-License-Identifier: LicenseRef-ONF-Member-1.0

package kpm2

import (
	"context"
	"fmt"
	"math/rand"
	"strconv"
	"time"

	"github.com/onosproject/ran-simulator/pkg/utils/e2sm/kpm2/measurments"

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

var measTypes = map[string]int32{
	"RRC.ConnEstabAtt.Tot":            1,
	"RRC.ConnEstabSucc.Tot":           2,
	"RRC.ConnReEstabAtt.Tot":          3,
	"RRC.ConnReEstabAtt.reconfigFail": 4,
	"RRC.ConnReEstabAtt.HOFail":       5,
	"RRC.ConnReEstabAtt.Other":        6,
	"RRC.Conn.Avg":                    7,
	"RRC.Conn.Max":                    8,
}

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

	cells := node.Cells
	cmol := make([]*e2smkpmv2.CellMeasurementObjectItem, 0)
	for _, cellEcgi := range cells {
		cellGlobalID, err := pdubuilder.CreateCellGlobalIDNRCGI(plmnID, uint64(cellEcgi)<<28)
		if err != nil {
			log.Error(err)
			return registry.ServiceModel{}, err
		}

		cellMeasObjItem := pdubuilder.CreateCellMeasurementObjectItem(strconv.FormatUint(uint64(cellEcgi), 10), cellGlobalID)
		cmol = append(cmol, cellMeasObjItem)
	}

	// TODO - Fix hardcoded IDs
	var gnbCuUpID int64 = 12345
	var gnbDuID int64 = 6789
	globalKpmnodeID, err := pdubuilder.CreateGlobalKpmnodeIDgNBID(&bs, plmnID, gnbCuUpID, gnbDuID)
	if err != nil {
		log.Error(err)
		return registry.ServiceModel{}, err
	}

	kpmNodeItem := pdubuilder.CreateRicKpmnodeItem(globalKpmnodeID, cmol)

	rknl := make([]*e2smkpmv2.RicKpmnodeItem, 0)
	rknl = append(rknl, kpmNodeItem)

	retsi := pdubuilder.CreateRicEventTriggerStyleItem(ricStyleType, ricStyleName, ricFormatType)

	retsl := make([]*e2smkpmv2.RicEventTriggerStyleItem, 0)
	retsl = append(retsl, retsi)

	measInfoActionList := e2smkpmv2.MeasurementInfoActionList{
		Value: make([]*e2smkpmv2.MeasurementInfoActionItem, 0),
	}

	for measTypeName, measTypeID := range measTypes {
		log.Debug("Measurement Name and ID:", measTypeName, measTypeID)
		measInfoActionItem := pdubuilder.CreateMeasurementInfoActionItem(measTypeName, measTypeID)
		measInfoActionList.Value = append(measInfoActionList.Value, measInfoActionItem)

	}

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

func (sm *Client) createMeasurementInfoList() (*e2smkpmv2.MeasurementInfoList, error) {
	// Creates measurement info list
	measInfoList := e2smkpmv2.MeasurementInfoList{
		Value: make([]*e2smkpmv2.MeasurementInfoItem, 0),
	}
	labelInfoList, err := sm.createInfoLabelList()
	if err != nil {
		return nil, err
	}
	for measTypeName := range measTypes {
		measTypeName, _ := measurments.NewMeasurementTypeMeasName(
			measurments.WithMeasurementName(measTypeName)).
			Build()
		measInfoItem, _ := measurments.NewMeasurementInfoItem(
			measurments.WithMeasType(measTypeName),
			measurments.WithLabelInfoList(labelInfoList)).Build()

		measInfoList.Value = append(measInfoList.Value, measInfoItem)
	}

	return &measInfoList, nil

}

func (sm *Client) createMeasurementData() (*e2smkpmv2.MeasurementData, error) {

	measData := e2smkpmv2.MeasurementData{
		Value: make([]*e2smkpmv2.MeasurementDataItem, 0),
	}

	for measTypeName := range measTypes {
		log.Debug("Creating measurement data for:", measTypeName)

		var integerValue = rand.Int63()
		// Creates meas record
		measRecord := e2smkpmv2.MeasurementRecord{
			Value: make([]*e2smkpmv2.MeasurementRecordItem, 0),
		}

		measRecordInteger := measurments.NewMeasurementRecordItemInteger(
			measurments.WithIntegerValue(integerValue)).
			Build()

		measRecord.Value = append(measRecord.Value, measRecordInteger)

		measDataItem, err := measurments.NewMeasurementDataItem(
			measurments.WithMeasurementRecord(&measRecord),
			measurments.WithIncompleteFlag(e2smkpmv2.IncompleteFlag_INCOMPLETE_FLAG_TRUE)).
			Build()
		if err != nil {
			log.Warn(err)
			return nil, err
		}

		measData.Value = append(measData.Value, measDataItem)
	}
	return &measData, nil

}

func (sm *Client) createIndicationMessageFormat1(cellECGI ransimtypes.ECGI) ([]byte, error) {

	measInfoList, err := sm.createMeasurementInfoList()
	if err != nil {
		return nil, err
	}

	measData, err := sm.createMeasurementData()
	if err != nil {
		return nil, err
	}

	var granularity int32 = 21
	// Creating an indication message format 1
	indicationMessage := kpm2MessageFormat1.NewIndicationMessage(
		kpm2MessageFormat1.WithCellObjID(strconv.FormatUint(uint64(cellECGI), 10)),
		kpm2MessageFormat1.WithGranularity(granularity),
		kpm2MessageFormat1.WithSubscriptionID(123456),
		kpm2MessageFormat1.WithMeasData(measData),
		kpm2MessageFormat1.WithMeasInfoList(measInfoList))

	kpmModelPlugin, err := sm.ServiceModel.ModelPluginRegistry.GetPlugin(e2smtypes.OID(sm.ServiceModel.OID))
	if err != nil {
		return nil, err
	}
	indicationMessageBytes, err := indicationMessage.ToAsn1Bytes(kpmModelPlugin)
	if err != nil {
		log.Warn(err)
		return nil, err
	}

	return indicationMessageBytes, nil
}

func (sm *Client) createIndicationHeaderBytes() ([]byte, error) {
	gNbID, err := strconv.ParseUint(fmt.Sprintf("%d", sm.ServiceModel.Node.EnbID), 10, 64)
	if err != nil {
		log.Warn(err)
		return nil, err
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
		return nil, err
	}

	header := kpm2IndicationHeader.NewIndicationHeader(
		kpm2IndicationHeader.WithGlobalKpmNodeID(kpmNodeID),
		kpm2IndicationHeader.WithFileFormatVersion(fileFormatVersion),
		kpm2IndicationHeader.WithSenderName(senderName),
		kpm2IndicationHeader.WithSenderType(senderType),
		kpm2IndicationHeader.WithVendorName(vendorName),
		kpm2IndicationHeader.WithTimeStamp([]byte{0x21, 0x22, 0x23, 0x24}))

	kpmModelPlugin, err := sm.ServiceModel.ModelPluginRegistry.GetPlugin(e2smtypes.OID(sm.ServiceModel.OID))
	if err != nil {
		return nil, err
	}
	indicationHeaderAsn1Bytes, err := header.ToAsn1Bytes(kpmModelPlugin)
	if err != nil {
		log.Warn(err)
		return nil, err
	}

	return indicationHeaderAsn1Bytes, nil

}

func (sm *Client) createRicIndication(ctx context.Context, ecgi ransimtypes.ECGI, subscription *subutils.Subscription) (*e2appducontents.Ricindication, error) {
	indicationMessageBytes, err := sm.createIndicationMessageFormat1(ecgi)
	if err != nil {
		log.Warn(err)
		return nil, err
	}
	indicationHeaderAsn1Bytes, err := sm.createIndicationHeaderBytes()
	if err != nil {
		log.Warn(err)
		return nil, err
	}

	indication := e2apIndicationUtils.NewIndication(
		e2apIndicationUtils.WithRicInstanceID(subscription.GetRicInstanceID()),
		e2apIndicationUtils.WithRanFuncID(subscription.GetRanFuncID()),
		e2apIndicationUtils.WithRequestID(subscription.GetReqID()),
		e2apIndicationUtils.WithIndicationHeader(indicationHeaderAsn1Bytes),
		e2apIndicationUtils.WithIndicationMessage(indicationMessageBytes))

	ricIndication, err := indication.Build()
	if err != nil {
		log.Error("creating indication message is failed for Cell with ID", ecgi, err)
		return nil, err
	}
	return ricIndication, nil
}

func (sm *Client) sendRicIndication(ctx context.Context, subscription *subutils.Subscription) error {
	subID := subscriptions.NewID(subscription.GetRicInstanceID(), subscription.GetReqID(), subscription.GetRanFuncID())
	sub, err := sm.ServiceModel.Subscriptions.Get(subID)
	if err != nil {
		return err
	}

	node := sm.ServiceModel.Node
	// Creates and sends an indication message for each cell in the node
	for _, ecgi := range node.Cells {
		ricIndication, err := sm.createRicIndication(ctx, ecgi, subscription)
		if err != nil {
			log.Error(err)
			return err
		}
		err = sub.E2Channel.RICIndication(ctx, ricIndication)
		if err != nil {
			log.Error(err)
			return err
		}
	}
	return nil
}

func (sm *Client) reportIndication(ctx context.Context, interval int32, subscription *subutils.Subscription) error {
	subID := subscriptions.NewID(subscription.GetRicInstanceID(), subscription.GetReqID(), subscription.GetRanFuncID())
	// Creates an indication header

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
			err = sm.sendRicIndication(ctx, subscription)
			if err != nil {
				log.Error("creating indication message is failed", err)
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
