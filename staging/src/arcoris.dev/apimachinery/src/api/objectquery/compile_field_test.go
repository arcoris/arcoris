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

// TestValidateFieldLiteralsEnforcesOperatorAndKind verifies field compilation
// rejects operator/type combinations before runtime matching.
func TestValidateFieldLiteralsEnforcesOperatorAndKind(t *testing.T) {
	field := selectable(fieldRef("spec.replicas"), value.KindInteger, Operators(OperatorEquals, OperatorLessThan))

	requireNoError(t, validateFieldLiterals(field, OperatorLessThan, []value.Value{value.Int64Value(2)}))

	err := validateFieldLiterals(field, OperatorLessThan, []value.Value{value.StringValue("2")})
	requireErrorIs(t, err, ErrInvalidField)
	requireErrorIs(t, err, ErrInvalidTerm)

	err = validateFieldLiterals(field, OperatorLessThan, []value.Value{value.NullValue()})
	requireErrorIs(t, err, ErrUnsupportedOperator)
}

// TestRequireFieldMatchComparable reports runtime kind mismatches without
// panicking if a malformed field declaration slips through.
func TestRequireFieldMatchComparable(t *testing.T) {
	requireNoError(t, requireFieldMatchComparable(value.StringValue("a"), value.StringValue("b")))

	if err := requireFieldMatchComparable(value.StringValue("a"), value.Int64Value(1)); err == nil {
		t.Fatal("requireFieldMatchComparable mismatch error = nil; want error")
	}
}
