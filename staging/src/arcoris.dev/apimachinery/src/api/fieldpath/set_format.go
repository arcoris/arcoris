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

package fieldpath

import "strings"

// String returns a deterministic diagnostic representation of s.
//
// The format is intended for debugging and logs. It is not a wire format and
// does not imply JSON serialization.
func (s Set) String() string {
	if len(s.paths) == 0 {
		return "{}"
	}

	var builder strings.Builder
	builder.WriteByte('{')
	for i, p := range s.paths {
		if i > 0 {
			builder.WriteString(", ")
		}

		builder.WriteString(p.String())
	}
	builder.WriteByte('}')

	return builder.String()
}
