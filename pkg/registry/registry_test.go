// SPDX-FileCopyrightText: 2020-present Open Networking Foundation <info@opennetworking.org>
//
// SPDX-License-Identifier: LicenseRef-ONF-Member-1.0

package registry

import (
	"testing"

	"github.com/onosproject/onos-e2t/api/e2ap/v1beta1/e2appducontents"

	"github.com/onosproject/ran-simulator/pkg/servicemodels"
)

type mockServiceModel struct {
}

func (sm mockServiceModel) ProcessSubscriptionDelete(request *e2appducontents.RicsubscriptionDeleteRequest) (response *e2appducontents.RicsubscriptionDeleteResponse, failure *e2appducontents.RicsubscriptionDeleteFailure, err error) {
	return nil, nil, nil
}

func (sm mockServiceModel) ProcessSubscription(request *e2appducontents.RicsubscriptionRequest) (response *e2appducontents.RicsubscriptionResponse, failure *e2appducontents.RicsubscriptionFailure, err error) {
	return nil, nil, nil
}

func (sm mockServiceModel) ProcessControl(request *e2appducontents.RiccontrolRequest) (response *e2appducontents.RiccontrolAcknowledge, failure *e2appducontents.RiccontrolFailure, err error) {
	return nil, nil, nil
}

func TestRegisterServiceModel(t *testing.T) {
	registry := &ServiceModelRegistry{
		serviceModels: make(map[ID]servicemodels.ServiceModel),
	}

	m := &mockServiceModel{}

	if err := registry.RegisterServiceModel(0, m); err != nil {
		t.Fatalf("failed to register the service model")
	}

	if err := registry.GetServiceModel(0, m); err != nil {
		t.Fatal("the service model does not exist", err)
	}

}
