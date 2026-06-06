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

func TestVectorSnapshotFitDiagnostics(t *testing.T) {
	snap := capacity.NewVectorSnapshot(
		vector(t, entry("memory_bytes", 4), entry("worker_slots", 4)),
		vector(t, entry("memory_bytes", 6), entry("worker_slots", 1)),
	)

	if !snap.IsValid() || !snap.HasDebt() {
		t.Fatalf("snapshot invalid or missing debt: %#v", snap)
	}
	if !snap.CanReserve(demand(t, entry("worker_slots", 2))) {
		t.Fatal("worker demand was blocked by unrelated memory debt")
	}

	memory := snap.Fit(demand(t, entry("memory_bytes", 1)))
	if memory.Refusal != capacity.RefusalDebt {
		t.Fatalf("memory refusal = %s, want debt", memory.Refusal)
	}
	requireVector(t, memory.Debt, entry("memory_bytes", 2))

	workers := snap.Fit(demand(t, entry("worker_slots", 4)))
	if workers.Refusal != capacity.RefusalInsufficient {
		t.Fatalf("worker refusal = %s, want insufficient", workers.Refusal)
	}
	requireVector(t, workers.Missing, entry("worker_slots", 1))

	unknown := snap.Fit(demand(t, entry("queue_slots", 1)))
	if unknown.Refusal != capacity.RefusalUnknownResource {
		t.Fatalf("unknown refusal = %s, want unknown_resource", unknown.Refusal)
	}
}
