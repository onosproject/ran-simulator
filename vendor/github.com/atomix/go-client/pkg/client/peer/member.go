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

package peer

import (
	"fmt"
	"google.golang.org/grpc"
	"net"
)

// NewMember returns a new local group member
func NewMember(id ID, host string, port int, services ...Service) *Member {
	return &Member{
		services: services,
		Peer: &Peer{
			ID:   id,
			Host: host,
			Port: port,
		},
	}
}

// Service is a peer-to-peer primitive service
type Service func(ID, *grpc.Server)

// Member is a local group member
type Member struct {
	*Peer
	services []Service
	stopCh   chan struct{}
}

// serve begins serving the local member
func (m *Member) serve(opts ...grpc.ServerOption) error {
	server := grpc.NewServer(opts...)
	for _, service := range m.services {
		service(m.ID, server)
	}

	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", m.Port))
	if err != nil {
		return err
	}

	go func() {
		err := server.Serve(lis)
		if err != nil {
			fmt.Println(err)
		}
	}()
	go func() {
		<-m.stopCh
		server.Stop()
	}()
	return nil
}

// Stop stops the local member serving
func (m *Member) Stop() {
	close(m.stopCh)
}
