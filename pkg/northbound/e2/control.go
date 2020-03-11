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
	"github.com/onosproject/ran-simulator/pkg/utils"
	"io"

	"github.com/onosproject/ran-simulator/api/e2"
	"github.com/onosproject/ran-simulator/api/trafficsim"
	"github.com/onosproject/ran-simulator/api/types"
	"github.com/onosproject/ran-simulator/pkg/manager"
)

// SendControl ...
func (s *Server) SendControl(stream e2.InterfaceService_SendControlServer) error {
	c := make(chan e2.ControlUpdate)
	defer close(c)
	go recvControlLoop(s.GetPort(), s.GetEcID(), stream, c)
	return sendControlLoop(s.GetPort(), stream, c)
}

func sendControlLoop(port int, stream e2.InterfaceService_SendControlServer, c chan e2.ControlUpdate) error {
	for {
		select {
		case msg := <-c:
			if err := stream.Send(&msg); err != nil {
				log.Infof("send error %v", err)
				return err
			}
		case <-stream.Context().Done():
			log.Infof("Controller on Port %d has disconnected", port)
			return nil
		}
	}
}

func recvControlLoop(port int, towerID types.EcID, stream e2.InterfaceService_SendControlServer, c chan e2.ControlUpdate) {
	for {
		in, err := stream.Recv()
		if err == io.EOF || err != nil {
			log.Errorf("Unexpectedly ended when receiving Control responses on Port %d %s", port, err.Error())
			return
		}
		//log.Infof("Recv messageType %d", in.MessageType)
		switch x := in.S.(type) {
		case *e2.ControlResponse_CellConfigRequest:
			handleCellConfigRequest(port, towerID, c)
			go handleUeAdmissions(towerID, stream, c)
		case *e2.ControlResponse_HORequest:
			UpdateControlMetrics(in)
			err = handleHORequest(x.HORequest)
			if err != nil {
				log.Error(err)
			}
		case *e2.ControlResponse_RRMConfig:
			handleRRMConfig(x.RRMConfig)
		default:
			log.Errorf("ControlResponse has unexpected type %T", x)
		}
	}
}

func handleRRMConfig(req *e2.RRMConfig) {
	var powerAdjust float32
	switch req.PA[0] {
	case e2.XICICPA_XICIC_PA_DB_MINUS6:
		powerAdjust = -6
	case e2.XICICPA_XICIC_PA_DB_MINUX4DOT77:
		powerAdjust = -4.77
	case e2.XICICPA_XICIC_PA_DB_MINUS3:
		powerAdjust = -3
	case e2.XICICPA_XICIC_PA_DB_MINUS1DOT77:
		powerAdjust = -1.77
	case e2.XICICPA_XICIC_PA_DB_0:
		//Nothing to do
	case e2.XICICPA_XICIC_PA_DB_1:
		powerAdjust = 1
	case e2.XICICPA_XICIC_PA_DB_2:
		powerAdjust = 2
	case e2.XICICPA_XICIC_PA_DB_3:
		powerAdjust = 3
	}
	trafficSimMgr := manager.GetManager()
	err := trafficSimMgr.UpdateTower(types.EcID(req.Ecgi.Ecid), powerAdjust)
	if err != nil {
		log.Warn(err.Error())
	}
}

func handleHORequest(req *e2.HORequest) error {
	log.Infof("handleHORequest:  %s/%s -> %s", req.EcgiS.Ecid, req.Crnti, req.EcgiT.Ecid)
	m := manager.GetManager()
	ueName, err := m.CrntiToName(types.Crnti(req.Crnti), types.EcID(req.EcgiS.Ecid))
	if err != nil {
		log.Error(err)
		return fmt.Errorf("handleHORequest: ue %s/%s not found", req.EcgiS.Ecid, req.Crnti)
	}
	m.UeHandover(ueName, types.EcID(req.EcgiT.Ecid))
	return err
}

func handleCellConfigRequest(port int, ecID types.EcID, c chan e2.ControlUpdate) {
	log.Infof("handleCellConfigRequest on Port %d", port)

	trafficSimMgr := manager.GetManager()
	trafficSimMgr.TowersLock.RLock()
	defer trafficSimMgr.TowersLock.RUnlock()
	tower, ok := trafficSimMgr.Towers[ecID]
	if !ok {
		log.Warnf("Tower %s not found for handleCellConfigRequest on Port %d", ecID, port)
		return
	}
	cells := make([]*e2.CandScell, 0, 8)
	for _, neighbor := range tower.Neighbors {
		t := trafficSimMgr.Towers[neighbor]
		cell := e2.CandScell{
			Ecgi: &e2.ECGI{
				PlmnId: string(t.PlmnID),
				Ecid:   string(t.EcID),
			}}
		cells = append(cells, &cell)
	}
	cellConfigReport := e2.ControlUpdate{
		MessageType: e2.MessageType_CELL_CONFIG_REPORT,
		S: &e2.ControlUpdate_CellConfigReport{
			CellConfigReport: &e2.CellConfigReport{
				Ecgi: &e2.ECGI{
					PlmnId: string(tower.PlmnID),
					Ecid:   string(tower.EcID),
				},
				MaxNumConnectedUes: tower.MaxUEs,
				CandScells:         cells,
			},
		},
	}

	c <- cellConfigReport
	log.Infof("handleCellConfigReport eci: %s", tower.EcID)
}

