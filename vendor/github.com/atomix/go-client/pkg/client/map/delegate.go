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

package _map //nolint:golint

import (
	"context"
	"github.com/atomix/go-client/pkg/client/primitive"
)

// newDelegatingMap returns a Map that delegates all method calls to the given Map
func newDelegatingMap(_map Map) *delegatingMap {
	return &delegatingMap{
		delegate: _map,
	}
}

// delegatingMap is a Map that delegates method calls to an underlying Map
type delegatingMap struct {
	delegate Map
}

func (m *delegatingMap) Name() primitive.Name {
	return m.delegate.Name()
}

func (m *delegatingMap) Put(ctx context.Context, key string, value []byte, opts ...PutOption) (*Entry, error) {
	return m.delegate.Put(ctx, key, value, opts...)
}

func (m *delegatingMap) Get(ctx context.Context, key string, opts ...GetOption) (*Entry, error) {
	return m.delegate.Get(ctx, key, opts...)
}

func (m *delegatingMap) Remove(ctx context.Context, key string, opts ...RemoveOption) (*Entry, error) {
	return m.delegate.Remove(ctx, key, opts...)
}

func (m *delegatingMap) Len(ctx context.Context) (int, error) {
	return m.delegate.Len(ctx)
}

func (m *delegatingMap) Clear(ctx context.Context) error {
	return m.delegate.Clear(ctx)
}

func (m *delegatingMap) Entries(ctx context.Context, ch chan<- *Entry) error {
	return m.delegate.Entries(ctx, ch)
}

func (m *delegatingMap) Watch(ctx context.Context, ch chan<- *Event, opts ...WatchOption) error {
	return m.delegate.Watch(ctx, ch, opts...)
}

func (m *delegatingMap) Close(ctx context.Context) error {
	return m.delegate.Close(ctx)
}

func (m *delegatingMap) Delete(ctx context.Context) error {
	return m.delegate.Delete(ctx)
}
