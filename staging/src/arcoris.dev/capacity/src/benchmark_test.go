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

var (
	benchmarkBoolSink           bool
	benchmarkReservationSink    *capacity.Reservation
	benchmarkObservationSink    capacity.Observation
	benchmarkScalarSnapshotSink snapshot.Snapshot[capacity.Snapshot]
	benchmarkVectorSnapshotSink snapshot.Snapshot[capacity.VectorSnapshot]
)

func BenchmarkLedgerRawReserveRelease(b *testing.B) {
	ledger := capacity.NewLedger(capacity.Amount(b.N + 1))

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		if !ledger.TryReserve(1) {
			b.Fatal("reserve refused")
		}
		ledger.Release(1)
	}
}

func BenchmarkLedgerRawReserveObservedReleaseObserved(b *testing.B) {
	ledger := capacity.NewLedger(capacity.Amount(b.N + 1))

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		observation, ok := ledger.TryReserveObserved(1)
		if !ok {
			b.Fatal("reserve refused")
		}

		benchmarkObservationSink = observation
		benchmarkScalarSnapshotSink = ledger.ReleaseObserved(1)
	}
}

func BenchmarkLedgerRawDeniedReserve(b *testing.B) {
	ledger := capacity.NewLedger(1)

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		if ledger.TryReserve(2) {
			b.Fatal("reserve unexpectedly succeeded")
		}
	}
}

func BenchmarkLedgerRawDebtRefusal(b *testing.B) {
	ledger := capacity.NewLedger(1)
	if !ledger.TryReserve(1) {
		b.Fatal("initial reserve refused")
	}
	ledger.SetLimit(0)
	defer ledger.Release(1)

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		if ledger.TryReserve(1) {
			b.Fatal("reserve unexpectedly succeeded")
		}
	}
}

func BenchmarkLedgerRawRelease(b *testing.B) {
	ledger := capacity.NewLedger(capacity.Amount(b.N + 1))
	for i := 0; i < b.N; i++ {
		if !ledger.TryReserve(1) {
			b.Fatal("reserve refused during setup")
		}
	}

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		ledger.Release(1)
	}
}

func BenchmarkLedgerOwnedAcquireRelease(b *testing.B) {
	ledger := capacity.NewLedger(capacity.Amount(b.N + 1))

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		reservation, ok := ledger.TryAcquire(1)
		if !ok {
			b.Fatal("acquire refused")
		}
		reservation.Release()
		benchmarkReservationSink = reservation
	}
}

func BenchmarkLedgerOwnedAcquireFailure(b *testing.B) {
	ledger := capacity.NewLedger(1)

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		reservation, ok := ledger.TryAcquire(2)
		if ok || reservation != nil {
			b.Fatal("acquire unexpectedly succeeded")
		}
	}
}

func BenchmarkLedgerOwnedTryRelease(b *testing.B) {
	ledger := capacity.NewLedger(capacity.Amount(b.N + 1))

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		reservation, ok := ledger.TryAcquire(1)
		if !ok {
			b.Fatal("acquire refused")
		}
		benchmarkBoolSink = reservation.TryRelease()
	}
}

func BenchmarkLedgerSnapshot(b *testing.B) {
	ledger := capacity.NewLedger(1024)

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		benchmarkScalarSnapshotSink = ledger.Snapshot()
	}
}

func BenchmarkLedgerParallelRawReserveRelease(b *testing.B) {
	ledger := capacity.NewLedger(1024)

	b.ReportAllocs()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			if ledger.TryReserve(1) {
				ledger.Release(1)
			}
		}
	})
}

func BenchmarkLedgerParallelOwnedAcquireRelease(b *testing.B) {
	ledger := capacity.NewLedger(1024)

	b.ReportAllocs()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			reservation, ok := ledger.TryAcquire(1)
			if ok {
				reservation.Release()
			}
		}
	})
}

func BenchmarkLedgerParallelMixedSnapshot(b *testing.B) {
	ledger := capacity.NewLedger(1024)

	b.ReportAllocs()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			if ledger.TryReserve(1) {
				ledger.Release(1)
				continue
			}
			benchmarkScalarSnapshotSink = ledger.Snapshot()
		}
	})
}

func BenchmarkVectorLedgerReserveRelease(b *testing.B) {
	ledger := capacity.NewVectorLedger(vector(b, entry("memory_bytes", uint64(b.N+1)), entry("worker_slots", uint64(b.N+1))))
	demand := demand(b, entry("memory_bytes", 1), entry("worker_slots", 1))

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		reservation, ok := ledger.TryReserve(demand)
		if !ok {
			b.Fatal("reserve refused")
		}
		reservation.Release()
	}
}

func BenchmarkVectorLedgerDeniedInsufficient(b *testing.B) {
	ledger := capacity.NewVectorLedger(vector(b, entry("worker_slots", 1)))
	demand := demand(b, entry("worker_slots", 2))

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		if reservation, ok := ledger.TryReserve(demand); ok || reservation != nil {
			b.Fatal("reserve unexpectedly succeeded")
		}
	}
}

func BenchmarkVectorLedgerDeniedDebt(b *testing.B) {
	ledger := capacity.NewVectorLedger(vector(b, entry("memory_bytes", 1)))
	reservation, ok := ledger.TryReserve(demand(b, entry("memory_bytes", 1)))
	if !ok {
		b.Fatal("initial reserve refused")
	}
	ledger.SetLimits(capacity.Vector{})
	defer reservation.Release()
	demand := demand(b, entry("memory_bytes", 1))

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		if reservation, ok := ledger.TryReserve(demand); ok || reservation != nil {
			b.Fatal("reserve unexpectedly succeeded")
		}
	}
}

func BenchmarkVectorLedgerDeniedUnknownResource(b *testing.B) {
	ledger := capacity.NewVectorLedger(vector(b, entry("worker_slots", 1)))
	demand := demand(b, entry("memory_bytes", 1))

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		if reservation, ok := ledger.TryReserve(demand); ok || reservation != nil {
			b.Fatal("reserve unexpectedly succeeded")
		}
	}
}

func BenchmarkVectorLedgerSnapshot(b *testing.B) {
	ledger := capacity.NewVectorLedger(vector(b, entry("memory_bytes", 8), entry("worker_slots", 4)))

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		benchmarkVectorSnapshotSink = ledger.Snapshot()
	}
}

func BenchmarkManyLedgersDistributed(b *testing.B) {
	ledgers := make([]*capacity.Ledger, 64)
	for i := range ledgers {
		ledgers[i] = capacity.NewLedger(1024)
	}

	b.ReportAllocs()
	b.RunParallel(func(pb *testing.PB) {
		i := 0
		for pb.Next() {
			ledger := ledgers[i%len(ledgers)]
			i++

			if ledger.TryReserve(1) {
				ledger.Release(1)
			}
		}
	})
}
