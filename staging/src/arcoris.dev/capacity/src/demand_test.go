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
	"errors"
	"testing"

	"arcoris.dev/capacity"
)

func TestNewDemandRequiresNonEmptyVector(t *testing.T) {
	t.Parallel()

	_, err := capacity.NewDemand()
	if !errors.Is(err, capacity.ErrEmptyDemand) {
		t.Fatalf("NewDemand() error = %v, want ErrEmptyDemand", err)
	}
}

func TestDemandAccessorsAreCopySafe(t *testing.T) {
	t.Parallel()

	d := demand(t, entry("worker_slots", 2), entry("memory_bytes", 8))
	if !d.IsValid() || d.Len() != 2 {
		t.Fatalf("demand valid=%v len=%d, want true/2", d.IsValid(), d.Len())
	}

	entries := d.Entries()
	entries[0] = entry("memory_bytes", 99)
	requireEntries(t, d.Entries(), entry("memory_bytes", 8), entry("worker_slots", 2))

	v := d.Vector()
	if got := v.Amount(capacity.MustResource("worker_slots")); got != 2 {
		t.Fatalf("Vector().Amount(worker_slots) = %d, want 2", got)
	}
}

func TestDemandRejectsInvalidVectorInput(t *testing.T) {
	t.Parallel()

	_, err := capacity.NewDemand(capacity.Entry{Resource: "bad-name", Amount: 1})
	if !errors.Is(err, capacity.ErrInvalidResource) {
		t.Fatalf("NewDemand() error = %v, want ErrInvalidResource", err)
	}
}

func TestZeroDemandInvalid(t *testing.T) {
	t.Parallel()

	var d capacity.Demand
	if d.IsValid() {
		t.Fatal("zero Demand is valid")
	}
}
