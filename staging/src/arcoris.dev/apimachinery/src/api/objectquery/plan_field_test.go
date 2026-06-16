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
	positive := term{kind: termField, fieldRef: ref, operator: OperatorGreaterThan, values: []value.Value{value.Int64Value(1)}}
	negative := term{kind: termField, fieldRef: ref, operator: OperatorNotEquals, values: []value.Value{value.Int64Value(1)}}

	got := constraintsForFieldTerm(positive)
	if len(got) != 1 || got[0].Kind != ConstraintField || got[0].Op != OperatorGreaterThan {
		t.Fatalf("positive field constraints = %#v", got)
	}
	if got := constraintsForFieldTerm(negative); got != nil {
		t.Fatalf("negative field constraints = %#v; want nil", got)
	}
}
