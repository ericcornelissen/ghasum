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
	"bytes"
	"testing"
	"testing/quick"

	"github.com/liamg/memoryfs"
)

func TestActionsInWorkflows(t *testing.T) {
	t.Parallel()

	t.Run("Valid examples", func(t *testing.T) {
		t.Parallel()

		type TestCase struct {
			name string
			in   []workflow
			want int
		}

		testCases := []TestCase{
			{
				name: "no jobs",
				in: []workflow{
					{
						Jobs: map[string]job{},
					},
				},
				want: 0,
			},
			{
				name: "one job without steps",
				in: []workflow{
					{
						Jobs: map[string]job{
							"example": {
								Steps: []step{},
							},
						},
					},
				},
				want: 0,
			},
			{
				name: "multiple jobs without steps",
				in: []workflow{
					{
						Jobs: map[string]job{
							"example-a": {
								Steps: []step{},
							},
							"example-b": {
								Steps: []step{},
							},
						},
					},
				},
				want: 0,
			},
			{
				name: "one job with a step without uses",
				in: []workflow{
					{
						Jobs: map[string]job{
							"example": {
								Steps: []step{
									{},
								},
							},
						},
					},
				},
				want: 0,
			},
			{
				name: "one job with one step",
				in: []workflow{
					{
						Jobs: map[string]job{
							"example": {
								Steps: []step{
									{
										Uses: "foo/bar@v1",
									},
								},
							},
						},
					},
				},
				want: 1,
			},
			{
				name: "multiple jobs with one unique step each",
				in: []workflow{
					{
						Jobs: map[string]job{
							"example-a": {
								Steps: []step{
									{
										Uses: "foo/bar@v1",
									},
								},
							},
							"example-b": {
								Steps: []step{
									{
										Uses: "foo/baz@v1",
									},
								},
							},
						},
					},
				},
				want: 2,
			},
			{
				name: "one job with multiple unique steps",
				in: []workflow{
					{
						Jobs: map[string]job{
							"example": {
								Steps: []step{
									{
										Uses: "foo/bar@v1",
									},
									{
										Uses: "foo/baz@v2",
									},
								},
							},
						},
					},
				},
				want: 2,
			},
			{
				name: "multiple jobs with multiple unique steps",
				in: []workflow{
					{
						Jobs: map[string]job{
							"example-a": {
								Steps: []step{
									{
										Uses: "foo/bar@v1",
									},
									{
										Uses: "hello/world@v2",
									},
								},
							},
							"example-b": {
								Steps: []step{
									{
										Uses: "foo/baz@v1",
									},
									{
										Uses: "hallo/wereld@v2",
									},
								},
							},
						},
					},
				},
				want: 4,
			},
			{
				name: "one jobs with duplicate steps",
				in: []workflow{
					{
						Jobs: map[string]job{
							"example": {
								Steps: []step{
									{
										Uses: "foo/bar@v1",
									},
									{
										Uses: "foo/bar@v1",
									},
								},
							},
						},
					},
				},
				want: 1,
			},
			{
				name: "multiple jobs with duplicate step between them",
				in: []workflow{
					{

						Jobs: map[string]job{
							"example-a": {
								Steps: []step{
									{
										Uses: "foo/bar@v1",
									},
								},
							},
							"example-b": {
								Steps: []step{
									{
										Uses: "foo/bar@v1",
									},
								},
							},
						},
					},
				},
				want: 1,
			},
		}

		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				t.Parallel()

				got, err := actionsInWorkflows(tc.in)
				if err != nil {
					t.Fatalf("Unexpected error: %+v", err)
				}

				if got, want := len(got), tc.want; got != want {
					t.Errorf("Incorrect result length (got %d, want %d)", got, want)
				}
			})
		}
	})

	t.Run("Invalid examples", func(t *testing.T) {
		t.Parallel()

		type TestCase struct {
			name string
			in   []workflow
		}

		testCases := []TestCase{
			{
				name: "invalid uses value",
				in: []workflow{
					{
						Jobs: map[string]job{
							"example": {
								Steps: []step{
									{
										Uses: "this isn't an action",
									},
								},
							},
						},
					},
				},
			},
		}

		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				if _, err := actionsInWorkflows(tc.in); err == nil {
					t.Fatal("Unexpected success")
				}
			})
		}
	})

	t.Run("Arbitrary", func(t *testing.T) {
		t.Parallel()

		unique := func(workflows []workflow) bool {
			actions, err := actionsInWorkflows(workflows)
			if err != nil {
				return true
			}

			seen := make(map[GitHubAction]struct{}, 0)
			for _, action := range actions {
				if _, ok := seen[action]; ok {
					return false
				}

				seen[action] = struct{}{}
			}

			return true
		}

		if err := quick.Check(unique, nil); err != nil {
			t.Errorf("Duplicate value detected for: %v", err)
		}
	})
}

func TestWorkflowsInRepo(t *testing.T) {
	t.Parallel()

	t.Run("Valid examples", func(t *testing.T) {
		t.Parallel()

		type TestCase struct {
			name      string
			workflows map[string]mockFsEntry
			want      []workflowFile
		}

		testCases := []TestCase{
			{
				name: ".yml workflow",
				workflows: map[string]mockFsEntry{
					"example.yml": {
						Content: []byte(workflowWithJobsWithSteps),
					},
				},
				want: []workflowFile{
					{
						content: []byte(workflowWithJobsWithSteps),
						path:    ".github/workflows/example.yml",
					},
				},
			},
			{
				name: ".yaml workflow",
				workflows: map[string]mockFsEntry{
					"example.yaml": {
						Content: []byte(workflowWithJobsWithSteps),
					},
				},
				want: []workflowFile{
					{
						content: []byte(workflowWithJobsWithSteps),
						path:    ".github/workflows/example.yaml",
					},
				},
			},
			{
				name: "non-workflow file",
				workflows: map[string]mockFsEntry{
					"greeting.txt": {
						Content: []byte("Hello world!"),
					},
				},
				want: []workflowFile{},
			},
			{
				name: "nested directory",
				workflows: map[string]mockFsEntry{
					"greeting.txt": {
						Dir: true,
						Children: map[string]mockFsEntry{
							"workflow.yml": {},
						},
					},
				},
				want: []workflowFile{},
			},
		}

		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				t.Parallel()

				repo, err := mockRepo(tc.workflows)
				if err != nil {
					t.Fatalf("Could not initialize file system: %+v", err)
				}

				got, err := workflowsInRepo(repo)
				if err != nil {
					t.Fatalf("Unexpected error: %+v", err)
				}

				if got, want := len(got), len(tc.want); got != want {
					t.Fatalf("Incorrect result length (got %d, want %d)", got, want)
				}

				for i, got := range got {
					if got, want := got.content, tc.want[i].content; !bytes.Equal(got, want) {
						t.Errorf("Incorrect content for workflow %d (got %s, want %s)", i, got, want)
					}

					if got, want := got.path, tc.want[i].path; got != want {
						t.Errorf("Incorrect path for workflow %d (got %s, want %s)", i, got, want)
					}
				}
			})
		}
	})

	t.Run("No actions", func(t *testing.T) {
		t.Parallel()

		repo := memoryfs.New()
		if _, err := workflowsInRepo(repo); err == nil {
			t.Fatal("Unexpected success")
		}
	})
}
