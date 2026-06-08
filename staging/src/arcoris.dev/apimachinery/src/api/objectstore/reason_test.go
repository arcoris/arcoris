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

package objectstore

import "testing"

func TestReasonStringAndValidity(t *testing.T) {
	tests := []struct {
		name  string
		in    Reason
		text  string
		valid bool
	}{
		{name: "not found", in: ReasonNotFound, text: "not_found", valid: true},
		{name: "already exists", in: ReasonAlreadyExists, text: "already_exists", valid: true},
		{name: "conflict", in: ReasonConflict, text: "conflict", valid: true},
		{name: "stale revision", in: ReasonStaleRevision, text: "stale_revision", valid: true},
		{name: "invalid key", in: ReasonInvalidKey, text: "invalid_key", valid: true},
		{name: "invalid state", in: ReasonInvalidState, text: "invalid_state", valid: true},
		{name: "invalid revision", in: ReasonInvalidRevision, text: "invalid_revision", valid: true},
		{name: "uninitialized store", in: ReasonUninitializedStore, text: "uninitialized_store", valid: true},
		{name: "unknown", in: 0, text: "unknown", valid: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.in.String(); got != tt.text {
				t.Fatalf("String() = %q; want %q", got, tt.text)
			}
			if got := tt.in.IsValid(); got != tt.valid {
				t.Fatalf("IsValid() = %v; want %v", got, tt.valid)
			}
		})
	}
}
