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
	"github.com/onosproject/ran-simulator/api/e2"
	"github.com/onosproject/ran-simulator/api/trafficsim"
	"github.com/onosproject/ran-simulator/api/types"
	"github.com/onosproject/ran-simulator/pkg/manager"
)

const e2TelemetryNbi = "e2TelemetryNbi"

func makeCqi(distance float32, txPowerdB float32) uint32 {
	cqi := uint32(0.0001 * txPowerdB / (distance * distance))
	if cqi > 15 {
		cqi = 15
	}
	return cqi
}

// SendTelemetry ...
func (s *Server) SendTelemetry(req *e2.L2MeasConfig, stream e2.InterfaceService_SendTelemetryServer) error {
	c := make(chan e2.TelemetryMessage)
	defer close(c)
	go func() {
		err := radioMeasReportPerUE(s.GetPort(), s.GetEcID(), stream, c)
		if err != nil {
			log.Errorf("Unable to send radioMeasReportPerUE on Port %d %s", s.GetPort(), err.Error())
		}
	}()
	return sendTelemetryLoop(s.GetPort(), stream, c)
}

func sendTelemetryLoop(port int, stream e2.InterfaceService_SendTelemetryServer, c chan e2.TelemetryMessage) error {
	for {
		select {
		case msg := <-c:
			UpdateTelemetryMetrics(&msg)
			if err := stream.Send(&msg); err != nil {
				log.Infof("send error on Port %d %v", port, err)
				return err
			}
		case <-stream.Context().Done():
			log.Infof("Controller has disconnected on Port %d", port)
			return nil
		}
	}
}

func radioMeasReportPerUE(port int, towerID types.EcID, stream e2.InterfaceService_SendTelemetryServer, c chan e2.TelemetryMessage) error {
	trafficSimMgr := manager.GetManager()

	// replay any existing UE's
	for _, ue := range trafficSimMgr.UserEquipments {
		if ue.ServingTower != towerID {
			continue
		}
		if ue.Admitted {
			c <- generateReport(ue)
		}
	}

	streamID := fmt.Sprintf("%s-%p", e2TelemetryNbi, stream)
	ueChangeChannel, err := trafficSimMgr.Dispatcher.RegisterUeListener(streamID)
	defer trafficSimMgr.Dispatcher.UnregisterUeListener(streamID)
	if err != nil {
		return err
	}
	log.Infof("Listening for changes on UEs with ServingTower=%s for Port %d", towerID, port)
	for {
		select {
		// block here and listen out for any updates to UEs
		case ueUpdate := <-ueChangeChannel:
			if ueUpdate.Type == trafficsim.Type_UPDATED && ueUpdate.UpdateType == trafficsim.UpdateType_TOWER {
				ue, ok := ueUpdate.Object.(*types.Ue)
				if !ok {
					log.Fatalf("Object %v could not be converted to UE", ueUpdate)
				}
				if ue.ServingTower != towerID {
					continue
				}
				c <- generateReport(ue)
			}
		case <-stream.Context().Done():
			log.Infof("Controller has disconnected on Port %d", port)
			return nil
		}
	}
}

func generateReport(ue *types.Ue) e2.TelemetryMessage {
	trafficSimMgr := manager.GetManager()

	trafficSimMgr.TowersLock.RLock()
	defer trafficSimMgr.TowersLock.RUnlock()

	servingTower := trafficSimMgr.Towers[ue.ServingTower]
	tower1 := trafficSimMgr.Towers[ue.Tower1]
	tower2 := trafficSimMgr.Towers[ue.Tower2]
	tower3 := trafficSimMgr.Towers[ue.Tower3]

	reports := make([]*e2.RadioRepPerServCell, 3)

	reports[0] = new(e2.RadioRepPerServCell)
	reports[0].Ecgi = &e2.ECGI{
		PlmnId: string(tower1.PlmnID),
		Ecid:   string(tower1.EcID),
	}
	reports[0].CqiHist = make([]uint32, 1)
	reports[0].CqiHist[0] = makeCqi(ue.Tower1Dist, tower1.GetTxPowerdB())

	reports[1] = new(e2.RadioRepPerServCell)
	reports[1].Ecgi = &e2.ECGI{
		PlmnId: string(tower2.PlmnID),
		Ecid:   string(tower2.EcID),
	}
	reports[1].CqiHist = make([]uint32, 1)
	reports[1].CqiHist[0] = makeCqi(ue.Tower2Dist, tower2.GetTxPowerdB())

	reports[2] = new(e2.RadioRepPerServCell)
	reports[2].Ecgi = &e2.ECGI{
		PlmnId: string(tower3.PlmnID),
		Ecid:   string(tower3.EcID),
	}
	reports[2].CqiHist = make([]uint32, 1)
	reports[2].CqiHist[0] = makeCqi(ue.Tower3Dist, tower3.GetTxPowerdB())

	log.Infof("RadioMeasReport %s %s cqi:%d(%s),%d(%s),%d(%s)", servingTower.EcID, ue.Name,
		reports[0].CqiHist[0], reports[0].Ecgi.Ecid,
		reports[1].CqiHist[0], reports[1].Ecgi.Ecid,
		reports[2].CqiHist[0], reports[2].Ecgi.Ecid)

	return e2.TelemetryMessage{
		MessageType: e2.MessageType_RADIO_MEAS_REPORT_PER_UE,
		S: &e2.TelemetryMessage_RadioMeasReportPerUE{
			RadioMeasReportPerUE: &e2.RadioMeasReportPerUE{
				Ecgi: &e2.ECGI{
					PlmnId: string(servingTower.PlmnID),
					Ecid:   string(servingTower.EcID),
				},
				Crnti:                string(ue.Crnti),
				RadioReportServCells: reports,
			},
		},
	}
}
