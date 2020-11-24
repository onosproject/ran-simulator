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

package util

import (
	"sort"
	"sync"
)

// IterAsync executes the given function f up to n times concurrently.
// Each call is done in a separate goroutine. On each iteration, the function f
// will be called with a unique sequential index i such that the index can be
// used to reference an element in an array or slice. If an error is returned
// by the function f for any index, an error will be returned. Otherwise,
// a nil result will be returned once all function calls have completed.
func IterAsync(n int, f func(i int) error) error {
	wg := sync.WaitGroup{}
	asyncErrors := make(chan error, n)

	wg.Add(n)
	for i := 0; i < n; i++ {
		go func(j int) {
			err := f(j)
			if err != nil {
				asyncErrors <- err
			}
			wg.Done()
		}(i)
	}

	go func() {
		wg.Wait()
		close(asyncErrors)
	}()

	for err := range asyncErrors {
		return err
	}
	return nil
}

// ExecuteAsync executes the given function f up to n times concurrently, populating
// the given results slice with the results of each function call.
// Each call is done in a separate goroutine. On each iteration, the function f
// will be called with a unique sequential index i such that the index can be
// used to reference an element in an array or slice. If an error is returned
// by the function f for any index, an error will be returned. Otherwise,
// a nil result will be returned once all function calls have completed.
func ExecuteAsync(n int, f func(i int) (interface{}, error)) ([]interface{}, error) {
	wg := sync.WaitGroup{}
	asyncErrors := make(chan error, n)
	asyncResults := make(chan interface{}, n)

	wg.Add(n)
	for i := 0; i < n; i++ {
		go func(j int) {
			result, err := f(j)
			if err != nil {
				asyncErrors <- err
			} else {
				asyncResults <- result
			}
			wg.Done()
		}(i)
	}

	go func() {
		wg.Wait()
		close(asyncErrors)
		close(asyncResults)
	}()

	for err := range asyncErrors {
		return nil, err
	}

	results := make([]interface{}, 0, n)
	for result := range asyncResults {
		results = append(results, result)
	}
	return results, nil
}

// ExecuteOrderedAsync executes the given function f up to n times concurrently, populating
// the given results slice with the results of each function call.
// Each call is done in a separate goroutine. On each iteration, the function f
// will be called with a unique sequential index i such that the index can be
// used to reference an element in an array or slice. If an error is returned
// by the function f for any index, an error will be returned. Otherwise,
// a nil result will be returned once all function calls have completed.
func ExecuteOrderedAsync(n int, f func(i int) (interface{}, error)) ([]interface{}, error) {
	wg := sync.WaitGroup{}
	asyncErrors := make(chan error, n)
	asyncResults := make(chan *asyncResult, n)

	wg.Add(n)
	for i := 0; i < n; i++ {
		go func(j int) {
			result, err := f(j)
			if err != nil {
				asyncErrors <- err
			} else {
				asyncResults <- &asyncResult{
					i:      j,
					result: result,
				}
			}
			wg.Done()
		}(i)
	}

	go func() {
		wg.Wait()
		close(asyncErrors)
		close(asyncResults)
	}()

	for err := range asyncErrors {
		return nil, err
	}

	sortedResults := make([]*asyncResult, 0, n)
	for result := range asyncResults {
		sortedResults = append(sortedResults, result)
	}

	sort.Slice(sortedResults, func(i, j int) bool {
		return sortedResults[i].i < sortedResults[j].i
	})

	results := make([]interface{}, n)
	for i, result := range sortedResults {
		results[i] = result.result
	}

	return results, nil
}

type asyncResult struct {
	i      int
	result interface{}
}
