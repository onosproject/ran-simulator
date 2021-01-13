// SPDX-FileCopyrightText: 2020-present Open Networking Foundation <info@opennetworking.org>
//
// SPDX-License-Identifier: LicenseRef-ONF-Member-1.0

package e2agent

import "github.com/onosproject/ran-simulator/api/types"

// E2Agents represents a collection of E2 agents to allow centralized management
type E2Agents struct {
	Agents map[types.ECGI]E2Agent
}
