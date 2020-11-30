// Copyright 2019-present Open Networking Foundation.
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

package list

import (
	"context"
	"github.com/atomix/api/proto/atomix/headers"
	api "github.com/atomix/api/proto/atomix/list"
	"github.com/atomix/go-client/pkg/client/primitive"
	"google.golang.org/grpc"
)

type primitiveHandler struct{}

func (h *primitiveHandler) Create(ctx context.Context, s *primitive.Instance) error {
	return s.DoCreate(ctx, func(ctx context.Context, conn *grpc.ClientConn, header *headers.RequestHeader) (*headers.ResponseHeader, interface{}, error) {
		request := &api.CreateRequest{
			Header: header,
		}
		client := api.NewListServiceClient(conn)
		response, err := client.Create(ctx, request)
		if err != nil {
			return nil, nil, err
		}
		return response.Header, response, nil
	})
}

func (h *primitiveHandler) Close(ctx context.Context, s *primitive.Instance) error {
	return s.DoClose(ctx, func(ctx context.Context, conn *grpc.ClientConn, header *headers.RequestHeader) (*headers.ResponseHeader, interface{}, error) {
		request := &api.CloseRequest{
			Header: header,
		}
		client := api.NewListServiceClient(conn)
		response, err := client.Close(ctx, request)
		if err != nil {
			return nil, nil, err
		}
		return response.Header, response, nil
	})
}

func (h *primitiveHandler) Delete(ctx context.Context, s *primitive.Instance) error {
	return s.DoClose(ctx, func(ctx context.Context, conn *grpc.ClientConn, header *headers.RequestHeader) (*headers.ResponseHeader, interface{}, error) {
		request := &api.CloseRequest{
			Header: header,
			Delete: true,
		}
		client := api.NewListServiceClient(conn)
		response, err := client.Close(ctx, request)
		if err != nil {
			return nil, nil, err
		}
		return response.Header, response, nil
	})
}
