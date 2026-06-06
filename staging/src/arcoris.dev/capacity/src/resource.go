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

package capacity

import "fmt"

// Resource is a stable accounting dimension.
//
// Resource identifies what is counted, not a runtime object instance. Valid
// values are dot-separated lower_snake_case segments such as "worker_slots" or
// "resilience.bulkhead.slots".
type Resource string

// MustResource returns value as a Resource or panics when it is invalid.
func MustResource(value string) Resource {
	resource := Resource(value)
	if !resource.IsValid() {
		panic(errorAt(
			"resource",
			ErrInvalidResource,
			fmt.Sprintf("resource %q must be dot-separated lower_snake_case", resource),
		))
	}
	return resource
}

// IsValid reports whether r follows the capacity resource grammar.
func (r Resource) IsValid() bool {
	value := string(r)
	if value == "" {
		return false
	}

	startSegment := true
	previousUnderscore := false
	for i := 0; i < len(value); i++ {
		c := value[i]

		switch {
		case c == '.':
			if startSegment || previousUnderscore {
				return false
			}
			startSegment = true
			previousUnderscore = false

		case isLowerASCII(c):
			startSegment = false
			previousUnderscore = false

		case isDigitASCII(c):
			if startSegment {
				return false
			}
			previousUnderscore = false

		case c == '_':
			if startSegment || previousUnderscore {
				return false
			}
			previousUnderscore = true

		default:
			return false
		}
	}

	return !startSegment && !previousUnderscore
}

// String returns r as its stable identifier string.
func (r Resource) String() string {
	return string(r)
}

// isLowerASCII reports whether c is an ASCII lowercase letter.
func isLowerASCII(c byte) bool {
	return c >= 'a' && c <= 'z'
}

// isDigitASCII reports whether c is an ASCII decimal digit.
func isDigitASCII(c byte) bool {
	return c >= '0' && c <= '9'
}
