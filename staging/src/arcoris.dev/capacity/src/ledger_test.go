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

func TestLedgerImplementsSnapshotSources(t *testing.T) {
	var _ snapshot.Source[capacity.Snapshot] = (*capacity.Ledger)(nil)
	var _ snapshot.RevisionSource = (*capacity.Ledger)(nil)
}

func TestLedgerSnapshotRevisionAndSetLimit(t *testing.T) {
	ledger := capacity.NewLedger(4)
	initial := ledger.Snapshot()

	if initial.Value != capacity.NewSnapshot(4, 0) {
		t.Fatalf("initial snapshot = %#v", initial.Value)
	}

	ledger.SetLimit(4)
	noChange := ledger.Snapshot()

	if noChange.Revision != initial.Revision {
		t.Fatalf("same SetLimit advanced revision: %d -> %d", initial.Revision, noChange.Revision)
	}

	changed := ledger.SetLimitObserved(2)

	if changed.Value != capacity.NewSnapshot(2, 0) {
		t.Fatalf("changed snapshot = %#v", changed.Value)
	}
	if changed.Revision == initial.Revision {
		t.Fatal("changed SetLimit did not advance revision")
	}
}
