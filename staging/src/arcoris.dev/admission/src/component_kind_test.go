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

package admission

import "testing"

func TestComponentKindIsValid(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		kind ComponentKind
		want bool
	}{
		{name: "bulkhead", kind: ComponentKind("bulkhead"), want: true},
		{name: "retry budget", kind: ComponentKind("retry_budget"), want: true},
		{name: "deadline", kind: ComponentKind("deadline"), want: true},
		{name: "custom", kind: "custom_guard", want: true},
		{name: "empty", kind: "", want: false},
		{name: "uppercase", kind: "Bulkhead", want: false},
		{name: "hyphen", kind: "bulkhead-kind", want: false},
		{name: "leading underscore", kind: "_bulkhead", want: false},
		{name: "trailing underscore", kind: "bulkhead_", want: false},
		{name: "repeated underscore", kind: "bulkhead__kind", want: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			if got := tt.kind.IsValid(); got != tt.want {
				t.Fatalf("%q IsValid = %v, want %v", tt.kind, got, tt.want)
			}
		})
	}
}

func TestComponentKindString(t *testing.T) {
	t.Parallel()

	if got := ComponentKind("bulkhead").String(); got != "bulkhead" {
		t.Fatalf("String = %q, want bulkhead", got)
	}
}

func TestMustComponentKind(t *testing.T) {
	t.Parallel()

	if got := MustComponentKind("bulkhead"); got != "bulkhead" {
		t.Fatalf("MustComponentKind() = %q, want bulkhead", got)
	}

	requirePanicValue(t, "admission.ComponentKind: invalid component kind", func() {
		_ = MustComponentKind("bad-kind")
	})
}
