// SPDX-FileCopyrightText: 2020-present Open Networking Foundation <info@opennetworking.org>
//
// SPDX-License-Identifier: LicenseRef-ONF-Member-1.0

package kpm

import (
	"context"
	"strconv"
	"time"

	kpmutils "github.com/onosproject/ran-simulator/pkg/utils/indication/kpm"

	"github.com/onosproject/ran-simulator/pkg/model"

	"github.com/onosproject/ran-simulator/pkg/modelplugins"

	"github.com/onosproject/onos-e2-sm/servicemodels/e2sm_kpm/pdubuilder"
	"github.com/onosproject/onos-e2t/pkg/protocols/e2"
	indicationutils "github.com/onosproject/ran-simulator/pkg/utils/indication"
	subutils "github.com/onosproject/ran-simulator/pkg/utils/subscription"
	subdeleteutils "github.com/onosproject/ran-simulator/pkg/utils/subscriptiondelete"

	"github.com/onosproject/onos-lib-go/pkg/logging"

	"github.com/onosproject/onos-lib-go/pkg/errors"
	"github.com/onosproject/ran-simulator/pkg/servicemodel/registry"

	"github.com/onosproject/onos-e2t/api/e2ap/v1beta1/e2apies"
	"github.com/onosproject/onos-e2t/api/e2ap/v1beta1/e2appducontents"
	"github.com/onosproject/onos-e2t/pkg/southbound/e2ap/types"
	"github.com/onosproject/ran-simulator/pkg/servicemodel"
	"google.golang.org/protobuf/proto"
)

var _ servicemodel.Client = &Client{}

var log = logging.GetLogger("sm", "kpm")

const (
	modelFullName = "e2sm_kpm-v1beta1"
	version       = "v1beta1"
)

// Client kpm service model client
type Client struct {
	Channel      e2.ClientChannel
	ServiceModel *registry.ServiceModel
}

// NewServiceModel creates a new service model
func NewServiceModel(node model.Node, model *model.Model, modelPluginRegistry *modelplugins.ModelPluginRegistry) (registry.ServiceModel, error) {
	modelFullName := modelplugins.ModelFullName(modelFullName)
	kpmSm := registry.ServiceModel{
		RanFunctionID:       registry.Kpm,
		ModelFullName:       modelFullName,
		Client:              &Client{},
		Revision:            1,
		Version:             version,
		ModelPluginRegistry: modelPluginRegistry,
		Node:                node,
		Model:               model,
	}
	var ranFunctionShortName = "“ORAN-E2SM-KPM”"
	var ranFunctionE2SmOid = "OID123"
	var ranFunctionDescription = "KPM Monitor"
	var ranFunctionInstance int32 = 1
	var ricEventStyleType int32 = 1
	var ricEventStyleName = "Periodic report"
	var ricEventFormatType int32 = 5
	var ricReportStyleType int32 = 1
	var ricReportStyleName = "O-CU-CP Measurement Container for the 5GC connected deployment"
	var ricIndicationHeaderFormatType int32 = 1
	var ricIndicationMessageFormatType int32 = 1
	ranFuncDescPdu, err := pdubuilder.CreateE2SmKpmRanfunctionDescriptionMsg(ranFunctionShortName, ranFunctionE2SmOid, ranFunctionDescription,
		ranFunctionInstance, ricEventStyleType, ricEventStyleName, ricEventFormatType, ricReportStyleType, ricReportStyleName,
		ricIndicationHeaderFormatType, ricIndicationMessageFormatType)
	if err != nil {
		log.Error(err)
		return registry.ServiceModel{}, err
	}

	protoBytes, err := proto.Marshal(ranFuncDescPdu)
	log.Debug("Proto bytes of KPM Ran Function Description:", protoBytes)
	if err != nil {
		log.Error(err)
		return registry.ServiceModel{}, err
	}
	kpmModelPlugin := modelPluginRegistry.ModelPlugins[modelFullName]
	if kpmModelPlugin == nil {
		return registry.ServiceModel{}, errors.New(errors.Invalid, "model plugin is nil")
	}
	// TODO it panics and it should be fixed in kpm service model otherwise it panics
	/*ranFuncDescBytes, err = kpmModelPlugin.RanFuncDescriptionProtoToASN1(protoBytes)
	if err != nil {
		log.Error(err)
		return registry.ServiceModel{}, err
	}*/

	kpmSm.Description = ranFuncDescBytes
	return kpmSm, nil
}

