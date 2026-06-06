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

func TestVectorStateWithReservedIsAllOrNothing(t *testing.T) {
	state := capacity.MustVectorState(
		vector(t, entry("memory_bytes", 4), entry("worker_slots", 2)),
		vector(t, entry("memory_bytes", 1)),
	)

	next, fit := state.WithReserved(demand(t, entry("memory_bytes", 2), entry("worker_slots", 2)))
	if !fit.Fits() {
		t.Fatalf("fit = %#v, want success", fit)
	}
	requireVector(t, next.Reserved, entry("memory_bytes", 3), entry("worker_slots", 2))

	unchanged, fit := state.WithReserved(demand(t, entry("memory_bytes", 2), entry("worker_slots", 3)))
	if fit.Refusal != capacity.RefusalInsufficient {
		t.Fatalf("fit refusal = %s, want insufficient", fit.Refusal)
	}
	if !unchanged.Reserved.Equal(state.Reserved) {
		t.Fatal("refused WithReserved mutated state")
	}
}

func TestVectorStateWithoutReserved(t *testing.T) {
	state := capacity.MustVectorState(
		vector(t, entry("worker_slots", 4)),
		vector(t, entry("worker_slots", 3)),
	)

	next, ok := state.WithoutReserved(demand(t, entry("worker_slots", 2)))
	if !ok {
		t.Fatal("WithoutReserved() returned false")
	}
	requireVector(t, next.Reserved, entry("worker_slots", 1))

	unchanged, ok := state.WithoutReserved(demand(t, entry("worker_slots", 4)))
	if ok {
		t.Fatal("WithoutReserved() underflow returned true")
	}
	if !unchanged.Reserved.Equal(state.Reserved) {
		t.Fatal("underflow mutated state")
	}
}
