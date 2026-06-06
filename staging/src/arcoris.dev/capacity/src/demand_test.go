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

func TestDemandValidationAndCopySafety(t *testing.T) {
	if _, err := capacity.NewDemand(); !errors.Is(err, capacity.ErrEmptyDemand) {
		t.Fatalf("NewDemand() error = %v, want ErrEmptyDemand", err)
	}
	if _, err := capacity.NewDemand(capacity.Entry{Resource: "bad-name", Amount: 1}); !errors.Is(err, capacity.ErrInvalidResource) {
		t.Fatalf("NewDemand() error = %v, want ErrInvalidResource", err)
	}

	demand := demand(t, entry("worker_slots", 2))
	if !demand.IsValid() {
		t.Fatal("Demand was invalid")
	}

	entries := demand.Entries()
	entries[0] = entry("queue_slots", 1)
	requireVector(t, demand.Vector(), entry("worker_slots", 2))

	var zero capacity.Demand
	if zero.IsValid() {
		t.Fatal("zero-value Demand was valid")
	}
}
