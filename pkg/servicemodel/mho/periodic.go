// SPDX-FileCopyrightText: 2020-present Open Networking Foundation <info@opennetworking.org>
//
// SPDX-License-Identifier: LicenseRef-ONF-Member-1.0

package mho

import (
	"github.com/onosproject/ran-simulator/pkg/store/subscriptions"
	"time"
)

func (m *Mho) reportPeriodicIndication(interval int32) {
	log.Debugf("Starting periodic report with interval %d ms", interval)
	subID := subscriptions.NewID(m.subscription.GetRicInstanceID(), m.subscription.GetReqID(), m.subscription.GetRanFuncID())
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
			err = m.sendRicIndication()
			if err != nil {
				log.Error("Failure sending indication message: ", err)
			}

		case <-sub.E2Channel.Context().Done():
			sub.Ticker.Stop()
			return
		}
	}
}
