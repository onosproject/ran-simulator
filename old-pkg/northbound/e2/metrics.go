// SPDX-FileCopyrightText: 2020-present Open Networking Foundation <info@opennetworking.org>
//
// SPDX-License-Identifier: LicenseRef-ONF-Member-1.0
//

package e2

import (
	"time"

	"github.com/onosproject/ran-simulator/api/types"

	e2 "github.com/onosproject/onos-ric/api/sb"
	"github.com/onosproject/ran-simulator/pkg/manager"
)

// UpdateTelemetryMetrics ...
func UpdateTelemetryMetrics(msg *e2.RadioMeasReportPerUE, t time.Time) {
	trafficSimMgr := manager.GetManager()
	towerID := toTypesEcgi(msg.GetEcgi())
	name, err := trafficSimMgr.CrntiToName(types.Crnti(msg.GetCrnti()), towerID)
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

	trafficSimMgr.UserEquipmentsLock.RLock()
	servingTowerID := toE2Ecgi(ue.ServingTower)
	trafficSimMgr.UserEquipmentsLock.RUnlock()

	if servingTowerID.Ecid != bestStationID.Ecid || servingTowerID.PlmnId != bestStationID.PlmnId {
		trafficSimMgr.UserEquipmentsLock.Lock()
		m, err := generateReport(ue)
		if err == nil {
			msg := m.GetMsg().GetRadioMeasReportPerUE()
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

			servingTowerID := toE2Ecgi(ue.ServingTower)
			if servingTowerID.Ecid != bestStationID.Ecid || servingTowerID.PlmnId != bestStationID.PlmnId {
				if ue.Metrics.HoReportTimestamp == 0 {
					ue.Metrics.HoReportTimestamp = t.UnixNano()
				}
			} else {
				ue.Metrics.HoReportTimestamp = 0
			}
		}
		trafficSimMgr.UserEquipmentsLock.Unlock()
	}
}
