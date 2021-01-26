// SPDX-FileCopyrightText: 2020-present Open Networking Foundation <info@opennetworking.org>
//
// SPDX-License-Identifier: LicenseRef-ONF-Member-1.0

package e2agent

import (
	"fmt"
	"github.com/onosproject/onos-e2t/api/e2ap/v1beta1/e2apies"
	"github.com/onosproject/onos-e2t/api/e2ap/v1beta1/e2appducontents"
	"sync"
)

// ID is an alias for string subscription ID
type ID string

// Subscription is an auxiliary wrapper for tracking subscriptions by each E2 agent
type Subscription struct {
	ID      ID
	ReqID   *e2apies.RicrequestId
	FnID    *e2apies.RanfunctionId
	Details *e2appducontents.RicsubscriptionDetails
}

// GenID returns the locally unique ID for the specified subscription add/delete request
func GenID(instID int32, rqID int32, fnID int32) ID {
	return ID(fmt.Sprintf("%d/%d/%d", instID, rqID, fnID))
}

// NewSubscription generates a subscription record from the E2AP subscription request
func NewSubscription(e2apsub *e2appducontents.RicsubscriptionRequest) *Subscription {
	id := GenID(e2apsub.ProtocolIes.E2ApProtocolIes29.Value.RicInstanceId,
		e2apsub.ProtocolIes.E2ApProtocolIes29.Value.RicRequestorId,
		e2apsub.ProtocolIes.E2ApProtocolIes5.Value.Value)
	return &Subscription{
		ID:      id,
		ReqID:   e2apsub.ProtocolIes.E2ApProtocolIes29.Value,
		FnID:    e2apsub.ProtocolIes.E2ApProtocolIes5.Value,
		Details: e2apsub.ProtocolIes.E2ApProtocolIes30.Value,
	}
}

type subscriptions struct {
	subs map[ID]*Subscription
	mu   sync.RWMutex
}

// Add adds the specified subscription
func (s *subscriptions) Add(sub *Subscription) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.subs[sub.ID] = sub
}

// Remove removes the specified subscription
func (s *subscriptions) Remove(id ID) {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.subs, id)
}

// Get returns the subscription with the specified ID
func (s *subscriptions) Get(id ID) *Subscription {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.subs[id]
}

// List returns slice containing all current subscriptions
func (s *subscriptions) List() []*Subscription {
	s.mu.RLock()
	defer s.mu.RUnlock()

	resp := make([]*Subscription, 0, len(s.subs))
	for _, sub := range s.subs {
		resp = append(resp, sub)
	}
	return resp
}
