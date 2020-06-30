// SPDX-FileCopyrightText: 2020-present Open Networking Foundation <info@opennetworking.org>
//
// SPDX-License-Identifier: LicenseRef-ONF-Member-1.0
//

package e2

import (
	"fmt"
	"math"
	"time"

	e2 "github.com/onosproject/onos-ric/api/sb"
	"github.com/onosproject/onos-ric/api/sb/e2ap"
	"github.com/onosproject/onos-ric/api/sb/e2sm"
	"github.com/onosproject/ran-simulator/api/trafficsim"
	"github.com/onosproject/ran-simulator/api/types"
	"github.com/onosproject/ran-simulator/pkg/manager"
)

const e2TelemetryNbi = "e2TelemetryNbi"

// Conversion of signal strength in dB to CQI
// Here we just map 0 dB to the middle of the scale 0-15 CQI scale
func makeCqi(strengthdB float64) uint32 {
	cqi := math.Round(strengthdB) + 7
	if cqi > 15 {
		return 15
	} else if cqi < 0 {
		return 0
	}

	return uint32(cqi)
}

func (s *Server) radioMeasReportPerUE(indChan chan e2ap.RicIndication, stream e2ap.E2AP_RicChanServer) error {
	trafficSimMgr := manager.GetManager()

	streamID := fmt.Sprintf("%s-%p", e2TelemetryNbi, stream)
	ueChangeChannel, err := trafficSimMgr.Dispatcher.RegisterUeListener(streamID, int(trafficSimMgr.MapLayout.MaxUes))
	if err != nil {
		log.Errorf("RegisterUeListener failed for ServingTower=%s for Port %d", s.GetECGI(), s.GetPort())
		return err
	}
	defer trafficSimMgr.Dispatcher.UnregisterUeListener(streamID)

	log.Infof("Waiting for l2MeasConfig for ServingTower=%s for Port %d", s.GetECGI(), s.GetPort())
	configDone := make(chan bool)
	go s.waitForConfig(configDone)
	<-configDone

	log.Infof("Listening for changes on UEs with ServingTower=%s for Port %d", s.GetECGI(), s.GetPort())

	for {
		select {
		case ueUpdate := <-ueChangeChannel:
			if ueUpdate.Type == trafficsim.Type_UPDATED && ueUpdate.UpdateType == trafficsim.UpdateType_TOWER {
				ue, ok := ueUpdate.Object.(*types.Ue)
				if !ok {
					log.Fatalf("Object %v could not be converted to UE", ueUpdate)
				}
				if ue.ServingTower.EcID != s.GetECGI().EcID || ue.ServingTower.PlmnID != s.GetECGI().PlmnID {
					continue
				}
				ind, err := generateReport(ue)
				if err != nil {
					log.Warnf("generateReport returned error %v", err)
				} else {
					indChan <- ind
				}
			}
		case <-stream.Context().Done():
			log.Infof("Controller has disconnected on Port %d", s.GetPort())
			return nil
		}
	}
}

func (s *Server) waitForConfig(configDone chan bool) {
	ticker := time.NewTicker(500 * time.Millisecond)
	for range ticker.C {
		if s.l2MeasConfig.RadioMeasReportPerUe != 0 {
			ticker.Stop()
			configDone <- true
			return
		}
	}
}

func generateReport(ue *types.Ue) (e2ap.RicIndication, error) {
	trafficSimMgr := manager.GetManager()
	if ue == nil {
		return e2ap.RicIndication{}, fmt.Errorf("ue is empty when generating RicIndication")
	}

	trafficSimMgr.CellsLock.RLock()
	defer trafficSimMgr.CellsLock.RUnlock()

	servingTower, servingOk := trafficSimMgr.Cells[*ue.ServingTower]
	if !servingOk {
		return e2ap.RicIndication{}, fmt.Errorf("serving tower not found %s", *ue.ServingTower)
	}
	tower1, t1ok := trafficSimMgr.Cells[*ue.Tower1]
	tower2, t2ok := trafficSimMgr.Cells[*ue.Tower2]
	tower3, t3ok := trafficSimMgr.Cells[*ue.Tower3]

	reports := make([]*e2.RadioRepPerServCell, 4)

	reports[0] = new(e2.RadioRepPerServCell)
	sTowerEcgi := toE2Ecgi(servingTower.Ecgi)
	reports[0].Ecgi = &sTowerEcgi
	reports[0].CqiHist = make([]uint32, 1)
	reports[0].CqiHist[0] = makeCqi(ue.ServingTowerStrength)

	if t1ok {
		reports[1] = new(e2.RadioRepPerServCell)
		tower1Ecgi := toE2Ecgi(tower1.Ecgi)
		reports[1].Ecgi = &tower1Ecgi
		reports[1].CqiHist = make([]uint32, 1)
		reports[1].CqiHist[0] = makeCqi(ue.Tower1Strength)
	}

	if t2ok {
		reports[2] = new(e2.RadioRepPerServCell)
		tower2Ecgi := toE2Ecgi(tower2.Ecgi)
		reports[2].Ecgi = &tower2Ecgi
		reports[2].CqiHist = make([]uint32, 1)
		reports[2].CqiHist[0] = makeCqi(ue.Tower2Strength)
	}

	if t3ok {
		reports[3] = new(e2.RadioRepPerServCell)
		tower3Ecgi := toE2Ecgi(tower3.Ecgi)
		reports[3].Ecgi = &tower3Ecgi
		reports[3].CqiHist = make([]uint32, 1)
		reports[3].CqiHist[0] = makeCqi(ue.Tower3Strength)
	}

	log.Infof("RadioMeasReportPerUE %s [cqi:%d] %s(%d) cqi:%d(%s),%d(%s),%d(%s)",
		servingTower.Ecgi.EcID, reports[0].CqiHist[0], ue.Crnti, ue.Imsi,
		reports[1].CqiHist[0], reports[1].Ecgi.Ecid,
		reports[2].CqiHist[0], reports[2].Ecgi.Ecid,
		reports[3].CqiHist[0], reports[3].Ecgi.Ecid)

	servingTower2Ecgi := toE2Ecgi(servingTower.Ecgi)
	return e2ap.RicIndication{
		Hdr: &e2sm.RicIndicationHeader{
			MessageType: e2.MessageType_RADIO_MEAS_REPORT_PER_UE,
		},
		Msg: &e2sm.RicIndicationMessage{
			S: &e2sm.RicIndicationMessage_RadioMeasReportPerUE{
				RadioMeasReportPerUE: &e2.RadioMeasReportPerUE{
					Ecgi:                 &servingTower2Ecgi,
					Crnti:                string(ue.Crnti),
					RadioReportServCells: reports,
				},
			},
		},
	}, nil
}
