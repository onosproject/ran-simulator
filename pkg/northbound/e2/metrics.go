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
	"github.com/onosproject/onos-ric/api/sb/e2ap"
	"github.com/onosproject/ran-simulator/pkg/manager"
	"github.com/onosproject/ran-simulator/pkg/northbound/metrics"
)

// UpdateTelemetryMetrics ...
func UpdateTelemetryMetrics(m *e2ap.RicIndication) {
	trafficSimMgr := manager.GetManager()
	switch m.GetHdr().GetMessageType() {
	case e2.MessageType_RADIO_MEAS_REPORT_PER_UE:
		msg := m.GetMsg().GetRadioMeasReportPerUE()
		towerID := toTypesEcgi(msg.GetEcgi())
		name, err := trafficSimMgr.CrntiToName(types.Crnti(msg.GetCrnti()), &towerID)
		if err != nil {
			log.Errorf("ue %s/%s not found", msg.GetEcgi().GetEcid(), msg.GetCrnti())
			return
		}
		var ue *types.Ue
		var ok bool
		trafficSimMgr.UserEquipmentsMapLock.RLock()
		if ue, ok = trafficSimMgr.UserEquipments[name]; !ok {
			trafficSimMgr.UserEquipmentsMapLock.RUnlock()
			return
		}
		trafficSimMgr.UserEquipmentsMapLock.RUnlock()

		reports := msg.RadioReportServCells

		bestCQI := reports[0].CqiHist[0]
		bestStationID := reports[0].Ecgi

		for i := 1; i < len(reports); i++ {
			temp := reports[i].CqiHist[0]
			if bestCQI < temp {
				bestCQI = temp
				bestStationID = reports[i].Ecgi
			}
		}
		trafficSimMgr.CellsLock.RLock()
		servingTower := trafficSimMgr.Cells[*ue.ServingTower]
		trafficSimMgr.CellsLock.RUnlock()

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
	trafficSimMgr.UserEquipmentsMapLock.RLock()
	var ok bool
	var ue *types.Ue
	if ue, ok = trafficSimMgr.UserEquipments[imsi]; !ok {
		log.Errorf("ue %s not found", imsi)
		trafficSimMgr.UserEquipmentsMapLock.RUnlock()
		return
	}
	trafficSimMgr.UserEquipmentsMapLock.RUnlock()
	if ue.Metrics.IsFirst {
		// Discard the first one as it may have been waiting for onos-ric-ho to startup
		trafficSimMgr.UserEquipmentsLock.Lock()
		ue.Metrics.HoReportTimestamp = 0
		ue.Metrics.IsFirst = false
		trafficSimMgr.UserEquipmentsLock.Unlock()
	} else if ue.Metrics.HoReportTimestamp != 0 {
		trafficSimMgr.UserEquipmentsLock.Lock()
		ue.Metrics.HoLatency = time.Now().UnixNano() - ue.Metrics.HoReportTimestamp
		ue.Metrics.HoReportTimestamp = 0
		tmpHOEvent := metrics.HOEvent{
			Timestamp:    time.Now(),
			Crnti:        ue.GetCrnti(),
			ServingTower: *ue.ServingTower,
			HOLatency:    ue.Metrics.HoLatency,
		}
		trafficSimMgr.UserEquipmentsLock.Unlock()
		trafficSimMgr.LatencyChannel <- tmpHOEvent
		log.Infof("%d Hand-over latency: %d µs", ue.Imsi, ue.Metrics.HoLatency/1000)
	}
}
