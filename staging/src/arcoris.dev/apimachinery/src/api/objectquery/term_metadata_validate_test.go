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

// TestValidateMetadataTermErrors verifies metadata key, value, operator, and
// arity failures keep stable query sentinels.
func TestValidateMetadataTermErrors(t *testing.T) {
	tests := []struct {
		name string
		op   Operator
		key  string
		vals []string
	}{
		{name: "invalid operator", op: Operator(99), key: "env"},
		{name: "unsupported operator", op: OperatorContains, key: "env", vals: []string{"p"}},
		{name: "empty key", op: OperatorEquals, key: "", vals: []string{"prod"}},
		{name: "exists with value", op: OperatorExists, key: "env", vals: []string{"prod"}},
		{name: "equals no value", op: OperatorEquals, key: "env"},
		{name: "in no values", op: OperatorIn, key: "env"},
		{name: "invalid value", op: OperatorEquals, key: "env", vals: []string{"bad value"}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := validateMetadataTerm(metadataLabels, tt.op, tt.key, tt.vals)
			requireErrorIs(t, err, ErrInvalidQuery)
		})
	}
}

// TestValidateMetadataKeyAndValueUnknownDomain verifies unknown internal
// metadata domains fail closed.
func TestValidateMetadataKeyAndValueUnknownDomain(t *testing.T) {
	requireErrorIs(t, validateMetadataKey(0, "env"), ErrInvalidTerm)
	requireErrorIs(t, validateMetadataValue(0, "prod"), ErrInvalidTerm)
}
