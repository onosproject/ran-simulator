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

	e2 "github.com/onosproject/onos-ric/api/sb"
	"github.com/onosproject/onos-ric/api/sb/e2ap"
	"github.com/onosproject/onos-ric/api/sb/e2sm"
	"github.com/onosproject/ran-simulator/api/trafficsim"
	"github.com/onosproject/ran-simulator/api/types"
	"github.com/onosproject/ran-simulator/pkg/manager"
)

// RicControl ...
func (s *Server) RicControl(stream e2ap.E2AP_RicControlServer) error {
	go ricControlRequest(s.GetPort(), s.GetECGI(), stream)
	c := make(chan e2ap.RicControlResponse)
	defer close(c)
	return ricControlResponse(s.GetPort(), stream, c)
}

// SendControl ...
func (s *Server) SendControl(stream e2ap.E2AP_SendControlServer) error {
	c := make(chan e2.ControlUpdate)
	defer close(c)
	go recvControlLoop(s.GetPort(), s.GetECGI(), stream, c)
	return sendControlLoop(s.GetPort(), stream, c)
}

func ricControlResponse(port int, stream e2ap.E2AP_RicControlServer, c chan e2ap.RicControlResponse) error {
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

func sendControlLoop(port int, stream e2ap.E2AP_SendControlServer, c chan e2.ControlUpdate) error {
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

func ricControlRequest(port int, towerID types.ECGI, stream e2ap.E2AP_RicControlServer) {
	for {
		in, err := stream.Recv()
		if err == io.EOF || err != nil {
			log.Errorf("Unexpectedly ended when receiving Control responses on Port %d %s", port, err.Error())
			return
		}
		//log.Infof("Recv messageType %d", in.MessageType)
		switch in.Hdr.MessageType {
		case e2.MessageType_HO_REQUEST:
			if x, ok := in.Msg.S.(*e2sm.RicControlMessage_HORequest); ok {
				err = handleHORequest(towerID, x.HORequest)
				if err != nil {
					log.Error(err)
				}
			} else {
				log.Fatalf("Unexpected payload in MessageType_HO_REQUEST %v", in)
			}
		case e2.MessageType_RRM_CONFIG:
			if x, ok := in.Msg.S.(*e2sm.RicControlMessage_RRMConfig); ok {
				handleRRMConfig(x.RRMConfig)
			} else {
				log.Fatalf("Unexpected payload in MessageType_RRM_CONFIG %v", in)
			}
		default:
			log.Errorf("ControlResponse has unexpected type %T", in.Hdr.MessageType)
		}
	}
}

func recvControlLoop(port int, towerID types.ECGI, stream e2ap.E2AP_SendControlServer, c chan e2.ControlUpdate) {
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
	err := trafficSimMgr.UpdateTower(toTypesEcgi(req.Ecgi), powerAdjust)
	if err != nil {
		log.Warn(err.Error())
	}
}

func handleHORequest(towerID types.ECGI, req *e2.HORequest) error {
	sourceEcgi := toTypesEcgi(req.EcgiS)
	targetEcgi := toTypesEcgi(req.EcgiT)

	if towerID.EcID == sourceEcgi.EcID && towerID.PlmnID == sourceEcgi.PlmnID {
		log.Infof("Source handleHORequest:  %s/%s -> %s", req.EcgiS.Ecid, req.Crnti, req.EcgiT.Ecid)
		m := manager.GetManager()
		imsi, err := m.CrntiToName(types.Crnti(req.Crnti), &towerID)
		if err != nil {
			log.Error(err)
			return fmt.Errorf("handleHORequest: ue %s/%s not found", req.EcgiS.Ecid, req.Crnti)
		}
		UpdateControlMetrics(imsi)
		return m.UeHandover(imsi, &targetEcgi)
	} else if towerID.EcID == targetEcgi.EcID && towerID.PlmnID == targetEcgi.PlmnID {
		log.Infof("Target handleHORequest:  %s/%s -> %s", req.EcgiS.Ecid, req.Crnti, req.EcgiT.Ecid)
		return nil
	}
	return fmt.Errorf("unexpected handleHORequest on tower: %s %s/%s -> %s", towerID, req.EcgiS.Ecid, req.Crnti, req.EcgiT.Ecid)
}

func handleCellConfigRequest(port int, ecgi types.ECGI, c chan e2.ControlUpdate) {
	log.Infof("handleCellConfigRequest on Port %d", port)

	trafficSimMgr := manager.GetManager()
	trafficSimMgr.TowersLock.RLock()
	defer trafficSimMgr.TowersLock.RUnlock()
	tower, ok := trafficSimMgr.Towers[ecgi]
	if !ok {
		log.Warnf("Tower %s not found for handleCellConfigRequest on Port %d", ecgi, port)
		return
	}
	cells := make([]*e2.CandScell, 0, 8)
	for _, neighbor := range tower.Neighbors {
		t := trafficSimMgr.Towers[*neighbor]
		e2Ecgi := toE2Ecgi(t.Ecgi)
		cell := e2.CandScell{
			Ecgi: &e2Ecgi,
		}
		cells = append(cells, &cell)
	}
	e2Ecgi := toE2Ecgi(tower.Ecgi)
	cellConfigReport := e2.ControlUpdate{
		MessageType: e2.MessageType_CELL_CONFIG_REPORT,
		S: &e2.ControlUpdate_CellConfigReport{
			CellConfigReport: &e2.CellConfigReport{
				Ecgi:               &e2Ecgi,
				MaxNumConnectedUes: tower.MaxUEs,
				CandScells:         cells,
			},
		},
	}

	c <- cellConfigReport
	log.Infof("handleCellConfigReport eci: %v", tower.Ecgi)
}

func handleUeAdmissions(towerID types.ECGI, stream e2ap.E2AP_SendControlServer, c chan e2.ControlUpdate) {
	trafficSimMgr := manager.GetManager()
	// Initiate UE admissions - handle what's currently here and listen for others
	for _, ue := range trafficSimMgr.UserEquipments {
		trafficSimMgr.UserEquipmentsLock.Lock()
		if ue.GetServingTower().EcID != towerID.EcID || ue.GetServingTower().PlmnID != towerID.PlmnID {
			trafficSimMgr.UserEquipmentsLock.Unlock()
			continue
		}
		ueAdmReq := formatUeAdmissionReq(ue.ServingTower, ue.Crnti, ue.Imsi)
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
			if ue.ServingTower.EcID != towerID.EcID || ue.ServingTower.PlmnID != towerID.PlmnID {
				continue // listen for the next event
			}
			if event.Type == trafficsim.Type_ADDED ||
				event.Type == trafficsim.Type_UPDATED && event.UpdateType == trafficsim.UpdateType_HANDOVER {
				ueAdmReq := formatUeAdmissionReq(ue.ServingTower, ue.Crnti, ue.Imsi)
				c <- *ueAdmReq
				log.Infof("ueAdmissionRequest eci:%s crnti:%s", ue.ServingTower, ue.Crnti)
				ue.Admitted = true
			} else if event.Type == trafficsim.Type_REMOVED {
				err = trafficSimMgr.DelCrnti(ue.ServingTower, ue.Crnti)
				if err != nil {
					log.Error(err.Error())
					continue
				}
				ueRelInd := formatUeReleaseInd(ue.ServingTower, ue.Crnti)
				c <- *ueRelInd
				log.Infof("ueReleaseInd eci:%s crnti:%s", ue.ServingTower, ue.Crnti)
				ue.Crnti = manager.InvalidCrnti
			}
			// Nothing to be done for trafficsim.Type_UPDATED - they are handled by Telemetry
		case <-stream.Context().Done():
			log.Infof("Controller has disconnected")
			return
		}
	}
}

func formatUeAdmissionReq(eci *types.ECGI, crnti types.Crnti, imsi types.Imsi) *e2.ControlUpdate {
	e2Ecgi := toE2Ecgi(eci)
	return &e2.ControlUpdate{
		MessageType: e2.MessageType_UE_ADMISSION_REQUEST,
		S: &e2.ControlUpdate_UEAdmissionRequest{
			UEAdmissionRequest: &e2.UEAdmissionRequest{
				Ecgi:              &e2Ecgi,
				Crnti:             string(crnti),
				AdmissionEstCause: e2.AdmEstCause_MO_SIGNALLING,
				Imsi:              uint64(imsi),
			},
		},
	}
}

func formatUeReleaseInd(eci *types.ECGI, crnti types.Crnti) *e2.ControlUpdate {
	e2Ecgi := toE2Ecgi(eci)
	return &e2.ControlUpdate{
		MessageType: e2.MessageType_UE_RELEASE_IND,
		S: &e2.ControlUpdate_UEReleaseInd{
			UEReleaseInd: &e2.UEReleaseInd{
				Ecgi:         &e2Ecgi,
				Crnti:        string(crnti),
				ReleaseCause: e2.ReleaseCause_RELEASE_INACTIVITY,
			},
		},
	}
}
