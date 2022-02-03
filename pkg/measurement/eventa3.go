// SPDX-FileCopyrightText: 2020-present Open Networking Foundation <info@opennetworking.org>
//
// SPDX-License-Identifier: Apache-2.0

package measurement

import (
	"github.com/onosproject/onos-lib-go/pkg/logging"
	"github.com/onosproject/rrm-son-lib/pkg/measurement"
	"github.com/onosproject/rrm-son-lib/pkg/model/device"
)

var logEventA3 = logging.GetLogger("measurement", "eventa3")

// MeasEventA3 is an abstraction of measurement Event A3
type MeasEventA3 interface {

	// Start starts the Event A3 module
	Start()

	// GetInputChan returns the channel to push all measurements
	GetInputChan() chan device.UE

	// GetOutputChan returns the channel to get measurements for Event A3
	GetOutputChan() chan device.UE

	// PushMeasurement pushes measurements to MeasEventA3 handler
	PushMeasurement(device.UE)
}

type measEventA3 struct {
	eventA3Handler *measurement.MeasEventA3Handler
}

// NewMeasEventA3 returns the measurement Event A3 object
func NewMeasEventA3() MeasEventA3 {
	return &measEventA3{
		eventA3Handler: measurement.NewMeasEventA3Handler(),
	}
}

func (m *measEventA3) Start() {
	logEventA3.Info("Measurement event A3 handler starting")
	go m.eventA3Handler.Run()
}

func (m *measEventA3) GetInputChan() chan device.UE {
	return m.eventA3Handler.Chans.InputChan
}

func (m *measEventA3) GetOutputChan() chan device.UE {
	return m.eventA3Handler.Chans.OutputChan
}

func (m *measEventA3) PushMeasurement(ue device.UE) {
	m.eventA3Handler.Chans.InputChan <- ue
}
