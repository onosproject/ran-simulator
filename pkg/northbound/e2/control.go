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
	"io"
	"strconv"

	"github.com/onosproject/ran-simulator/api/e2"
	"github.com/onosproject/ran-simulator/api/trafficsim"
	"github.com/onosproject/ran-simulator/api/types"
	"github.com/onosproject/ran-simulator/pkg/manager"
)

func crntiToName(crnti string) string {
	return "Ue-" + crnti
}

// TODO return the error OR refactor altogether
func eciToName(eci string) string {
	id, _ := strconv.Atoi(eci)
	return fmt.Sprintf("Tower-%d", id)
}

// SendControl ...
func (s *Server) SendControl(stream e2.InterfaceService_SendControlServer) error {
	c := make(chan e2.ControlUpdate)
	defer close(c)
	go recvControlLoop(stream, c)
	return sendControlLoop(stream, c)
}

func sendControlLoop(stream e2.InterfaceService_SendControlServer, c chan e2.ControlUpdate) error {
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

func recvControlLoop(stream e2.InterfaceService_SendControlServer, c chan e2.ControlUpdate) {
	for {
		in, err := stream.Recv()
		if err == io.EOF || err != nil {
			log.Errorf("Unexpectedly ended when receiving Control responses %s", err.Error())
			return
		}
		//log.Infof("Recv messageType %d", in.MessageType)
		switch x := in.S.(type) {
		case *e2.ControlResponse_CellConfigRequest:
			go func() {
				err := handleCellConfigRequest(c)
				if err != nil {
					log.Errorf("error listening for UE updates %s", err.Error())
				}
			}()
		case *e2.ControlResponse_HORequest:
			handleHORequest(x.HORequest)
		case *e2.ControlResponse_RRMConfig:
			handleRRMConfig(x.RRMConfig)
		default:
			log.Errorf("ControlResponse has unexpected type %T", x)
		}
		UpdateControlMetrics(in)
	}
}

func handleRRMConfig(req *e2.RRMConfig) {
	trafficSimMgr := manager.GetManager()
	tower := trafficSimMgr.GetTower(eciToName(req.Ecgi.Ecid))
	switch req.PA[0] {
	case e2.XICICPA_XICIC_PA_DB_MINUS6:
		tower.TxPowerdB -= 6
	case e2.XICICPA_XICIC_PA_DB_MINUX4DOT77:
		tower.TxPowerdB -= 4.77
	case e2.XICICPA_XICIC_PA_DB_MINUS3:
		tower.TxPowerdB -= 3
	case e2.XICICPA_XICIC_PA_DB_MINUS1DOT77:
		tower.TxPowerdB -= 1.77
	case e2.XICICPA_XICIC_PA_DB_0:
		tower.TxPowerdB -= 0
	case e2.XICICPA_XICIC_PA_DB_1:
		tower.TxPowerdB++
	case e2.XICICPA_XICIC_PA_DB_2:
		tower.TxPowerdB += 2
	case e2.XICICPA_XICIC_PA_DB_3:
		tower.TxPowerdB += 3
	}
	trafficSimMgr.UpdateTower(tower)
}

func handleHORequest(req *e2.HORequest) {
	//log.Infof("handleHORequest crnti:%s, name:%s serving:%s, target:%s", req.Crnti, crntiToName(req.Crnti), req.EcgiS.Ecid, req.EcgiT.Ecid)

	trafficSimMgr := manager.GetManager()

	log.Infof("hand-over %s from %s to %s", crntiToName(req.Crnti), eciToName(req.EcgiS.Ecid), eciToName(req.EcgiT.Ecid))
	trafficSimMgr.UeHandover(crntiToName(req.Crnti), eciToName(req.EcgiT.Ecid))
}

func handleCellConfigRequest(c chan e2.ControlUpdate) error {
	log.Infof("handleCellConfigRequest")

	trafficSimMgr := manager.GetManager()

	for _, tower := range trafficSimMgr.Towers {
		cells := make([]*e2.CandScell, 0, 8)
		for _, neighbor := range tower.Neighbors {
			t := trafficSimMgr.Towers[neighbor]
			cell := e2.CandScell{
				Ecgi: &e2.ECGI{
					PlmnId: t.PlmnID,
					Ecid:   t.EcID,
				}}
			cells = append(cells, &cell)
		}
		cellConfigReport := e2.ControlUpdate{
			MessageType: e2.MessageType_CELL_CONFIG_REPORT,
			S: &e2.ControlUpdate_CellConfigReport{
				CellConfigReport: &e2.CellConfigReport{
					Ecgi: &e2.ECGI{
						PlmnId: tower.PlmnID,
						Ecid:   tower.EcID,
					},
					MaxNumConnectedUes: tower.MaxUEs,
					CandScells:         cells,
				},
			},
		}

		c <- cellConfigReport
		log.Infof("handleCellConfigReport eci: %s", tower.EcID)
	}

	// Initate UE admissions
	ueUpdatesLsnr, err := trafficSimMgr.Dispatcher.RegisterUeListener("handleCellConfigReport")
	if err != nil {
		log.Errorf("could not register for UE events")
		return err
	}
	for event := range ueUpdatesLsnr {
		ue, ok := event.Object.(*types.Ue)
		if !ok {
			return fmt.Errorf("invalid event type for %v", event.Object)
		}
		if event.Type == trafficsim.Type_ADDED || event.Type == trafficsim.Type_NONE {
			eci := trafficSimMgr.GetTowerByName(ue.ServingTower).EcID
			ueAdmReq := e2.ControlUpdate{
				MessageType: e2.MessageType_UE_ADMISSION_REQUEST,
				S: &e2.ControlUpdate_UEAdmissionRequest{
					UEAdmissionRequest: &e2.UEAdmissionRequest{
						Ecgi: &e2.ECGI{
							PlmnId: manager.TestPlmnID,
							Ecid:   eci,
						},
						Crnti:             ue.Crnti,
						AdmissionEstCause: e2.AdmEstCause_MO_SIGNALLING,
					},
				},
			}
			c <- ueAdmReq
			log.Infof("ueAdmissionRequest eci:%s crnti:%s", eci, ue.Crnti)
			trafficSimMgr.UeAdmitted(ue)
		} else if event.Type == trafficsim.Type_REMOVED {
			log.Warnf("removal of UE %s - not yet supported", ue.GetName())
		}
	}
	return nil
}
