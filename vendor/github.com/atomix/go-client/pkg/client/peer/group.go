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
	"context"
	"errors"
	"fmt"
	membershipapi "github.com/atomix/api/proto/atomix/membership"
	"github.com/atomix/go-client/pkg/client/util"
	"google.golang.org/grpc"
	"io"
	"sync"
	"time"
)

// NewGroup creates a new peer group
func NewGroup(address string, opts ...Option) (*Group, error) {
	ctx := context.Background()
	options := applyOptions(opts...)
	if options.joinTimeout != nil {
		c, cancel := context.WithTimeout(context.Background(), *options.joinTimeout)
		defer cancel()
		ctx = c
	}
	return NewGroupWithContext(ctx, address, opts...)
}

// NewGroupWithContext creates a new peer group
func NewGroupWithContext(ctx context.Context, address string, opts ...Option) (*Group, error) {
	options := applyOptions(opts...)

	conn, err := grpc.DialContext(ctx, address, grpc.WithInsecure(), grpc.WithBlock(), grpc.WithUnaryInterceptor(util.RetryingUnaryClientInterceptor()), grpc.WithStreamInterceptor(util.RetryingStreamClientInterceptor(time.Second)))
	if err != nil {
		return nil, err
	}

	var member *Member
	if options.memberID != "" {
		services := options.services
		if services == nil {
			services = []Service{}
		}
		member = NewMember(ID(options.memberID), options.peerHost, options.peerPort, services...)
	}

	group := &Group{
		Namespace: options.namespace,
		Name:      options.scope,
		member:    member,
		peers:     make(Set),
		conn:      conn,
		options:   *options,
		leaveCh:   make(chan struct{}),
		watchers:  make([]chan<- Set, 0),
	}

	err = group.join(ctx)
	if err != nil {
		return nil, err
	}
	return group, nil
}

// Group manages the peer group for a client
type Group struct {
	Namespace string
	Name      string
	member    *Member
	conn      *grpc.ClientConn
	options   options
	peers     Set
	watchers  []chan<- Set
	closer    context.CancelFunc
	leaveCh   chan struct{}
	mu        sync.RWMutex
}

// Member returns the local group member
func (c *Group) Member() *Peer {
	return c.member.Peer
}

// Peer returns a peer by ID
func (c *Group) Peer(id ID) *Peer {
	return c.peers[id]
}

// Peers returns the current group peers
func (c *Group) Peers() Set {
	c.mu.RLock()
	defer c.mu.RUnlock()
	if c.peers != nil {
		return c.peers
	}
	return Set{}
}

// join joins the group
func (c *Group) join(ctx context.Context) error {
	if c.member != nil {
		err := c.member.serve()
		if err != nil {
			return nil
		}
	}

	var member *membershipapi.Member
	if c.member != nil {
		member = &membershipapi.Member{
			ID: membershipapi.MemberId{
				Namespace: c.options.namespace,
				Name:      c.options.memberID,
			},
			Host: c.member.Host,
			Port: int32(c.member.Port),
		}
	}

	client := membershipapi.NewMembershipServiceClient(c.conn)
	request := &membershipapi.JoinGroupRequest{
		Member: member,
		GroupID: membershipapi.GroupId{
			Namespace: c.options.namespace,
			Name:      c.options.scope,
		},
	}
	streamCtx, cancel := context.WithCancel(context.Background())
	stream, err := client.JoinGroup(streamCtx, request)
	if err != nil {
		cancel()
		return err
	}

	c.mu.Lock()
	joinCh := make(chan struct{})
	c.closer = cancel
	c.mu.Unlock()

	go func() {
		joined := false
		for {
			response, err := stream.Recv()
			if err == io.EOF {
				close(c.leaveCh)
				return
			} else if err != nil {
				fmt.Println(err)
				close(c.leaveCh)
				return
			} else {
				c.mu.Lock()
				members := make(map[ID]membershipapi.Member)
				for _, member := range response.Members {
					members[ID(member.ID.Name)] = member
				}

				for id := range c.peers {
					_, ok := members[id]
					if !ok {
						delete(c.peers, id)
					}
				}

				for id, member := range members {
					_, ok := c.peers[id]
					if !ok {
						c.peers[id] = NewPeer(id, member.Host, int(member.Port))
					}
				}
				peers := c.peers
				c.mu.Unlock()

				if !joined {
					close(joinCh)
					joined = true
				}

				c.mu.RLock()
				for _, watcher := range c.watchers {
					watcher <- peers
				}
				c.mu.RUnlock()
			}
		}
	}()

	select {
	case <-joinCh:
		return nil
	case <-ctx.Done():
		return errors.New("join timed out")
	}
}

// Watch watches the peers for changes
func (c *Group) Watch(ctx context.Context, ch chan<- Set) error {
	c.mu.Lock()
	watcher := make(chan Set)
	peers := c.peers
	go func() {
		if peers != nil {
			ch <- peers
		}
		for {
			select {
			case peers, ok := <-watcher:
				if !ok {
					return
				}
				ch <- peers
			case <-ctx.Done():
				c.mu.Lock()
				watchers := make([]chan<- Set, 0)
				for _, ch := range c.watchers {
					if ch != watcher {
						watchers = append(watchers, ch)
					}
				}
				c.watchers = watchers
				c.mu.Unlock()
				close(watcher)
			}
		}
	}()
	c.watchers = append(c.watchers, watcher)
	c.mu.Unlock()
	return nil
}

// Close closes the group
func (c *Group) Close() error {
	c.mu.RLock()
	closer := c.closer
	leaveCh := c.leaveCh
	c.mu.RUnlock()
	timeout := time.Minute
	if c.options.joinTimeout != nil {
		timeout = *c.options.joinTimeout
	}
	if closer != nil {
		closer()
		select {
		case <-leaveCh:
			return nil
		case <-time.After(timeout):
			return errors.New("leave timed out")
		}
	}
	return nil
}

// Set is a set of peers
type Set map[ID]*Peer
