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
	"archive/tar"
	"errors"
	"fmt"
	"github.com/onosproject/helmit/pkg/kubernetes"
	"io"
	"io/ioutil"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/tools/remotecommand"
	"os"
	"path"
	"strings"
)

// Copy returns a new copier
func Copy(client kubernetes.Client) *CopyOptions {
	return &CopyOptions{
		client:    client,
		namespace: client.Namespace(),
	}
}

// CopyOptions is options for copying files from a source to a destination
type CopyOptions struct {
	client    kubernetes.Client
	source    string
	dest      string
	namespace string
	pod       string
	container string
}

// From sets the copy source
func (c *CopyOptions) From(src string) *CopyOptions {
	c.source = src
	return c
}

// To sets the copy destination path
func (c *CopyOptions) To(dest string) *CopyOptions {
	c.dest = dest
	return c
}

// On sets the copy destination pod
func (c *CopyOptions) On(pod string, container ...string) *CopyOptions {
	c.pod = pod
	if len(container) > 0 {
		c.container = container[0]
	}
	return c
}

// Do executes the copy to the pod
func (c *CopyOptions) Do() error {
	if c.source == "" || c.pod == "" {
		return errors.New("source and destination cannot be empty")
	}

	pod, err := c.client.Clientset().CoreV1().Pods(c.client.Namespace()).Get(c.pod, metav1.GetOptions{})
	if err != nil {
		return err
	}

	containerName := c.container
	if len(containerName) == 0 {
		if len(pod.Spec.Containers) > 1 {
			return errors.New("destination container is ambiguous")
		}
		containerName = pod.Spec.Containers[0].Name
	}

	reader, writer := io.Pipe()

	if c.dest == "" {
		c.dest = c.source
	}

	// strip trailing slash (if any)
	if c.source != "/" && strings.HasSuffix(string(c.source[len(c.source)-1]), "/") {
		c.source = c.source[:len(c.source)-1]
	}
	if c.dest != "/" && strings.HasSuffix(string(c.dest[len(c.dest)-1]), "/") {
		c.dest = c.dest[:len(c.dest)-1]
	}

	go func() {
		defer writer.Close()
		err := makeTar(c.source, c.dest, writer)
		if err != nil {
			fmt.Println(err)
		}
	}()

	cmd := []string{"tar", "-xf", "-"}
	req := c.client.Clientset().CoreV1().RESTClient().
		Post().
		Resource("pods").
		Name(c.pod).
		Namespace(c.namespace).
		SubResource("exec").
		VersionedParams(&corev1.PodExecOptions{
			Container: containerName,
			Command:   cmd,
			Stdin:     true,
			Stdout:    true,
			Stderr:    true,
			TTY:       false,
		}, scheme.ParameterCodec)

	exec, err := remotecommand.NewSPDYExecutor(c.client.Config(), "POST", req.URL())
	if err != nil {
		return err
	}
	err = exec.Stream(remotecommand.StreamOptions{
		Stdin:  reader,
		Stdout: os.Stdout,
		Stderr: os.Stderr,
		Tty:    false,
	})
	if err != nil {
		return err
	}
	return nil
}

func makeTar(srcPath, destPath string, writer io.Writer) error {
	// TODO: use compression here?
	tarWriter := tar.NewWriter(writer)
	defer tarWriter.Close()

	srcPath = path.Clean(srcPath)
	destPath = path.Clean(destPath)
	return recursiveTar(path.Dir(srcPath), path.Base(srcPath), path.Dir(destPath), path.Base(destPath), tarWriter)
}

func recursiveTar(srcBase, srcFile, destBase, destFile string, tw *tar.Writer) error {
	filepath := path.Join(srcBase, srcFile)
	stat, err := os.Lstat(filepath)
	if err != nil {
		return err
	}
	if stat.IsDir() {
		files, err := ioutil.ReadDir(filepath)
		if err != nil {
			return err
		}
		if len(files) == 0 {
			//case empty directory
			hdr, _ := tar.FileInfoHeader(stat, filepath)
			hdr.Name = destFile
			if err := tw.WriteHeader(hdr); err != nil {
				return err
			}
		}
		for _, f := range files {
			if err := recursiveTar(srcBase, path.Join(srcFile, f.Name()), destBase, path.Join(destFile, f.Name()), tw); err != nil {
				return err
			}
		}
		return nil
	} else if stat.Mode()&os.ModeSymlink != 0 {
		//case soft link
		hdr, _ := tar.FileInfoHeader(stat, filepath)
		target, err := os.Readlink(filepath)
		if err != nil {
			return err
		}

		hdr.Linkname = target
		hdr.Name = destFile
		if err := tw.WriteHeader(hdr); err != nil {
			return err
		}
	} else {
		//case regular file or other file type like pipe
		hdr, err := tar.FileInfoHeader(stat, filepath)
		if err != nil {
			return err
		}
		hdr.Name = destFile

		if err := tw.WriteHeader(hdr); err != nil {
			return err
		}

		f, err := os.Open(filepath)
		if err != nil {
			return err
		}
		defer f.Close()

		if _, err := io.Copy(tw, f); err != nil {
			return err
		}
		return f.Close()
	}
	return nil
}
