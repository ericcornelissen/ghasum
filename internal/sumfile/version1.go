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
	"errors"
	"fmt"
	"sort"
	"strings"
)

func decodeV1(lines []string) ([]Entry, error) {
	entries := make([]Entry, len(lines))
	for i, line := range lines {
		// split "line" into "id[@id..]" "sum"
		j := strings.IndexRune(line, ' ')
		if j <= 0 || j >= len(line)-1 {
			err := fmt.Errorf("syntax error on line %d", i+2)
			return nil, errors.Join(ErrCorrupted, err)
		}

		entries[i] = Entry{
			ID:       strings.Split(line[:j], "@"),
			Checksum: line[j+1:],
		}
	}

	if !validV1(entries) {
		return nil, ErrInvalid
	}

	return entries, nil
}

func encodeV1(entries []Entry) (string, error) {
	if !validV1(entries) {
		return "", ErrInvalid
	}

	var sb strings.Builder
	lines := make([]string, len(entries))
	for i, entry := range entries {
		for i, part := range entry.ID {
			if i != 0 {
				sb.WriteRune('@')
			}
			sb.WriteString(part)
		}

		sb.WriteRune(' ')
		sb.WriteString(entry.Checksum)
		sb.WriteRune('\n')

		lines[i] = sb.String()
		sb.Reset()
	}

	sort.Strings(lines)
	return strings.Join(lines, ""), nil
}

func validV1(entries []Entry) bool {
	if hasDuplicates(entries) || hasMissing(entries) {
		return false
	}

	for _, entry := range entries {
		if strings.ContainsAny(entry.Checksum, "\n ") {
			return false
		}

		if strings.ContainsAny(strings.Join(entry.ID, ""), "\n @") {
			return false
		}
	}

	return true
}
