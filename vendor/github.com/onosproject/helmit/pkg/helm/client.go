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

// nolint
package helm

import (
	"log"

	"github.com/onosproject/helmit/pkg/kubernetes/config"
	"helm.sh/helm/v3/pkg/action"
	"k8s.io/client-go/kubernetes"
)

var clients = make(map[string]HelmClient)

// Namespace returns the Helm namespace
func Namespace() string {
	return config.GetNamespaceFromEnv()
}

// Client returns the Helm client
func Client() HelmClient {
	return getClient(Namespace())
}

// getClient returns the client for the given namespace
func getClient(namespace string) HelmClient {
	client, ok := clients[namespace]
	if !ok {
		configuration, err := getConfig(namespace)
		if err != nil {
			panic(err)
		}
		client = &helmClient{
			namespace: namespace,
			client:    kubernetes.NewForConfigOrDie(config.GetRestConfigOrDie()),
			charts:    make(map[string]*HelmChart),
			config:    configuration,
		}
		clients[namespace] = client
	}
	return client
}

// getConfig gets the Helm configuration for the given namespace
func getConfig(namespace string) (*action.Configuration, error) {
	config := &action.Configuration{}
	if err := config.Init(settings.RESTClientGetter(), namespace, "memory", log.Printf); err != nil {
		return nil, err
	}
	return config, nil
}

// HelmClient is a Helm client
type HelmClient interface {
	HelmChartClient
	HelmReleaseClient

	// Namespace returns the client for the given namespace
	Namespace(namespace string) HelmClient
}

// helmClient is an implementation of the HelmClient interface
type helmClient struct {
	namespace string
	client    *kubernetes.Clientset
	charts    map[string]*HelmChart
	config    *action.Configuration
}

func (c *helmClient) Namespace(namespace string) HelmClient {
	return getClient(namespace)
}

// Charts returns a list of charts in the cluster
func (c *helmClient) Charts() []*HelmChart {
	charts := make([]*HelmChart, 0, len(c.charts))
	for _, chart := range c.charts {
		charts = append(charts, chart)
	}
	return charts
}

// HelmChart returns a chart
func (c *helmClient) Chart(name string, repository ...string) *HelmChart {
	chart, ok := c.charts[name]
	if !ok {
		chart = newChart(name, repository, c.namespace, c.client, c.config)
		c.charts[name] = chart
	}
	return chart
}

// Releases returns a list of releases
func (c *helmClient) Releases() []*HelmRelease {
	releases := make([]*HelmRelease, 0)
	for _, chart := range c.charts {
		releases = append(releases, chart.Releases()...)
	}
	return releases
}

// Release returns the release with the given name
func (c *helmClient) Release(name string) *HelmRelease {
	for _, chart := range c.charts {
		for _, release := range chart.Releases() {
			if release.Name() == name {
				return release
			}
		}
	}
	return nil
}
