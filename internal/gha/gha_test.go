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
	"slices"
	"testing"

	"github.com/liamg/memoryfs"
)

func TestNoWorkflows(t *testing.T) {
	t.Parallel()

	repo := memoryfs.New()
	if _, err := RepoActions(repo); err == nil {
		t.Fatal("Unexpected success")
	}
}

func TestFaultyWorkflow(t *testing.T) {
	t.Parallel()

	workflows := map[string]mockFsEntry{
		"workflow.yaml": {
			Content: []byte(workflowWithJobWithSteps),
		},
		"syntax-error.yml": {
			Content: []byte(workflowWithSyntaxError),
		},
	}

	repo, err := mockRepo(workflows)
	if err != nil {
		t.Fatalf("Could not initialize file system: %+v", err)
	}

	if _, err := RepoActions(repo); err == nil {
		t.Fatal("Unexpected success")
	}
}

func TestFaultyUses(t *testing.T) {
	t.Parallel()

	workflows := map[string]mockFsEntry{
		"workflow.yaml": {
			Content: []byte(workflowWithJobWithSteps),
		},
		"invalid-uses.yml": {
			Content: []byte(workflowWithInvalidUses),
		},
	}

	repo, err := mockRepo(workflows)
	if err != nil {
		t.Fatalf("Could not initialize file system: %+v", err)
	}

	if _, err := RepoActions(repo); err == nil {
		t.Fatal("Unexpected success")
	}
}

func TestRealisticRepository(t *testing.T) {
	t.Parallel()

	workflows := map[string]mockFsEntry{
		"nested": {
			Dir: true,
			Children: map[string]mockFsEntry{
				"foo.bar": {
					Content: []byte("foobar"),
				},
			},
		},
		"not-a-workflow.txt": {
			Content: []byte("Hello world!"),
		},
		"one-job.yaml": {
			Content: []byte(workflowWithJobWithSteps),
		},
		"multiple-jobs.yml": {
			Content: []byte(workflowWithJobsWithSteps),
		},
		"nested-action.yml": {
			Content: []byte(workflowWithNestedActions),
		},
	}

	repo, err := mockRepo(workflows)
	if err != nil {
		t.Fatalf("Could not initialize file system: %+v", err)
	}

	got, err := RepoActions(repo)
	if err != nil {
		t.Fatalf("Unexpected error: %+v", err)
	}

	want := []GitHubAction{
		{
			Owner:   "foo",
			Project: "bar",
			Ref:     "v1",
		},
		{
			Owner:   "foo",
			Project: "baz",
			Ref:     "v2",
		},
		{
			Owner:   "nested",
			Project: "action",
			Ref:     "v1",
		},
	}

	if got, want := len(got), len(want); got != want {
		t.Errorf("Incorrect result length (got %d, want %d)", got, want)
	}

	for _, got := range got {
		if !slices.Contains(want, got) {
			t.Errorf("Unwanted value found %v", got)
		}
	}

	for _, want := range want {
		if !slices.Contains(got, want) {
			t.Errorf("Wanted value missing %v", want)
		}
	}
}

func TestWorkflowActions(t *testing.T) {
	t.Parallel()

	type TestCase struct {
		workflow []byte
		wantErr  bool
	}

	testCases := []TestCase{
		{
			workflow: []byte(workflowWithNoJobs),
			wantErr:  false,
		},
		{
			workflow: []byte(workflowWithJobNoSteps),
			wantErr:  false,
		},
		{
			workflow: []byte(workflowWithJobWithSteps),
			wantErr:  false,
		},
		{
			workflow: []byte(workflowWithJobsWithSteps),
			wantErr:  false,
		},
		{
			workflow: []byte(workflowWithNestedActions),
			wantErr:  false,
		},
		{
			workflow: []byte(workflowWithSyntaxError),
			wantErr:  true,
		},
		{
			workflow: []byte(workflowWithInvalidUses),
			wantErr:  true,
		},
	}

	for i, tc := range testCases {
		t.Run(fmt.Sprintf("#%d", i), func(t *testing.T) {
			t.Parallel()

			_, err := WorkflowActions(tc.workflow)
			if err == nil && tc.wantErr {
				t.Error("Unexpected success")
			} else if err != nil && !tc.wantErr {
				t.Error("Unexpected failure")
			}
		})
	}
}
