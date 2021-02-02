// SPDX-FileCopyrightText: 2020-present Open Networking Foundation <info@opennetworking.org>
//
// SPDX-License-Identifier: LicenseRef-ONF-Member-1.0

package rc

import (
	"context"

	"github.com/onosproject/ran-simulator/pkg/model"
	"github.com/onosproject/ran-simulator/pkg/modelplugins"

	"github.com/onosproject/onos-e2t/api/e2ap/v1beta1/e2appducontents"
	"github.com/onosproject/onos-lib-go/pkg/logging"
	"github.com/onosproject/ran-simulator/pkg/servicemodel"
	"github.com/onosproject/ran-simulator/pkg/servicemodel/registry"
	"github.com/onosproject/ran-simulator/pkg/store/subscriptions"
)

var _ servicemodel.Client = &Client{}

var log = logging.GetLogger("sm", "rc")

const (
	modelFullName = "e2sm_rc-pre"
	version       = "v1beta1"
)

// Client kpm service model client
type Client struct {
	Subscriptions *subscriptions.Subscriptions
	ServiceModel  *registry.ServiceModel
}

// NewServiceModel creates a new service model
func NewServiceModel(node model.Node, model *model.Model, modelPluginRegistry *modelplugins.ModelPluginRegistry) (registry.ServiceModel, error) {
	modelFullName := modelplugins.ModelFullName(modelFullName)
	RcSm := registry.ServiceModel{
		RanFunctionID:       registry.Rc,
		ModelFullName:       modelFullName,
		Client:              &Client{},
		Revision:            1,
		Version:             version,
		ModelPluginRegistry: modelPluginRegistry,
		Node:                node,
		Model:               model,
	}
	// TODO Creates ran function description and update the value in RcSm
	return RcSm, nil
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
