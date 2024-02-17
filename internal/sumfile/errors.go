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

	// ErrHeaders is the error for when checksum headers are invalid.
	ErrHeaders = errors.New("checksum headers are invalid")

	// ErrInvalid is the error for when checksums are invalid.
	ErrInvalid = errors.New("checksums are invalid")

	// ErrSyntax is the error when a checksum file has a syntax error.
	ErrSyntax = errors.New("syntax error")
)
