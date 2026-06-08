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

package bulkhead

import "testing"

func BenchmarkTryAcquireRelease(b *testing.B) {
	b.ReportAllocs()

	bulkhead := New(1)
	for b.Loop() {
		lease, observation, ok := bulkhead.TryAcquire()
		if !ok {
			b.Fatal("TryAcquire failed")
		}
		benchmarkObservation = observation
		benchmarkSnapshot = lease.Release()
	}
}

func BenchmarkTryAcquireDenied(b *testing.B) {
	b.ReportAllocs()

	bulkhead := New(1)
	held, _, ok := bulkhead.TryAcquire()
	if !ok {
		b.Fatal("initial TryAcquire failed")
	}
	defer held.Release()

	for b.Loop() {
		lease, observation, ok := bulkhead.TryAcquire()
		if ok {
			b.Fatal("TryAcquire succeeded, want denied")
		}
		benchmarkLease = lease
		benchmarkObservation = observation
		benchmarkOK = ok
	}
}

func BenchmarkTryAcquireAmountWeightedRelease(b *testing.B) {
	b.ReportAllocs()

	bulkhead := New(8)
	for b.Loop() {
		lease, observation, ok := bulkhead.TryAcquireAmount(4)
		if !ok {
			b.Fatal("TryAcquireAmount failed")
		}
		benchmarkObservation = observation
		benchmarkSnapshot = lease.Release()
	}
}

func BenchmarkTryAcquireAmountDeniedDebt(b *testing.B) {
	b.ReportAllocs()

	bulkhead := New(4)
	held, _, ok := bulkhead.TryAcquireAmount(4)
	if !ok {
		b.Fatal("initial TryAcquireAmount failed")
	}
	defer held.Release()
	bulkhead.SetLimit(2)

	for b.Loop() {
		lease, observation, ok := bulkhead.TryAcquireAmount(1)
		if ok {
			b.Fatal("TryAcquireAmount succeeded, want debt denial")
		}
		benchmarkLease = lease
		benchmarkObservation = observation
		benchmarkOK = ok
	}
}
