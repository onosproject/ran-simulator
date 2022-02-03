// SPDX-FileCopyrightText: 2020-present Open Networking Foundation <info@opennetworking.org>
//
// SPDX-License-Identifier: Apache-2.0

package registry

import (
	"context"
	"testing"

	"github.com/onosproject/ran-simulator/pkg/servicemodel"
	"github.com/stretchr/testify/assert"

	e2appducontents "github.com/onosproject/onos-e2t/api/e2ap/v2/e2ap-pdu-contents"
)

var _ servicemodel.Client = &mockServiceModel{}

type mockServiceModel struct {
	t *testing.T
}

func (sm mockServiceModel) E2ConnectionUpdate(ctx context.Context, request *e2appducontents.E2ConnectionUpdate) (response *e2appducontents.E2ConnectionUpdateAcknowledge, failure *e2appducontents.E2ConnectionUpdateFailure, err error) {
	panic("implement me")
}

func (sm mockServiceModel) RICControl(ctx context.Context, request *e2appducontents.RiccontrolRequest) (response *e2appducontents.RiccontrolAcknowledge, failure *e2appducontents.RiccontrolFailure, err error) {
	panic("implement me")
}

func (sm mockServiceModel) RICSubscription(ctx context.Context, request *e2appducontents.RicsubscriptionRequest) (response *e2appducontents.RicsubscriptionResponse, failure *e2appducontents.RicsubscriptionFailure, err error) {
	sm.t.Log("Test Ric Subscription")
	return &e2appducontents.RicsubscriptionResponse{}, &e2appducontents.RicsubscriptionFailure{}, nil
}

func (sm mockServiceModel) RICSubscriptionDelete(ctx context.Context, request *e2appducontents.RicsubscriptionDeleteRequest) (response *e2appducontents.RicsubscriptionDeleteResponse, failure *e2appducontents.RicsubscriptionDeleteFailure, err error) {
	panic("implement me")
}

func TestRegisterServiceModel(t *testing.T) {

	registry := NewServiceModelRegistry()

	m := &mockServiceModel{
		t: t,
	}

	testServiceModelConfig := ServiceModel{
		RanFunctionID: Internal,
		Client:        m,
		Description:   []byte{0x01, 0x02, 0x03, 0x04},
		Revision:      1,
	}

	if err := registry.RegisterServiceModel(testServiceModelConfig); err != nil {
		t.Fatalf("failed to register the service model")
	}

	sm, err := registry.GetServiceModel(Internal)
	if err != nil {
		t.Fatal("the service model does not exist", err)
	}

	testSm := sm

	_, _, err = testSm.Client.RICSubscription(context.Background(), &e2appducontents.RicsubscriptionRequest{})
	assert.NoError(t, err)

	ranFunctions := registry.GetRanFunctions()
	assert.Equal(t, len(ranFunctions), 1)

}
