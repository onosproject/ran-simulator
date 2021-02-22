// SPDX-FileCopyrightText: 2020-present Open Networking Foundation <info@opennetworking.org>
//
// SPDX-License-Identifier: LicenseRef-ONF-Member-1.0

package procedures

import (
	"context"
	e2appducontents "github.com/onosproject/onos-e2t/api/e2ap/v1beta2/e2ap-pdu-contents"
	e2appdudescriptions "github.com/onosproject/onos-e2t/api/e2ap/v1beta2/e2ap-pdu-descriptions"
	"github.com/onosproject/onos-lib-go/pkg/errors"
	"sync"
)

// RICSubscriptionDelete is a RIC subscription delete procedure
type RICSubscriptionDelete interface {
	RICSubscriptionDelete(ctx context.Context, request *e2appducontents.RicsubscriptionDeleteRequest) (response *e2appducontents.RicsubscriptionDeleteResponse, failure *e2appducontents.RicsubscriptionDeleteFailure, err error)
}

func NewRICSubscriptionDeleteInitiator(dispatcher Dispatcher) *RICSubscriptionDeleteInitiator {
	return &RICSubscriptionDeleteInitiator{
		dispatcher:  dispatcher,
		responseChs: make(map[int32]chan e2appdudescriptions.E2ApPdu),
	}
}

type RICSubscriptionDeleteInitiator struct {
	dispatcher  Dispatcher
	responseChs map[int32]chan e2appdudescriptions.E2ApPdu
	mu          sync.RWMutex
}

func (p *RICSubscriptionDeleteInitiator) Initiate(ctx context.Context, request *e2appducontents.RicsubscriptionDeleteRequest) (*e2appducontents.RicsubscriptionDeleteResponse, *e2appducontents.RicsubscriptionDeleteFailure, error) {
	requestPDU := &e2appdudescriptions.E2ApPdu{
		E2ApPdu: &e2appdudescriptions.E2ApPdu_InitiatingMessage{
			InitiatingMessage: &e2appdudescriptions.InitiatingMessage{
				ProcedureCode: &e2appdudescriptions.E2ApElementaryProcedures{
					RicSubscriptionDelete: &e2appdudescriptions.RicSubscriptionDelete{
						InitiatingMessage: request,
					},
				},
			},
		},
	}
	if err := requestPDU.Validate(); err != nil {
		return nil, nil, errors.NewInvalid("E2AP PDU validation failed: %v", err)
	}

	if err := p.dispatcher(requestPDU); err != nil {
		return nil, nil, errors.NewUnavailable("RIC Subscription Delete initiation failed: %v", err)
	}

	responseCh := make(chan e2appdudescriptions.E2ApPdu, 1)
	requestID := request.ProtocolIes.E2ApProtocolIes29.Value.RicRequestorId
	p.mu.Lock()
	p.responseChs[requestID] = responseCh
	p.mu.Unlock()

	defer func() {
		p.mu.Lock()
		delete(p.responseChs, requestID)
		p.mu.Unlock()
	}()

	select {
	case responsePDU, ok := <-responseCh:
		if !ok {
			return nil, nil, errors.NewUnavailable("connection closed")
		}

		switch response := responsePDU.E2ApPdu.(type) {
		case *e2appdudescriptions.E2ApPdu_SuccessfulOutcome:
			return response.SuccessfulOutcome.ProcedureCode.RicSubscriptionDelete.SuccessfulOutcome, nil, nil
		case *e2appdudescriptions.E2ApPdu_UnsuccessfulOutcome:
			return nil, response.UnsuccessfulOutcome.ProcedureCode.RicSubscriptionDelete.UnsuccessfulOutcome, nil
		default:
			return nil, nil, errors.NewInternal("received unexpected outcome")
		}
	case <-ctx.Done():
		return nil, nil, ctx.Err()
	}
}

func (p *RICSubscriptionDeleteInitiator) Matches(pdu *e2appdudescriptions.E2ApPdu) bool {
	switch msg := pdu.E2ApPdu.(type) {
	case *e2appdudescriptions.E2ApPdu_SuccessfulOutcome:
		return msg.SuccessfulOutcome.ProcedureCode.RicSubscriptionDelete != nil
	case *e2appdudescriptions.E2ApPdu_UnsuccessfulOutcome:
		return msg.UnsuccessfulOutcome.ProcedureCode.RicSubscriptionDelete != nil
	default:
		return false
	}
}

