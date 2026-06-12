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

package objectlifecycle

import "testing"

func TestOperationStringAndValidity(t *testing.T) {
	tests := []struct {
		name  string
		op    Operation
		text  string
		valid bool
	}{
		{name: "zero", op: 0, text: "unknown", valid: false},
		{name: "get", op: OperationGet, text: "get", valid: true},
		{name: "create", op: OperationCreate, text: "create", valid: true},
		{name: "apply", op: OperationApply, text: "apply", valid: true},
		{name: "update observed", op: OperationUpdateObserved, text: "update_observed", valid: true},
		{name: "patch metadata", op: OperationPatchMetadata, text: "patch_metadata", valid: true},
		{name: "delete", op: OperationDelete, text: "delete", valid: true},
		{name: "unknown", op: Operation(99), text: "unknown", valid: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.op.String(); got != tt.text {
				t.Fatalf("String() = %q; want %q", got, tt.text)
			}
			if got := tt.op.IsValid(); got != tt.valid {
				t.Fatalf("IsValid() = %v; want %v", got, tt.valid)
			}
		})
	}
}
