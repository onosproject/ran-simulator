// SPDX-FileCopyrightText: 2020-present Open Networking Foundation <info@opennetworking.org>
//
// SPDX-License-Identifier: LicenseRef-ONF-Member-1.0

package connection

import (
	e2 "github.com/onosproject/onos-e2t/pkg/protocols/e2ap"
	"github.com/onosproject/ran-simulator/pkg/e2agent/addressing"
	"github.com/onosproject/ran-simulator/pkg/model"
	"github.com/onosproject/ran-simulator/pkg/servicemodel/registry"
	"github.com/onosproject/ran-simulator/pkg/store/connections"
	"github.com/onosproject/ran-simulator/pkg/store/subscriptions"
)

// InstanceOptions e2 channel instance options
type InstanceOptions struct {
	node            model.Node
	model           *model.Model
	ricAddress      addressing.RICAddress
	channel         e2.ClientConn
	registry        *registry.ServiceModelRegistry
	subStore        *subscriptions.Subscriptions
	connectionStore connections.Store
}

// InstanceOption instance option
type InstanceOption func(*InstanceOptions)

// WithNode sets model node
func WithNode(node model.Node) func(options *InstanceOptions) {
	return func(options *InstanceOptions) {
		options.node = node
	}
}

// WithModel sets model
func WithModel(model *model.Model) func(options *InstanceOptions) {
	return func(options *InstanceOptions) {
		options.model = model
	}
}

// WithRICAddress sets RIC address
func WithRICAddress(ricAddress addressing.RICAddress) func(options *InstanceOptions) {
	return func(options *InstanceOptions) {
		options.ricAddress = ricAddress
	}
}

// WithChannel sets E2 channel
func WithChannel(channel e2.ClientConn) func(options *InstanceOptions) {
	return func(options *InstanceOptions) {
		options.channel = channel
	}
}

// WithSMRegistry sets service model registry
func WithSMRegistry(registry *registry.ServiceModelRegistry) func(options *InstanceOptions) {
	return func(options *InstanceOptions) {
		options.registry = registry
	}
}

// WithSubStore sets subscription store
func WithSubStore(subStore *subscriptions.Subscriptions) func(options *InstanceOptions) {
	return func(options *InstanceOptions) {
		options.subStore = subStore
	}
}

// WithConnectionStore sets connection store
func WithConnectionStore(connectionStore connections.Store) func(options *InstanceOptions) {
	return func(options *InstanceOptions) {
		options.connectionStore = connectionStore
	}
}
