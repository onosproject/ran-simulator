// SPDX-FileCopyrightText: 2019-present Open Networking Foundation <info@opennetworking.org>
//
// SPDX-License-Identifier: Apache-2.0
//

package gnmi

import (
	"fmt"
	"github.com/openconfig/gnmi/proto/gnmi"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// Subscribe implements gNMI Subscribe
func (s *Server) Subscribe(stream gnmi.GNMI_SubscribeServer) error {
	// TODO: Implement Subscribe if we add any state attributes to the model
	return status.Error(codes.Unimplemented, fmt.Sprintf("gNMI Subscribe not yet supported on Port %d", s.port))
}
