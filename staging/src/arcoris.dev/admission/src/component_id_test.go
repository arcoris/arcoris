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

func TestComponentIDIsValid(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		id   ComponentID
		want bool
	}{
		{name: "single segment", id: "scheduler", want: true},
		{name: "two segments", id: "resilience.bulkhead", want: true},
		{name: "snake segment", id: "scheduler.tenant_fairness", want: true},
		{name: "digit after first character", id: "queue.priority_v2", want: true},
		{name: "empty", id: "", want: false},
		{name: "uppercase", id: "Bulkhead", want: false},
		{name: "slash", id: "resilience/Bulkhead", want: false},
		{name: "empty middle segment", id: "resilience..bulkhead", want: false},
		{name: "leading dot", id: ".bulkhead", want: false},
		{name: "trailing dot", id: "bulkhead.", want: false},
		{name: "punctuation", id: "resilience.bulkhead#1", want: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			if got := tt.id.IsValid(); got != tt.want {
				t.Fatalf("%q IsValid = %v, want %v", tt.id, got, tt.want)
			}
		})
	}
}

func TestMustComponentID(t *testing.T) {
	t.Parallel()

	if got := MustComponentID("resilience.bulkhead"); got != "resilience.bulkhead" {
		t.Fatalf("got %q, want resilience.bulkhead", got)
	}

	defer func() {
		if recover() == nil {
			t.Fatal("expected panic")
		}
	}()

	_ = MustComponentID("bad/id")
}

func TestComponentIDString(t *testing.T) {
	t.Parallel()

	if got := ComponentID("resilience.bulkhead").String(); got != "resilience.bulkhead" {
		t.Fatalf("String = %q, want resilience.bulkhead", got)
	}
}
