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
	"errors"
	"github.com/atomix/go-client/pkg/client/primitive"
)

// slicedList is a slice of a list
type slicedList struct {
	from *int
	to   *int
	list List
}

func (l *slicedList) Name() primitive.Name {
	return l.list.Name()
}

func (l *slicedList) inRangeIndex(index int) bool {
	return (l.from == nil || index >= *l.from) && (l.to == nil || index < *l.to)
}

func (l *slicedList) Append(ctx context.Context, value []byte) error {
	return errors.New("cannot append to list slice")
}

func (l *slicedList) Insert(ctx context.Context, index int, value []byte) error {
	if l.from != nil {
		index += *l.from
	}
	if !l.inRangeIndex(index) {
		return errors.New("index out of slice range")
	}
	return l.list.Insert(ctx, index, value)
}

func (l *slicedList) Set(ctx context.Context, index int, value []byte) error {
	if l.from != nil {
		index += *l.from
	}
	if !l.inRangeIndex(index) {
		return errors.New("index out of slice range")
	}
	return l.list.Set(ctx, index, value)
}

func (l *slicedList) Get(ctx context.Context, index int) ([]byte, error) {
	if l.from != nil {
		index += *l.from
	}
	if !l.inRangeIndex(index) {
		return nil, errors.New("index out of slice range")
	}
	return l.list.Get(ctx, index)
}

func (l *slicedList) Remove(ctx context.Context, index int) ([]byte, error) {
	if l.from != nil {
		index += *l.from
	}
	if !l.inRangeIndex(index) {
		return nil, errors.New("index out of slice range")
	}
	return l.list.Remove(ctx, index)
}

func (l *slicedList) Len(ctx context.Context) (int, error) {
	size, err := l.list.Len(ctx)
	if err != nil {
		return 0, err
	}
	if l.to != nil && *l.to < size {
		size = *l.to
	}
	if l.from != nil {
		if *l.from > size {
			return 0, nil
		}
		size -= *l.from
	}
	return size, nil
}

func (l *slicedList) Slice(ctx context.Context, from int, to int) (List, error) {
	if l.from != nil {
		from += *l.from
		to += *l.from
	}
	return &slicedList{
		from: &from,
		to:   &to,
		list: l.list,
	}, nil
}

func (l *slicedList) SliceFrom(ctx context.Context, from int) (List, error) {
	if l.from != nil {
		from += *l.from
	}
	return &slicedList{
		from: &from,
		list: l.list,
	}, nil
}

func (l *slicedList) SliceTo(ctx context.Context, to int) (List, error) {
	if l.from != nil {
		to += *l.from
	}
	return &slicedList{
		to:   &to,
		list: l.list,
	}, nil
}

func (l *slicedList) Items(ctx context.Context, ch chan<- []byte) error {
	itemsCh := make(chan []byte)
	go func() {
		// TODO: This method should not have to filter all the elements in the list to get to the relevant elements!
		i := 0
		for item := range itemsCh {
			if l.inRangeIndex(i) {
				ch <- item
			}
			i++
			if l.to != nil && i == *l.to {
				close(ch)
				return
			}
		}
		close(ch)
	}()
	return l.list.Items(ctx, itemsCh)
}

func (l *slicedList) Watch(ctx context.Context, ch chan<- *Event, opts ...WatchOption) error {
	eventCh := make(chan *Event)
	go func() {
		for event := range eventCh {
			if (l.from == nil || *l.from >= event.Index) && (l.to == nil || event.Index < *l.to) {
				ch <- event
			}
		}
	}()
	return l.list.Watch(ctx, eventCh, opts...)
}

func (l *slicedList) Clear(ctx context.Context) error {
	return errors.New("cannot clear list slice")
}

func (l *slicedList) Close(ctx context.Context) error {
	return l.list.Close(ctx)
}

func (l *slicedList) Delete(ctx context.Context) error {
	return errors.New("cannot delete list slice")
}
