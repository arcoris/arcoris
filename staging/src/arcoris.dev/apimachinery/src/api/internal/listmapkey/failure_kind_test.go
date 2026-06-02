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

package listmapkey

import "testing"

func TestFailureKindNamesAreStable(t *testing.T) {
	tests := []struct {
		name string
		kind FailureKind
		want string
	}{
		{name: "invalid descriptor", kind: FailureInvalidDescriptor, want: "invalid_descriptor"},
		{name: "unresolved ref", kind: FailureUnresolvedRef, want: "unresolved_ref"},
		{name: "reference cycle", kind: FailureReferenceCycle, want: "reference_cycle"},
		{name: "item kind mismatch", kind: FailureItemKindMismatch, want: "item_kind_mismatch"},
		{name: "missing key", kind: FailureMissingKey, want: "missing_key"},
		{name: "null key", kind: FailureNullKey, want: "null_key"},
		{name: "key kind mismatch", kind: FailureKeyKindMismatch, want: "key_kind_mismatch"},
		{name: "key integer range", kind: FailureKeyIntegerRange, want: "key_integer_range"},
		{name: "invalid selector", kind: FailureInvalidSelector, want: "invalid_selector"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			requireEqual(t, string(tt.kind), tt.want)
		})
	}
}
