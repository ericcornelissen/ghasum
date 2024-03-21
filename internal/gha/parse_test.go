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
	"strings"
	"testing"
	"testing/quick"
)

func TestParseUses(t *testing.T) {
	t.Parallel()

	t.Run("Valid examples", func(t *testing.T) {
		t.Parallel()

		type TestCase struct {
			in   string
			want GitHubAction
		}

		testCases := []TestCase{
			{
				in: "foo/bar@v1",
				want: GitHubAction{
					Owner:   "foo",
					Project: "bar",
					Ref:     "v1",
				},
			},
			{
				in: "foo/baz@v3.1.4",
				want: GitHubAction{
					Owner:   "foo",
					Project: "baz",
					Ref:     "v3.1.4",
				},
			},
			{
				in: "hello/world@random-ref",
				want: GitHubAction{
					Owner:   "hello",
					Project: "world",
					Ref:     "random-ref",
				},
			},
			{
				in: "hallo/wereld@35dd46a3b3dfbb14198f8d19fb083ce0832dce4a",
				want: GitHubAction{
					Owner:   "hallo",
					Project: "wereld",
					Ref:     "35dd46a3b3dfbb14198f8d19fb083ce0832dce4a",
				},
			},
			{
				in: "foo/bar/baz@v2",
				want: GitHubAction{
					Owner:   "foo",
					Project: "bar",
					Ref:     "v2",
				},
			},
		}

		for _, tc := range testCases {
			t.Run(tc.in, func(t *testing.T) {
				t.Parallel()

				got, err := parseUses(tc.in)
				if err != nil {
					t.Fatalf("Unexpected error: %+v", err)
				}

				if got, want := got.Owner, tc.want.Owner; got != want {
					t.Errorf("Incorrect owner (got %q, want %q)", got, want)
				}

				if got, want := got.Project, tc.want.Project; got != want {
					t.Errorf("Incorrect project (got %q, want %q)", got, want)
				}

				if got, want := got.Ref, tc.want.Ref; got != want {
					t.Errorf("Incorrect ref (got %q, want %q)", got, want)
				}
			})
		}
	})

	t.Run("Invalid examples", func(t *testing.T) {
		t.Parallel()

		type TestCase struct {
			in   string
			want string
		}

		testCases := []TestCase{
			{
				in:   "foobar",
				want: "invalid uses value",
			},
			{
				in:   "foo/bar",
				want: "invalid uses value",
			},
			{
				in:   "f@o/bar@baz",
				want: "invalid uses value",
			},
			{
				in:   "foo/b@r@baz",
				want: "invalid uses value",
			},
			{
				in:   "foo/bar/b@z@ref",
				want: "invalid uses value",
			},
			{
				in:   "foo@bar",
				want: "invalid repository in uses",
			},
			{
				in:   "foo/@bar",
				want: "invalid repository in uses",
			},
			{
				in:   "foo/bar/@baz",
				want: "invalid repository path in uses",
			},
			{
				in:   "foo//@bar",
				want: "invalid repository path in uses",
			},
			{
				in:   "foo//bar@baz",
				want: "invalid repository path in uses",
			},
		}

		for _, tc := range testCases {
			t.Run(tc.in, func(t *testing.T) {
				t.Parallel()

				_, err := parseUses(tc.in)
				if err == nil {
					t.Fatal("Unexpected success")
				}

				if got, want := err.Error(), tc.want; got != want {
					t.Errorf("Incorrect error message (got %q, want %q)", got, want)
				}
			})
		}
	})

	t.Run("Arbitrary values", func(t *testing.T) {
		t.Parallel()

		constructive := func(owner, project, path, ref string) bool {
			if len(owner) == 0 || len(project) == 0 || len(ref) == 0 {
				return true
			}

			repo := fmt.Sprintf("%s/%s", owner, project)
			if strings.Count(repo, "/") != 1 {
				return true
			}

			if len(path) > 0 {
				repo = fmt.Sprintf("%s/%s", repo, path)
			}

			if strings.ContainsRune(repo, '@') {
				return true
			}

			uses := fmt.Sprintf("%s@%s", repo, ref)

			action, err := parseUses(uses)
			if err != nil {
				return false
			}

			return action.Owner == owner && action.Project == project && action.Ref == ref
		}

		if err := quick.Check(constructive, nil); err != nil {
			t.Errorf("Parsing failed for: %v", err)
		}

		noPanic := func(uses string) bool {
			_, _ = parseUses(uses)
			return true
		}

		if err := quick.Check(noPanic, nil); err != nil {
			t.Errorf("Parsing failed for: %v", err)
		}
	})
}

func TestParseWorkflow(t *testing.T) {
	t.Parallel()

	t.Run("Valid examples", func(t *testing.T) {
		t.Parallel()

		type TestCase struct {
			in   string
			want workflow
		}

		testCases := []TestCase{
			{
				in: workflowWithNoJobs,
				want: workflow{
					Jobs: map[string]job{},
				},
			},
			{
				in: workflowWithJobNoSteps,
				want: workflow{
					Jobs: map[string]job{
						"no-steps": {},
					},
				},
			},
			{
				in: workflowWithJobWithSteps,
				want: workflow{
					Jobs: map[string]job{
						"only-job": {
							Steps: []step{
								{
									Uses: "foo/bar@v1",
								},
								{
									Uses: "",
								},
								{
									Uses: "foo/baz@v2",
								},
							},
						},
					},
				},
			},
			{
				in: workflowWithJobsWithSteps,
				want: workflow{
					Jobs: map[string]job{
						"job-a": {
							Steps: []step{
								{
									Uses: "foo/bar@v1",
								},
							},
						},
						"job-b": {
							Steps: []step{
								{
									Uses: "",
								},
								{
									Uses: "foo/baz@v2",
								},
							},
						},
					},
				},
			},
		}

		for _, tc := range testCases {
			t.Run(strings.Split(tc.in, "\n")[0], func(t *testing.T) {
				t.Parallel()

				got, err := parseWorkflow([]byte(tc.in))
				if err != nil {
					t.Fatalf("Unexpected error: %+v", err)
				}

				if got, want := len(got.Jobs), len(tc.want.Jobs); got != want {
					t.Fatalf("Incorrect jobs length (got %d, want %d)", got, want)
				}

				for name, job := range got.Jobs {
					want, ok := tc.want.Jobs[name]
					if !ok {
						t.Errorf("Got unwanted job %q", name)
						continue
					}

					if got, want := len(job.Steps), len(want.Steps); got != want {
						t.Errorf("Incorrect steps length for job %q (got %d, want %d)", name, got, want)
						continue
					}

					for i, step := range job.Steps {
						want := want.Steps[i]

						if got, want := step.Uses, want.Uses; got != want {
							t.Errorf("Incorrect uses for step %d of job %q (got %q, want %q)", i, name, got, want)
						}
					}
				}
			})
		}
	})

	t.Run("Invalid examples", func(t *testing.T) {
		t.Parallel()

		cases := []string{
			workflowWithSyntaxError,
		}

		for _, tc := range cases {
			t.Run(tc, func(t *testing.T) {
				t.Parallel()

				if _, err := parseWorkflow([]byte(tc)); err == nil {
					t.Fatal("Unexpected success")
				}
			})
		}
	})

	t.Run("Arbitrary values", func(t *testing.T) {
		t.Parallel()

		noPanic := func(w []byte) bool {
			_, _ = parseWorkflow(w)
			return true
		}

		if err := quick.Check(noPanic, nil); err != nil {
			t.Errorf("Parsing failed for: %v", err)
		}
	})
}
