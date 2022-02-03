// SPDX-FileCopyrightText: 2020-present Open Networking Foundation <info@opennetworking.org>
//
// SPDX-License-Identifier: Apache-2.0

package utils

import (
	meastype "github.com/onosproject/rrm-son-lib/pkg/model/measurement/type"
	"sort"
)

// TimeToTriggerRange is a struct for TimeToTriggerRanges element
type TimeToTriggerRange struct {
	Min   int32
	Max   int32
	Value meastype.TimeToTriggerRange
}

// TimeToTriggerRanges is the list type of TimeToTriggerRange
type TimeToTriggerRanges []TimeToTriggerRange

// Len returns length
func (t TimeToTriggerRanges) Len() int {
	return len(t)
}

// Less checks if i is less than j
func (t TimeToTriggerRanges) Less(i, j int) bool {
	return t[i].Min < t[j].Min
}

// Swap swaps two values
func (t TimeToTriggerRanges) Swap(i, j int) {
	t[i], t[j] = t[j], t[i]
}

// Sort sorts the list
func (t TimeToTriggerRanges) Sort() {
	sort.Sort(t)
}

// Search returns an appropriate value satisfying a condition
func (t TimeToTriggerRanges) Search(v int32) meastype.TimeToTriggerRange {
	length := t.Len()
	if i := sort.Search(length, func(i int) bool { return v < t[i].Max }); i < length {
		if it := &t[i]; v >= it.Min && v < it.Max {
			return it.Value
		}
	}
	return meastype.TTT0ms
}
