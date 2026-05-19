/*
  Copyright 2026 The ARCORIS Authors

  Licensed under the Apache License, Version 2.0 (the "License");
  you may not use this file except in compliance with the License.
  You may obtain a copy of the License at

      http://www.apache.org/licenses/LICENSE-2.0

  Unless required by applicable law or agreed to in writing, software
  distributed under the License is distributed on an "AS IS" BASIS,
  WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
  See the License for the specific language governing permissions and
  limitations under the License.
*/

package bulkhead

import "testing"

func TestNewCapacitySnapshot(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		limit     uint64
		inFlight  uint64
		available uint64
		full      bool
	}{
		{name: "empty", limit: 3, inFlight: 0, available: 3, full: false},
		{name: "partial", limit: 3, inFlight: 2, available: 1, full: false},
		{name: "full", limit: 3, inFlight: 3, available: 0, full: true},
		{name: "overfull defensive", limit: 3, inFlight: 4, available: 0, full: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got := newCapacitySnapshot(tt.limit, tt.inFlight)
			if got.Available != tt.available || got.Full != tt.full {
				t.Fatalf("newCapacitySnapshot() = %+v", got)
			}
		})
	}
}

func TestCapacitySnapshotIsValid(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		snap CapacitySnapshot
		want bool
	}{
		{name: "valid empty", snap: CapacitySnapshot{Limit: 2, InFlight: 0, Available: 2, Full: false}, want: true},
		{name: "valid full", snap: CapacitySnapshot{Limit: 2, InFlight: 2, Available: 0, Full: true}, want: true},
		{name: "zero limit", snap: CapacitySnapshot{}, want: false},
		{name: "over limit", snap: CapacitySnapshot{Limit: 1, InFlight: 2, Available: 0, Full: true}, want: false},
		{name: "bad available", snap: CapacitySnapshot{Limit: 2, InFlight: 1, Available: 2, Full: false}, want: false},
		{name: "bad full flag", snap: CapacitySnapshot{Limit: 2, InFlight: 2, Available: 0, Full: false}, want: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			if got := tt.snap.IsValid(); got != tt.want {
				t.Fatalf("IsValid() = %v, want %v", got, tt.want)
			}
		})
	}
}
