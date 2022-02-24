// SPDX-FileCopyrightText: 2020-present Open Networking Foundation <info@opennetworking.org>
//
// SPDX-License-Identifier: Apache-2.0

package handover

import (
	"github.com/onosproject/onos-lib-go/pkg/logging"
	"github.com/onosproject/rrm-son-lib/pkg/handover"
	"github.com/onosproject/rrm-son-lib/pkg/model/device"
)

var logA3ho = logging.GetLogger()

// A3Handover is an abstraction of A3 handover
type A3Handover interface {
	// Start starts the A3 handover module
	Start()

	// GetInputChan returns the channel to push measurement
	GetInputChan() chan device.UE

	// GetOutputChan returns the channel to get handover event
	GetOutputChan() chan handover.A3HandoverDecision

	// PushMeasurementEventA3 pushes measurement to the input channel
	PushMeasurementEventA3(device.UE)
}

type a3Handover struct {
	a3HandoverHandler *handover.A3HandoverHandler
}

// NewA3Handover returns an A3 handover object
func NewA3Handover() A3Handover {
	return &a3Handover{
		a3HandoverHandler: handover.NewA3HandoverHandler(),
	}
}

func (h *a3Handover) Start() {
	logA3ho.Info("A3 handover handler starting")
	go h.a3HandoverHandler.Run()
}

func (h *a3Handover) GetInputChan() chan device.UE {
	return h.a3HandoverHandler.Chans.InputChan
}

func (h *a3Handover) GetOutputChan() chan handover.A3HandoverDecision {
	return h.a3HandoverHandler.Chans.OutputChan
}

func (h *a3Handover) PushMeasurementEventA3(ue device.UE) {
	h.a3HandoverHandler.Chans.InputChan <- ue
}
