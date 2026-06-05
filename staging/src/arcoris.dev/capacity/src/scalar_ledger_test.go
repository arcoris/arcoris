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

func TestScalarLedgerImplementsSnapshotSources(t *testing.T) {
	t.Parallel()

	var _ snapshot.Source[capacity.ScalarSnapshot] = (*capacity.ScalarLedger)(nil)
	var _ snapshot.RevisionSource = (*capacity.ScalarLedger)(nil)
}

func TestScalarLedgerSetLimitCreatesDebt(t *testing.T) {
	t.Parallel()

	ledger := capacity.NewScalarLedger(4)
	result := ledger.TryReserve(3)
	if !result.Reserved() {
		t.Fatalf("TryReserve() status = %s, want reserved", result.Status)
	}

	debt := ledger.SetLimit(2)
	if debt.Value.Debt != 1 || debt.Value.Available != 0 {
		t.Fatalf("debt snapshot = %+v, want debt=1 available=0", debt.Value)
	}
}
