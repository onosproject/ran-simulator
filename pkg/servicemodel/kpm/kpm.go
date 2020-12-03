// SPDX-FileCopyrightText: 2020-present Open Networking Foundation <info@opennetworking.org>
//
// SPDX-License-Identifier: LicenseRef-ONF-Member-1.0

package kpm

import (
	"context"

	"github.com/onosproject/onos-e2t/api/e2ap/v1beta1/e2appducontents"
	"github.com/onosproject/ran-simulator/pkg/servicemodel"
)

var _ servicemodel.ServiceModel = &ServiceModel{}

// ServiceModel kpm service model struct
type ServiceModel struct {
}

// RICControl ....
func (sm *ServiceModel) RICControl(ctx context.Context, request *e2appducontents.RiccontrolRequest) (response *e2appducontents.RiccontrolAcknowledge, failure *e2appducontents.RiccontrolFailure, err error) {
	panic("implement me")
}

// RICSubscription ...
func (sm *ServiceModel) RICSubscription(ctx context.Context, request *e2appducontents.RicsubscriptionRequest) (response *e2appducontents.RicsubscriptionResponse, failure *e2appducontents.RicsubscriptionFailure, err error) {
	panic("implement me")
}

// RICSubscriptionDelete ...
func (sm *ServiceModel) RICSubscriptionDelete(ctx context.Context, request *e2appducontents.RicsubscriptionDeleteRequest) (response *e2appducontents.RicsubscriptionDeleteResponse, failure *e2appducontents.RicsubscriptionDeleteFailure, err error) {
	panic("implement me")
}
