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

package files

import (
	"errors"
	"fmt"
	"os"

	"github.com/onosproject/helmit/pkg/kubernetes"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/tools/remotecommand"
)

// Echo returns a new echo client
func Echo(client kubernetes.Client) *EchoOptions {
	return &EchoOptions{
		client:    client,
		namespace: client.Namespace(),
	}
}

// EchoOptions is options for echoing output to a file
type EchoOptions struct {
	client    kubernetes.Client
	namespace string
	pod       string
	container string
	file      string
	bytes     []byte
}

// Bytes sets the bytes to echo
func (o *EchoOptions) Bytes(bytes []byte) *EchoOptions {
	o.bytes = bytes
	return o
}

// Contents sets the contents to echo
func (o *EchoOptions) String(s string) *EchoOptions {
	return o.Bytes([]byte(s))
}

// To configures the file to which to echo
func (o *EchoOptions) To(filename string) *EchoOptions {
	o.file = filename
	return o
}

// On configures the copy destination pod
func (o *EchoOptions) On(pod string, container ...string) *EchoOptions {
	o.pod = pod
	if len(container) > 0 {
		o.container = container[0]
	}
	return o
}

// Do executes the copy to the pod
func (o *EchoOptions) Do() error {
	if o.pod == "" || o.file == "" {
		return errors.New("target file cannot be empty")
	}

	pod, err := o.client.Clientset().CoreV1().Pods(o.client.Namespace()).Get(o.pod, metav1.GetOptions{})
	if err != nil {
		return err
	}

	containerName := o.container
	if len(containerName) == 0 {
		if len(pod.Spec.Containers) > 1 {
			return errors.New("destination container is ambiguous")
		}
		containerName = pod.Spec.Containers[0].Name
	}

	cmd := []string{"/bin/sh", "-c", fmt.Sprintf("echo \"%s\" > %s", string(o.bytes), o.file)}
	req := o.client.Clientset().CoreV1().RESTClient().
		Post().
		Resource("pods").
		Name(o.pod).
		Namespace(o.namespace).
		SubResource("exec").
		VersionedParams(&corev1.PodExecOptions{
			Container: containerName,
			Command:   cmd,
			Stdin:     false,
			Stdout:    true,
			Stderr:    true,
			TTY:       false,
		}, scheme.ParameterCodec)

	exec, err := remotecommand.NewSPDYExecutor(o.client.Config(), "POST", req.URL())
	if err != nil {
		return err
	}
	err = exec.Stream(remotecommand.StreamOptions{
		Stdout: os.Stdout,
		Stderr: os.Stderr,
		Tty:    false,
	})
	if err != nil {
		return err
	}
	return nil
}
