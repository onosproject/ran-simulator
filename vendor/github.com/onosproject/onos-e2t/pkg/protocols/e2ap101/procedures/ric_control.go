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

// RICControl is a RIC control procedure
type RICControl interface {
	RICControl(ctx context.Context, request *e2appducontents.RiccontrolRequest) (response *e2appducontents.RiccontrolAcknowledge, failure *e2appducontents.RiccontrolFailure, err error)
}

func NewRICControlInitiator(dispatcher Dispatcher) *RICControlInitiator {
	return &RICControlInitiator{
		dispatcher:  dispatcher,
		responseChs: make(map[int32]chan e2appdudescriptions.E2ApPdu),
	}
}

type RICControlInitiator struct {
	dispatcher  Dispatcher
	responseChs map[int32]chan e2appdudescriptions.E2ApPdu
	mu          sync.RWMutex
}

func (p *RICControlInitiator) Initiate(ctx context.Context, request *e2appducontents.RiccontrolRequest) (*e2appducontents.RiccontrolAcknowledge, *e2appducontents.RiccontrolFailure, error) {
	requestPDU := &e2appdudescriptions.E2ApPdu{
		E2ApPdu: &e2appdudescriptions.E2ApPdu_InitiatingMessage{
			InitiatingMessage: &e2appdudescriptions.InitiatingMessage{
				ProcedureCode: &e2appdudescriptions.E2ApElementaryProcedures{
					RicControl: &e2appdudescriptions.RicControl{
						InitiatingMessage: request,
					},
				},
			},
		},
	}
	if err := requestPDU.Validate(); err != nil {
		return nil, nil, errors.NewInvalid("E2AP PDU validation failed: %v", err)
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

	if err := p.dispatcher(requestPDU); err != nil {
		return nil, nil, errors.NewUnavailable("RIC Control initiation failed: %v", err)
	}

	select {
	case responsePDU, ok := <-responseCh:
		if !ok {
			return nil, nil, errors.NewUnavailable("connection closed")
		}

		switch response := responsePDU.E2ApPdu.(type) {
		case *e2appdudescriptions.E2ApPdu_SuccessfulOutcome:
			return response.SuccessfulOutcome.ProcedureCode.RicControl.SuccessfulOutcome, nil, nil
		case *e2appdudescriptions.E2ApPdu_UnsuccessfulOutcome:
			return nil, response.UnsuccessfulOutcome.ProcedureCode.RicControl.UnsuccessfulOutcome, nil
		default:
			return nil, nil, errors.NewInternal("received unexpected outcome")
		}
	case <-ctx.Done():
		return nil, nil, ctx.Err()
	}
}

func (p *RICControlInitiator) Matches(pdu *e2appdudescriptions.E2ApPdu) bool {
	switch msg := pdu.E2ApPdu.(type) {
	case *e2appdudescriptions.E2ApPdu_SuccessfulOutcome:
		return msg.SuccessfulOutcome.ProcedureCode.RicControl != nil
	case *e2appdudescriptions.E2ApPdu_UnsuccessfulOutcome:
		return msg.UnsuccessfulOutcome.ProcedureCode.RicControl != nil
	default:
		return false
	}
}

func (p *RICControlInitiator) Handle(pdu *e2appdudescriptions.E2ApPdu) {
	var requestID int32
	switch response := pdu.E2ApPdu.(type) {
	case *e2appdudescriptions.E2ApPdu_SuccessfulOutcome:
		requestID = response.SuccessfulOutcome.ProcedureCode.RicControl.SuccessfulOutcome.ProtocolIes.E2ApProtocolIes29.Value.RicRequestorId
	case *e2appdudescriptions.E2ApPdu_UnsuccessfulOutcome:
		requestID = response.UnsuccessfulOutcome.ProcedureCode.RicControl.UnsuccessfulOutcome.ProtocolIes.E2ApProtocolIes29.Value.RicRequestorId
	}

	p.mu.RLock()
	responseCh, ok := p.responseChs[requestID]
	p.mu.RUnlock()
	if ok {
		responseCh <- *pdu
		close(responseCh)
	} else {
		log.Errorf("Received RIC Control response for unknown request %d", requestID)
	}
}

func (p *RICControlInitiator) Close() error {
	p.mu.Lock()
	for _, responseCh := range p.responseChs {
		close(responseCh)
	}
	p.mu.Unlock()
	return nil
}

var _ ElementaryProcedure = &RICControlInitiator{}

func NewRICControlProcedure(dispatcher Dispatcher, handler RICControl) *RICControlProcedure {
	return &RICControlProcedure{
		dispatcher: dispatcher,
		handler:    handler,
	}
}

type RICControlProcedure struct {
	dispatcher Dispatcher
	handler    RICControl
}

func (p *RICControlProcedure) Matches(pdu *e2appdudescriptions.E2ApPdu) bool {
	switch msg := pdu.E2ApPdu.(type) {
	case *e2appdudescriptions.E2ApPdu_InitiatingMessage:
		return msg.InitiatingMessage.ProcedureCode.RicControl != nil
	default:
		return false
	}
}

func (p *RICControlProcedure) Handle(requestPDU *e2appdudescriptions.E2ApPdu) {
	response, failure, err := p.handler.RICControl(context.Background(), requestPDU.GetInitiatingMessage().ProcedureCode.RicControl.InitiatingMessage)
	if err != nil {
		log.Errorf("RIC Control procedure failed: %v", err)
	} else if response != nil {
		responsePDU := &e2appdudescriptions.E2ApPdu{
			E2ApPdu: &e2appdudescriptions.E2ApPdu_SuccessfulOutcome{
				SuccessfulOutcome: &e2appdudescriptions.SuccessfulOutcome{
					ProcedureCode: &e2appdudescriptions.E2ApElementaryProcedures{
						RicControl: &e2appdudescriptions.RicControl{
							SuccessfulOutcome: response,
						},
					},
				},
			},
		}
		if err := requestPDU.Validate(); err != nil {
			log.Errorf("RIC Control response validation failed: %v", err)
		} else {
			err := p.dispatcher(responsePDU)
			if err != nil {
				log.Errorf("RIC Control response failed: %v", err)
			}
		}
	} else if failure != nil {
		responsePDU := &e2appdudescriptions.E2ApPdu{
			E2ApPdu: &e2appdudescriptions.E2ApPdu_UnsuccessfulOutcome{
				UnsuccessfulOutcome: &e2appdudescriptions.UnsuccessfulOutcome{
					ProcedureCode: &e2appdudescriptions.E2ApElementaryProcedures{
						RicControl: &e2appdudescriptions.RicControl{
							UnsuccessfulOutcome: failure,
						},
					},
				},
			},
		}
		if err := requestPDU.Validate(); err != nil {
			log.Errorf("RIC Control response validation failed: %v", err)
		} else {
			err := p.dispatcher(responsePDU)
			if err != nil {
				log.Errorf("RIC Control response failed: %v", err)
			}
		}
	} else {
		log.Errorf("RIC Control function returned invalid output: no response message found")
	}
}

func (p *RICControlProcedure) Close() error {
	return nil
}

var _ ElementaryProcedure = &RICControlProcedure{}
