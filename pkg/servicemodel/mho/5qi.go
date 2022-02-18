// SPDX-FileCopyrightText: 2020-present Open Networking Foundation <info@opennetworking.org>
//
// SPDX-License-Identifier: Apache-2.0

package mho

//ToDo - probably this is redundant - test it out
//func (m *Mho) processFiveQiUpdate(ctx context.Context, subscription *subutils.Subscription) {
//	log.Info("Start processing FiveQi updates")
//	for update := range m.fiveQiUpdateChan {
//		log.Debugf("Received FiveQI Update, IMSI:%v, GnbID:%v, NCGI:%v, FiveQI:%v", update.IMSI, update.Cell.ID, update.Cell.NCGI, update.FiveQi)
//
//		ue, err := m.ServiceModel.UEs.Get(ctx, update.IMSI)
//		if err != nil {
//			log.Warn(err)
//			continue
//		}
//
//		err = m.sendRicIndicationFormat1(ctx, update.Cell.NCGI, ue, subscription)
//		if err != nil {
//			log.Warn(err)
//			continue
//		}
//	}
//}
