// Copyright 2024 Eric Cornelissen
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

package gha

import (
	"fmt"
	"io/fs"
	"path"

	"github.com/liamg/memoryfs"
)

type mockFsEntry struct {
	/* File */
	Content []byte

	/* Directory */
	Dir      bool
	Children map[string]mockFsEntry
}

const (
	workflowWithNoJobs     = `name: no jobs`
	workflowWithJobNoSteps = `name: job without steps
jobs:
  no-steps: ~
`
	workflowWithJobWithSteps = `name: job with steps
jobs:
  only-job:
    steps:
      - uses: foo/bar@v1
      - name: no uses
      - uses: foo/baz@v2
`
	workflowWithJobsWithSteps = `name: jobs with steps
jobs:
  job-a:
    steps:
      - uses: foo/bar@v1
  job-b:
    steps:
      - name: no uses
      - uses: foo/baz@v2
`
	workflowWithNestedActions = `name: job using an action that is not at the root
jobs:
  only-job:
    steps:
      - uses: nested/action/1@v1
      - uses: nested/action/2@v1
`
	workflowWithSyntaxError = `Hello world!`
	workflowWithInvalidUses = `name: invalid 'uses' value
jobs:
  job:
    steps:
      - uses: this-is-not-an-action
`
)

func mockRepo(entries map[string]mockFsEntry) (fs.FS, error) {
	repo := memoryfs.New()

	err := repo.MkdirAll(WorkflowsPath, 0o700)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize workflows directory: %v", err)
	}

	err = mockRepoInternal(repo, WorkflowsPath, entries)
	return repo, err
}

func mockRepoInternal(fsys *memoryfs.FS, base string, entries map[string]mockFsEntry) error {
	for name, entry := range entries {
		entryPath := path.Join(base, name)
		if entry.Dir {
			err := fsys.MkdirAll(entryPath, 0o700)
			if err != nil {
				return fmt.Errorf("failed to create dir %q: %v", entryPath, err)
			}

			err = mockRepoInternal(fsys, entryPath, entry.Children)
			if err != nil {
				return err
			}
		} else {
			err := fsys.WriteFile(entryPath, entry.Content, 0o600)
			if err != nil {
				return fmt.Errorf("failed to create %q: %v", entryPath, err)
			}
		}
	}

	return nil
}
