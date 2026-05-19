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

func TestValidLowerSnakeIdentifier(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name  string
		value string
		limit int
		want  bool
	}{
		{name: "single word", value: "bulkhead", limit: 64, want: true},
		{name: "snake case", value: "retry_budget", limit: 64, want: true},
		{name: "digit after first character", value: "queue_v2", limit: 64, want: true},
		{name: "empty", value: "", limit: 64},
		{name: "too long", value: "bulkhead", limit: 3},
		{name: "leading digit", value: "2queue", limit: 64},
		{name: "uppercase", value: "Bulkhead", limit: 64},
		{name: "hyphen", value: "retry-budget", limit: 64},
		{name: "leading underscore", value: "_budget", limit: 64},
		{name: "trailing underscore", value: "budget_", limit: 64},
		{name: "repeated underscore", value: "retry__budget", limit: 64},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			if got := validLowerSnakeIdentifier(tt.value, tt.limit); got != tt.want {
				t.Fatalf("validLowerSnakeIdentifier(%q, %d) = %v, want %v", tt.value, tt.limit, got, tt.want)
			}
		})
	}
}

func TestValidDotPathIdentifier(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name  string
		value string
		limit int
		want  bool
	}{
		{name: "single segment", value: "bulkhead", limit: 128, want: true},
		{name: "multiple segments", value: "resilience.retry_budget", limit: 128, want: true},
		{name: "empty", value: "", limit: 128},
		{name: "too long", value: "resilience.bulkhead", limit: 3},
		{name: "empty segment", value: "resilience..bulkhead", limit: 128},
		{name: "leading dot", value: ".bulkhead", limit: 128},
		{name: "trailing dot", value: "bulkhead.", limit: 128},
		{name: "invalid segment", value: "resilience.bulkhead-v2", limit: 128},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			if got := validDotPathIdentifier(tt.value, tt.limit); got != tt.want {
				t.Fatalf("validDotPathIdentifier(%q, %d) = %v, want %v", tt.value, tt.limit, got, tt.want)
			}
		})
	}
}
