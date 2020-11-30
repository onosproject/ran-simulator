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
	"github.com/atomix/go-framework/pkg/atomix/primitive"
	streams "github.com/atomix/go-framework/pkg/atomix/stream"
	"github.com/golang/protobuf/proto"
	log "github.com/sirupsen/logrus"
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

// Size gets the number of elements in the list
func (s *Server) Size(ctx context.Context, request *api.SizeRequest) (*api.SizeResponse, error) {
	log.Tracef("Received SizeRequest %+v", request)
	in, err := proto.Marshal(&SizeRequest{})
	if err != nil {
		return nil, err
	}

	out, header, err := s.DoQuery(ctx, opSize, in, request.Header)
	if err != nil {
		return nil, err
	}

	sizeResponse := &SizeResponse{}
	if err = proto.Unmarshal(out, sizeResponse); err != nil {
		return nil, err
	}

	response := &api.SizeResponse{
		Header: header,
		Size_:  sizeResponse.Size_,
	}
	log.Tracef("Sending SizeResponse %+v", response)
	return response, nil
}

// Contains checks whether the list contains a value
func (s *Server) Contains(ctx context.Context, request *api.ContainsRequest) (*api.ContainsResponse, error) {
	log.Tracef("Received ContainsRequest %+v", request)
	in, err := proto.Marshal(&ContainsRequest{
		Value: request.Value,
	})
	if err != nil {
		return nil, err
	}

	out, header, err := s.DoQuery(ctx, opContains, in, request.Header)
	if err != nil {
		return nil, err
	}

	containsResponse := &ContainsResponse{}
	if err = proto.Unmarshal(out, containsResponse); err != nil {
		return nil, err
	}

	response := &api.ContainsResponse{
		Header:   header,
		Contains: containsResponse.Contains,
	}
	log.Tracef("Sending ContainsResponse %+v", response)
	return response, nil
}

// Append adds a value to the end of the list
func (s *Server) Append(ctx context.Context, request *api.AppendRequest) (*api.AppendResponse, error) {
	log.Tracef("Received AppendRequest %+v", request)
	in, err := proto.Marshal(&AppendRequest{
		Value: request.Value,
	})
	if err != nil {
		return nil, err
	}

	out, header, err := s.DoCommand(ctx, opAppend, in, request.Header)
	if err != nil {
		return nil, err
	}

	appendResponse := &AppendResponse{}
	if err = proto.Unmarshal(out, appendResponse); err != nil {
		return nil, err
	}

	response := &api.AppendResponse{
		Header: header,
		Status: getResponseStatus(appendResponse.Status),
	}
	log.Tracef("Sending AppendResponse %+v", response)
	return response, nil
}

// Insert inserts a value at a specific index
func (s *Server) Insert(ctx context.Context, request *api.InsertRequest) (*api.InsertResponse, error) {
	log.Tracef("Received InsertRequest %+v", request)
	in, err := proto.Marshal(&InsertRequest{
		Index: request.Index,
		Value: request.Value,
	})
	if err != nil {
		return nil, err
	}

	out, header, err := s.DoCommand(ctx, opInsert, in, request.Header)
	if err != nil {
		return nil, err
	}

	insertResponse := &InsertResponse{}
	if err = proto.Unmarshal(out, insertResponse); err != nil {
		return nil, err
	}

	response := &api.InsertResponse{
		Header: header,
		Status: getResponseStatus(insertResponse.Status),
	}
	log.Tracef("Sending InsertResponse %+v", response)
	return response, nil
}

// Set sets the value at a specific index
func (s *Server) Set(ctx context.Context, request *api.SetRequest) (*api.SetResponse, error) {
	log.Tracef("Received SetRequest %+v", request)
	in, err := proto.Marshal(&SetRequest{
		Index: request.Index,
		Value: request.Value,
	})
	if err != nil {
		return nil, err
	}

	out, header, err := s.DoCommand(ctx, opSet, in, request.Header)
	if err != nil {
		return nil, err
	}

	setResponse := &SetResponse{}
	if err = proto.Unmarshal(out, setResponse); err != nil {
		return nil, err
	}

	response := &api.SetResponse{
		Header: header,
		Status: getResponseStatus(setResponse.Status),
	}
	log.Tracef("Sending SetResponse %+v", response)
	return response, nil
}

// Get gets the value at a specific index
func (s *Server) Get(ctx context.Context, request *api.GetRequest) (*api.GetResponse, error) {
	log.Tracef("Received GetRequest %+v", request)
	in, err := proto.Marshal(&GetRequest{
		Index: request.Index,
	})
	if err != nil {
		return nil, err
	}

	out, header, err := s.DoQuery(ctx, opGet, in, request.Header)
	if err != nil {
		return nil, err
	}

	getResponse := &GetResponse{}
	if err = proto.Unmarshal(out, getResponse); err != nil {
		return nil, err
	}

	response := &api.GetResponse{
		Header: header,
		Status: getResponseStatus(getResponse.Status),
		Value:  getResponse.Value,
	}
	log.Tracef("Sending GetResponse %+v", response)
	return response, nil
}

