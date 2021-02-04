// SPDX-FileCopyrightText: 2020-present Open Networking Foundation <info@opennetworking.org>
//
// SPDX-License-Identifier: LicenseRef-ONF-Member-1.0

package e2agent

import (
	"bytes"
	"encoding/gob"
	"hash/fnv"
	"time"

	"github.com/cenkalti/backoff"
	"github.com/onosproject/ran-simulator/pkg/model"
)

const (
	backoffInterval = 10 * time.Millisecond
	maxBackoffTime  = 5 * time.Second
)

func newExpBackoff() *backoff.ExponentialBackOff {
	b := backoff.NewExponentialBackOff()
	b.InitialInterval = backoffInterval
	// MaxInterval caps the RetryInterval
	b.MaxInterval = maxBackoffTime
	// Never stops retrying
	b.MaxElapsedTime = 0
	return b
}

func nodeID(plmndID model.PlmnID, enbID model.EnbID) (uint64, error) {
	gEnbID := model.GEnbID{
		PlmnID: plmndID,
		EnbID:  enbID,
	}
	buf := bytes.Buffer{}
	enc := gob.NewEncoder(&buf)
	err := enc.Encode(gEnbID)
	if err != nil {
		return 0, err
	}

	h := fnv.New64a()
	_, _ = h.Write(buf.Bytes())
	return h.Sum64(), nil
}
