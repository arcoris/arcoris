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
	"sync"
	"testing"

	"arcoris.dev/capacity"
)

func TestScalarLedgerReserveSuccess(t *testing.T) {
	t.Parallel()

	ledger := capacity.NewScalarLedger(4)
	result := ledger.TryReserve(3)
	if !result.Reserved() || result.Reservation == nil {
		t.Fatalf("TryReserve() status = %s reservation=%v", result.Status, result.Reservation)
	}
	if result.Snapshot.Value != capacity.NewScalarSnapshot(4, 3) {
		t.Fatalf("snapshot = %+v, want limit=4 reserved=3", result.Snapshot.Value)
	}
}

func TestScalarLedgerReserveRefusesDebt(t *testing.T) {
	t.Parallel()

	ledger := capacity.NewScalarLedger(4)
	result := ledger.TryReserve(3)
	if !result.Reserved() {
		t.Fatalf("TryReserve() status = %s, want reserved", result.Status)
	}
	_ = ledger.SetLimit(2)

	denied := ledger.TryReserve(1)
	if denied.Status != capacity.ReserveStatusDebt {
		t.Fatalf("TryReserve() while in debt status = %s, want debt", denied.Status)
	}
}

func TestScalarLedgerInsufficientAndZeroAmount(t *testing.T) {
	t.Parallel()

	ledger := capacity.NewScalarLedger(2)
	denied := ledger.TryReserve(3)
	if denied.Status != capacity.ReserveStatusInsufficient || denied.Reservation != nil {
		t.Fatalf("denied result = %+v, want insufficient nil reservation", denied)
	}
	requirePanicIs(t, capacity.ErrZeroAmount, func() { _ = ledger.TryReserve(0) })
}

func TestScalarConcurrentReserveDoesNotOverspend(t *testing.T) {
	t.Parallel()

	ledger := capacity.NewScalarLedger(32)
	var wg sync.WaitGroup
	for i := 0; i < 128; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			_ = ledger.TryReserve(1)
		}()
	}
	wg.Wait()

	snap := ledger.Snapshot()
	if snap.Value.Reserved != 32 {
		t.Fatalf("reserved = %d, want 32", snap.Value.Reserved)
	}
}
