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

func TestOperatorValidity(t *testing.T) {
	tests := []struct {
		name  string
		op    Operator
		valid bool
	}{
		{name: "zero", op: 0},
		{name: "exists", op: OperatorExists, valid: true},
		{name: "does not exist", op: OperatorDoesNotExist, valid: true},
		{name: "equals", op: OperatorEquals, valid: true},
		{name: "not equals", op: OperatorNotEquals, valid: true},
		{name: "in", op: OperatorIn, valid: true},
		{name: "not in", op: OperatorNotIn, valid: true},
		{name: "unknown", op: Operator(255)},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.op.IsValid(); got != tt.valid {
				t.Fatalf("IsValid() = %v; want %v", got, tt.valid)
			}
		})
	}
}