func (p *RICSubscriptionDeleteInitiator) Handle(pdu *e2appdudescriptions.E2ApPdu) {
	var requestID int32
	switch response := pdu.E2ApPdu.(type) {
	case *e2appdudescriptions.E2ApPdu_SuccessfulOutcome:
		requestID = response.SuccessfulOutcome.ProcedureCode.RicSubscriptionDelete.SuccessfulOutcome.ProtocolIes.E2ApProtocolIes29.Value.RicRequestorId
	case *e2appdudescriptions.E2ApPdu_UnsuccessfulOutcome:
		requestID = response.UnsuccessfulOutcome.ProcedureCode.RicSubscriptionDelete.UnsuccessfulOutcome.ProtocolIes.E2ApProtocolIes29.Value.RicRequestorId
	}

	p.mu.RLock()
	responseCh, ok := p.responseChs[requestID]
	p.mu.RUnlock()
	if ok {
		responseCh <- *pdu
		close(responseCh)
	} else {
		log.Errorf("Received RIC Subscription Delete response for unknown request %d", requestID)
	}
}

func (p *RICSubscriptionDeleteInitiator) Close() error {
	p.mu.Lock()
	for _, responseCh := range p.responseChs {
		close(responseCh)
	}
	p.mu.Unlock()
	return nil
}

var _ ElementaryProcedure = &RICSubscriptionDeleteInitiator{}

func NewRICSubscriptionDeleteProcedure(dispatcher Dispatcher, handler RICSubscriptionDelete) *RICSubscriptionDeleteProcedure {
	return &RICSubscriptionDeleteProcedure{
		dispatcher: dispatcher,
		handler:    handler,
	}
}

type RICSubscriptionDeleteProcedure struct {
	dispatcher Dispatcher
	handler    RICSubscriptionDelete
}

func (p *RICSubscriptionDeleteProcedure) Matches(pdu *e2appdudescriptions.E2ApPdu) bool {
	switch msg := pdu.E2ApPdu.(type) {
	case *e2appdudescriptions.E2ApPdu_InitiatingMessage:
		return msg.InitiatingMessage.ProcedureCode.RicSubscriptionDelete != nil
	default:
		return false
	}
}

func (p *RICSubscriptionDeleteProcedure) Handle(requestPDU *e2appdudescriptions.E2ApPdu) {
	response, failure, err := p.handler.RICSubscriptionDelete(context.Background(), requestPDU.GetInitiatingMessage().ProcedureCode.RicSubscriptionDelete.InitiatingMessage)
	if err != nil {
		log.Errorf("RIC Subscription Delete procedure failed: %v", err)
	} else if response != nil {
		responsePDU := &e2appdudescriptions.E2ApPdu{
			E2ApPdu: &e2appdudescriptions.E2ApPdu_SuccessfulOutcome{
				SuccessfulOutcome: &e2appdudescriptions.SuccessfulOutcome{
					ProcedureCode: &e2appdudescriptions.E2ApElementaryProcedures{
						RicSubscriptionDelete: &e2appdudescriptions.RicSubscriptionDelete{
							SuccessfulOutcome: response,
						},
					},
				},
			},
		}
		if err := requestPDU.Validate(); err != nil {
			log.Errorf("RIC Subscription Delete response validation failed: %v", err)
		} else {
			err := p.dispatcher(responsePDU)
			if err != nil {
				log.Errorf("RIC Subscription Delete response failed: %v", err)
			}
		}
	} else if failure != nil {
		responsePDU := &e2appdudescriptions.E2ApPdu{
			E2ApPdu: &e2appdudescriptions.E2ApPdu_UnsuccessfulOutcome{
				UnsuccessfulOutcome: &e2appdudescriptions.UnsuccessfulOutcome{
					ProcedureCode: &e2appdudescriptions.E2ApElementaryProcedures{
						RicSubscriptionDelete: &e2appdudescriptions.RicSubscriptionDelete{
							UnsuccessfulOutcome: failure,
						},
					},
				},
			},
		}
		if err := requestPDU.Validate(); err != nil {
			log.Errorf("RIC Subscription Delete response validation failed: %v", err)
		} else {
			err := p.dispatcher(responsePDU)
			if err != nil {
				log.Errorf("RIC Subscription Delete response failed: %v", err)
			}
		}
	} else {
		log.Errorf("RIC Subscription Delete function returned invalid output: no response message found")
	}
}

func (p *RICSubscriptionDeleteProcedure) Close() error {
	return nil
}

var _ ElementaryProcedure = &RICSubscriptionDeleteProcedure{}
