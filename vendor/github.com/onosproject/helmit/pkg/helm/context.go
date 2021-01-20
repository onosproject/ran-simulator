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

package helm

import "path/filepath"

var context = &Context{}

// SetContext sets the Helm context
func SetContext(ctx *Context) error {
	ctxWorkDir := ctx.WorkDir
	if ctxWorkDir != "" {
		absDir, err := filepath.Abs(ctxWorkDir)
		if err != nil {
			return err
		}
		ctxWorkDir = absDir
	}

	ctxValueFiles := make(map[string][]string)
	for release, valueFiles := range ctx.ValueFiles {
		cleanValueFiles := make([]string, 0)
		for _, valueFile := range valueFiles {
			absPath, err := filepath.Abs(valueFile)
			if err != nil {
				return err
			}
			cleanValueFiles = append(cleanValueFiles, absPath)
		}
		ctxValueFiles[release] = cleanValueFiles
	}

	context = &Context{
		WorkDir:    ctxWorkDir,
		Values:     ctx.Values,
		ValueFiles: ctxValueFiles,
	}
	return nil
}

// Context is a Helm context
type Context struct {
	// WorkDir is the Helm working directory
	WorkDir string

	// Values is a mapping of release values
	Values map[string][]string

	// ValueFiles is a mapping of release value files
	ValueFiles map[string][]string
}

// Release returns the context for the given release
func (c *Context) Release(name string) *ReleaseContext {
	return &ReleaseContext{
		Values:     c.Values[name],
		ValueFiles: c.ValueFiles[name],
	}
}

// ReleaseContext is a Helm release context
type ReleaseContext struct {
	// ValueFiles is the release value files
	ValueFiles []string

	// Values is the release values
	Values []string
}
