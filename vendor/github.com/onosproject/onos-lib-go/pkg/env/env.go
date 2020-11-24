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

package env

import "os"

// PodNamespace is the name of the environment variable containing the pod namespace
const PodNamespace = "POD_NAMESPACE"

// GetPodNamespace gets the pod namespace from the environment
func GetPodNamespace() string {
	return os.Getenv(PodNamespace)
}

// PodName is the name of the environment variable containing the pod name
const PodName = "POD_NAME"

// GetPodName gets the pod name from the environment
func GetPodName() string {
	return os.Getenv(PodName)
}

// PodID is the name of the environment variable containing the pod network identifier
const PodID = "POD_ID"

// GetPodID gets the pod network identifier from the environment
func GetPodID() string {
	return os.Getenv(PodID)
}

// PodIP is the name of the environment variable containing the pod IP address
const PodIP = "POD_IP"

// GetPodIP gets the pod IP address from the environment
func GetPodIP() string {
	return os.Getenv(PodIP)
}

// ServiceNamespace is the name of the environment variable containing the service namespace
const ServiceNamespace = "SERVICE_NAMESPACE"

// GetServiceNamespace gets the service namespace from the environment
func GetServiceNamespace() string {
	return os.Getenv(ServiceNamespace)
}

// ServiceName is the name of the environment variable containing the service name
const ServiceName = "SERVICE_NAME"

// GetServiceName gets the service name from the environment
func GetServiceName() string {
	return os.Getenv(ServiceName)
}
