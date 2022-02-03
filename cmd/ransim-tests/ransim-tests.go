// SPDX-FileCopyrightText: 2020-present Open Networking Foundation <info@opennetworking.org>
//
// SPDX-License-Identifier: Apache-2.0

package main

import (
	"github.com/onosproject/helmit/pkg/registry"
	"github.com/onosproject/helmit/pkg/test"
	"github.com/onosproject/ran-simulator/tests/e2t"
	"github.com/onosproject/ran-simulator/tests/ransim"
)

func main() {
	registry.RegisterTestSuite("e2t", &e2t.TestSuite{})
	registry.RegisterTestSuite("ransim", &ransim.TestSuite{})
	test.Main()
}
