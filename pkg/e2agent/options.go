// SPDX-FileCopyrightText: ${year}-present Open Networking Foundation <info@opennetworking.org>
// SPDX-License-Identifier: Apache-2.0

package e2agent

import (
	e2 "github.com/onosproject/onos-e2t/pkg/protocols/e2ap101"
	"github.com/onosproject/ran-simulator/pkg/model"
	"github.com/onosproject/ran-simulator/pkg/servicemodel/registry"
	"github.com/onosproject/ran-simulator/pkg/store/channels"
	"github.com/onosproject/ran-simulator/pkg/store/subscriptions"
)

// InstanceOptions e2 instance options
type InstanceOptions struct {
	node         model.Node
	model        *model.Model
	ricAddress   RICAddress
	channel      e2.ClientChannel
	registry     *registry.ServiceModelRegistry
	subStore     *subscriptions.Subscriptions
	channelStore channels.Store
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
func WithRICAddress(ricAddress RICAddress) func(options *InstanceOptions) {
	return func(options *InstanceOptions) {
		options.ricAddress = ricAddress
	}
}

// WithChannel sets E2 channel
func WithChannel(channel e2.ClientChannel) func(options *InstanceOptions) {
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

// WithChannelStore sets channel store
func WithChannelStore(channelStore channels.Store) func(options *InstanceOptions) {
	return func(options *InstanceOptions) {
		options.channelStore = channelStore
	}
}
