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
	"github.com/onosproject/helmit/pkg/helm"
	"github.com/onosproject/helmit/pkg/registry"
	"google.golang.org/grpc"
	"net"
	"os"
	"testing"
)

// newWorker returns a new test worker
func newWorker(config *Config) (*Worker, error) {
	return &Worker{
		config: config,
	}, nil
}

// Worker runs a test job
type Worker struct {
	config *Config
}

// Run runs a benchmark
func (w *Worker) Run() error {
	err := helm.SetContext(&helm.Context{
		WorkDir:    w.config.Context,
		Values:     w.config.Values,
		ValueFiles: w.config.ValueFiles,
	})
	if err != nil {
		return err
	}

	lis, err := net.Listen("tcp", ":5000")
	if err != nil {
		return err
	}
	server := grpc.NewServer()
	RegisterWorkerServiceServer(server, w)
	return server.Serve(lis)
}

// RunTests runs a suite of tests
func (w *Worker) RunTests(ctx context.Context, request *TestRequest) (*TestResponse, error) {
	go w.runTests(request)
	return &TestResponse{}, nil
}

func (w *Worker) runTests(request *TestRequest) {
	test := registry.GetTestSuite(request.Suite)
	if test == nil {
		fmt.Println(fmt.Errorf("unknown test suite %s", request.Suite))
		os.Exit(1)
	}

	tests := []testing.InternalTest{
		{
			Name: request.Suite,
			F: func(t *testing.T) {
				RunTests(t, test, request.Tests)
			},
		},
	}

	// Hack to enable verbose testing.
	os.Args = []string{
		os.Args[0],
		"-test.v",
	}

	testing.Main(func(_, _ string) (bool, error) { return true, nil }, tests, nil, nil)
}
