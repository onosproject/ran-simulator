// SPDX-FileCopyrightText: 2020-present Open Networking Foundation <info@opennetworking.org>
//
// SPDX-License-Identifier: LicenseRef-ONF-Member-1.0

package e2t

import (
	"context"
	"testing"
	"time"

	"github.com/onosproject/onos-api/go/onos/e2sub/subscription"
	"google.golang.org/protobuf/proto"

	"github.com/onosproject/onos-ric-sdk-go/pkg/e2/indication"

	"gotest.tools/assert"

	e2client "github.com/onosproject/onos-ric-sdk-go/pkg/e2"

	"github.com/onosproject/onos-e2-sm/servicemodels/e2sm_kpm/pdubuilder"
	"github.com/onosproject/onos-e2t/test/utils"
)

const (
	subscriptionServiceHost = "onos-e2sub"
	subscriptionServicePort = 5150
)

func createEventTriggerTestData() ([]byte, error) {

	var rtPeriod int32 = 12
	e2SmKpmEventTriggerDefinition, err := pdubuilder.CreateE2SmKpmEventTriggerDefinition(rtPeriod)
	if err != nil {
		return []byte{}, err
	}
	err = e2SmKpmEventTriggerDefinition.Validate()
	if err != nil {
		return []byte{}, err
	}
	protoBytes, err := proto.Marshal(e2SmKpmEventTriggerDefinition)
	if err != nil {
		return []byte{}, err
	}
	return protoBytes, nil
}

func createSubscriptionRequest(nodeID string) (subscription.SubscriptionDetails, error) {
	eventTriggerTestData, err := createEventTriggerTestData()
	if err != nil {
		return subscription.SubscriptionDetails{}, err
	}

	subReq := subscription.SubscriptionDetails{
		E2NodeID: subscription.E2NodeID(nodeID),
		ServiceModel: subscription.ServiceModel{
			Name:    KpmServiceModelName,
			Version: KpmServiceModelVersion1,
		},
		EventTrigger: subscription.EventTrigger{
			Payload: subscription.Payload{
				Encoding: subscription.Encoding_ENCODING_PROTO,
				Data:     eventTriggerTestData,
			},
		},
		Actions: []subscription.Action{
			{
				ID:   100,
				Type: subscription.ActionType_ACTION_TYPE_REPORT,
				SubsequentAction: &subscription.SubsequentAction{
					Type:       subscription.SubsequentActionType_SUBSEQUENT_ACTION_TYPE_CONTINUE,
					TimeToWait: subscription.TimeToWait_TIME_TO_WAIT_ZERO,
				},
			},
		},
	}

	return subReq, nil
}

// TestSubscription issues a subscription request to the simulator
func (s *TestSuite) TestSubscription(t *testing.T) {
	clientConfig := e2client.Config{
		AppID: "subscription-test",
		SubscriptionService: e2client.ServiceConfig{
			Host: subscriptionServiceHost,
			Port: subscriptionServicePort,
		},
	}
	client, err := e2client.NewClient(clientConfig)
	assert.NilError(t, err)

	ch := make(chan indication.Indication)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	nodeIDs, err := utils.GetNodeIDs()
	assert.NilError(t, err)

	subReq, err := createSubscriptionRequest(nodeIDs[0])
	assert.NilError(t, err)

	subCtx, err := client.Subscribe(ctx, subReq, ch)
	assert.NilError(t, err)

	select {
	case indicationMsg := <-ch:
		t.Log(indicationMsg)

	case <-time.After(20 * time.Second):
		t.Fatal("failed to receive an indication")
	}
	assert.NilError(t, err)

	// TODO: enable once subscription removal works
	if false {
		err = subCtx.Close()
		assert.NilError(t, err)

		select {
		case <-ch:
			t.Fatal("received an extraneous indication")

		case <-time.After(4 * time.Second):
			t.Log("all quiet on the western front")
		}

		assert.NilError(t, err)
	}
}
