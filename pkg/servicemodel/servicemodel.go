// SPDX-FileCopyrightText: 2020-present Open Networking Foundation <info@opennetworking.org>
//
// SPDX-License-Identifier: LicenseRef-ONF-Member-1.0

package servicemodel

import "github.com/onosproject/onos-e2t/pkg/protocols/e2"

// ServiceModel service model interface
type ServiceModel interface {
	e2.ClientInterface
}
