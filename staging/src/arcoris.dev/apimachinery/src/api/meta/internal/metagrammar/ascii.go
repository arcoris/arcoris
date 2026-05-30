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

package metagrammar

// HasControl reports whether s contains an ASCII control byte.
func HasControl(s string) bool {
	for i := 0; i < len(s); i++ {
		if s[i] < 0x20 || s[i] == 0x7f {
			return true
		}
	}
	return false
}

// HasWhitespace reports whether s contains ASCII whitespace.
func HasWhitespace(s string) bool {
	for i := 0; i < len(s); i++ {
		switch s[i] {
		case ' ', '\t', '\n', '\r', '\v', '\f':
			return true
		}
	}
	return false
}

// HasUnsafeScalarChar reports whether s contains bytes that make opaque
// metadata scalars ambiguous in line-oriented logs, route fragments, or compact
// diagnostics.
func HasUnsafeScalarChar(s string) bool {
	for i := 0; i < len(s); i++ {
		switch s[i] {
		case '/', '\\':
			return true
		}
		if s[i] < 0x20 || s[i] == 0x7f {
			return true
		}
	}
	return false
}
