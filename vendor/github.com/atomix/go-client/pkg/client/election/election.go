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

package election

import (
	"context"
	api "github.com/atomix/api/proto/atomix/election"
	"github.com/atomix/api/proto/atomix/headers"
	"github.com/atomix/go-client/pkg/client/primitive"
	"github.com/atomix/go-client/pkg/client/util"
	"github.com/google/uuid"
	"google.golang.org/grpc"
)

// Option is an election option
type Option interface {
	apply(options *options)
}

// options is election options
type options struct {
	id string
}

// idOption is an identifier option
type idOption struct {
	id string
}

func (o *idOption) apply(options *options) {
	options.id = o.id
}

// WithID sets the election instance identifier
func WithID(id string) Option {
	return &idOption{
		id: id,
	}
}

// Type is the election type
const Type primitive.Type = "Election"

// Client provides an API for creating Elections
type Client interface {
	// GetElection gets the Election instance of the given name
	GetElection(ctx context.Context, name string, opts ...Option) (Election, error)
}

// Election provides distributed leader election
type Election interface {
	primitive.Primitive

	// ID returns the ID of the instance of the election
	ID() string

	// GetTerm gets the current election term
	GetTerm(ctx context.Context) (*Term, error)

	// Enter enters the instance into the election
	Enter(ctx context.Context) (*Term, error)

	// Leave removes the instance from the election
	Leave(ctx context.Context) (*Term, error)

	// Anoint assigns leadership to the instance with the given ID
	Anoint(ctx context.Context, id string) (*Term, error)

	// Promote increases the priority of the instance with the given ID in the election queue
	Promote(ctx context.Context, id string) (*Term, error)

	// Evict removes the instance with the given ID from the election
	Evict(ctx context.Context, id string) (*Term, error)

	// Watch watches the election for changes
	Watch(ctx context.Context, c chan<- *Event) error
}

// newTerm returns a new term from the response term
func newTerm(term *api.Term) *Term {
	if term == nil {
		return nil
	}
	return &Term{
		ID:         term.ID,
		Leader:     term.Leader,
		Candidates: term.Candidates,
	}
}

// Term is a leadership term
// A term is guaranteed to have a monotonically increasing, globally unique ID.
type Term struct {
	// ID is a globally unique, monotonically increasing term number
	ID uint64

	// Leader is the ID of the leader that was elected
	Leader string

	// Candidates is a list of candidates currently participating in the election
	Candidates []string
}

// EventType is the type of an Election event
type EventType string

const (
	// EventChanged indicates the election term changed
	EventChanged EventType = "changed"
)

// Event is an election event
type Event struct {
	// Type is the type of the event
	Type EventType

	// Term is the term that occurs as a result of the election event
	Term Term
}

// New creates a new election primitive
func New(ctx context.Context, name primitive.Name, partitions []*primitive.Session, opts ...Option) (Election, error) {
	options := &options{
		id: uuid.New().String(),
	}
	for _, opt := range opts {
		opt.apply(options)
	}

	i, err := util.GetPartitionIndex(name.Name, len(partitions))
	if err != nil {
		return nil, err
	}

	instance, err := primitive.NewInstance(ctx, name, partitions[i], &primitiveHandler{})
	if err != nil {
		return nil, err
	}

	return &election{
		id:       options.id,
		name:     name,
		instance: instance,
	}, nil
}

// election is the default single-partition implementation of Election
type election struct {
	id       string
	name     primitive.Name
	instance *primitive.Instance
}

func (e *election) Name() primitive.Name {
	return e.name
}

func (e *election) ID() string {
	return e.id
}

func (e *election) GetTerm(ctx context.Context) (*Term, error) {
	response, err := e.instance.DoQuery(ctx, func(ctx context.Context, conn *grpc.ClientConn, header *headers.RequestHeader) (*headers.ResponseHeader, interface{}, error) {
		client := api.NewLeaderElectionServiceClient(conn)
		request := &api.GetTermRequest{
			Header: header,
		}
		response, err := client.GetTerm(ctx, request)
		if err != nil {
			return nil, nil, err
		}
		return response.Header, response, nil
	})
	if err != nil {
		return nil, err
	}
	return newTerm(response.(*api.GetTermResponse).Term), nil
}

