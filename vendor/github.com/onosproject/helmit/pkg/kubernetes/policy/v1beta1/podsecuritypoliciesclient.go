// Code generated by helmit-generate. DO NOT EDIT.

package v1beta1

import (
	"github.com/onosproject/helmit/pkg/kubernetes/resource"
)

type PodSecurityPoliciesClient interface {
	PodSecurityPolicies() PodSecurityPoliciesReader
}

func NewPodSecurityPoliciesClient(resources resource.Client, filter resource.Filter) PodSecurityPoliciesClient {
	return &podSecurityPoliciesClient{
		Client: resources,
		filter: filter,
	}
}

type podSecurityPoliciesClient struct {
	resource.Client
	filter resource.Filter
}

func (c *podSecurityPoliciesClient) PodSecurityPolicies() PodSecurityPoliciesReader {
	return NewPodSecurityPoliciesReader(c.Client, c.filter)
}
