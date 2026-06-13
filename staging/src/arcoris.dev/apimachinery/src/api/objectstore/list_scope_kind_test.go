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

func TestListScopeKindStringAndValidity(t *testing.T) {
	tests := []struct {
		name  string
		kind  ListScopeKind
		text  string
		valid bool
	}{
		{name: "all", kind: ListScopeAll, text: "all", valid: true},
		{name: "namespace", kind: ListScopeNamespace, text: "namespace", valid: true},
		{name: "unknown", kind: 0, text: "unknown", valid: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.kind.String(); got != tt.text {
				t.Fatalf("String() = %q; want %q", got, tt.text)
			}
			if got := tt.kind.IsValid(); got != tt.valid {
				t.Fatalf("IsValid() = %v; want %v", got, tt.valid)
			}
		})
	}
}