func (sm *Client) reportIndication(ctx context.Context, interval int32, subscription *subutils.Subscription) error {
	gNbID, err := strconv.ParseUint(string(sm.ServiceModel.Node.EnbID), 10, 64)
	if err != nil {
		log.Error(err)
		return err
	}
	newIndicationHeader, _ := kpmutils.NewIndicationHeader(
		kpmutils.WithPlmnID(string(sm.ServiceModel.Model.PlmnID)),
		kpmutils.WithGnbID(gNbID),
		kpmutils.WithSst("1"),
		kpmutils.WithSd("SD1"),
		kpmutils.WithPlmnIDnrcgi(string(sm.ServiceModel.Model.PlmnID)))

	// Creating an indication header
	indicationHeader, err := kpmutils.CreateIndicationHeader(newIndicationHeader)
	if err != nil {
		log.Error(err)
		return err
	}

	indicationHeaderProtoBytes, err := proto.Marshal(indicationHeader)
	if err != nil {
		log.Error("Error in creating indication header proto bytes")
		return err
	}
	kpmModelPlugin := sm.ServiceModel.ModelPluginRegistry.ModelPlugins[sm.ServiceModel.ModelFullName]
	indicationHeaderAsn1Bytes, err := kpmModelPlugin.IndicationHeaderProtoToASN1(indicationHeaderProtoBytes)
	if err != nil {
		log.Error("Error in creating indication header ASN1 bytes", err)
		return err
	}
	// Creating an indication message
	newIndicationMessage, err := kpmutils.NewIndicationMessage(
		kpmutils.WithNumberOfActiveUes(10))
	if err != nil {
		log.Error(err)
		return err
	}
	indicationMessage, err := kpmutils.CreateIndicationMessage(newIndicationMessage)
	if err != nil {
		log.Error(err)
		return err
	}
	_, err = proto.Marshal(indicationMessage)
	if err != nil {
		log.Error("Error in creating indication header proto bytes")
		return err
	}

	// TODO model plugin bug should be fixed to call this function otherwise it panics
	/*indicationMessageAsn1Bytes, err := kpmModelPlugin.IndicationMessageProtoToASN1(indicationMessageProtoBytes)
	if err != nil {
		log.Error(err)
		return err
	}*/

	indication, _ := indicationutils.NewIndication(
		indicationutils.WithRicInstanceID(subscription.GetRicInstanceID()),
		indicationutils.WithRanFuncID(subscription.GetRanFuncID()),
		indicationutils.WithRequestID(subscription.GetReqID()),
		indicationutils.WithIndicationHeader(indicationHeaderAsn1Bytes),
		indicationutils.WithIndicationMessage(indicationMessageBytes))

	ricIndication := indicationutils.CreateIndication(indication)

	intervalDuration := time.Duration(interval)
	ticker := time.NewTicker(intervalDuration * time.Millisecond)
	for range ticker.C {
		log.Info("Sending indication")
		err := sm.Channel.RICIndication(ctx, ricIndication)
		if err != nil {
			log.Error("Sending indication report is failed:", err)
			return err
		}
	}
	return nil
}

// RICControl implements control handler for kpm service model
func (sm *Client) RICControl(ctx context.Context, request *e2appducontents.RiccontrolRequest) (response *e2appducontents.RiccontrolAcknowledge, failure *e2appducontents.RiccontrolFailure, err error) {
	return nil, nil, errors.New(errors.NotSupported, "Control operation is not supported")
}

