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

import (
	"testing"

	"arcoris.dev/apimachinery/api/value"
)

// TestConstraintsForFieldTermEmitsPositiveOperatorsOnly verifies field
// planning is conservative and leaves negative/missing-sensitive operators
// residual.
func TestConstraintsForFieldTermEmitsPositiveOperatorsOnly(t *testing.T) {
	ref := fieldRef("spec.replicas")
	field := selectable(ref, value.KindInteger, Operators(OperatorGreaterThan, OperatorNotEquals))
	field.Index = IndexRange
	positive := term{kind: termField, fieldRef: ref, field: field, operator: OperatorGreaterThan, values: []value.Value{value.Int64Value(1)}}
	negative := term{kind: termField, fieldRef: ref, field: field, operator: OperatorNotEquals, values: []value.Value{value.Int64Value(1)}}

	got := constraintsForFieldTerm(positive)
	if len(got) != 1 || got[0].Kind != ConstraintField || got[0].Op != OperatorGreaterThan {
		t.Fatalf("positive field constraints = %#v", got)
	}
	if got := constraintsForFieldTerm(negative); got != nil {
		t.Fatalf("negative field constraints = %#v; want nil", got)
	}
}

// TestFieldIndexHintControlsPlanning verifies SelectableField.Index is honored
// by field constraint extraction.
func TestFieldIndexHintControlsPlanning(t *testing.T) {
	ref := fieldRef("spec.replicas")
	tests := []struct {
		name string
		hint IndexHint
		op   Operator
		want bool
	}{
		{name: "none equality", hint: IndexNone, op: OperatorEquals},
		{name: "equality equals", hint: IndexEquality, op: OperatorEquals, want: true},
		{name: "equality in", hint: IndexEquality, op: OperatorIn, want: true},
		{name: "equality range", hint: IndexEquality, op: OperatorGreaterThan},
		{name: "range equals", hint: IndexRange, op: OperatorEquals, want: true},
		{name: "range greater", hint: IndexRange, op: OperatorGreaterThan, want: true},
		{name: "range negative", hint: IndexRange, op: OperatorNotEquals},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			field := selectable(ref, value.KindInteger, Operators(tt.op))
			field.Index = tt.hint
			term := term{kind: termField, fieldRef: ref, field: field, operator: tt.op, values: []value.Value{value.Int64Value(1)}}
			if tt.op == OperatorExists {
				term.values = nil
			}

			got := constraintsForFieldTerm(term)
			if (len(got) > 0) != tt.want {
				t.Fatalf("has constraints = %v; want %v (%#v)", len(got) > 0, tt.want, got)
			}
		})
	}
}