func (e *election) Enter(ctx context.Context) (*Term, error) {
	response, err := e.instance.DoCommand(ctx, func(ctx context.Context, conn *grpc.ClientConn, header *headers.RequestHeader) (*headers.ResponseHeader, interface{}, error) {
		client := api.NewLeaderElectionServiceClient(conn)
		request := &api.EnterRequest{
			Header:      header,
			CandidateID: e.ID(),
		}
		response, err := client.Enter(ctx, request)
		if err != nil {
			return nil, nil, err
		}
		return response.Header, response, nil
	})
	if err != nil {
		return nil, err
	}
	return newTerm(response.(*api.EnterResponse).Term), nil
}

func (e *election) Leave(ctx context.Context) (*Term, error) {
	response, err := e.instance.DoCommand(ctx, func(ctx context.Context, conn *grpc.ClientConn, header *headers.RequestHeader) (*headers.ResponseHeader, interface{}, error) {
		client := api.NewLeaderElectionServiceClient(conn)
		request := &api.WithdrawRequest{
			Header:      header,
			CandidateID: e.ID(),
		}
		response, err := client.Withdraw(ctx, request)
		if err != nil {
			return nil, nil, err
		}
		return response.Header, response, nil
	})
	if err != nil {
		return nil, err
	}
	return newTerm(response.(*api.WithdrawResponse).Term), nil
}

func (e *election) Anoint(ctx context.Context, id string) (*Term, error) {
	response, err := e.instance.DoCommand(ctx, func(ctx context.Context, conn *grpc.ClientConn, header *headers.RequestHeader) (*headers.ResponseHeader, interface{}, error) {
		client := api.NewLeaderElectionServiceClient(conn)
		request := &api.AnointRequest{
			Header:      header,
			CandidateID: id,
		}
		response, err := client.Anoint(ctx, request)
		if err != nil {
			return nil, nil, err
		}
		return response.Header, response, nil
	})
	if err != nil {
		return nil, err
	}
	return newTerm(response.(*api.AnointResponse).Term), nil
}

func (e *election) Promote(ctx context.Context, id string) (*Term, error) {
	response, err := e.instance.DoCommand(ctx, func(ctx context.Context, conn *grpc.ClientConn, header *headers.RequestHeader) (*headers.ResponseHeader, interface{}, error) {
		client := api.NewLeaderElectionServiceClient(conn)
		request := &api.PromoteRequest{
			Header:      header,
			CandidateID: id,
		}
		response, err := client.Promote(ctx, request)
		if err != nil {
			return nil, nil, err
		}
		return response.Header, response, nil
	})
	if err != nil {
		return nil, err
	}
	return newTerm(response.(*api.PromoteResponse).Term), nil
}

func (e *election) Evict(ctx context.Context, id string) (*Term, error) {
	response, err := e.instance.DoCommand(ctx, func(ctx context.Context, conn *grpc.ClientConn, header *headers.RequestHeader) (*headers.ResponseHeader, interface{}, error) {
		client := api.NewLeaderElectionServiceClient(conn)
		request := &api.EvictRequest{
			Header:      header,
			CandidateID: id,
		}
		response, err := client.Evict(ctx, request)
		if err != nil {
			return nil, nil, err
		}
		return response.Header, response, nil
	})
	if err != nil {
		return nil, err
	}
	return newTerm(response.(*api.EvictResponse).Term), nil
}

func (e *election) Watch(ctx context.Context, ch chan<- *Event) error {
	stream, err := e.instance.DoCommandStream(ctx, func(ctx context.Context, conn *grpc.ClientConn, header *headers.RequestHeader) (interface{}, error) {
		client := api.NewLeaderElectionServiceClient(conn)
		request := &api.EventRequest{
			Header: header,
		}
		return client.Events(ctx, request)
	}, func(responses interface{}) (*headers.ResponseHeader, interface{}, error) {
		response, err := responses.(api.LeaderElectionService_EventsClient).Recv()
		if err != nil {
			return nil, nil, err
		}
		return response.Header, response, nil
	})
	if err != nil {
		return err
	}

	go func() {
		defer close(ch)
		for event := range stream {
			response := event.(*api.EventResponse)
			ch <- &Event{
				Type: EventChanged,
				Term: *newTerm(response.Term),
			}
		}
	}()
	return nil
}

func (e *election) Close(ctx context.Context) error {
	return e.instance.Close(ctx)
}

func (e *election) Delete(ctx context.Context) error {
	return e.instance.Delete(ctx)
}
