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
		err := radioMeasReportPerUE(stream, c)
		if err != nil {
			log.Errorf("Unable to send radioMeasReportPerUE %s", err.Error())
		}
	}()
	return sendTelemetryLoop(stream, c)
}

func sendTelemetryLoop(stream e2.InterfaceService_SendTelemetryServer, c chan e2.TelemetryMessage) error {
	for {
		select {
		case msg := <-c:
			if err := stream.Send(&msg); err != nil {
				log.Infof("send error %v", err)
				return err
			}
		case <-stream.Context().Done():
			log.Infof("Controller has disconnected")
			return nil
		}
	}
}

func radioMeasReportPerUE(stream e2.InterfaceService_SendTelemetryServer, c chan e2.TelemetryMessage) error {
	trafficSimMgr := manager.GetManager()

	// replay any existing UE's
	for _, ue := range trafficSimMgr.UserEquipments {
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
	// block here and listen out for any updates to UEs
	for {
		select {
		case ueUpdate := <-ueChangeChannel:
			if ueUpdate.Type == trafficsim.Type_UPDATED && ueUpdate.UpdateType == trafficsim.UpdateType_TOWER {
				ue, ok := ueUpdate.Object.(*types.Ue)
				if !ok {
					log.Fatalf("Object %v could not be converted to UE", ueUpdate)
				}
				c <- generateReport(ue)
			}
		case <-stream.Context().Done():
			log.Infof("Controller has disconnected")
			return nil
		}
	}
}

func generateReport(ue *types.Ue) e2.TelemetryMessage {
	trafficSimMgr := manager.GetManager()

	servingTower := trafficSimMgr.GetTowerByName(ue.ServingTower)
	tower1 := trafficSimMgr.GetTowerByName(ue.Tower1)
	tower2 := trafficSimMgr.GetTowerByName(ue.Tower2)
	tower3 := trafficSimMgr.GetTowerByName(ue.Tower3)

	reports := make([]*e2.RadioRepPerServCell, 3)

	trafficSimMgr.TowersLock.RLock()
	reports[0] = new(e2.RadioRepPerServCell)
	reports[0].Ecgi = &e2.ECGI{
		PlmnId: tower1.PlmnID,
		Ecid:   tower1.EcID,
	}
	reports[0].CqiHist = make([]uint32, 1)
	reports[0].CqiHist[0] = makeCqi(ue.Tower1Dist, tower1.GetTxPowerdB())

	reports[1] = new(e2.RadioRepPerServCell)
	reports[1].Ecgi = &e2.ECGI{
		PlmnId: tower2.PlmnID,
		Ecid:   tower2.EcID,
	}
	reports[1].CqiHist = make([]uint32, 1)
	reports[1].CqiHist[0] = makeCqi(ue.Tower2Dist, tower2.GetTxPowerdB())

	reports[2] = new(e2.RadioRepPerServCell)
	reports[2].Ecgi = &e2.ECGI{
		PlmnId: tower3.PlmnID,
		Ecid:   tower3.EcID,
	}
	reports[2].CqiHist = make([]uint32, 1)
	reports[2].CqiHist[0] = makeCqi(ue.Tower3Dist, tower3.GetTxPowerdB())
	trafficSimMgr.TowersLock.RUnlock()

	log.Infof("RadioMeasReport %s %s cqi:%d(%s),%d(%s),%d(%s)", servingTower.EcID, ue.Name,
		reports[0].CqiHist[0], reports[0].Ecgi.Ecid,
		reports[1].CqiHist[0], reports[1].Ecgi.Ecid,
		reports[2].CqiHist[0], reports[2].Ecgi.Ecid)

	return e2.TelemetryMessage{
		MessageType: e2.MessageType_RADIO_MEAS_REPORT_PER_UE,
		S: &e2.TelemetryMessage_RadioMeasReportPerUE{
			RadioMeasReportPerUE: &e2.RadioMeasReportPerUE{
				Ecgi: &e2.ECGI{
					PlmnId: servingTower.PlmnID,
					Ecid:   servingTower.EcID,
				},
				Crnti:                ue.Crnti,
				RadioReportServCells: reports,
			},
		},
	}
}
