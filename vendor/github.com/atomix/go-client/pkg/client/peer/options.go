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
	"google.golang.org/grpc"
	"os"
	"time"
)

func applyOptions(opts ...Option) *options {
	options := &options{
		namespace:     os.Getenv("ATOMIX_NAMESPACE"),
		scope:         os.Getenv("ATOMIX_SCOPE"),
		peerPort:      8080,
		services:      make([]Service, 0),
		serverOptions: make([]grpc.ServerOption, 0),
	}
	for _, opt := range opts {
		opt.apply(options)
	}
	return options
}

type options struct {
	memberID      string
	peerHost      string
	peerPort      int
	services      []Service
	joinTimeout   *time.Duration
	scope         string
	namespace     string
	serverOptions []grpc.ServerOption
}

// Option provides a peer option
type Option interface {
	apply(options *options)
}

// WithMemberID configures the peer's member ID
func WithMemberID(memberID string) Option {
	return &memberIDOption{id: memberID}
}

type memberIDOption struct {
	id string
}

func (o *memberIDOption) apply(options *options) {
	options.memberID = o.id
	if options.peerHost == "" {
		options.peerHost = o.id
	}
}

// WithHost configures the peer's host
func WithHost(host string) Option {
	return &hostOption{host: host}
}

type hostOption struct {
	host string
}

func (o *hostOption) apply(options *options) {
	options.peerHost = o.host
}

// WithPort configures the peer's port
func WithPort(port int) Option {
	return &portOption{port: port}
}

type portOption struct {
	port int
}

func (o *portOption) apply(options *options) {
	options.peerPort = o.port
}

// WithService configures a peer-to-peer service
func WithService(service Service) Option {
	return &serviceOption{
		service: service,
	}
}

type serviceOption struct {
	service Service
}

func (o *serviceOption) apply(options *options) {
	options.services = append(options.services, o.service)
}

// WithServices configures peer-to-peer services
func WithServices(services ...Service) Option {
	return &servicesOption{
		services: services,
	}
}

type servicesOption struct {
	services []Service
}

func (o *servicesOption) apply(options *options) {
	options.services = append(options.services, o.services...)
}

// WithJoinTimeout configures the client's join timeout
func WithJoinTimeout(timeout time.Duration) Option {
	return &joinTimeoutOption{timeout: timeout}
}

type joinTimeoutOption struct {
	timeout time.Duration
}

func (o *joinTimeoutOption) apply(options *options) {
	options.joinTimeout = &o.timeout
}

type scopeOption struct {
	scope string
}

func (o *scopeOption) apply(options *options) {
	options.scope = o.scope
}

// WithScope configures the application scope for the client
func WithScope(scope string) Option {
	return &scopeOption{scope: scope}
}

type namespaceOption struct {
	namespace string
}

func (o *namespaceOption) apply(options *options) {
	options.namespace = o.namespace
}

// WithNamespace configures the client's partition group namespace
func WithNamespace(namespace string) Option {
	return &namespaceOption{namespace: namespace}
}

// WithServerOption configures a server option
func WithServerOption(option grpc.ServerOption) Option {
	return &serverOptionOption{
		option: option,
	}
}

type serverOptionOption struct {
	option grpc.ServerOption
}

func (o *serverOptionOption) apply(options *options) {
	options.serverOptions = append(options.serverOptions, o.option)
}

// WithServerOptions configures server options
func WithServerOptions(options ...grpc.ServerOption) Option {
	return &serverOptionsOption{
		options: options,
	}
}

type serverOptionsOption struct {
	options []grpc.ServerOption
}

func (o *serverOptionsOption) apply(options *options) {
	options.serverOptions = append(options.serverOptions, o.options...)
}

func applyConnectOptions(opts ...ConnectOption) *connectOptions {
	options := &connectOptions{}
	for _, opt := range opts {
		opt.apply(options)
	}
	if options.dialOptions == nil {
		options.dialOptions = []grpc.DialOption{
			grpc.WithInsecure(),
		}
	}
	return options
}

type connectOptions struct {
	dialOptions []grpc.DialOption
}

// ConnectOption is an option for connecting to a peer
type ConnectOption interface {
	apply(options *connectOptions)
}

// WithDialOption creates a dial option for the gRPC connection
func WithDialOption(option grpc.DialOption) ConnectOption {
	return &dialOptionOption{
		option: option,
	}
}

type dialOptionOption struct {
	option grpc.DialOption
}

func (o *dialOptionOption) apply(options *connectOptions) {
	if options.dialOptions == nil {
		options.dialOptions = make([]grpc.DialOption, 0)
	}
	options.dialOptions = append(options.dialOptions, o.option)
}

// WithDialOptions creates a dial option for the gRPC connection
func WithDialOptions(options ...grpc.DialOption) ConnectOption {
	return &dialOptionsOption{
		options: options,
	}
}

type dialOptionsOption struct {
	options []grpc.DialOption
}

func (o *dialOptionsOption) apply(options *connectOptions) {
	if options.dialOptions == nil {
		options.dialOptions = make([]grpc.DialOption, 0)
	}
	options.dialOptions = append(options.dialOptions, o.options...)
}
