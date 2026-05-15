/*
  Copyright 2026 The ARCORIS Authors

  Licensed under the Apache License, Version 2.0 (the "License");
  you may not use this file except in compliance with the License.
  You may obtain a copy of the License at

      http://www.apache.org/licenses/LICENSE-2.0

  Unless required by applicable law or agreed to in writing, software
  distributed under the License is distributed on an "AS IS" BASIS,
  WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
  See the License for the specific language governing permissions and
  limitations under the License.
*/

package health

// validLowerSnakeIdentifier reports whether value is a stable lower_snake_case
// identifier within maxLength.
//
// Health check names and reasons share the same syntax so they remain safe for
// reports, diagnostics, metrics labels, tests, and transport adapters. The
// helper deliberately rejects empty values; callers that own an empty sentinel,
// such as ReasonNone, must handle that case before calling it.
func validLowerSnakeIdentifier(val string, maxLength int) bool {
	if len(val) == 0 || len(val) > maxLength {
		return false
	}

	previousUnderscore := false

	for i := 0; i < len(val); i++ {
		c := val[i]

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
