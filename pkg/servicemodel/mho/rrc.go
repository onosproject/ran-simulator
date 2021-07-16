// SPDX-FileCopyrightText: 2020-present Open Networking Foundation <info@opennetworking.org>
//
// SPDX-License-Identifier: LicenseRef-ONF-Member-1.0

package mho

import "context"

func (m *Mho) processRrcUpdate(ctx context.Context) {
	log.Info("Start processing RRC updates")
	for update := range m.rrcUpdateChan {
		log.Debugf("Received RRC Update, IMSI:%v, GnbID:%v, NCGI:%v", update.IMSI, update.Cell.ID, update.Cell.NCGI)

		ue, err := m.ServiceModel.UEs.Get(ctx, update.IMSI)
		if err != nil {
			log.Warn(err)
			continue
		}
		err = m.sendRicIndicationFormat2(ctx, update.Cell.NCGI, ue)
		if err != nil {
			log.Warn(err)
			continue
		}
	}
}
