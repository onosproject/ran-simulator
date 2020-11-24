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
	"github.com/hashicorp/golang-lru"
	"sync"
)

// newCachingMap returns a decorated map that caches updates to the given map
func newCachingMap(_map Map, size int) (Map, error) {
	cache, err := lru.New(size)
	if err != nil {
		return nil, err
	}
	cachingMap := &cachingMap{
		delegatingMap: newDelegatingMap(_map),
		pending:       make(map[string]*cachedEntry),
		cache:         cache,
	}
	if err := cachingMap.open(); err != nil {
		return nil, err
	}
	return cachingMap, nil
}

// cachingMap is an implementation of the Map interface that caches entries
type cachingMap struct {
	*delegatingMap
	cancel       context.CancelFunc
	pending      map[string]*cachedEntry
	cache        *lru.Cache
	cacheVersion Version
	mu           sync.RWMutex
}

// open opens the map listeners
func (m *cachingMap) open() error {
	ch := make(chan *Event)
	ctx, cancel := context.WithCancel(context.Background())
	m.mu.Lock()
	m.cancel = cancel
	m.mu.Unlock()
	if err := m.delegatingMap.Watch(ctx, ch, WithReplay()); err != nil {
		return err
	}
	go func() {
		for event := range ch {
			m.cacheUpdate(event.Entry, event.Type == EventRemoved)
		}
	}()
	return nil
}

// cacheUpdate caches the given updated entry
func (m *cachingMap) cacheUpdate(update *Entry, tombstone bool) {
	m.mu.Lock()
	defer m.mu.Unlock()

	// If the update version is less than the cache version, the cache contains
	// more recent updates. Ignore the update.
	if update.Version <= m.cacheVersion {
		return
	}

	// If the pending entry is newer than the update entry, the update can be ignored.
	// Otherwise, remove the entry from the pending cache if present.
	if pending, ok := m.pending[update.Key]; ok {
		if pending.Version > update.Version {
			return
		}
		delete(m.pending, update.Key)
	}

	// If the entry is a tombstone, remove it from the cache, otherwise insert it.
	if tombstone {
		m.cache.Remove(update.Key)
	} else {
		m.cache.Add(update.Key, update)
	}

	// Update the cache version.
	m.cacheVersion = update.Version
}

// cacheRead caches the given read entry
func (m *cachingMap) cacheRead(read *Entry, tombstone bool) {
	m.mu.Lock()
	defer m.mu.Unlock()

	// If the entry version is less than the cache version, ignore the update. The entry will
	// have been cached as an update.
	if read.Version <= m.cacheVersion {
		return
	}

	// The pending cache contains the most recent known state for the entry.
	// If the read entry is newer than the pending entry for the key, update
	// the pending cache.
	if pending, ok := m.pending[read.Key]; !ok || read.Version > pending.Version {
		m.pending[read.Key] = &cachedEntry{
			Entry:     read,
			tombstone: tombstone,
		}
	}
}

// getCache gets a cached entry
func (m *cachingMap) getCache(key string) (*Entry, bool) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	// The pending cache contains the most recent known states. If the entry is present
	// in the pending cache, return it rather than using the LRU cache.
	if entry, ok := m.pending[key]; ok {
		if entry.tombstone {
			return nil, true
		}
		return entry.Entry, true
	}

	// If the entry is present in the LRU cache, return it.
	if entry, ok := m.cache.Get(key); ok {
		return entry.(*Entry), true
	}
	return nil, false
}

func (m *cachingMap) Get(ctx context.Context, key string, opts ...GetOption) (*Entry, error) {
	// If the entry is already in the cache, return it
	if entry, ok := m.getCache(key); ok {
		return entry, nil
	}

	// Otherwise, fetch the entry from the underlying map
	entry, err := m.delegatingMap.Get(ctx, key, opts...)
	if err != nil {
		return nil, err
	}

	// Update the cache if necessary
	if err != nil {
		return nil, err
	}
	m.cacheRead(entry, entry.Value == nil)
	return entry, nil
}

func (m *cachingMap) Put(ctx context.Context, key string, value []byte, opts ...PutOption) (*Entry, error) {
	// Put the entry in the map using the underlying map delegate
	entry, err := m.delegatingMap.Put(ctx, key, value, opts...)
	if err != nil {
		return nil, err
	}

	// Update the cache if necessary
	if err != nil {
		return nil, err
	}
	m.cacheRead(entry, false)
	return entry, nil
}

func (m *cachingMap) Remove(ctx context.Context, key string, opts ...RemoveOption) (*Entry, error) {
	// Remove the entry from the map using the underlying map delegate
	entry, err := m.delegatingMap.Remove(ctx, key, opts...)
	if err != nil {
		return nil, err
	}

	// Update the cache if necessary
	if err != nil {
		return nil, err
	}
	m.cacheRead(entry, true)
	return entry, nil
}

func (m *cachingMap) Close(ctx context.Context) error {
	m.mu.Lock()
	if m.cancel != nil {
		m.cancel()
	}
	m.mu.Unlock()
	return m.delegatingMap.Close(ctx)
}

// cachedEntry is a cached entry
type cachedEntry struct {
	*Entry
	tombstone bool
}
