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

package config

import (
	"github.com/onosproject/helmit/pkg/util/random"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"os"
)

// NamespaceEnv is the environment variable for setting the k8s namespace
const NamespaceEnv = "POD_NAMESPACE"

// GetNamespaceFromEnv gets the Kubernetes namespace from the environment
func GetNamespaceFromEnv() string {
	namespace := os.Getenv(NamespaceEnv)
	if namespace == "" {
		namespace = random.NewPetName(2)
	}
	return namespace
}

// GetRestConfigOrDie returns the Kubernetes REST API configuration
func GetRestConfigOrDie() *rest.Config {
	config, err := GetRestConfig()
	if err != nil {
		panic(err)
	}
	return config
}

// GetRestConfig returns the Kubernetes REST API configuration
func GetRestConfig() (*rest.Config, error) {
	restconfig, err := rest.InClusterConfig()
	if err == nil {
		return restconfig, nil
	}

	kubeconfig := clientcmd.NewNonInteractiveDeferredLoadingClientConfig(
		clientcmd.NewDefaultClientConfigLoadingRules(),
		&clientcmd.ConfigOverrides{},
	)
	return kubeconfig.ClientConfig()
}
