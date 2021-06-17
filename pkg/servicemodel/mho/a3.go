// SPDX-FileCopyrightText: 2020-present Open Networking Foundation <info@opennetworking.org>
//
// SPDX-License-Identifier: LicenseRef-ONF-Member-1.0

package mho

import (
	"github.com/onosproject/onos-api/go/onos/ransim/types"
	ransimtypes "github.com/onosproject/onos-api/go/onos/ransim/types"
	"github.com/onosproject/rrm-son-lib/pkg/model/id"
)

func (m *Mho) processEventA3MeasReport() {
	log.Info("Start processing event a3 measurement report")
	for report := range m.ServiceModel.MeasChan {
		log.Debugf("received event a3 measurement report: %v", report)
		log.Debugf("Send upon-rcv-meas-report indication for cell ecgi:%d, IMSI:%s",
			report.GetSCell().GetID().GetID().(id.ECGI), report.GetID().String())
		ecgi := report.GetSCell().GetID().GetID().(id.ECGI)
		imsi := report.GetID().GetID().(id.UEID).IMSI
		ue, err := m.ServiceModel.UEs.Get(m.context, types.IMSI(imsi))
		if err != nil {
			log.Warn(err)
			continue
		}
		err = m.sendRicIndicationFormat1(ransimtypes.NCGI(ecgi), ue)
		if err != nil {
			log.Warn(err)
			continue
		}
	}
}
