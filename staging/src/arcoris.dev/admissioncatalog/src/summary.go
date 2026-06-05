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

package admissioncatalog

import "unicode/utf8"

// validSummary reports whether summary is acceptable descriptor prose.
//
// Summaries are optional human-facing metadata. They are intentionally kept as
// single-line text so catalogs do not become a place for stack traces, logs, or
// dynamic request data. The package cannot detect every kind of dynamic data,
// but it can reject invalid UTF-8 and control characters.
func validSummary(summary string) bool {
	if !utf8.ValidString(summary) {
		return false
	}
	for _, r := range summary {
		if r < ' ' || r == 0x7f {
			return false
		}
	}
	return true
}
