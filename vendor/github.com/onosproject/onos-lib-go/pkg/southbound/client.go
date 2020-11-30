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

package southbound

import (
	"context"
	"crypto/tls"
	"fmt"
	"github.com/onosproject/onos-lib-go/pkg/certs"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

// Connect establishes a client-side connection to the gRPC end-point.
func Connect(ctx context.Context, address string, certPath string, keyPath string, opts ...grpc.DialOption) (*grpc.ClientConn, error) {
	var tlsOpts []grpc.DialOption
	if certPath != "" && keyPath != "" {
		cert, err := tls.LoadX509KeyPair(certPath, keyPath)
		if err != nil {
			return nil, err
		}
		tlsOpts = []grpc.DialOption{
			grpc.WithTransportCredentials(credentials.NewTLS(&tls.Config{
				Certificates:       []tls.Certificate{cert},
				InsecureSkipVerify: true,
			})),
		}
	} else {
		// Load default Certificates
		cert, err := tls.X509KeyPair([]byte(certs.DefaultClientCrt), []byte(certs.DefaultClientKey))
		if err != nil {
			return nil, err
		}
		tlsOpts = []grpc.DialOption{
			grpc.WithTransportCredentials(credentials.NewTLS(&tls.Config{
				Certificates:       []tls.Certificate{cert},
				InsecureSkipVerify: true,
			})),
		}
	}

	opts = append(tlsOpts, opts...)
	conn, err := grpc.DialContext(ctx, address, opts...)
	if err != nil {
		fmt.Println("Can't connect", err)
		return nil, err
	}
	return conn, nil
}
