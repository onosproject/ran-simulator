// SPDX-FileCopyrightText: 2020-present Open Networking Foundation <info@opennetworking.org>
//
// SPDX-License-Identifier: Apache-2.0

package measurement

import (
	"context"
	"github.com/onosproject/onos-lib-go/pkg/logging"
	"github.com/onosproject/ran-simulator/pkg/model"
	"github.com/onosproject/ran-simulator/pkg/store/cells"
	"github.com/onosproject/ran-simulator/pkg/store/ues"
	"github.com/onosproject/rrm-son-lib/pkg/model/device"
)

var logMeasCtrl = logging.GetLogger("measurement", "controller")

// NewMeasController returns the measurement controller object
func NewMeasController(measType MeasEventType, cellStore cells.Store, ueStore ues.Store) MeasController {
	return &measController{
		measType:   measType,
		cellStore:  cellStore,
		ueStore:    ueStore,
		inputChan:  make(chan *model.UE),
		outputChan: make(chan device.UE),
	}
}

// MeasController is an abstraction of the measurement controller
type MeasController interface {
	// Start starts measurement controller
	Start(ctx context.Context)

	// GetInputChan returns input channel
	GetInputChan() chan *model.UE

	// GetOutputChan returns output channel
	GetOutputChan() chan device.UE
}

// MeasEventType is the type for measurement event - currently it is string
// ToDo: define enumerated measurement type into rrm-son-lib
type MeasEventType string

type measController struct {
	cellStore  cells.Store
	ueStore    ues.Store
	measType   MeasEventType
	inputChan  chan *model.UE
	outputChan chan device.UE
}

func (m *measController) Start(ctx context.Context) {
	switch m.measType {
	case "EventA3":
		m.startMeasEventA3Handler(ctx)
	}
}

func (m *measController) startMeasEventA3Handler(ctx context.Context) {
	logMeasCtrl.Info("Measurement controller starting with EventA3Handler")
	handler := NewMeasEventA3()
	converter := NewMeasReportConverter(m.cellStore, m.ueStore)

	go handler.Start()
	// for input
	go m.forwardReportToEventA3Handler(ctx, handler, converter)
	// for output
	go m.forwardReportFromEventA3Handler(handler)
}

func (m *measController) forwardReportToEventA3Handler(ctx context.Context, handler MeasEventA3, converter MeasReportConverter) {
	for ue := range m.inputChan {
		logMeasCtrl.Debugf("[input] Measurement report to Event A3 handler: %v", *ue)
		report := converter.Convert(ctx, ue)
		handler.PushMeasurement(report)
	}
}

func (m *measController) forwardReportFromEventA3Handler(handler MeasEventA3) {
	for report := range handler.GetOutputChan() {
		logMeasCtrl.Debugf("[output] Measurement report for Event A3: %v", report)
		m.outputChan <- report
	}
}

func (m *measController) GetInputChan() chan *model.UE {
	return m.inputChan
}

func (m *measController) GetOutputChan() chan device.UE {
	return m.outputChan
}
