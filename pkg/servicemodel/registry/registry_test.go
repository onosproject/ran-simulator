// SPDX-FileCopyrightText: 2020-present Open Networking Foundation <info@opennetworking.org>
//
// SPDX-License-Identifier: LicenseRef-ONF-Member-1.0

package registry

import (
	"context"
	"testing"

	"github.com/magiconair/properties/assert"

	"github.com/onosproject/ran-simulator/pkg/servicemodel"

	"github.com/onosproject/onos-e2t/api/e2ap/v1beta1/e2appducontents"
)

var _ servicemodel.ServiceModel = &mockServiceModel{}

type mockServiceModel struct {
	t *testing.T
}

func (sm mockServiceModel) RICControl(ctx context.Context, request *e2appducontents.RiccontrolRequest) (response *e2appducontents.RiccontrolAcknowledge, failure *e2appducontents.RiccontrolFailure, err error) {
	panic("implement me")
}

func (sm mockServiceModel) RICSubscription(ctx context.Context, request *e2appducontents.RicsubscriptionRequest) (response *e2appducontents.RicsubscriptionResponse, failure *e2appducontents.RicsubscriptionFailure, err error) {
	sm.t.Log("Test Ric Subscription")
	return nil, nil, nil
}

func (sm mockServiceModel) RICSubscriptionDelete(ctx context.Context, request *e2appducontents.RicsubscriptionDeleteRequest) (response *e2appducontents.RicsubscriptionDeleteResponse, failure *e2appducontents.RicsubscriptionDeleteFailure, err error) {
	panic("implement me")
}

func TestRegisterServiceModel(t *testing.T) {

	registry := NewServiceModelRegistry()

	m := &mockServiceModel{
		t: t,
	}

	testServiceModelConfig := ServiceModelConfig{
		ID:           Internal,
		ServiceModel: m,
		Description:  "Test Service model",
		Revision:     1,
	}

	if err := registry.RegisterServiceModel(testServiceModelConfig); err != nil {
		t.Fatalf("failed to register the service model")
	}

	if err := registry.GetServiceModel(Internal, m); err != nil {
		t.Fatal("the service model does not exist", err)
	}

	ranFunctions := registry.GetRanFunctions()
	assert.Equal(t, len(ranFunctions), 1)

}
