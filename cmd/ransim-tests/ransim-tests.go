// SPDX-FileCopyrightText: 2020-present Open Networking Foundation <info@opennetworking.org>
//
// SPDX-License-Identifier: LicenseRef-ONF-Member-1.0

package main

import (
	"github.com/onosproject/helmit/pkg/registry"
	"github.com/onosproject/helmit/pkg/test"
	"github.com/onosproject/ran-simulator/tests/e2t"
)

func main() {
	registry.RegisterTestSuite("e2t", &e2t.TestSuite{})
	test.Main()
}
