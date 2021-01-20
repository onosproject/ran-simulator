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

package v1

import (
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/wait"
	"time"
)

// Wait waits for the Pod to be ready
func (p *Pod) Wait(timeout time.Duration) error {
	return wait.Poll(time.Second, timeout, func() (bool, error) {
		pod, err := p.Clientset().CoreV1().Pods(p.Namespace).Get(p.Name, metav1.GetOptions{})
		if err != nil {
			return false, err
		}
		for _, c := range pod.Status.Conditions {
			if c.Type == corev1.PodReady && c.Status == corev1.ConditionTrue {
				return true, nil
			}
		}
		return false, nil
	})
}

// Wait waits for the Service to be ready
func (s *Service) Wait(timeout time.Duration) error {
	return wait.Poll(time.Second, timeout, func() (bool, error) {
		service, err := s.Clientset().CoreV1().Services(s.Namespace).Get(s.Name, metav1.GetOptions{})
		if err != nil {
			return false, err
		}
		if service.Spec.Type == corev1.ServiceTypeExternalName {
			return true, nil
		}
		if service.Spec.ClusterIP == "" {
			return false, nil
		}
		if service.Spec.Type == corev1.ServiceTypeLoadBalancer {
			if len(service.Spec.ExternalIPs) > 0 {
				return true, nil
			}
			if service.Status.LoadBalancer.Ingress == nil {
				return false, nil
			}
		}
		return true, nil
	})
}
