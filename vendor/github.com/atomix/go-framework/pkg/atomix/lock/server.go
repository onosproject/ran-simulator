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

package lock

import (
	"context"
	"github.com/atomix/api/proto/atomix/headers"
	api "github.com/atomix/api/proto/atomix/lock"
	"github.com/atomix/go-framework/pkg/atomix/primitive"
	streams "github.com/atomix/go-framework/pkg/atomix/stream"
	"github.com/golang/protobuf/proto"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// Server is an implementation of MapServiceServer for the map primitive
type Server struct {
	*primitive.Server
}

// Create opens a new session
func (s *Server) Create(ctx context.Context, request *api.CreateRequest) (*api.CreateResponse, error) {
	log.Tracef("Received CreateRequest %+v", request)
	header, err := s.DoCreateService(ctx, request.Header)
	if err != nil {
		return nil, err
	}
	response := &api.CreateResponse{
		Header: header,
	}
	log.Tracef("Sending CreateResponse %+v", response)
	return response, nil
}

// Close closes a session
func (s *Server) Close(ctx context.Context, request *api.CloseRequest) (*api.CloseResponse, error) {
	log.Tracef("Received CloseRequest %+v", request)
	if request.Delete {
		header, err := s.DoDeleteService(ctx, request.Header)
		if err != nil {
			return nil, err
		}
		response := &api.CloseResponse{
			Header: header,
		}
		log.Tracef("Sending CloseResponse %+v", response)
		return response, nil
	}

	header, err := s.DoCloseService(ctx, request.Header)
	if err != nil {
		return nil, err
	}
	response := &api.CloseResponse{
		Header: header,
	}
	log.Tracef("Sending CloseResponse %+v", response)
	return response, nil
}

// Lock acquires a lock
func (s *Server) Lock(ctx context.Context, request *api.LockRequest) (*api.LockResponse, error) {
	log.Tracef("Received LockRequest %+v", request)

	in, err := proto.Marshal(&LockRequest{
		Timeout: request.Timeout,
	})
	if err != nil {
		return nil, err
	}

	stream := streams.NewBufferedStream()
	if err := s.DoCommandStream(ctx, opLock, in, request.Header, stream); err != nil {
		return nil, err
	}

	for {
		result, ok := stream.Receive()
		if !ok {
			return nil, status.Error(codes.Canceled, "stream closed")
		}

		if result.Failed() {
			return nil, result.Error
		}

		output := result.Value.(primitive.SessionOutput)

		if output.Header.Type == headers.ResponseType_RESPONSE {
			lockResponse := &LockResponse{}
			if err = proto.Unmarshal(output.Value.([]byte), lockResponse); err != nil {
				return nil, err
			}
			response := &api.LockResponse{
				Header:  output.Header,
				Version: uint64(lockResponse.Index),
			}
			log.Tracef("Sending LockResponse %+v", response)
			return response, nil
		}
	}
}

// Unlock releases the lock
func (s *Server) Unlock(ctx context.Context, request *api.UnlockRequest) (*api.UnlockResponse, error) {
	log.Tracef("Received UnlockRequest %+v", request)
	in, err := proto.Marshal(&UnlockRequest{
		Index: int64(request.Version),
	})
	if err != nil {
		return nil, err
	}

	out, header, err := s.DoCommand(ctx, opUnlock, in, request.Header)
	if err != nil {
		return nil, err
	}

	unlockResponse := &UnlockResponse{}
	if err = proto.Unmarshal(out, unlockResponse); err != nil {
		return nil, err
	}

	response := &api.UnlockResponse{
		Header:   header,
		Unlocked: unlockResponse.Succeeded,
	}
	log.Tracef("Sending UnlockResponse %+v", response)
	return response, nil
}

// IsLocked checks whether the lock is held by any session
func (s *Server) IsLocked(ctx context.Context, request *api.IsLockedRequest) (*api.IsLockedResponse, error) {
	log.Tracef("Received IsLockedRequest %+v", request)
	in, err := proto.Marshal(&IsLockedRequest{
		Index: int64(request.Version),
	})
	if err != nil {
		return nil, err
	}

	out, header, err := s.DoQuery(ctx, opIsLocked, in, request.Header)
	if err != nil {
		return nil, err
	}

	isLockedResponse := &IsLockedResponse{}
	if err = proto.Unmarshal(out, isLockedResponse); err != nil {
		return nil, err
	}

	response := &api.IsLockedResponse{
		Header:   header,
		IsLocked: isLockedResponse.Locked,
	}
	log.Tracef("Sending IsLockedResponse %+v", response)
	return response, nil
}
