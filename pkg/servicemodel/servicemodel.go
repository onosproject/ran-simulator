// SPDX-FileCopyrightText: 2020-present Open Networking Foundation <info@opennetworking.org>
//
// SPDX-License-Identifier: Apache-2.0

package servicemodel

import e2 "github.com/onosproject/onos-e2t/pkg/protocols/e2ap"

// Client service model client interface
type Client interface {
	e2.ClientInterface
}
