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

package test

import (
	"context"
	"fmt"
	"strconv"

	"github.com/onosproject/helmit/pkg/job"
	"github.com/onosproject/helmit/pkg/registry"
	"google.golang.org/grpc"
)

// newCoordinator returns a new test coordinator
func newCoordinator(config *Config) (*Coordinator, error) {
	return &Coordinator{
		config: config,
		runner: job.NewNamespace(config.Namespace),
	}, nil
}

// Coordinator coordinates workers for suites of tests
type Coordinator struct {
	config *Config
	runner *job.Runner
}

// Run runs the tests
func (c *Coordinator) Run() (int, error) {
	var returnCode int
	for iteration := 1; iteration <= c.config.Iterations || c.config.Iterations < 0; iteration++ {
		suites := c.config.Suites
		if len(suites) == 0 || suites[0] == "" {
			suites = registry.GetTestSuites()
		}
		returnCode = 0
		for _, suite := range suites {
			jobID := newJobID(c.config.ID+"-"+strconv.Itoa(iteration), suite)
			env := c.config.Env
			if env == nil {
				env = make(map[string]string)
			}
			env[testTypeEnv] = string(testTypeWorker)
			config := &Config{
				Config: &job.Config{
					ID:              jobID,
					Namespace:       c.config.Config.Namespace,
					ServiceAccount:  c.config.Config.ServiceAccount,
					Image:           c.config.Config.Image,
					ImagePullPolicy: c.config.Config.ImagePullPolicy,
					Executable:      c.config.Config.Executable,
					Context:         c.config.Config.Context,
					Values:          c.config.Config.Values,
					ValueFiles:      c.config.Config.ValueFiles,
					Env:             env,
					Timeout:         c.config.Config.Timeout,
					NoTeardown:      c.config.Config.NoTeardown,
					Secrets:         c.config.Config.Secrets,
				},
				Suites:     []string{suite},
				Tests:      c.config.Tests,
				Iterations: c.config.Iterations,
			}
			task := &WorkerTask{
				runner: c.runner,
				config: config,
			}
			status, err := task.Run()
			if err != nil {
				return status, err
			} else if returnCode == 0 {
				returnCode = status
			}
		}
		if returnCode == 0 {
			return 0, nil
		}
	}
	return returnCode, nil
}

// newJobID returns a new unique test job ID
func newJobID(testID, suite string) string {
	return fmt.Sprintf("%s-%s", testID, suite)
}

// WorkerTask manages a single test job for a test worker
type WorkerTask struct {
	runner *job.Runner
	config *Config
}

// Run runs the worker job
func (t *WorkerTask) Run() (int, error) {
	job := &job.Job{
		Config:    t.config.Config,
		JobConfig: t.config,
		Type:      testJobType,
	}

	err := t.runner.StartJob(job)
	if err != nil {
		return 0, err
	}

	address := fmt.Sprintf("%s:5000", job.ID)
	conn, err := grpc.Dial(address, grpc.WithInsecure())
	if err != nil {
		return 0, err
	}
	client := NewWorkerServiceClient(conn)
	_, err = client.RunTests(context.Background(), &TestRequest{
		Suite: t.config.Suites[0],
		Tests: t.config.Tests,
	})

	if err != nil {
		return 0, err
	}

	status, err := t.runner.WaitForExit(job)
	if err != nil {
		return 0, err
	}
	return status, err
}
