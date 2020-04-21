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

	e2 "github.com/onosproject/onos-ric/api/sb"
	"github.com/onosproject/onos-ric/api/sb/e2ap"
	"github.com/onosproject/onos-ric/api/sb/e2sm"
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

func radioMeasReportPerUE(port int, towerID types.ECGI, stream e2ap.E2AP_RicChanServer, c chan e2ap.RicIndication) error {
	trafficSimMgr := manager.GetManager()

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
				if ue.ServingTower.EcID != towerID.EcID || ue.ServingTower.PlmnID != towerID.PlmnID {
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

func generateReport(ue *types.Ue) e2ap.RicIndication {
	trafficSimMgr := manager.GetManager()

	trafficSimMgr.TowersLock.RLock()
	defer trafficSimMgr.TowersLock.RUnlock()

	servingTower := trafficSimMgr.Towers[*ue.ServingTower]
	tower1 := trafficSimMgr.Towers[*ue.Tower1]
	tower2 := trafficSimMgr.Towers[*ue.Tower2]
	tower3 := trafficSimMgr.Towers[*ue.Tower3]

	reports := make([]*e2.RadioRepPerServCell, 3)

	reports[0] = new(e2.RadioRepPerServCell)
	tower1Ecgi := toE2Ecgi(tower1.Ecgi)
	reports[0].Ecgi = &tower1Ecgi
	reports[0].CqiHist = make([]uint32, 1)
	reports[0].CqiHist[0] = makeCqi(ue.Tower1Dist, tower1.GetTxPowerdB())

	reports[1] = new(e2.RadioRepPerServCell)
	tower2Ecgi := toE2Ecgi(tower2.Ecgi)
	reports[1].Ecgi = &tower2Ecgi
	reports[1].CqiHist = make([]uint32, 1)
	reports[1].CqiHist[0] = makeCqi(ue.Tower2Dist, tower2.GetTxPowerdB())

	reports[2] = new(e2.RadioRepPerServCell)
	tower3Ecgi := toE2Ecgi(tower2.Ecgi)
	reports[2].Ecgi = &tower3Ecgi
	reports[2].CqiHist = make([]uint32, 1)
	reports[2].CqiHist[0] = makeCqi(ue.Tower3Dist, tower3.GetTxPowerdB())

	log.Infof("RadioMeasReport %s %d cqi:%d(%s),%d(%s),%d(%s)", servingTower.Ecgi.EcID, ue.Imsi,
		reports[0].CqiHist[0], reports[0].Ecgi.Ecid,
		reports[1].CqiHist[0], reports[1].Ecgi.Ecid,
		reports[2].CqiHist[0], reports[2].Ecgi.Ecid)

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
