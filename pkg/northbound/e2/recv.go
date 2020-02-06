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
	"regexp"
	"strconv"

	"github.com/onosproject/ran-simulator/api/e2"
	"github.com/onosproject/ran-simulator/pkg/manager"
	"github.com/prometheus/common/log"
)

// TestPlmnID - https://en.wikipedia.org/wiki/Mobile_country_code#Test_networks
const TestPlmnID = "001001"

func makeEci(towerName string) string {
	re := regexp.MustCompile("[0-9]+")
	id, _ := strconv.Atoi(re.FindAllString(towerName, 1)[0])
	return fmt.Sprintf("%07X", id)
}

func makeCrnti(ueName string) string {
	re := regexp.MustCompile("[0-9]+")
	id, _ := strconv.Atoi(re.FindAllString(ueName, 1)[0])
	return fmt.Sprintf("%04X", id+1)
}

func recv(stream e2.InterfaceService_SendControlServer, c chan e2.ControlUpdate) {
	for {
		in, err := stream.Recv()
		if err == io.EOF || err != nil {
			return
		}
		log.Infof("Recv messageType %d", in.MessageType)
		switch x := in.S.(type) {
		case *e2.ControlResponse_CellConfigRequest:
			handleCellConfigRequest(stream, x.CellConfigRequest, c)
		default:
			log.Errorf("ControlResponse has unexpected type %T", x)
		}
	}
}

func handleCellConfigRequest(stream e2.InterfaceService_SendControlServer, req *e2.CellConfigRequest, c chan e2.ControlUpdate) {
	log.Infof("handleCellConfigRequest")

	mgr := manager.GetManager()

	for _, tower := range mgr.Towers {
		tower.Eci = makeEci(tower.Name)
		cellConfigReport := e2.ControlUpdate{
			MessageType: e2.MessageType_CELL_CONFIG_REPORT,
			S: &e2.ControlUpdate_CellConfigReport{
				CellConfigReport: &e2.CellConfigReport{
					Ecgi: &e2.ECGI{
						PlmnId: TestPlmnID,
						Ecid:   tower.Eci,
					},
				},
			},
		}

		c <- cellConfigReport
		log.Infof("handleCellConfigReport eci: %s", tower.Eci)
	}

	// Initate UE admissions
	for _, ue := range mgr.UserEquipments {
		log.Infof("Ue: %s", ue.Name)
		ue.Crnti = makeCrnti(ue.Name)
		log.Infof("Ue: %s", ue.Crnti)

		eci := mgr.GetTowerByName(ue.Tower).Eci
		ueAdmReq := e2.ControlUpdate{
			MessageType: e2.MessageType_UE_ADMISSION_REQUEST,
			S: &e2.ControlUpdate_UEAdmissionRequest{
				UEAdmissionRequest: &e2.UEAdmissionRequest{
					Ecgi: &e2.ECGI{
						PlmnId: TestPlmnID,
						Ecid:   eci,
					},
					Crnti:             ue.Crnti,
					AdmissionEstCause: e2.AdmEstCause_MO_SIGNALLING,
				},
			},
		}
		c <- ueAdmReq
		log.Infof("ueAdmissionRequest eci:%s crnti:%s", eci, ue.Crnti)
	}
}
