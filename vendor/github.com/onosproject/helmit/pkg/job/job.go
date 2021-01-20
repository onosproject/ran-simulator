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

package job

import (
	"encoding/json"
	"io/ioutil"
	corev1 "k8s.io/api/core/v1"
	"os"
	"path/filepath"
	"time"
)

const configPath = "/etc/helmit"
const configFile = "job.json"
const readyFile = "/tmp/job-ready"

// Config is a job configuration
type Config struct {
	ID              string
	Namespace       string
	ServiceAccount  string
	Image           string
	ImagePullPolicy corev1.PullPolicy
	Executable      string
	Context         string
	Values          map[string][]string
	ValueFiles      map[string][]string
	Args            []string
	Env             map[string]string
	Timeout         time.Duration
	NoTeardown      bool
	Secrets         map[string]string
}

// Job is a job configuration
type Job struct {
	*Config
	JobConfig interface{}
	Type      string
}

// Bootstrap bootstraps the job
func Bootstrap(config interface{}) error {
	awaitReady()
	return LoadConfig(config)
}

// awaitReady waits for the job to become ready
func awaitReady() {
	for {
		if isReady() {
			return
		}
		time.Sleep(time.Second)
	}
}

// isReady checks if the job is ready
func isReady() bool {
	info, err := os.Stat(readyFile)
	return err == nil && !info.IsDir()
}

// LoadConfig returns the job configuration
func LoadConfig(config interface{}) error {
	file, err := os.Open(filepath.Join(configPath, configFile))
	if err != nil {
		return err
	}
	defer file.Close()
	bytes, err := ioutil.ReadAll(file)
	if err != nil {
		return err
	}
	err = json.Unmarshal(bytes, config)
	if err != nil {
		return err
	}
	return nil
}
