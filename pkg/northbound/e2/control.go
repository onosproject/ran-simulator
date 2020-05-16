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

// RicChan ...
func (s *Server) RicChan(stream e2ap.E2AP_RicChanServer) error {
	s.indChan = make(chan e2ap.RicIndication)
	s.stream = stream
	defer close(s.indChan)
	go s.ricControlRequest()
	go s.recvControlLoop()
	go func() {
		err := s.radioMeasReportPerUE()
		if err != nil {
			log.Errorf("Unable to send radioMeasReportPerUE on Port %d %s", s.GetPort(), err.Error())
		}
	}()
	return s.ricControlResponse()
}

func (s *Server) ricControlResponse() error {
	for {
		select {
		case msg := <-s.indChan:
			go UpdateTelemetryMetrics(&msg)
			if err := s.stream.Send(&msg); err != nil {
				log.Infof("send error %v", err)
				return err
			}
		case <-s.stream.Context().Done():
			log.Infof("Controller on Port %d has disconnected", s.GetPort())
			return nil
		}
	}
}

func (s *Server) ricControlRequest() {
	for {
		in, err := s.stream.Recv()
		if err == io.EOF || err != nil {
			log.Errorf("Unexpectedly ended when receiving Control responses on Port %d %s", s.GetPort(), err.Error())
			return
		}
		//log.Infof("Recv messageType %d", in.GetHdr().GetMessageType())
		if in == nil || in.Hdr == nil || in.Msg == nil {
			log.Errorf("Unexpected empty Control request message on Port %d %v", s.GetPort(), in)
			return
		}
		switch in.Hdr.MessageType {
		case e2.MessageType_CELL_CONFIG_REQUEST:
			err = s.handleCellConfigRequest()
			if err != nil {
				log.Error(err)
			}
		case e2.MessageType_HO_REQUEST:
			if x, ok := in.Msg.S.(*e2sm.RicControlMessage_HORequest); ok {
				err = s.handleHORequest(x.HORequest)
				if err != nil {
					log.Error(err)
				}
			} else {
				log.Fatalf("Unexpected payload in MessageType_HO_REQUEST %v", in)
			}
		case e2.MessageType_RRM_CONFIG:
			if x, ok := in.Msg.S.(*e2sm.RicControlMessage_RRMConfig); ok {
				s.handleRRMConfig(x.RRMConfig)
			} else {
				log.Fatalf("Unexpected payload in MessageType_RRM_CONFIG %v", in)
			}
		case e2.MessageType_L2_MEAS_CONFIG:
			if x, ok := in.Msg.S.(*e2sm.RicControlMessage_L2MeasConfig); ok {
				s.handleL2MeasConfig(x.L2MeasConfig)
			} else {
				log.Fatalf("Unexpected payload in MessageType_RRM_CONFIG %v", in)
			}
		default:
			log.Errorf("ControlRequest has unexpected type %d", in.Hdr.MessageType)
		}
	}
}

func (s *Server) recvControlLoop() {
	s.handleUeAdmissions()
}

func (s *Server) handleL2MeasConfig(req *e2.L2MeasConfig) {
	log.Infof("handleL2MeasConfig radioMeasReportPerUe=%d", req.RadioMeasReportPerUe)
	s.l2MeasConfig = *req
}

func (s *Server) handleRRMConfig(req *e2.RRMConfig) {
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
	err := trafficSimMgr.UpdateCell(toTypesEcgi(req.Ecgi), powerAdjust)
	if err != nil {
		log.Warn(err.Error())
	}
}

func (s *Server) handleHORequest(req *e2.HORequest) error {
	sourceEcgi := toTypesEcgi(req.EcgiS)
	targetEcgi := toTypesEcgi(req.EcgiT)

	for _, crnti := range req.Crntis {
		if s.GetECGI().EcID == sourceEcgi.EcID && s.GetECGI().PlmnID == sourceEcgi.PlmnID {
			log.Infof("Source handleHORequest:  %s/%s -> %s", req.EcgiS.Ecid, crnti, req.EcgiT.Ecid)
			m := manager.GetManager()
			imsi, err := m.CrntiToName(types.Crnti(crnti), s.GetECGI())
			if err != nil {
				log.Error(err)
				continue
			}
			UpdateControlMetrics(imsi)
			err = m.UeHandover(imsi, &targetEcgi)
			if err != nil {
				log.Error(err)
				continue
			}
		} else if s.GetECGI().EcID == targetEcgi.EcID && s.GetECGI().PlmnID == targetEcgi.PlmnID {
			log.Infof("Target handleHORequest:  %s/%s -> %s", req.EcgiS.Ecid, crnti, req.EcgiT.Ecid)
		}
		log.Errorf("unexpected handleHORequest on tower: %s %s/%s -> %s", s.GetECGI(), req.EcgiS.Ecid, crnti, req.EcgiT.Ecid)
	}
	return nil
}

