// Copyright 2020-present Open Networking Foundation.
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

package onostest

const (
	AtomixChartRepo                 = "https://charts.atomix.io"
	OnosChartRepo                   = "https://charts.onosproject.org"
	SdranChartRepo                  = "https://sdrancharts.onosproject.org"
	AtomixControllerPort            = "5679"
	SecretsName                     = "helmit-secrets"
	ControllerChartName             = "atomix-controller"
	RaftStorageControllerChartName  = "raft-storage-controller"
	CacheStorageControllerChartName = "cache-storage-controller"
)

func AtomixName(testName string, componentName string) string {
	return testName + "-" + componentName + "-atomix"
}

func AtomixControllerName(testName string, componentName string) string {
	return AtomixName(testName, componentName) + "-atomix-controller"
}

func AtomixController(testName string, componentName string) string {
	return AtomixControllerName(testName, componentName) + ":" + AtomixControllerPort
}

func RaftReleaseName(componentName string) string {
	return componentName + "-raft"
}

func CacheReleaseName(componentName string) string {
	return componentName + "-cache"
}
