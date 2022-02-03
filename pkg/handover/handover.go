// SPDX-FileCopyrightText: 2020-present Open Networking Foundation <info@opennetworking.org>
//
// SPDX-License-Identifier: Apache-2.0

package handover

import (
	"context"
	"github.com/onosproject/onos-lib-go/pkg/logging"
	"github.com/onosproject/ran-simulator/pkg/store/cells"
	"github.com/onosproject/ran-simulator/pkg/store/ues"
	"github.com/onosproject/rrm-son-lib/pkg/handover"
	"github.com/onosproject/rrm-son-lib/pkg/model/device"
)

var logHoCtrl = logging.GetLogger("handover", "controller")

// NewHOController returns the hanover controller
func NewHOController(hoType HOType, cellStore cells.Store, ueStore ues.Store) HOController {
	return &hoController{
		hoType:     hoType,
		cellStore:  cellStore,
		ueStore:    ueStore,
		inputChan:  make(chan device.UE),
		outputChan: make(chan handover.A3HandoverDecision),
	}
}

// HOController is an abstraction of the handover controller
type HOController interface {
	// Start starts handover controller
	Start(ctx context.Context)

	// GetInputChan returns input channel
	GetInputChan() chan device.UE

	// GetOutputChan returns output channel
	GetOutputChan() chan handover.A3HandoverDecision
}

// HOType is the type of hanover - currently it is string
// ToDo: define enumerated handover type into rrm-son-lib
type HOType string

type hoController struct {
	cellStore  cells.Store
	ueStore    ues.Store
	hoType     HOType
	inputChan  chan device.UE
	outputChan chan handover.A3HandoverDecision
}

func (h *hoController) Start(ctx context.Context) {
	switch h.hoType {
	case "A3":
		h.startA3HandoverHandler(ctx)
	}
}

func (h *hoController) startA3HandoverHandler(ctx context.Context) {
	logHoCtrl.Info("Handover controller starting with A3HandoveHandler")
	handler := NewA3Handover()

	go handler.Start()
	// for input
	go h.forwardReportToA3HandoverHandler(handler)
	//for output
	go h.forwardHandoverDecision(handler)
}

func (h *hoController) forwardReportToA3HandoverHandler(handler A3Handover) {
	for ue := range h.inputChan {
		logHoCtrl.Debugf("[input] Measurement report for HO decision: %v", ue)
		handler.PushMeasurementEventA3(ue)
	}
}

func (h *hoController) forwardHandoverDecision(handler A3Handover) {
	for hoDecision := range handler.GetOutputChan() {
		logHoCtrl.Debugf("[output] Handover decision: %v", hoDecision)
		h.outputChan <- hoDecision
	}
}

func (h *hoController) GetInputChan() chan device.UE {
	return h.inputChan
}

func (h *hoController) GetOutputChan() chan handover.A3HandoverDecision {
	return h.outputChan
}
