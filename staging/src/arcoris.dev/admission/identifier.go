// Copyright 2026 The ARCORIS Authors.
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

package admission

// validLowerSnakeIdentifier reports whether value is a compact ASCII
// lower_snake_case identifier.
//
// The helper is deliberately stricter than a general-purpose identifier parser:
// values must be non-empty, bounded by maxLength, start with a lowercase letter,
// and contain only lowercase letters, digits, and single underscores between
// tokens. This keeps admission reasons and kinds stable enough for logs,
// documentation, and downstream dimensions without admitting dynamic payloads.
func validLowerSnakeIdentifier(value string, maxLength int) bool {
	if len(value) == 0 || len(value) > maxLength {
		return false
	}

	previousUnderscore := false
	for i := 0; i < len(value); i++ {
		c := value[i]
		switch {
		case c >= 'a' && c <= 'z':
			previousUnderscore = false
		case c >= '0' && c <= '9':
			if i == 0 {
				return false
			}
			previousUnderscore = false
		case c == '_':
			if i == 0 || previousUnderscore {
				return false
			}
			previousUnderscore = true
		default:
			return false
		}
	}

	return !previousUnderscore
}

// validDotPathIdentifier reports whether value is a dot-separated path of
// valid lower_snake_case segments.
//
// ComponentID uses this grammar so packages can express stable ownership paths
// such as "resilience.bulkhead" while rejecting empty segments, dynamic
// instance identifiers, and punctuation that would be hard to preserve across
// observability systems.
func validDotPathIdentifier(value string, maxLength int) bool {
	if len(value) == 0 || len(value) > maxLength {
		return false
	}

	segmentStart := 0
	for i := 0; i <= len(value); i++ {
		if i != len(value) && value[i] != '.' {
			continue
		}
		if i == segmentStart {
			return false
		}
		if !validLowerSnakeIdentifier(value[segmentStart:i], maxLength) {
			return false
		}
		segmentStart = i + 1
	}

	return true
}
