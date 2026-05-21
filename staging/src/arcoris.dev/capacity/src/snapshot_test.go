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


package capacity_test

import (
	"testing"

	"arcoris.dev/capacity"
)

func TestSnapshotIsValid(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name  string
		snap  capacity.Snapshot
		valid bool
	}{
		{name: "zero", snap: capacity.Snapshot{}, valid: true},
		{name: "empty limit", snap: capacity.Snapshot{Limit: 10, Reserved: 0, Available: 10, Debt: 0}, valid: true},
		{name: "partially reserved", snap: capacity.Snapshot{Limit: 10, Reserved: 4, Available: 6, Debt: 0}, valid: true},
		{name: "fully reserved", snap: capacity.Snapshot{Limit: 10, Reserved: 10, Available: 0, Debt: 0}, valid: true},
		{name: "overcommitted", snap: capacity.Snapshot{Limit: 10, Reserved: 12, Available: 0, Debt: 2}, valid: true},
		{name: "wrong available", snap: capacity.Snapshot{Limit: 10, Reserved: 4, Available: 5, Debt: 0}, valid: false},
		{name: "unexpected debt", snap: capacity.Snapshot{Limit: 10, Reserved: 4, Available: 6, Debt: 1}, valid: false},
		{name: "available while overcommitted", snap: capacity.Snapshot{Limit: 10, Reserved: 12, Available: 1, Debt: 2}, valid: false},
		{name: "wrong debt", snap: capacity.Snapshot{Limit: 10, Reserved: 12, Available: 0, Debt: 1}, valid: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			if got := tt.snap.IsValid(); got != tt.valid {
				t.Fatalf("IsValid() = %t, want %t for %+v", got, tt.valid, tt.snap)
			}
		})
	}
}

func TestSnapshotCanReserve(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name   string
		snap   capacity.Snapshot
		amount capacity.Amount
		want   bool
	}{
		{name: "valid snapshot", snap: capacity.Snapshot{Limit: 10, Reserved: 4, Available: 6}, amount: 5, want: true},
		{name: "zero amount", snap: capacity.Snapshot{Limit: 10, Reserved: 4, Available: 6}, amount: 0, want: false},
		{name: "available equal to amount", snap: capacity.Snapshot{Limit: 10, Reserved: 4, Available: 6}, amount: 6, want: true},
		{name: "insufficient available", snap: capacity.Snapshot{Limit: 10, Reserved: 4, Available: 6}, amount: 7, want: false},
		{name: "debt refuses reservation", snap: capacity.Snapshot{Limit: 5, Reserved: 8, Available: 0, Debt: 3}, amount: 1, want: false},
		{name: "invalid snapshot", snap: capacity.Snapshot{Limit: 10, Reserved: 4, Available: 100, Debt: 0}, amount: 1, want: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			if got := tt.snap.CanReserve(tt.amount); got != tt.want {
				t.Fatalf("CanReserve(%d) = %t, want %t", tt.amount, got, tt.want)
			}
		})
	}
}

func TestSnapshotStatePredicates(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name          string
		snap          capacity.Snapshot
		exhausted     bool
		overcommitted bool
	}{
		{name: "available", snap: capacity.Snapshot{Limit: 10, Reserved: 4, Available: 6}, exhausted: false, overcommitted: false},
		{name: "exhausted", snap: capacity.Snapshot{Limit: 10, Reserved: 10, Available: 0}, exhausted: true, overcommitted: false},
		{name: "overcommitted", snap: capacity.Snapshot{Limit: 5, Reserved: 8, Available: 0, Debt: 3}, exhausted: true, overcommitted: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			if got := tt.snap.Exhausted(); got != tt.exhausted {
				t.Fatalf("Exhausted() = %t, want %t", got, tt.exhausted)
			}
			if got := tt.snap.Overcommitted(); got != tt.overcommitted {
				t.Fatalf("Overcommitted() = %t, want %t", got, tt.overcommitted)
			}
		})
	}
}
