// SPDX-FileCopyrightText: 2020-present Open Networking Foundation <info@opennetworking.org>
//
// SPDX-License-Identifier: Apache-2.0

package subscriptions

import (
	"testing"

	e2apies "github.com/onosproject/onos-e2t/api/e2ap/v2/e2ap-ies"

	"github.com/stretchr/testify/assert"
)

// TestSubscriptions test subscriptions store interface
func TestSubscriptions(t *testing.T) {
	subStore := NewStore()
	numSubs, err := subStore.Len()
	assert.NoError(t, err)
	assert.Equal(t, 0, numSubs)
	sub1 := &Subscription{
		ID: "sub1",
		ReqID: &e2apies.RicrequestId{
			RicRequestorId: 1,
			RicInstanceId:  1,
		},
	}

	err = subStore.Add(sub1)
	assert.NoError(t, err)
	numSubs, err = subStore.Len()
	assert.NoError(t, err)
	assert.Equal(t, 1, numSubs)

	sub2 := &Subscription{
		ID: "sub2",
		ReqID: &e2apies.RicrequestId{
			RicRequestorId: 2,
			RicInstanceId:  2,
		},
	}
	err = subStore.Add(sub2)
	assert.NoError(t, err)
	numSubs, err = subStore.Len()
	assert.NoError(t, err)
	assert.Equal(t, 2, numSubs)

	_, err = subStore.Get("sub3")
	assert.Error(t, err, "subscription entry has not been found")

	sub1Entry, err := subStore.Get("sub1")
	assert.NoError(t, err)
	assert.Equal(t, ID("sub1"), sub1Entry.ID)

	err = subStore.Remove("sub1")
	assert.NoError(t, err)
	numSubs, err = subStore.Len()
	assert.NoError(t, err)
	assert.Equal(t, 1, numSubs)

	_, err = subStore.Get("sub1")
	assert.Error(t, err, "subscription entry has not been found")

	subscriptionList, err := subStore.List()
	assert.NoError(t, err)
	numSubs, err = subStore.Len()
	assert.Equal(t, 1, numSubs)
	assert.NoError(t, err)
	assert.Equal(t, 1, len(subscriptionList))

}
