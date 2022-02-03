// SPDX-FileCopyrightText: 2020-present Open Networking Foundation <info@opennetworking.org>
//
// SPDX-License-Identifier: Apache-2.0

package utils

import (
	meastype "github.com/onosproject/rrm-son-lib/pkg/model/measurement/type"
	"sort"
)

// QOffsetRange is a struct for QOffsetRanges element
type QOffsetRange struct {
	Min   int32
	Max   int32
	Value meastype.QOffsetRange
}

// QOffsetRanges is the list type of QOffsetRange
type QOffsetRanges []QOffsetRange

// Len returns length
func (q QOffsetRanges) Len() int {
	return len(q)
}

// Less checks if i is less than j
func (q QOffsetRanges) Less(i, j int) bool {
	return q[i].Min < q[j].Min
}

// Swap swaps two values
func (q QOffsetRanges) Swap(i, j int) {
	q[i], q[j] = q[j], q[i]
}

// Sort sorts the list
func (q QOffsetRanges) Sort() {
	sort.Sort(q)
}

// Search returns an appropriate value satisfying a condition
func (q QOffsetRanges) Search(v int32) meastype.QOffsetRange {
	length := q.Len()
	if i := sort.Search(length, func(i int) bool { return v < q[i].Max }); i < length {
		if it := &q[i]; v >= it.Min && v < it.Max {
			return it.Value
		}
	}
	return meastype.QOffset0dB
}
