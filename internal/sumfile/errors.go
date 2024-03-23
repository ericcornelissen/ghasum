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

import "errors"

var (
	// ErrCorrupted is the error when a checksum file is corrupted.
	ErrCorrupted = errors.New("checksums are corrupted")

	// ErrDuplicate is the error when a checksum file contains two entries for the
	// same identifier.
	ErrDuplicate = errors.New("duplicate entry found")

	// ErrHeaders is the error for when sumfile headers are invalid.
	ErrHeaders = errors.New("sumfile headers are invalid")

	// ErrMissing is the error when an entry is missing an id (part) or checksums.
	ErrMissing = errors.New("missing id or checksum")

	// ErrSyntax is the error when a checksum file has a syntax error.
	ErrSyntax = errors.New("syntax error")

	// ErrVersion is the error when the version is invalid or missing from the
	// checksum file.
	ErrVersion = errors.New("version error")
)
