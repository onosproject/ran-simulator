// SPDX-FileCopyrightText: 2020-present Open Networking Foundation <info@opennetworking.org>
//
// SPDX-License-Identifier: LicenseRef-ONF-Member-1.0

package e2smkpmies

type NrcgiBuilder interface {
	NewNrcgi()
	SetPlmnID(plmnID PlmnIdentity)
	SetNrcellIdentity(nrCellID NrcellIdentity)
	GetPLmnIdentity()
	GetNRcellIdentity()
	GetNrcgi()
}

func NewNrcgi() *Nrcgi {
	return &Nrcgi{}
}

func (b *Nrcgi) SetPlmnID(plmnID *PlmnIdentity) *Nrcgi {
	b.PLmnIdentity = plmnID
	return b
}

func (b *Nrcgi) SetNrcellIdentity(nrCellID *NrcellIdentity) *Nrcgi {
	b.NRcellIdentity = nrCellID
	return b
}

func (b *Nrcgi) GetNrcgi() *Nrcgi {
	return b
}
