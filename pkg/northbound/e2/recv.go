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

	"github.com/onosproject/ran-simulator/api/e2"
	"github.com/prometheus/common/log"
)

func recv(stream e2.InterfaceService_SendControlServer, c chan e2.ControlUpdate) error {
	for {
		in, err := stream.Recv()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			return err
		}
		log.Infof("Recv messageType %d", in.MessageType)
		switch x := in.S.(type) {
		case *e2.ControlResponse_CellConfigRequest:
			return handleCellConfigRequest(stream, x.CellConfigRequest, c)
		default:
			return fmt.Errorf("ControlResponse has unexpected type %T", x)
		}
	}
}

func handleCellConfigRequest(stream e2.InterfaceService_SendControlServer, req *e2.CellConfigRequest, c chan e2.ControlUpdate) error {
	log.Infof("handleCellConfigRequest")
	cellConfigReport := e2.ControlUpdate{
		MessageType: e2.MessageType_CELL_CONFIG_REPORT,
		S: &e2.ControlUpdate_CellConfigReport{
			CellConfigReport: &e2.CellConfigReport{
				Ecgi: &e2.ECGI{PlmnId: "test", Ecid: "test"},
			},
		},
	}

	c <- cellConfigReport
	log.Infof("handleCellConfigRequest sent")

	return nil
}
