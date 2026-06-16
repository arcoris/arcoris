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

package objectquery

import "testing"

// TestOperatorValidityAndString covers operator validity, diagnostics, and the
// metadata operator set.
func TestOperatorValidityAndString(t *testing.T) {
	tests := []struct {
		op      Operator
		valid   bool
		text    string
		support bool
	}{
		{op: 0, text: "unknown"},
		{op: OperatorExists, valid: true, text: "exists", support: true},
		{op: OperatorDoesNotExist, valid: true, text: "doesNotExist", support: true},
		{op: OperatorEquals, valid: true, text: "equals", support: true},
		{op: OperatorNotEquals, valid: true, text: "notEquals", support: true},
		{op: OperatorIn, valid: true, text: "in", support: true},
		{op: OperatorNotIn, valid: true, text: "notIn", support: true},
		{op: OperatorLessThan, valid: true, text: "lessThan"},
		{op: OperatorLessOrEqual, valid: true, text: "lessOrEqual"},
		{op: OperatorGreaterThan, valid: true, text: "greaterThan"},
		{op: OperatorGreaterOrEqual, valid: true, text: "greaterOrEqual"},
		{op: OperatorHasPrefix, valid: true, text: "hasPrefix"},
		{op: OperatorHasSuffix, valid: true, text: "hasSuffix"},
		{op: OperatorContains, valid: true, text: "contains"},
		{op: Operator(99), text: "unknown"},
	}

	for _, tt := range tests {
		t.Run(tt.text, func(t *testing.T) {
			if got := tt.op.IsValid(); got != tt.valid {
				t.Fatalf("IsValid() = %v; want %v", got, tt.valid)
			}
			if got := tt.op.String(); got != tt.text {
				t.Fatalf("String() = %q; want %q", got, tt.text)
			}
			if got := metadataOperators.Supports(tt.op); got != tt.support {
				t.Fatalf("metadata support = %v; want %v", got, tt.support)
			}
		})
	}
}
