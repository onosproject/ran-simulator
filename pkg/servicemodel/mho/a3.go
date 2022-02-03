// SPDX-FileCopyrightText: 2020-present Open Networking Foundation <info@opennetworking.org>
//
// SPDX-License-Identifier: Apache-2.0

package mho

import (
	"context"
	"github.com/onosproject/onos-api/go/onos/ransim/types"
	ransimtypes "github.com/onosproject/onos-api/go/onos/ransim/types"
	"github.com/onosproject/ran-simulator/pkg/store/subscriptions"
	subutils "github.com/onosproject/ran-simulator/pkg/utils/e2ap/subscription"
	"github.com/onosproject/rrm-son-lib/pkg/model/id"
)

func (m *Mho) processEventA3MeasReport(ctx context.Context, subscription *subutils.Subscription) {
	log.Info("Start processing event a3 measurement report")
	subID := subscriptions.NewID(subscription.GetRicInstanceID(), subscription.GetReqID(), subscription.GetRanFuncID())
	sub, err := m.ServiceModel.Subscriptions.Get(subID)
	if err != nil {
		log.Error(err)
		return
	}
	for {
		select {
		case report := <-m.ServiceModel.A3Chan:
			log.Debugf("received event a3 measurement report: %v", report)
			log.Debugf("Send upon-rcv-meas-report indication for cell ecgi:%d, IMSI:%s",
				report.UE.GetSCell().GetID().GetID().(id.ECGI), report.UE.GetID().String())
			ecgi := report.UE.GetSCell().GetID().GetID().(id.ECGI)
			imsi := report.UE.GetID().GetID().(id.UEID).IMSI
			ue, err := m.ServiceModel.UEs.Get(ctx, types.IMSI(imsi))
			if err != nil {
				log.Warn(err)
				continue
			}
			err = m.sendRicIndicationFormat1(ctx, ransimtypes.NCGI(ecgi), ue, subscription)
			if err != nil {
				log.Warn(err)
				continue
			}
		case <-sub.E2Channel.Context().Done():
			sub.Ticker.Stop()
			return
		}
	}
}
