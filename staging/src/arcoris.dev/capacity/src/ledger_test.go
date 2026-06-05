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
	t.Parallel()

	var _ snapshot.Source[capacity.Snapshot] = (*capacity.Ledger)(nil)
	var _ snapshot.RevisionSource = (*capacity.Ledger)(nil)
}

func TestLedgerSetSameLimitsDoesNotAdvanceRevision(t *testing.T) {
	t.Parallel()

	limits := vector(t, entry("worker_slots", 4))
	ledger := capacity.NewLedger(limits)
	before := ledger.Revision()
	after := ledger.SetLimits(limits)
	if after.Revision != before {
		t.Fatalf("revision = %d, want %d", after.Revision, before)
	}
}
