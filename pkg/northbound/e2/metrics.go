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
	switch m.MessageType {
	case e2.MessageType_RADIO_MEAS_REPORT_PER_UE:
		x, ok := m.S.(*e2.TelemetryMessage_RadioMeasReportPerUE)
		if !ok {
			log.Fatalf("Unexpected payload for RADIO_MEAS_REPORT_PER_UE message %v", m)
		}
		r := x.RadioMeasReportPerUE
		towerID := toTypesEcgi(r.Ecgi)
		name, err := trafficSimMgr.CrntiToName(types.Crnti(r.Crnti), &towerID)
		if err != nil {
			log.Errorf("ue %s/%s not found", r.Ecgi.Ecid, r.Crnti)
			return
		}
		var ue *types.Ue
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
		servingTower := trafficSimMgr.Towers[*ue.ServingTower]
		trafficSimMgr.TowersLock.RUnlock()

		if servingTower.Ecgi.EcID != types.EcID(bestStationID.Ecid) || servingTower.Ecgi.PlmnID != types.PlmnID(bestStationID.PlmnId) {
			trafficSimMgr.UserEquipmentsLock.Lock()
			if ue.Metrics.HoReportTimestamp == 0 {
				ue.Metrics.HoReportTimestamp = time.Now().UnixNano()
			}
			trafficSimMgr.UserEquipmentsLock.Unlock()
		}
	}
}

// UpdateControlMetrics ...
func UpdateControlMetrics(imsi types.Imsi) {
	trafficSimMgr := manager.GetManager()
	trafficSimMgr.UserEquipmentsLock.Lock()
	defer trafficSimMgr.UserEquipmentsLock.Unlock()
	var ok bool
	var ue *types.Ue
	if ue, ok = trafficSimMgr.UserEquipments[imsi]; !ok {
		log.Errorf("ue %s not found", imsi)
		return
	}
	if ue.Metrics.IsFirst {
		// Discard the first one as it may have been waiting for onos-ric-ho to startup
		ue.Metrics.HoReportTimestamp = 0
		ue.Metrics.IsFirst = false
	} else if ue.Metrics.HoReportTimestamp != 0 {
		ue.Metrics.HoLatency = time.Now().UnixNano() - ue.Metrics.HoReportTimestamp
		ue.Metrics.HoReportTimestamp = 0
		tmpHOEvent := metrics.HOEvent{
			Timestamp:    time.Now(),
			Crnti:        ue.GetCrnti(),
			ServingTower: *ue.ServingTower,
			HOLatency:    ue.Metrics.HoLatency,
		}
		trafficSimMgr.LatencyChannel <- tmpHOEvent
		log.Infof("%d Hand-over latency: %d µs", ue.Imsi, ue.Metrics.HoLatency/1000)
	}
}
