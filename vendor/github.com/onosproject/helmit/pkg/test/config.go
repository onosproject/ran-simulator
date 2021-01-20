// Copyright 2019-present Open Networking Foundation.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package test

import (
	"github.com/onosproject/helmit/pkg/job"
	"os"
)

type testType string

const (
	testTypeEnv = "TEST_TYPE"
	testJobType = "test"
)

const (
	testTypeCoordinator testType = "coordinator"
	testTypeWorker      testType = "worker"
)

// Config is a test configuration
type Config struct {
	*job.Config `json:",inline"`
	Suites      []string `json:"suites,omitempty"`
	Tests       []string `json:"tests,omitempty"`
	Iterations  int      `json:"iterations,omitempty"`
	Verbose     bool     `json:"verbose,omitempty"`
	NoTeardown  bool     `json:"verbose,omitempty"`
}

// getTestContext returns the current test context
func getTestType() testType {
	context := os.Getenv(testTypeEnv)
	if context != "" {
		return testType(context)
	}
	return testTypeCoordinator
}
