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
	"sync/atomic"
	"testing"

	"arcoris.dev/capacity"
)

func BenchmarkLedgerSnapshot(b *testing.B) {
	ledger := capacity.NewLedger(1024)
	_, _, _ = ledger.TryReserve(128)

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = ledger.Snapshot()
	}
}

func BenchmarkLedgerTryReserveRelease(b *testing.B) {
	ledger := capacity.NewLedger(capacity.Amount(b.N + 1))

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		reservation, _, ok := ledger.TryReserve(1)
		if !ok {
			b.Fatal("reservation failed")
		}
		reservation.Release()
	}
}

func BenchmarkLedgerTryReserveDenied(b *testing.B) {
	ledger := capacity.NewLedger(1)
	_, _, _ = ledger.TryReserve(1)

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _, _ = ledger.TryReserve(1)
	}
}

func BenchmarkLedgerTryReserveDeniedOvercommitted(b *testing.B) {
	ledger := capacity.NewLedger(10)
	_, _, _ = ledger.TryReserve(8)
	_ = ledger.SetLimit(5)

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _, _ = ledger.TryReserve(1)
	}
}

func BenchmarkReservationTryReleaseAlreadyReleased(b *testing.B) {
	ledger := capacity.NewLedger(1)
	reservation, _, ok := ledger.TryReserve(1)
	if !ok {
		b.Fatal("reservation failed")
	}
	reservation.Release()

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = reservation.TryRelease()
	}
}

func BenchmarkLedgerSetLimit(b *testing.B) {
	ledger := capacity.NewLedger(1024)

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = ledger.SetLimit(capacity.Amount(1024 + i%2))
	}
}

func BenchmarkLedgerSetLimitSameValue(b *testing.B) {
	ledger := capacity.NewLedger(1024)

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = ledger.SetLimit(1024)
	}
}

func BenchmarkLedgerConcurrentTryReserveRelease(b *testing.B) {
	ledger := capacity.NewLedger(1024)

	b.ReportAllocs()
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			reservation, _, ok := ledger.TryReserve(1)
			if ok {
				reservation.Release()
			}
		}
	})
}

func BenchmarkLedgerSnapshotParallel(b *testing.B) {
	ledger := capacity.NewLedger(1024)
	_, _, _ = ledger.TryReserve(128)

	b.ReportAllocs()
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			_ = ledger.Snapshot()
		}
	})
}

func BenchmarkLedgerSetLimitParallel(b *testing.B) {
	ledger := capacity.NewLedger(1024)

	b.ReportAllocs()
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			_ = ledger.SetLimit(1024)
		}
	})
}

func BenchmarkLedgerMixedReserveReleaseSnapshotParallel(b *testing.B) {
	ledger := capacity.NewLedger(1024)
	var counter atomic.Uint64

	b.ReportAllocs()
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			switch counter.Add(1) % 3 {
			case 0:
				if reservation, _, ok := ledger.TryReserve(1); ok {
					_, _ = reservation.TryRelease()
				}
			case 1:
				_ = ledger.Snapshot()
			default:
				_ = ledger.SetLimit(1024)
			}
		}
	})
}
