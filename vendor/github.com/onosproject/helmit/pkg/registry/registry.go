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

package registry

var tests = make(map[string]interface{})
var benchmarks = make(map[string]interface{})
var simulations = make(map[string]interface{})

// RegisterTestSuite registers a test suite
func RegisterTestSuite(name string, suite interface{}) {
	tests[name] = suite
}

// GetTestSuites returns a list of registered tests
func GetTestSuites() []string {
	names := make([]string, 0, len(tests))
	for name := range tests {
		names = append(names, name)
	}
	return names
}

// GetTestSuite gets a registered test suite by name
func GetTestSuite(name string) interface{} {
	return tests[name]
}

// RegisterBenchmarkSuite registers a benchmark suite
func RegisterBenchmarkSuite(name string, suite interface{}) {
	benchmarks[name] = suite
}

// GetBenchmarkSuites returns a list of registered benchmark suites
func GetBenchmarkSuites() []string {
	names := make([]string, 0, len(benchmarks))
	for name := range benchmarks {
		names = append(names, name)
	}
	return names
}

// GetBenchmarkSuite gets a registered simulation by name
func GetBenchmarkSuite(name string) interface{} {
	return benchmarks[name]
}

// RegisterSimulationSuite registers a simulation suite
func RegisterSimulationSuite(name string, suite interface{}) {
	simulations[name] = suite
}

// GetSimulationSuites returns a list of registered simulation suites
func GetSimulationSuites() []string {
	names := make([]string, 0, len(simulations))
	for name := range simulations {
		names = append(names, name)
	}
	return names
}

// GetSimulationSuite gets a registered simulation suite by name
func GetSimulationSuite(name string) interface{} {
	return simulations[name]
}
