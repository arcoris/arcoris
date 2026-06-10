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

import "testing"

func TestErrorReasonStrings(t *testing.T) {
	tests := []struct {
		name   string
		reason ErrorReason
		want   string
	}{
		{name: "invalid path", reason: ErrorReasonInvalidPath, want: "invalid_path"},
		{name: "invalid syntax", reason: ErrorReasonInvalidSyntax, want: "invalid_syntax"},
		{name: "empty field name", reason: ErrorReasonEmptyFieldName, want: "empty_field_name"},
		{name: "empty map key", reason: ErrorReasonEmptyMapKey, want: "empty_map_key"},
		{name: "duplicate selector field", reason: ErrorReasonDuplicateSelectorField, want: "duplicate_selector_field"},
		{name: "non canonical text", reason: ErrorReasonNonCanonicalText, want: "non_canonical_text"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			requireEqual(t, string(tt.reason), tt.want)
		})
	}
}
