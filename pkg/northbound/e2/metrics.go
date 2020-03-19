// Copyright 2020-present Open Networking Foundation.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package e2

import (
	"time"

	"github.com/onosproject/ran-simulator/api/types"

	e2 "github.com/onosproject/onos-ric/api/sb"
	"github.com/onosproject/ran-simulator/pkg/manager"
	"github.com/onosproject/ran-simulator/pkg/northbound/metrics"
)

// UpdateTelemetryMetrics ...
func UpdateTelemetryMetrics(m *e2.TelemetryMessage) {
	trafficSimMgr := manager.GetManager()
	switch x := m.S.(type) {
	case *e2.TelemetryMessage_RadioMeasReportPerUE:
		r := x.RadioMeasReportPerUE
		name, err := trafficSimMgr.CrntiToName(types.Crnti(r.Crnti), types.EcID(r.Ecgi.Ecid))
		if err != nil {
			log.Errorf("ue %s/%s not found", r.Ecgi.Ecid, r.Crnti)
		}
		var ue *types.Ue
		var ok bool
		trafficSimMgr.UserEquipmentsLock.RLock()
		if ue, ok = trafficSimMgr.UserEquipments[name]; !ok {
			trafficSimMgr.UserEquipmentsLock.RUnlock()
			return
		}
		trafficSimMgr.UserEquipmentsLock.RUnlock()

		reports := r.RadioReportServCells

		bestCQI := reports[0].CqiHist[0]
		bestStationID := reports[0].Ecgi

		for i := 1; i < len(reports); i++ {
			temp := reports[i].CqiHist[0]
			if bestCQI < temp {
				bestCQI = temp
				bestStationID = reports[i].Ecgi
			}
		}
		trafficSimMgr.TowersLock.RLock()
		servingTower := trafficSimMgr.Towers[ue.ServingTower]
		trafficSimMgr.TowersLock.RUnlock()

		if servingTower.EcID != types.EcID(bestStationID.Ecid) || servingTower.PlmnID != types.PlmnID(bestStationID.PlmnId) {
			trafficSimMgr.UserEquipmentsLock.Lock()
			if ue.Metrics.HoReportTimestamp == 0 {
				ue.Metrics.HoReportTimestamp = time.Now().UnixNano()
			}
			trafficSimMgr.UserEquipmentsLock.Unlock()
		}
	}
}

// UpdateControlMetrics ...
func UpdateControlMetrics(in *e2.ControlResponse) {
	trafficSimMgr := manager.GetManager()
	switch x := in.S.(type) {
	case *e2.ControlResponse_HORequest:
		m := x.HORequest
		trafficSimMgr.UserEquipmentsLock.Lock()
		defer trafficSimMgr.UserEquipmentsLock.Unlock()
		ueName, err := trafficSimMgr.CrntiToName(types.Crnti(m.Crnti), types.EcID(m.EcgiS.Ecid))
		if err != nil {
			log.Errorf("ue %s/%s not found", m.EcgiS.Ecid, m.Crnti)
			return
		}
		var ue *types.Ue
		var ok bool
		if ue, ok = trafficSimMgr.UserEquipments[ueName]; !ok {
			log.Errorf("ue %s not found", ueName)
			return
		}
		if ue.Metrics.HoReportTimestamp != 0 {
			ue.Metrics.HoLatency = time.Now().UnixNano() - ue.Metrics.HoReportTimestamp
			ue.Metrics.HoReportTimestamp = 0
			tmpHOEvent := metrics.HOEvent{
				Timestamp:    time.Now(),
				Crnti:        string(ue.GetCrnti()),
				ServingTower: string(ue.GetServingTower()),
				HOLatency:    ue.Metrics.HoLatency,
			}
			trafficSimMgr.LatencyChannel <- tmpHOEvent
			log.Infof("%s Hand-over latency: %d microsec", ue.Name, ue.Metrics.HoLatency/1000)
		}
	}
}
