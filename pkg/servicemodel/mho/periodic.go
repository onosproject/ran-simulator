// SPDX-FileCopyrightText: 2020-present Open Networking Foundation <info@opennetworking.org>
//
// SPDX-License-Identifier: Apache-2.0

package mho

import (
	"context"
	"github.com/onosproject/ran-simulator/pkg/store/subscriptions"
	subutils "github.com/onosproject/ran-simulator/pkg/utils/e2ap/subscription"
	"time"
)

func (m *Mho) reportPeriodicIndication(ctx context.Context, interval int32, subscription *subutils.Subscription) {
	log.Debugf("Starting periodic report with interval %d ms", interval)
	subID := subscriptions.NewID(subscription.GetRicInstanceID(), subscription.GetReqID(), subscription.GetRanFuncID())
	intervalDuration := time.Duration(interval)
	sub, err := m.ServiceModel.Subscriptions.Get(subID)
	if err != nil {
		return
	}
	sub.Ticker = time.NewTicker(intervalDuration * time.Millisecond)
	for {
		select {
		case <-sub.Ticker.C:
			log.Debug("Sending periodic indication report for subscription:", sub.ID)
			err = m.sendRicIndication(ctx, subscription)
			if err != nil {
				log.Error("Failure sending indication message: ", err)
			}

		case <-sub.E2Channel.Context().Done():
			sub.Ticker.Stop()
			return
		}
	}
}
