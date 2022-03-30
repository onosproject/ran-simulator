// SPDX-FileCopyrightText: 2020-present Open Networking Foundation <info@opennetworking.org>
//
// SPDX-License-Identifier: Apache-2.0

package subscriptions

import (
	"fmt"
	"sync"
	"time"

	v2 "github.com/onosproject/onos-e2t/api/e2ap/v2"

	"github.com/onosproject/onos-e2t/pkg/protocols/e2ap"

	"github.com/onosproject/onos-lib-go/pkg/errors"

	e2apies "github.com/onosproject/onos-e2t/api/e2ap/v2/e2ap-ies"
	e2appducontents "github.com/onosproject/onos-e2t/api/e2ap/v2/e2ap-pdu-contents"
)

// ID is an alias for string subscription ID
type ID string

// Subscription is an auxiliary wrapper for tracking subscriptions by each E2 agent
type Subscription struct {
	ID        ID
	ReqID     *e2apies.RicrequestId
	FnID      *e2apies.RanfunctionId
	Details   *e2appducontents.RicsubscriptionDetails
	E2Channel e2ap.ClientConn
	Ticker    *time.Ticker
}

// NewID returns the locally unique ID for the specified subscription add/delete request
func NewID(instID int32, rqID int32, fnID int32) ID {
	return ID(fmt.Sprintf("%d-%d-%d", instID, rqID, fnID))
}

// NewSubscription generates a subscription record from the E2AP subscription request
func NewSubscription(id ID, e2apsub *e2appducontents.RicsubscriptionRequest, ch e2ap.ClientConn) (*Subscription, error) {
	if id == "" {
		return nil, errors.New(errors.Forbidden, "id cannot be empty")
	}

	var rrID *e2apies.RicrequestId
	var rfID *e2apies.RanfunctionId
	var details *e2appducontents.RicsubscriptionDetails
	for _, v := range e2apsub.GetProtocolIes() {
		if v.Id == int32(v2.ProtocolIeIDRanfunctionID) {
			rfID = v.GetValue().GetRanfunctionId()
		}
		if v.Id == int32(v2.ProtocolIeIDRicrequestID) {
			rrID = v.GetValue().GetRicrequestId()
		}
		if v.Id == int32(v2.ProtocolIeIDRicsubscriptionDetails) {
			details = v.GetValue().GetRicsubscriptionDetails()
		}
	}

	return &Subscription{
		ID:        id,
		ReqID:     rrID,
		FnID:      rfID,
		Details:   details,
		E2Channel: ch,
	}, nil
}

// NewStore creates a new subscription store
func NewStore() *Subscriptions {
	return &Subscriptions{
		subscriptions: make(map[ID]*Subscription),
		mu:            sync.RWMutex{},
	}
}

// Store store interface
type Store interface {
	// Add   adds the specified subscription
	Add(subscription *Subscription) error
	// Remove removes the specified subscription
	Remove(id ID) error
	// Get gets a subscription based on a given ID
	Get(id ID) (*Subscription, error)
	// List lists subscriptions
	List() ([]*Subscription, error)
	// Len number of subscriptions
	Len() (int, error)
}

// Subscriptions data structure for storing subscriptions
type Subscriptions struct {
	subscriptions map[ID]*Subscription
	mu            sync.RWMutex
}

// Len number of subscriptions
func (s *Subscriptions) Len() (int, error) {
	return len(s.subscriptions), nil
}

// Add adds the specified subscription
func (s *Subscriptions) Add(sub *Subscription) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if sub.ID == "" {
		return errors.New(errors.Invalid, "Subscription ID cannot be empty")
	}
	s.subscriptions[sub.ID] = sub
	return nil
}

// Remove removes the specified subscription
func (s *Subscriptions) Remove(id ID) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if id == "" {
		return errors.New(errors.Invalid, "ID cannot be empty")
	}
	delete(s.subscriptions, id)
	return nil
}

// Get returns the subscription with the specified ID
func (s *Subscriptions) Get(id ID) (*Subscription, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	if val, ok := s.subscriptions[id]; ok {
		return val, nil
	}
	return nil, errors.New(errors.NotFound, "subscription entry has not been found")
}

// List returns slice containing all current subscriptions
func (s *Subscriptions) List() ([]*Subscription, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	resp := make([]*Subscription, 0, len(s.subscriptions))
	for _, sub := range s.subscriptions {
		resp = append(resp, sub)
	}
	return resp, nil
}

var _ Store = &Subscriptions{}