func handleUeAdmissions(towerID types.EcID, stream e2.InterfaceService_SendControlServer, c chan e2.ControlUpdate) {
	trafficSimMgr := manager.GetManager()
	// Initiate UE admissions - handle what's currently here and listen for others
	for _, ue := range trafficSimMgr.UserEquipments {
		trafficSimMgr.UserEquipmentsLock.Lock()
		if ue.GetServingTower() != towerID {
			trafficSimMgr.UserEquipmentsLock.Unlock()
			continue
		}
		ueAdmReq := formatUeAdmissionReq(ue.ServingTower, ue.Crnti)
		c <- *ueAdmReq
		log.Infof("ueAdmissionRequest eci:%s crnti:%s", ue.ServingTower, ue.Crnti)
		ue.Admitted = true
		trafficSimMgr.UserEquipmentsLock.Unlock()
		trafficSimMgr.UeAdmitted(ue)
	}

	streamID := fmt.Sprintf("handleUeAdmissions-%p", stream)
	ueUpdatesLsnr, err := trafficSimMgr.Dispatcher.RegisterUeListener(streamID)
	if err != nil {
		log.Fatalf("could not register for UE events")
	}
	defer trafficSimMgr.Dispatcher.UnregisterUeListener(streamID)
	for {
		// block here and listen for updates on UEs
		select {
		case event := <-ueUpdatesLsnr:
			ue, ok := event.Object.(*types.Ue)
			if !ok {
				log.Fatalf("Object %v could not be converted to UE", ue)
			}
			if ue.ServingTower != towerID {
				continue // listen for the next event
			}
			if event.Type == trafficsim.Type_ADDED {
				ueAdmReq := formatUeAdmissionReq(ue.ServingTower, ue.Crnti)
				c <- *ueAdmReq
				log.Infof("ueAdmissionRequest eci:%s crnti:%s", ue.ServingTower, ue.Crnti)
				ue.Admitted = true
			} else if event.Type == trafficsim.Type_REMOVED {
				err = trafficSimMgr.DelCrnti(ue.ServingTower, ue.Crnti)
				if err != nil {
					log.Error(err.Error())
					continue
				}
				ue.Crnti = manager.InvalidCrnti
				ueRelInd := formatUeReleaseInd(ue.ServingTower, ue.Crnti)
				c <- *ueRelInd
				log.Infof("ueReleaseInd eci:%s crnti:%s", ue.ServingTower, ue.Crnti)
			}
			// Nothing to be done for trafficsim.Type_UPDATED - they are handled by Telemetry
		case <-stream.Context().Done():
			log.Infof("Controller has disconnected")
			return
		}
	}
}

func formatUeAdmissionReq(eci types.EcID, crnti types.Crnti) *e2.ControlUpdate {
	return &e2.ControlUpdate{
		MessageType: e2.MessageType_UE_ADMISSION_REQUEST,
		S: &e2.ControlUpdate_UEAdmissionRequest{
			UEAdmissionRequest: &e2.UEAdmissionRequest{
				Ecgi: &e2.ECGI{
					PlmnId: utils.TestPlmnID,
					Ecid:   string(eci),
				},
				Crnti:             string(crnti),
				AdmissionEstCause: e2.AdmEstCause_MO_SIGNALLING,
			},
		},
	}
}

func formatUeReleaseInd(eci types.EcID, crnti types.Crnti) *e2.ControlUpdate {
	return &e2.ControlUpdate{
		MessageType: e2.MessageType_UE_RELEASE_IND,
		S: &e2.ControlUpdate_UEReleaseInd{
			UEReleaseInd: &e2.UEReleaseInd{
				Ecgi: &e2.ECGI{
					PlmnId: utils.TestPlmnID,
					Ecid:   string(eci),
				},
				Crnti:        string(crnti),
				ReleaseCause: e2.ReleaseCause_RELEASE_INACTIVITY,
			},
		},
	}
}
