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

package indexedmap

import (
	api "github.com/atomix/api/proto/atomix/indexedmap"
)

// SetOption is an option for the Put method
type SetOption interface {
	beforePut(request *api.PutRequest)
	afterPut(response *api.PutResponse)
}

// ReplaceOption is an option for the Replace method
type ReplaceOption interface {
	beforeReplace(request *api.ReplaceRequest)
	afterReplace(response *api.ReplaceResponse)
}

// RemoveOption is an option for the Remove method
type RemoveOption interface {
	beforeRemove(request *api.RemoveRequest)
	afterRemove(response *api.RemoveResponse)
}

// IfVersion sets the required version for optimistic concurrency control
func IfVersion(version Version) VersionOption {
	return VersionOption{version: version}
}

// VersionOption is an implementation of SetOption and RemoveOption to specify the version for concurrency control
type VersionOption struct {
	SetOption
	ReplaceOption
	RemoveOption
	version Version
}

func (o VersionOption) beforePut(request *api.PutRequest) {
	request.Version = uint64(o.version)
}

func (o VersionOption) afterPut(response *api.PutResponse) {

}

func (o VersionOption) beforeReplace(request *api.ReplaceRequest) {
	request.PreviousVersion = uint64(o.version)
}

func (o VersionOption) afterReplace(response *api.ReplaceResponse) {

}

func (o VersionOption) beforeRemove(request *api.RemoveRequest) {
	request.Version = uint64(o.version)
}

func (o VersionOption) afterRemove(response *api.RemoveResponse) {

}

// IfNotSet sets the value if the entry is not yet set
func IfNotSet() SetOption {
	return &NotSetOption{}
}

// NotSetOption is a SetOption that sets the value only if it's not already set
type NotSetOption struct {
}

func (o NotSetOption) beforePut(request *api.PutRequest) {
	request.IfEmpty = true
}

func (o NotSetOption) afterPut(response *api.PutResponse) {

}

// GetOption is an option for the Get method
type GetOption interface {
	beforeGet(request *api.GetRequest)
	afterGet(response *api.GetResponse)
}

// WithDefault sets the default value to return if the key is not present
func WithDefault(def []byte) GetOption {
	return defaultOption{def: def}
}

type defaultOption struct {
	def []byte
}

func (o defaultOption) beforeGet(request *api.GetRequest) {
}

func (o defaultOption) afterGet(response *api.GetResponse) {
	if response.Version == 0 {
		response.Value = o.def
	}
}

// WatchOption is an option for the Watch method
type WatchOption interface {
	beforeWatch(request *api.EventRequest)
	afterWatch(response *api.EventResponse)
}

// WithReplay returns a watch option that enables replay of watch events
func WithReplay() WatchOption {
	return replayOption{}
}

type replayOption struct{}

func (o replayOption) beforeWatch(request *api.EventRequest) {
	request.Replay = true
}

func (o replayOption) afterWatch(response *api.EventResponse) {

}

type filterOption struct {
	filter Filter
}

func (o filterOption) beforeWatch(request *api.EventRequest) {
	if o.filter.Key != "" {
		request.Key = o.filter.Key
	}
	if o.filter.Index > 0 {
		request.Index = uint64(o.filter.Index)
	}
}

func (o filterOption) afterWatch(response *api.EventResponse) {
}

// WithFilter returns a watch option that filters the watch events
func WithFilter(filter Filter) WatchOption {
	return filterOption{filter: filter}
}

// Filter is a watch filter configuration
type Filter struct {
	Key   string
	Index Index
}
