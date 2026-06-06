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
	"arcoris.dev/snapshot"
)

func TestVectorLedgerImplementsSnapshotSources(t *testing.T) {
	var _ snapshot.Source[capacity.VectorSnapshot] = (*capacity.VectorLedger)(nil)
	var _ snapshot.RevisionSource = (*capacity.VectorLedger)(nil)
}

func TestVectorLedgerSetLimitsAndSnapshot(t *testing.T) {
	ledger := capacity.NewVectorLedger(vector(t, entry("memory_bytes", 8), entry("worker_slots", 4)))
	initial := ledger.Snapshot()

	ledger.SetLimits(vector(t, entry("memory_bytes", 8), entry("worker_slots", 4)))
	same := ledger.Snapshot()

	if same.Revision != initial.Revision {
		t.Fatalf("same SetLimits advanced revision: %d -> %d", initial.Revision, same.Revision)
	}

	changed := ledger.SetLimitsObserved(vector(t, entry("worker_slots", 4)))
	if changed.Revision == initial.Revision {
		t.Fatal("changed SetLimits did not advance revision")
	}
	if !ledger.Snapshot().Value.Limits.Equal(vector(t, entry("worker_slots", 4))) {
		t.Fatalf("limits after SetLimits = %#v", ledger.Snapshot().Value.Limits)
	}
}
