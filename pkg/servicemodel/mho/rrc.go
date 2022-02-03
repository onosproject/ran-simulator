// SPDX-FileCopyrightText: 2020-present Open Networking Foundation <info@opennetworking.org>
//
// SPDX-License-Identifier: Apache-2.0

package mho

import (
	"context"
	subutils "github.com/onosproject/ran-simulator/pkg/utils/e2ap/subscription"
)

func (m *Mho) processRrcUpdate(ctx context.Context, subscription *subutils.Subscription) {
	log.Info("Start processing RRC updates")
	for update := range m.rrcUpdateChan {
		log.Debugf("Received RRC Update, IMSI:%v, GnbID:%v, NCGI:%v", update.IMSI, update.Cell.ID, update.Cell.NCGI)

		ue, err := m.ServiceModel.UEs.Get(ctx, update.IMSI)
		if err != nil {
			log.Warn(err)
			continue
		}
		err = m.sendRicIndicationFormat2(ctx, update.Cell.NCGI, ue, subscription)
		if err != nil {
			log.Warn(err)
			continue
		}
	}
}