// Remove removes an index from the list
func (s *Server) Remove(ctx context.Context, request *api.RemoveRequest) (*api.RemoveResponse, error) {
	log.Tracef("Received RemoveRequest %+v", request)
	in, err := proto.Marshal(&RemoveRequest{
		Index: request.Index,
	})
	if err != nil {
		return nil, err
	}

	out, header, err := s.DoCommand(ctx, opRemove, in, request.Header)
	if err != nil {
		return nil, err
	}

	removeResponse := &RemoveResponse{}
	if err = proto.Unmarshal(out, removeResponse); err != nil {
		return nil, err
	}

	response := &api.RemoveResponse{
		Header: header,
		Status: getResponseStatus(removeResponse.Status),
		Value:  removeResponse.Value,
	}
	log.Tracef("Sending RemoveResponse %+v", response)
	return response, nil
}

// Clear removes all indexes from the list
func (s *Server) Clear(ctx context.Context, request *api.ClearRequest) (*api.ClearResponse, error) {
	log.Tracef("Received ClearRequest %+v", request)
	in, err := proto.Marshal(&ClearRequest{})
	if err != nil {
		return nil, err
	}

	out, header, err := s.DoCommand(ctx, opClear, in, request.Header)
	if err != nil {
		return nil, err
	}

	clearResponse := &ClearResponse{}
	if err = proto.Unmarshal(out, clearResponse); err != nil {
		return nil, err
	}

	response := &api.ClearResponse{
		Header: header,
	}
	log.Tracef("Sending ClearResponse %+v", response)
	return response, nil
}

// Events listens for list change events
func (s *Server) Events(request *api.EventRequest, srv api.ListService_EventsServer) error {
	log.Tracef("Received EventRequest %+v", request)
	in, err := proto.Marshal(&ListenRequest{
		Replay: request.Replay,
	})
	if err != nil {
		return err
	}

	stream := streams.NewBufferedStream()
	if err := s.DoCommandStream(srv.Context(), opEvents, in, request.Header, stream); err != nil {
		return err
	}

	for {
		result, ok := stream.Receive()
		if !ok {
			break
		}

		if result.Failed() {
			return result.Error
		}

		response := &ListenResponse{}
		output := result.Value.(primitive.SessionOutput)
		if err = proto.Unmarshal(output.Value.([]byte), response); err != nil {
			return err
		}

		var eventResponse *api.EventResponse
		switch output.Header.Type {
		case headers.ResponseType_OPEN_STREAM:
			eventResponse = &api.EventResponse{
				Header: output.Header,
			}
		case headers.ResponseType_CLOSE_STREAM:
			eventResponse = &api.EventResponse{
				Header: output.Header,
			}
		default:
			eventResponse = &api.EventResponse{
				Header: output.Header,
				Type:   getEventType(response.Type),
				Index:  response.Index,
				Value:  response.Value,
			}
		}

		log.Tracef("Sending EventResponse %+v", eventResponse)
		if err = srv.Send(eventResponse); err != nil {
			return err
		}
	}

	log.Tracef("Finished EventRequest %+v", request)
	return nil
}

// Iterate lists all the value in the list
func (s *Server) Iterate(request *api.IterateRequest, srv api.ListService_IterateServer) error {
	log.Tracef("Received IterateRequest %+v", request)
	in, err := proto.Marshal(&IterateRequest{})
	if err != nil {
		return err
	}

	stream := streams.NewBufferedStream()
	if err := s.DoQueryStream(srv.Context(), opIterate, in, request.Header, stream); err != nil {
		return err
	}

	for {
		result, ok := stream.Receive()
		if !ok {
			break
		}

		if result.Failed() {
			return result.Error
		}

		response := &IterateResponse{}
		output := result.Value.(primitive.SessionOutput)
		if err = proto.Unmarshal(output.Value.([]byte), response); err != nil {
			return err
		}

		var iterateResponse *api.IterateResponse
		switch output.Header.Type {
		case headers.ResponseType_OPEN_STREAM:
			iterateResponse = &api.IterateResponse{
				Header: output.Header,
			}
		case headers.ResponseType_CLOSE_STREAM:
			iterateResponse = &api.IterateResponse{
				Header: output.Header,
			}
		default:
			iterateResponse = &api.IterateResponse{
				Header: output.Header,
				Value:  response.Value,
			}
		}

		log.Tracef("Sending IterateResponse %+v", iterateResponse)
		if err = srv.Send(iterateResponse); err != nil {
			return err
		}
	}

	log.Tracef("Finished IterateRequest %+v", request)
	return nil
}

func getResponseStatus(status ResponseStatus) api.ResponseStatus {
	switch status {
	case ResponseStatus_OK:
		return api.ResponseStatus_OK
	case ResponseStatus_NOOP:
		return api.ResponseStatus_NOOP
	case ResponseStatus_WRITE_LOCK:
		return api.ResponseStatus_WRITE_LOCK
	case ResponseStatus_OUT_OF_BOUNDS:
		return api.ResponseStatus_OUT_OF_BOUNDS
	}
	return api.ResponseStatus_OK
}

func getEventType(eventType ListenResponse_Type) api.EventResponse_Type {
	switch eventType {
	case ListenResponse_NONE:
		return api.EventResponse_NONE
	case ListenResponse_ADDED:
		return api.EventResponse_ADDED
	case ListenResponse_REMOVED:
		return api.EventResponse_REMOVED
	}
	return api.EventResponse_NONE
}
