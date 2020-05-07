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
	"fmt"
	"math"
	"time"

	e2 "github.com/onosproject/onos-ric/api/sb"
	"github.com/onosproject/onos-ric/api/sb/e2ap"
	"github.com/onosproject/onos-ric/api/sb/e2sm"
	"github.com/onosproject/ran-simulator/api/trafficsim"
	"github.com/onosproject/ran-simulator/api/types"
	"github.com/onosproject/ran-simulator/pkg/dispatcher"
	"github.com/onosproject/ran-simulator/pkg/manager"
)

const e2TelemetryNbi = "e2TelemetryNbi"

const ueChangeChannelLen = 1000

// Conversion of signal strength in dB to CQI
func makeCqi(strengthdB float64) uint32 {
	// TODO normalize this across the range
	cqi := uint32(math.Round(strengthdB) + 7)
	if cqi > 15 {
		cqi = 15
	}
	return cqi
}

func (s *Server) radioMeasReportPerUE() error {
	trafficSimMgr := manager.GetManager()

	streamID := fmt.Sprintf("%s-%p", e2TelemetryNbi, s.stream)
	ueChangeChannel, err := trafficSimMgr.Dispatcher.RegisterUeListener(streamID, ueChangeChannelLen)
	defer trafficSimMgr.Dispatcher.UnregisterUeListener(streamID)
	if err != nil {
		log.Errorf("RegisterUeListener failed for ServingTower=%s for Port %d", s.GetECGI(), s.GetPort())
		return err
	}

	log.Infof("Waiting for l2MeasConfig for ServingTower=%s for Port %d", s.GetECGI(), s.GetPort())
	configDone := make(chan bool)
	go s.waitForConfig(configDone)
	<-configDone

	s.telemetryTicker = time.NewTicker(time.Duration(s.l2MeasConfig.RadioMeasReportPerUe) * time.Millisecond)

	log.Infof("Listening for changes on UEs with ServingTower=%s for Port %d", s.GetECGI(), s.GetPort())

	for {
		select {
		case <-s.telemetryTicker.C:
			ues, err := processUeChange(ueChangeChannel, s.stream)
			if err != nil || ues == nil {
				continue
			}
			for _, ue := range ues {
				s.indChan <- generateReport(ue)
			}
		case <-s.stream.Context().Done():
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

func processUeChange(ueChangeChannel chan dispatcher.Event, stream e2ap.E2AP_RicChanServer) ([]*types.Ue, error) {
	var ues []*types.Ue
	num := len(ueChangeChannel)
	for i := 0; i < num; i++ {
		select {
		// block here and listen out for any updates to UEs
		case ueUpdate := <-ueChangeChannel:
			if ueUpdate.Type == trafficsim.Type_UPDATED && ueUpdate.UpdateType == trafficsim.UpdateType_TOWER {
				ues = append(ues, ueUpdate.Object.(*types.Ue))
			}
		case <-stream.Context().Done():
			return nil, fmt.Errorf("Controller has disconnected")
		}
	}
	return ues, nil
}

func generateReport(ue *types.Ue) e2ap.RicIndication {
	trafficSimMgr := manager.GetManager()

	trafficSimMgr.CellsLock.RLock()
	defer trafficSimMgr.CellsLock.RUnlock()

	servingTower := trafficSimMgr.Cells[*ue.ServingTower]
	tower1 := trafficSimMgr.Cells[*ue.Tower1]
	tower2 := trafficSimMgr.Cells[*ue.Tower2]
	tower3 := trafficSimMgr.Cells[*ue.Tower3]
	sTower := trafficSimMgr.Cells[*ue.ServingTower]

	reports := make([]*e2.RadioRepPerServCell, 4)

	reports[0] = new(e2.RadioRepPerServCell)
	sTowerEcgi := toE2Ecgi(sTower.Ecgi)
	reports[0].Ecgi = &sTowerEcgi
	reports[0].CqiHist = make([]uint32, 1)
	reports[0].CqiHist[0] = makeCqi(ue.ServingTowerStrength)

	reports[1] = new(e2.RadioRepPerServCell)
	tower1Ecgi := toE2Ecgi(tower1.Ecgi)
	reports[1].Ecgi = &tower1Ecgi
	reports[1].CqiHist = make([]uint32, 1)
	reports[1].CqiHist[0] = makeCqi(ue.Tower1Strength)

	reports[2] = new(e2.RadioRepPerServCell)
	tower2Ecgi := toE2Ecgi(tower2.Ecgi)
	reports[2].Ecgi = &tower2Ecgi
	reports[2].CqiHist = make([]uint32, 1)
	reports[2].CqiHist[0] = makeCqi(ue.Tower2Strength)

	reports[3] = new(e2.RadioRepPerServCell)
	tower3Ecgi := toE2Ecgi(tower3.Ecgi)
	reports[3].Ecgi = &tower3Ecgi
	reports[3].CqiHist = make([]uint32, 1)
	reports[3].CqiHist[0] = makeCqi(ue.Tower3Strength)

	log.Infof("RadioMeasReport %s [cqi:%d] %d cqi:%d(%s),%d(%s),%d(%s)", servingTower.Ecgi.EcID, reports[0].CqiHist[0], ue.Imsi,
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
	}
}
