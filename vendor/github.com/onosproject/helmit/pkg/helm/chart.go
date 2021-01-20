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

// nolint
package helm

import (
	"helm.sh/helm/v3/pkg/action"
	"k8s.io/client-go/kubernetes"
)

// HelmChartClient is a Helm chart client
type HelmChartClient interface {
	// Charts returns a list of charts in the namespace
	Charts() []*HelmChart

	// HelmChart gets a chart in the namespace
	Chart(name string, repository ...string) *HelmChart
}

// Chart returns a Helm chart
func Chart(name string, repository ...string) *HelmChart {
	return Client().Chart(name, repository...)
}

// Charts returns a list of Helm charts
func Charts() []*HelmChart {
	return Client().Charts()
}

func newChart(name string, repo []string, namespace string, client *kubernetes.Clientset, config *action.Configuration) *HelmChart {
	repository := ""
	if len(repo) > 0 {
		repository = repo[0]
	}
	return &HelmChart{
		name:       name,
		repository: repository,
		namespace:  namespace,
		client:     client,
		config:     config,
		releases:   make(map[string]*HelmRelease),
	}
}

// HelmChart is a Helm chart
type HelmChart struct {
	HelmReleaseClient
	namespace  string
	client     *kubernetes.Clientset
	config     *action.Configuration
	name       string
	repository string
	releases   map[string]*HelmRelease
}

// Name returns the chart name
func (c *HelmChart) Name() string {
	return c.name
}

// Repository returns the chart's repository URL
func (c *HelmChart) Repository() string {
	return c.repository
}

// Releases returns a list of releases of the chart
func (c *HelmChart) Releases() []*HelmRelease {
	releases := make([]*HelmRelease, 0, len(c.releases))
	for _, release := range c.releases {
		releases = append(releases, release)
	}
	return releases
}

// Release returns the release with the given name
func (c *HelmChart) Release(name string) *HelmRelease {
	release, ok := c.releases[name]
	if !ok {
		release = newRelease(name, c.namespace, c.client, c, c.config)
		c.releases[name] = release
	}
	return release
}
