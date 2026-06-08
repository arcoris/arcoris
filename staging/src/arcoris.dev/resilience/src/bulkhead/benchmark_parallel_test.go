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

func BenchmarkParallelAcquireRelease(b *testing.B) {
	b.ReportAllocs()

	bulkhead := New(1024)
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			lease, observation, ok := bulkhead.TryAcquire()
			if !ok {
				benchmarkOK = ok
				continue
			}
			benchmarkObservation = observation
			benchmarkSnapshot = lease.Release()
		}
	})
}

func BenchmarkParallelAcquireDenied(b *testing.B) {
	b.ReportAllocs()

	bulkhead := New(0)
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			lease, observation, ok := bulkhead.TryAcquire()
			if ok {
				b.Fatal("TryAcquire succeeded, want denied")
			}
			benchmarkLease = lease
			benchmarkObservation = observation
		}
	})
}

func BenchmarkParallelSetLimitAcquire(b *testing.B) {
	b.ReportAllocs()

	bulkhead := New(128)
	b.RunParallel(func(pb *testing.PB) {
		var local uint64
		for pb.Next() {
			local++
			if local%16 == 0 {
				benchmarkSnapshot = bulkhead.SetLimit(128)
				continue
			}

			lease, observation, ok := bulkhead.TryAcquire()
			if ok {
				benchmarkSnapshot = lease.Release()
			}
			benchmarkObservation = observation
		}
	})
}

func BenchmarkParallelTryReleaseSameLease(b *testing.B) {
	b.ReportAllocs()

	bulkhead := New(1)
	lease, _, ok := bulkhead.TryAcquire()
	if !ok {
		b.Fatal("TryAcquire failed")
	}
	b.Cleanup(func() {
		_, _ = lease.TryRelease()
	})

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			snap, ok := lease.TryRelease()
			benchmarkSnapshot = snap
			benchmarkOK = ok
		}
	})
}
