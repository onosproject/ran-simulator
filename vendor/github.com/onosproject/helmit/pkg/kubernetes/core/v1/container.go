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
	"bytes"
	"github.com/onosproject/helmit/pkg/kubernetes/config"
	"github.com/onosproject/helmit/pkg/kubernetes/resource"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/tools/remotecommand"
	executil "k8s.io/client-go/util/exec"
	"strings"
)

// Container is a Kubernetes container
type Container struct {
	resource.Client
	corev1.Container
	pod *corev1.Pod
}

// Exec executes the given command in the container
func (c *Container) Exec(command ...string) (output []string, code int, err error) {
	fullCommand := append([]string{"/bin/bash", "-c"}, command...)
	req := c.Clientset().CoreV1().RESTClient().Post().
		Resource("pods").
		Name(c.pod.Name).
		Namespace(c.pod.Namespace).
		SubResource("exec").
		Param("container", c.Name)
	req.VersionedParams(&corev1.PodExecOptions{
		Container: c.Name,
		Command:   fullCommand,
		Stdout:    true,
		Stderr:    true,
		Stdin:     false,
	}, scheme.ParameterCodec)

	restConfig, err := config.GetRestConfig()
	if err != nil {
		return nil, 0, err
	}

	exec, err := remotecommand.NewSPDYExecutor(restConfig, "POST", req.URL())
	if err != nil {
		if execErr, ok := err.(executil.ExitError); ok && execErr.Exited() {
			return []string{}, execErr.ExitStatus(), nil
		}
		return nil, 0, err
	}

	var stdout, stderr bytes.Buffer
	err = exec.Stream(remotecommand.StreamOptions{
		Stdout: &stdout,
		Stderr: &stderr,
		Tty:    false,
	})

	if err != nil {
		if execErr, ok := err.(executil.ExitError); ok && execErr.Exited() {
			return []string{}, execErr.ExitStatus(), nil
		}
		return nil, 0, err
	}

	return strings.Split(strings.Trim(stdout.String(), "\n"), "\n"), 0, nil
}

// Containers returns a list of containers in the pod
func (p *Pod) Containers() []*Container {
	containers := make([]*Container, len(p.Object.Spec.Containers))
	for i, container := range p.Object.Spec.Containers {
		containers[i] = &Container{
			Client:    p.Resource.Client,
			pod:       p.Object,
			Container: container,
		}
	}
	return containers
}

// Container returns a container by name
func (p *Pod) Container(name string) *Container {
	for _, container := range p.Object.Spec.Containers {
		if container.Name == name {
			return &Container{
				Client:    p.Resource.Client,
				pod:       p.Object,
				Container: container,
			}
		}
	}
	return nil
}
