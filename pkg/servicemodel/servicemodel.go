// SPDX-FileCopyrightText: 2020-present Open Networking Foundation <info@opennetworking.org>
//
// SPDX-License-Identifier: LicenseRef-ONF-Member-1.0

package servicemodel

import e2 "github.com/onosproject/onos-e2t/pkg/protocols/e2ap101"

// Client service model client interface
type Client interface {
	e2.ClientInterface
}
