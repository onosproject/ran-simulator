// SPDX-FileCopyrightText: 2020-present Open Networking Foundation <info@opennetworking.org>
//
// SPDX-License-Identifier: LicenseRef-ONF-Member-1.0

package utils

import (
	"context"

	subapi "github.com/onosproject/onos-api/go/onos/e2sub/subscription"
	"github.com/onosproject/onos-e2-sm/servicemodels/e2sm_kpm/pdubuilder"
	rcpdubuilder "github.com/onosproject/onos-e2-sm/servicemodels/e2sm_rc_pre/pdubuilder"
	"github.com/onosproject/onos-ric-sdk-go/pkg/e2/creds"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/protobuf/proto"
)

// Subscription subscription request for subscription SDK api
type Subscription struct {
	NodeID               string
	ServiceModelName     subapi.ServiceModelName
	ServiceModelVersion  subapi.ServiceModelVersion
	ActionType           subapi.ActionType
	ActionID             int32
	EncodingType         subapi.Encoding
	TimeToWait           subapi.TimeToWait
	SubSequentActionType subapi.SubsequentActionType
	EventTrigger         []byte
}

// CreateRcEventTrigger creates a rc service model event trigger
func CreateRcEventTrigger() ([]byte, error) {
	e2SmKpmEventTriggerDefinition, err := rcpdubuilder.CreateE2SmRcPreEventTriggerDefinitionUponChange()
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

// CreateKpmEventTrigger creates a kpm service model event trigger
func CreateKpmEventTrigger(rtPeriod int32) ([]byte, error) {
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

// Create creates a subscription request using SDK
func (subRequest *Subscription) Create() (subapi.SubscriptionDetails, error) {
	subReq := subapi.SubscriptionDetails{
		E2NodeID: subapi.E2NodeID(subRequest.NodeID),
		ServiceModel: subapi.ServiceModel{
			Name:    subRequest.ServiceModelName,
			Version: subRequest.ServiceModelVersion,
		},
		EventTrigger: subapi.EventTrigger{
			Payload: subapi.Payload{
				Encoding: subRequest.EncodingType,
				Data:     subRequest.EventTrigger,
			},
		},
		Actions: []subapi.Action{
			{
				ID:   subRequest.ActionID,
				Type: subRequest.ActionType,
				SubsequentAction: &subapi.SubsequentAction{
					Type:       subRequest.SubSequentActionType,
					TimeToWait: subRequest.TimeToWait,
				},
			},
		},
	}

	return subReq, nil
}

// ConnectSubscriptionServiceHost connects to subscription service
func ConnectSubscriptionServiceHost() (*grpc.ClientConn, error) {
	tlsConfig, err := creds.GetClientCredentials()
	if err != nil {
		return nil, err
	}
	opts := []grpc.DialOption{
		grpc.WithTransportCredentials(credentials.NewTLS(tlsConfig)),
	}

	return grpc.DialContext(context.Background(), SubscriptionServiceAddress, opts...)
}