// RICSubscription implements subscription handler for kpm service model
func (sm *Client) RICSubscription(ctx context.Context, request *e2appducontents.RicsubscriptionRequest) (response *e2appducontents.RicsubscriptionResponse, failure *e2appducontents.RicsubscriptionFailure, err error) {
	log.Info("RIC Subscription is called for service model:", sm.ServiceModel.ModelFullName)
	var ricActionsAccepted []*types.RicActionID
	ricActionsNotAdmitted := make(map[types.RicActionID]*e2apies.Cause)
	actionList := subutils.GetRicActionToBeSetupList(request)
	reqID := subutils.GetRequesterID(request)
	ranFuncID := subutils.GetRanFunctionID(request)
	ricInstanceID := subutils.GetRicInstanceID(request)

	for _, action := range actionList {
		actionID := types.RicActionID(action.Value.RicActionId.Value)
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
	subscription, _ := subutils.NewSubscription(
		subutils.WithRequestID(reqID),
		subutils.WithRanFuncID(ranFuncID),
		subutils.WithRicInstanceID(ricInstanceID),
		subutils.WithActionsAccepted(ricActionsAccepted),
		subutils.WithActionsNotAdmitted(ricActionsNotAdmitted))

	// At least one required action must be accepted otherwise sends a subscription failure response
	if len(ricActionsAccepted) == 0 {
		subscriptionFailure := subutils.CreateSubscriptionFailure(subscription)
		return nil, subscriptionFailure, errors.New(errors.Forbidden, "no required action is accepted")
	}

	reportInterval, err := sm.getReportPeriod(request)
	if err != nil {
		subscriptionFailure := subutils.CreateSubscriptionFailure(subscription)
		return nil, subscriptionFailure, err
	}

	subscriptionResponse := subutils.CreateSubscriptionResponse(subscription)
	go func() {
		err := sm.reportIndication(ctx, reportInterval, subscription)
		if err != nil {
			return
		}
	}()
	return subscriptionResponse, nil, nil

}

// RICSubscriptionDelete implements subscription delete handler for kpm service model
func (sm *Client) RICSubscriptionDelete(ctx context.Context, request *e2appducontents.RicsubscriptionDeleteRequest) (response *e2appducontents.RicsubscriptionDeleteResponse, failure *e2appducontents.RicsubscriptionDeleteFailure, err error) {
	reqID := subdeleteutils.GetRequesterID(request)
	ranFuncID := subdeleteutils.GetRanFunctionID(request)
	ricInstanceID := subdeleteutils.GetRicInstanceID(request)

	subscriptionDelete, _ := subdeleteutils.NewSubscriptionDelete(
		subdeleteutils.WithRequestID(reqID),
		subdeleteutils.WithRanFuncID(ranFuncID),
		subdeleteutils.WithRicInstanceID(ricInstanceID))
	subDeleteResponse := subdeleteutils.CreateSubscriptionDeleteResponse(subscriptionDelete)

	// TODO stop sending indication reports

	return subDeleteResponse, nil, errors.New(errors.NotSupported, "Ric subscription delete is not supported")
}

var indicationMessageBytes = []byte{0x40, 0x00, 0x00, 0x6c, 0x1a, 0x4f, 0x70, 0x65, 0x6e, 0x4e, 0x65, 0x74, 0x77, 0x6f, 0x72, 0x6b, 0x69, 0x6e, 0x67, 0x80, 0x00, 0x00, 0x0c, 0x72, 0x61, 0x6e, 0x43, 0x6f, 0x6e, 0x74, 0x61, 0x69, 0x6e, 0x65, 0x72}

var ranFuncDescBytes = []byte{
	0x20, 0xC0, 0x4F, 0x52, 0x41, 0x4E, 0x2D, 0x45, 0x32, 0x53, 0x4D, 0x2D, 0x4B, 0x50, 0x4D, 0x00, 0x00, 0x05, 0x4F, 0x49,
	0x44, 0x31, 0x32, 0x33, 0x05, 0x00, 0x4B, 0x50, 0x4D, 0x20, 0x6D, 0x6F, 0x6E, 0x69, 0x74, 0x6F, 0x72, 0x08, 0x93, 0x49,
	0xF4, 0x77, 0xF9, 0xE1, 0xAF, 0x00, 0x60, 0x00, 0x01, 0x01, 0x07, 0x00, 0x50, 0x65, 0x72, 0x69, 0x6F, 0x64, 0x69, 0x63,
	0x20, 0x72, 0x65, 0x70, 0x6F, 0x72, 0x74, 0x01, 0x05, 0x14, 0x01, 0x01, 0x1D, 0x00, 0x4F, 0x2D, 0x44, 0x55, 0x20, 0x4D,
	0x65, 0x61, 0x73, 0x75, 0x72, 0x65, 0x6D, 0x65, 0x6E, 0x74, 0x20, 0x43, 0x6F, 0x6E, 0x74, 0x61, 0x69, 0x6E, 0x65, 0x72,
	0x20, 0x66, 0x6F, 0x72, 0x20, 0x74, 0x68, 0x65, 0x20, 0x35, 0x47, 0x43, 0x20, 0x63, 0x6F, 0x6E, 0x6E, 0x65, 0x63, 0x74,
	0x65, 0x64, 0x20, 0x64, 0x65, 0x70, 0x6C, 0x6F, 0x79, 0x6D, 0x65, 0x6E, 0x74, 0x01, 0x01, 0x01, 0x01, 0x00, 0x01, 0x02,
	0x1D, 0x00, 0x4F, 0x2D, 0x44, 0x55, 0x20, 0x4D, 0x65, 0x61, 0x73, 0x75, 0x72, 0x65, 0x6D, 0x65, 0x6E, 0x74, 0x20, 0x43,
	0x6F, 0x6E, 0x74, 0x61, 0x69, 0x6E, 0x65, 0x72, 0x20, 0x66, 0x6F, 0x72, 0x20, 0x74, 0x68, 0x65, 0x20, 0x45, 0x50, 0x43,
	0x20, 0x63, 0x6F, 0x6E, 0x6E, 0x65, 0x63, 0x74, 0x65, 0x64, 0x20, 0x64, 0x65, 0x70, 0x6C, 0x6F, 0x79, 0x6D, 0x65, 0x6E,
	0x74, 0x01, 0x01, 0x01, 0x01, 0x00, 0x01, 0x03, 0x1E, 0x80, 0x4F, 0x2D, 0x43, 0x55, 0x2D, 0x43, 0x50, 0x20, 0x4D, 0x65,
	0x61, 0x73, 0x75, 0x72, 0x65, 0x6D, 0x65, 0x6E, 0x74, 0x20, 0x43, 0x6F, 0x6E, 0x74, 0x61, 0x69, 0x6E, 0x65, 0x72, 0x20,
	0x66, 0x6F, 0x72, 0x20, 0x74, 0x68, 0x65, 0x20, 0x35, 0x47, 0x43, 0x20, 0x63, 0x6F, 0x6E, 0x6E, 0x65, 0x63, 0x74, 0x65,
	0x64, 0x20, 0x64, 0x65, 0x70, 0x6C, 0x6F, 0x79, 0x6D, 0x65, 0x6E, 0x74, 0x01, 0x01, 0x01, 0x01, 0x00, 0x01, 0x04, 0x1E,
	0x80, 0x4F, 0x2D, 0x43, 0x55, 0x2D, 0x43, 0x50, 0x20, 0x4D, 0x65, 0x61, 0x73, 0x75, 0x72, 0x65, 0x6D, 0x65, 0x6E, 0x74,
	0x20, 0x43, 0x6F, 0x6E, 0x74, 0x61, 0x69, 0x6E, 0x65, 0x72, 0x20, 0x66, 0x6F, 0x72, 0x20, 0x74, 0x68, 0x65, 0x20, 0x45,
	0x50, 0x43, 0x20, 0x63, 0x6F, 0x6E, 0x6E, 0x65, 0x63, 0x74, 0x65, 0x64, 0x20, 0x64, 0x65, 0x70, 0x6C, 0x6F, 0x79, 0x6D,
	0x65, 0x6E, 0x74, 0x01, 0x01, 0x01, 0x01, 0x00, 0x01, 0x05, 0x1E, 0x80, 0x4F, 0x2D, 0x43, 0x55, 0x2D, 0x55, 0x50, 0x20,
	0x4D, 0x65, 0x61, 0x73, 0x75, 0x72, 0x65, 0x6D, 0x65, 0x6E, 0x74, 0x20, 0x43, 0x6F, 0x6E, 0x74, 0x61, 0x69, 0x6E, 0x65,
	0x72, 0x20, 0x66, 0x6F, 0x72, 0x20, 0x74, 0x68, 0x65, 0x20, 0x35, 0x47, 0x43, 0x20, 0x63, 0x6F, 0x6E, 0x6E, 0x65, 0x63,
	0x74, 0x65, 0x64, 0x20, 0x64, 0x65, 0x70, 0x6C, 0x6F, 0x79, 0x6D, 0x65, 0x6E, 0x74, 0x01, 0x01, 0x01, 0x01, 0x00, 0x01,
	0x06, 0x1E, 0x80, 0x4F, 0x2D, 0x43, 0x55, 0x2D, 0x55, 0x50, 0x20, 0x4D, 0x65, 0x61, 0x73, 0x75, 0x72, 0x65, 0x6D, 0x65,
	0x6E, 0x74, 0x20, 0x43, 0x6F, 0x6E, 0x74, 0x61, 0x69, 0x6E, 0x65, 0x72, 0x20, 0x66, 0x6F, 0x72, 0x20, 0x74, 0x68, 0x65,
	0x20, 0x45, 0x50, 0x43, 0x20, 0x63, 0x6F, 0x6E, 0x6E, 0x65, 0x63, 0x74, 0x65, 0x64, 0x20, 0x64, 0x65, 0x70, 0x6C, 0x6F,
	0x79, 0x6D, 0x65, 0x6E, 0x74, 0x01, 0x01, 0x01, 0x01}
