// SPDX-FileCopyrightText: 2020-present Open Networking Foundation <info@opennetworking.org>
//
// SPDX-License-Identifier: LicenseRef-ONF-Member-1.0

package registry

import (
	"context"
	"testing"

	"github.com/onosproject/ran-simulator/pkg/servicemodel"

	"github.com/onosproject/onos-e2t/api/e2ap/v1beta1/e2appducontents"
)

var _ servicemodel.ServiceModel = &mockServiceModel{}

type mockServiceModel struct {
}

func (sm mockServiceModel) RICControl(ctx context.Context, request *e2appducontents.RiccontrolRequest) (response *e2appducontents.RiccontrolAcknowledge, failure *e2appducontents.RiccontrolFailure, err error) {
	panic("implement me")
}

func (sm mockServiceModel) RICSubscription(ctx context.Context, request *e2appducontents.RicsubscriptionRequest) (response *e2appducontents.RicsubscriptionResponse, failure *e2appducontents.RicsubscriptionFailure, err error) {
	panic("implement me")
}

func (sm mockServiceModel) RICSubscriptionDelete(ctx context.Context, request *e2appducontents.RicsubscriptionDeleteRequest) (response *e2appducontents.RicsubscriptionDeleteResponse, failure *e2appducontents.RicsubscriptionDeleteFailure, err error) {
	panic("implement me")
}

func TestRegisterServiceModel(t *testing.T) {
	registry := &ServiceModelRegistry{
		serviceModels: make(map[ID]servicemodel.ServiceModel),
	}

	m := &mockServiceModel{}

	if err := registry.RegisterServiceModel(INTERNAL, m); err != nil {
		t.Fatalf("failed to register the service model")
	}

	if err := registry.GetServiceModel(INTERNAL, m); err != nil {
		t.Fatal("the service model does not exist", err)
	}

}
