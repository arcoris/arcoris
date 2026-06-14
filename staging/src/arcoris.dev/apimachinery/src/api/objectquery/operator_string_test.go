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

func TestOperatorString(t *testing.T) {
	tests := []struct {
		name string
		op   Operator
		want string
	}{
		{name: "zero", op: 0, want: "unknown"},
		{name: "exists", op: OperatorExists, want: "exists"},
		{name: "does not exist", op: OperatorDoesNotExist, want: "doesNotExist"},
		{name: "equals", op: OperatorEquals, want: "equals"},
		{name: "not equals", op: OperatorNotEquals, want: "notEquals"},
		{name: "in", op: OperatorIn, want: "in"},
		{name: "not in", op: OperatorNotIn, want: "notIn"},
		{name: "unknown", op: Operator(255), want: "unknown"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.op.String(); got != tt.want {
				t.Fatalf("String() = %q; want %q", got, tt.want)
			}
		})
	}
}