func (s *Server) handleCellConfigRequest() error {
	log.Infof("handleCellConfigRequest on Port %d", s.GetPort())

	trafficSimMgr := manager.GetManager()
	trafficSimMgr.CellsLock.RLock()
	cell, ok := trafficSimMgr.Cells[s.GetECGI()]
	if !ok {
		log.Warnf("Tower %s not found for handleCellConfigRequest on Port %d", s.GetECGI(), s.GetPort())
		trafficSimMgr.CellsLock.RUnlock()
		return nil
	}
	nCells := make([]*e2.CandScell, 0, 8)
	for _, neighbor := range cell.Neighbors {
		nc := trafficSimMgr.Cells[*neighbor]
		ncEcgi := toE2Ecgi(nc.Ecgi)
		nCell := e2.CandScell{
			Ecgi: &ncEcgi,
		}
		nCells = append(nCells, &nCell)
	}
	trafficSimMgr.CellsLock.RUnlock()
	e2Ecgi := toE2Ecgi(cell.Ecgi)
	cellConfigReport := e2ap.RicIndication{
		Hdr: &e2sm.RicIndicationHeader{
			MessageType: e2.MessageType_CELL_CONFIG_REPORT,
		},
		Msg: &e2sm.RicIndicationMessage{
			S: &e2sm.RicIndicationMessage_CellConfigReport{
				CellConfigReport: &e2.CellConfigReport{
					Ecgi:               &e2Ecgi,
					MaxNumConnectedUes: cell.MaxUEs,
					CandScells:         nCells,
				},
			},
		},
	}

	s.indChan <- cellConfigReport
	log.Infof("handleCellConfigReport eci: %s. CCR %v", cell.GetEcgi().String(), cellConfigReport)

	return nil
}

func (s *Server) handleUeAdmissions() {
	trafficSimMgr := manager.GetManager()
	// Initiate UE admissions - handle what's currently here and listen for others
	for _, ue := range trafficSimMgr.UserEquipments {
		trafficSimMgr.UserEquipmentsLock.Lock()
		if ue.GetServingTower().EcID != s.GetECGI().EcID || ue.GetServingTower().PlmnID != s.GetECGI().PlmnID {
			trafficSimMgr.UserEquipmentsLock.Unlock()
			continue
		}
		ueAdmReq := formatUeAdmissionReq(ue.ServingTower, ue.Crnti, ue.Imsi)
		s.indChan <- *ueAdmReq
		log.Infof("ueAdmissionRequest eci:%s crnti:%s", ue.ServingTower, ue.Crnti)
		ue.Admitted = true
		trafficSimMgr.UserEquipmentsLock.Unlock()
		trafficSimMgr.UeAdmitted(ue)
	}

	streamID := fmt.Sprintf("handleUeAdmissions-%p", s.stream)
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
			if ue.ServingTower.EcID != s.GetECGI().EcID || ue.ServingTower.PlmnID != s.GetECGI().PlmnID {
				continue // listen for the next event
			}
			if event.Type == trafficsim.Type_ADDED ||
				event.Type == trafficsim.Type_UPDATED && event.UpdateType == trafficsim.UpdateType_HANDOVER {
				ueAdmReq := formatUeAdmissionReq(ue.ServingTower, ue.Crnti, ue.Imsi)
				s.indChan <- *ueAdmReq
				log.Infof("ueAdmissionRequest eci:%s crnti:%s", ue.ServingTower, ue.Crnti)
				ue.Admitted = true
			} else if event.Type == trafficsim.Type_REMOVED {
				err = trafficSimMgr.DelCrnti(ue.ServingTower, ue.Crnti)
				if err != nil {
					log.Error(err.Error())
					continue
				}
				ueRelInd := formatUeReleaseInd(ue.ServingTower, ue.Crnti)
				s.indChan <- *ueRelInd
				log.Infof("ueReleaseInd eci:%s crnti:%s", ue.ServingTower, ue.Crnti)
				ue.Crnti = manager.InvalidCrnti
			}
			// Nothing to be done for trafficsim.Type_UPDATED - they are handled by Telemetry
		case <-s.stream.Context().Done():
			log.Infof("Controller has disconnected")
			return
		}
	}
}

func formatUeAdmissionReq(eci *types.ECGI, crnti types.Crnti, imsi types.Imsi) *e2ap.RicIndication {
	e2Ecgi := toE2Ecgi(eci)

	return &e2ap.RicIndication{
		Hdr: &e2sm.RicIndicationHeader{
			MessageType: e2.MessageType_UE_ADMISSION_REQUEST,
		},
		Msg: &e2sm.RicIndicationMessage{
			S: &e2sm.RicIndicationMessage_UEAdmissionRequest{
				UEAdmissionRequest: &e2.UEAdmissionRequest{
					Ecgi:              &e2Ecgi,
					Crnti:             string(crnti),
					AdmissionEstCause: e2.AdmEstCause_MO_SIGNALLING,
					Imsi:              uint64(imsi),
				},
			},
		},
	}
}

func formatUeReleaseInd(eci *types.ECGI, crnti types.Crnti) *e2ap.RicIndication {
	e2Ecgi := toE2Ecgi(eci)
	return &e2ap.RicIndication{
		Hdr: &e2sm.RicIndicationHeader{
			MessageType: e2.MessageType_UE_RELEASE_IND,
		},
		Msg: &e2sm.RicIndicationMessage{
			S: &e2sm.RicIndicationMessage_UEReleaseInd{
				UEReleaseInd: &e2.UEReleaseInd{
					Ecgi:         &e2Ecgi,
					Crnti:        string(crnti),
					ReleaseCause: e2.ReleaseCause_RELEASE_INACTIVITY,
				},
			},
		},
	}
}
