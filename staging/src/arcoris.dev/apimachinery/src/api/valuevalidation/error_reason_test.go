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

package valuevalidation_test

import (
	"testing"

	"arcoris.dev/apimachinery/api/valuevalidation"
)

func TestErrorReasonValuesAreStable(t *testing.T) {
	tests := []struct {
		name   string
		reason valuevalidation.ErrorReason
		want   string
	}{
		{name: "invalid zero", reason: valuevalidation.ErrorReasonInvalidZero, want: "invalid_zero"},
		{name: "invalid descriptor", reason: valuevalidation.ErrorReasonInvalidDescriptor, want: "invalid_descriptor"},
		{name: "invalid path", reason: valuevalidation.ErrorReasonInvalidPath, want: "invalid_path"},
		{name: "kind mismatch", reason: valuevalidation.ErrorReasonKindMismatch, want: "kind_mismatch"},
		{name: "null not allowed", reason: valuevalidation.ErrorReasonNullNotAllowed, want: "null_not_allowed"},
		{name: "missing field", reason: valuevalidation.ErrorReasonMissingField, want: "missing_field"},
		{name: "unknown field", reason: valuevalidation.ErrorReasonUnknownField, want: "unknown_field"},
		{name: "invalid field name", reason: valuevalidation.ErrorReasonInvalidFieldName, want: "invalid_field_name"},
		{name: "invalid map key", reason: valuevalidation.ErrorReasonInvalidMapKey, want: "invalid_map_key"},
		{name: "below minimum", reason: valuevalidation.ErrorReasonBelowMinimum, want: "below_minimum"},
		{name: "above maximum", reason: valuevalidation.ErrorReasonAboveMaximum, want: "above_maximum"},
		{name: "too short", reason: valuevalidation.ErrorReasonTooShort, want: "too_short"},
		{name: "too long", reason: valuevalidation.ErrorReasonTooLong, want: "too_long"},
		{name: "pattern mismatch", reason: valuevalidation.ErrorReasonPatternMismatch, want: "pattern_mismatch"},
		{name: "enum mismatch", reason: valuevalidation.ErrorReasonEnumMismatch, want: "enum_mismatch"},
		{name: "unresolved ref", reason: valuevalidation.ErrorReasonUnresolvedRef, want: "unresolved_ref"},
		{name: "reference cycle", reason: valuevalidation.ErrorReasonReferenceCycle, want: "reference_cycle"},
		{name: "missing list key", reason: valuevalidation.ErrorReasonMissingListKey, want: "missing_list_key"},
		{name: "invalid list key", reason: valuevalidation.ErrorReasonInvalidListKey, want: "invalid_list_key"},
		{name: "duplicate list key", reason: valuevalidation.ErrorReasonDuplicateListKey, want: "duplicate_list_key"},
		{name: "duplicate list set element", reason: valuevalidation.ErrorReasonDuplicateListSetElement, want: "duplicate_list_set_element"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := string(tt.reason); got != tt.want {
				t.Fatalf("reason = %q, want %q", got, tt.want)
			}
		})
	}
}
