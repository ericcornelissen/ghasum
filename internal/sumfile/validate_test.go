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

package sumfile

import (
	"testing"
)

func TestHasEmpty(t *testing.T) {
	t.Parallel()

	t.Run("Non-empty examples", func(t *testing.T) {
		t.Parallel()

		type TestCase struct {
			name    string
			entries []Entry
		}

		testCases := []TestCase{
			{
				name:    "no entries",
				entries: []Entry{},
			},
			{
				name: "one ID parts",
				entries: []Entry{
					{
						Checksum: "checksum",
						ID:       []string{"foobar"},
					},
				},
			},
			{
				name: "multiple ID parts",
				entries: []Entry{
					{
						Checksum: "checksum",
						ID:       []string{"foo", "bar"},
					},
				},
			},
		}

		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				t.Parallel()

				got := hasMissing(tc.entries)
				if got {
					t.Fatal("Unexpected positive result")
				}
			})
		}
	})

	t.Run("Empty examples", func(t *testing.T) {
		t.Parallel()

		type TestCase struct {
			name    string
			entries []Entry
		}

		testCases := []TestCase{
			{
				name: "empty checksum",
				entries: []Entry{
					{
						Checksum: "",
						ID:       []string{"foobar"},
					},
				},
			},
			{
				name: "empty id array",
				entries: []Entry{
					{
						Checksum: "not-empty",
						ID:       []string{},
					},
				},
			},
			{
				name: "empty id part",
				entries: []Entry{
					{
						Checksum: "not-empty",
						ID:       []string{""},
					},
				},
			},
		}

		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				t.Parallel()

				got := hasMissing(tc.entries)
				if !got {
					t.Fatal("Unexpected negative result")
				}
			})
		}
	})
}

func TestHasDuplicates(t *testing.T) {
	t.Parallel()

	t.Run("No duplicates examples", func(t *testing.T) {
		t.Parallel()

		type TestCase struct {
			name    string
			entries []Entry
		}

		testCases := []TestCase{
			{
				name:    "no entries",
				entries: []Entry{},
			},
			{
				name: "one part",
				entries: []Entry{
					{
						ID: []string{"foo"},
					},
					{
						ID: []string{"bar"},
					},
				},
			},
			{
				name: "two parts",
				entries: []Entry{
					{
						ID: []string{"foo", "bar"},
					},
					{
						ID: []string{"hello", "world"},
					},
				},
			},
			{
				name: "two parts, first differs",
				entries: []Entry{
					{
						ID: []string{"bar", "foo"},
					},
					{
						ID: []string{"baz", "foo"},
					},
				},
			},
			{
				name: "two parts, second differs",
				entries: []Entry{
					{
						ID: []string{"foo", "bar"},
					},
					{
						ID: []string{"foo", "baz"},
					},
				},
			},
		}

		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				t.Parallel()

				got := hasDuplicates(tc.entries)
				if got {
					t.Fatal("Unexpected positive result")
				}
			})
		}
	})

	t.Run("Duplicate examples", func(t *testing.T) {
		t.Parallel()

		type TestCase struct {
			name    string
			entries []Entry
		}

		testCases := []TestCase{
			{
				name: "one part",
				entries: []Entry{
					{
						ID: []string{"foobar"},
					},
					{
						ID: []string{"foobar"},
					},
				},
			},
			{
				name: "multiple parts",
				entries: []Entry{
					{
						ID: []string{"foo", "bar"},
					},
					{
						ID: []string{"foo", "bar"},
					},
				},
			},
		}

		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				t.Parallel()

				got := hasDuplicates(tc.entries)
				if !got {
					t.Fatal("Unexpected negative result")
				}
			})
		}
	})
}
