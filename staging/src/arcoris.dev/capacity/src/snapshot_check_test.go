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

func TestSnapshotCheckSuccessAndDiagnostics(t *testing.T) {
	t.Parallel()

	snap := capacity.NewSnapshot(
		vector(t, entry("memory_bytes", 8), entry("worker_slots", 4)),
		vector(t, entry("memory_bytes", 10), entry("worker_slots", 1)),
	)

	workerDemand := demand(t, entry("worker_slots", 2))
	if result := snap.Check(workerDemand); !result.Reserved() {
		t.Fatalf("worker demand status = %s, want reserved", result.Status)
	}

	memoryDemand := demand(t, entry("memory_bytes", 1))
	if result := snap.Check(memoryDemand); result.Status != capacity.ReserveStatusDebt {
		t.Fatalf("memory demand status = %s, want debt", result.Status)
	} else {
		requireVector(t, result.Debt, entry("memory_bytes", 2))
	}

	largeWorkerDemand := demand(t, entry("worker_slots", 4))
	if result := snap.Check(largeWorkerDemand); result.Status != capacity.ReserveStatusInsufficient {
		t.Fatalf("large worker demand status = %s, want insufficient", result.Status)
	} else {
		requireVector(t, result.Missing, entry("worker_slots", 1))
	}

	unknownDemand := demand(t, entry("queue_slots", 1))
	if result := snap.Check(unknownDemand); result.Status != capacity.ReserveStatusUnknownResource {
		t.Fatalf("unknown demand status = %s, want unknown_resource", result.Status)
	} else {
		requireVector(t, result.Missing, entry("queue_slots", 1))
	}
}
