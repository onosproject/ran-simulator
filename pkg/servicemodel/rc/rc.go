// SPDX-FileCopyrightText: 2020-present Open Networking Foundation <info@opennetworking.org>
//
// SPDX-License-Identifier: LicenseRef-ONF-Member-1.0

package rc

import (
	"context"

	"github.com/onosproject/onos-e2-sm/servicemodels/e2sm_rc_pre/pdubuilder"
	"github.com/onosproject/ran-simulator/pkg/model"
	"github.com/onosproject/ran-simulator/pkg/modelplugins"
	"google.golang.org/protobuf/proto"

	"github.com/onosproject/onos-e2t/api/e2ap/v1beta1/e2appducontents"
	"github.com/onosproject/onos-lib-go/pkg/errors"
	"github.com/onosproject/onos-lib-go/pkg/logging"
	"github.com/onosproject/ran-simulator/pkg/servicemodel"
	"github.com/onosproject/ran-simulator/pkg/servicemodel/registry"
	"github.com/onosproject/ran-simulator/pkg/store/subscriptions"
)

var _ servicemodel.Client = &Client{}

var log = logging.GetLogger("sm", "rc")

const (
	modelFullName = "e2sm_rc_pre-v1"
	version       = "v1"
)

// Client kpm service model client
type Client struct {
	Subscriptions *subscriptions.Subscriptions
	ServiceModel  *registry.ServiceModel
}

// NewServiceModel creates a new service model
func NewServiceModel(node model.Node, model *model.Model, modelPluginRegistry *modelplugins.ModelPluginRegistry) (registry.ServiceModel, error) {
	modelFullName := modelplugins.ModelFullName(modelFullName)
	rcSm := registry.ServiceModel{
		RanFunctionID:       registry.Rc,
		ModelFullName:       modelFullName,
		Client:              &Client{},
		Revision:            1,
		Version:             version,
		ModelPluginRegistry: modelPluginRegistry,
		Node:                node,
		Model:               model,
	}

	var ranFunctionShortName = "ORAN-E2SM-RC"
	var ranFunctionE2SmOid = "OID124"
	var ranFunctionDescription = "RC PRE"
	var ranFunctionInstance int32 = 3
	var ricEventStyleType int32 = 1
	var ricEventStyleName = "PeriodicReport"
	var ricEventFormatType int32 = 1
	var ricReportStyleType int32 = 1
	var ricReportStyleName = "PCI and NRT update for eNB"
	var ricIndicationHeaderFormatType int32 = 21
	var ricIndicationMessageFormatType int32 = 56
	ranFuncDescPdu, err := pdubuilder.CreateE2SmRcPreRanfunctionDescriptionMsg(ranFunctionShortName, ranFunctionE2SmOid, ranFunctionDescription,
		ranFunctionInstance, ricEventStyleType, ricEventStyleName, ricEventFormatType, ricReportStyleType, ricReportStyleName,
		ricIndicationHeaderFormatType, ricIndicationMessageFormatType)

	if err != nil {
		log.Error(err)
		return registry.ServiceModel{}, err
	}

	protoBytes, err := proto.Marshal(ranFuncDescPdu)
	log.Debug("Proto bytes of RC service model Ran Function Description:", protoBytes)
	if err != nil {
		log.Error(err)
		return registry.ServiceModel{}, err
	}
	rcModelPlugin := modelPluginRegistry.ModelPlugins[modelFullName]
	if rcModelPlugin == nil {
		log.Info("model plugin names:", modelPluginRegistry.ModelPlugins)
		return registry.ServiceModel{}, errors.New(errors.Invalid, "model plugin is nil")
	}
	// TODO it panics and it should be fixed in rc service model otherwise it panics
	/*ranFuncDescBytes, err := rcModelPlugin.RanFuncDescriptionProtoToASN1(protoBytes)
	if err != nil {
		log.Error(err)
		return registry.ServiceModel{}, err
	}*/

	rcSm.Description = ranFuncDescBytes

	return rcSm, nil
}

// RICControl implements control handler for RC service model
func (sm Client) RICControl(ctx context.Context, request *e2appducontents.RiccontrolRequest) (response *e2appducontents.RiccontrolAcknowledge, failure *e2appducontents.RiccontrolFailure, err error) {
	log.Info("Control Request is received for service model:", sm.ServiceModel.ModelFullName)
	// TODO implements handler for control requests
	return response, failure, err
}

// RICSubscription implements subscription handler for RC service model
func (sm Client) RICSubscription(ctx context.Context, request *e2appducontents.RicsubscriptionRequest) (response *e2appducontents.RicsubscriptionResponse, failure *e2appducontents.RicsubscriptionFailure, err error) {
	log.Info("Ric Subscription Request is received for service model:", sm.ServiceModel.ModelFullName)
	// TODO implements handler for subscription requests
	return response, failure, err
}

// RICSubscriptionDelete implements subscription delete handler for RC service model
func (sm Client) RICSubscriptionDelete(ctx context.Context, request *e2appducontents.RicsubscriptionDeleteRequest) (response *e2appducontents.RicsubscriptionDeleteResponse, failure *e2appducontents.RicsubscriptionDeleteFailure, err error) {
	log.Info("Ric subscription delete request is received for service model:", sm.ServiceModel.ModelFullName)
	// TODO implements handler for subscription delete requests
	return response, failure, err
}

var ranFuncDescBytes = []byte{
	0x20, 0xf0, 0x4f, 0x52, 0x41, 0x4e, 0x2d, 0x45, 0x32, 0x53, 0x4d, 0x2d, 0x52, 0x43, 0x2d, 0x50, 0x52,
	0x45, 0x00, 0x00, 0x05, 0x4f, 0x49, 0x44, 0x31, 0x32, 0x33, 0x02, 0x80, 0x52, 0x43, 0x20, 0x50, 0x52,
	0x45, 0x01, 0x03, 0x60, 0x00, 0x01, 0x01, 0x07, 0x00, 0x50, 0x65, 0x72, 0x69, 0x6f, 0x64, 0x69, 0x63,
	0x20, 0x72, 0x65, 0x70, 0x6f, 0x72, 0x74, 0x01, 0x05, 0x00, 0x01, 0x01, 0x0c, 0x80, 0x50, 0x43, 0x49,
	0x20, 0x61, 0x6e, 0x64, 0x20, 0x4e, 0x52, 0x54, 0x20, 0x75, 0x70, 0x64, 0x61, 0x74, 0x65, 0x20, 0x66,
	0x6f, 0x72, 0x20, 0x65, 0x4e, 0x42, 0x01, 0x01, 0x01, 0x01,
}
