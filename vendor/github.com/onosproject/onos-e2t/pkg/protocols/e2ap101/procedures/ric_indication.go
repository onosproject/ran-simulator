// SPDX-FileCopyrightText: 2020-present Open Networking Foundation <info@opennetworking.org>
//
// SPDX-License-Identifier: LicenseRef-ONF-Member-1.0

package procedures

import (
	"context"
	e2ap "github.com/onosproject/onos-e2t/api/e2ap/v1beta2"
	"github.com/onosproject/onos-e2t/api/e2ap/v1beta2/e2ap-commondatatypes"
	"github.com/onosproject/onos-e2t/api/e2ap/v1beta2/e2ap-constants"
	e2appducontents "github.com/onosproject/onos-e2t/api/e2ap/v1beta2/e2ap-pdu-contents"
	e2appdudescriptions "github.com/onosproject/onos-e2t/api/e2ap/v1beta2/e2ap-pdu-descriptions"
	"github.com/onosproject/onos-lib-go/pkg/errors"
)

// RICIndication is a RIC indication procedure
type RICIndication interface {
	RICIndication(ctx context.Context, request *e2appducontents.Ricindication) (err error)
}

func NewRICIndicationInitiator(dispatcher Dispatcher) *RICIndicationInitiator {
	return &RICIndicationInitiator{
		dispatcher: dispatcher,
	}
}

type RICIndicationInitiator struct {
	dispatcher Dispatcher
}

func (p *RICIndicationInitiator) Initiate(ctx context.Context, request *e2appducontents.Ricindication) (err error) {
	pdu := &e2appdudescriptions.E2ApPdu{
		E2ApPdu: &e2appdudescriptions.E2ApPdu_InitiatingMessage{
			InitiatingMessage: &e2appdudescriptions.InitiatingMessage{
				ProcedureCode: &e2appdudescriptions.E2ApElementaryProcedures{
					RicIndication: &e2appdudescriptions.RicIndication{
						InitiatingMessage: request,
						ProcedureCode: &e2ap_constants.IdRicindication{
							Value: int32(e2ap.ProcedureCodeIDRICindication),
						},
						Criticality: &e2ap_commondatatypes.CriticalityIgnore{
							Criticality: e2ap_commondatatypes.Criticality_CRITICALITY_IGNORE,
						},
					},
				},
			},
		},
	}
	if err := pdu.Validate(); err != nil {
		return errors.NewInvalid("E2AP PDU validation failed: %v", err)
	}
	return p.dispatcher(pdu)
}

func (p *RICIndicationInitiator) Matches(pdu *e2appdudescriptions.E2ApPdu) bool {
	return false
}

func (p *RICIndicationInitiator) Handle(pdu *e2appdudescriptions.E2ApPdu) {

}

func (p *RICIndicationInitiator) Close() error {
	return nil
}

var _ ElementaryProcedure = &RICIndicationInitiator{}

func NewRICIndicationProcedure(dispatcher Dispatcher, handler RICIndication) *RICIndicationProcedure {
	return &RICIndicationProcedure{
		dispatcher: dispatcher,
		handler:    handler,
	}
}

type RICIndicationProcedure struct {
	dispatcher Dispatcher
	handler    RICIndication
}

func (p *RICIndicationProcedure) Matches(pdu *e2appdudescriptions.E2ApPdu) bool {
	switch msg := pdu.E2ApPdu.(type) {
	case *e2appdudescriptions.E2ApPdu_InitiatingMessage:
		return msg.InitiatingMessage.ProcedureCode.RicIndication != nil
	default:
		return false
	}
}

func (p *RICIndicationProcedure) Handle(pdu *e2appdudescriptions.E2ApPdu) {
	err := p.handler.RICIndication(context.Background(), pdu.GetInitiatingMessage().ProcedureCode.RicIndication.InitiatingMessage)
	if err != nil {
		log.Errorf("RIC Indication procedure failed: %v", err)
	}
}

func (p *RICIndicationProcedure) Close() error {
	return nil
}

var _ ElementaryProcedure = &RICIndicationProcedure{}
