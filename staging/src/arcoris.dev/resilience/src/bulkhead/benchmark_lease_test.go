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

func BenchmarkLeaseAcquireAndRelease(b *testing.B) {
	b.ReportAllocs()

	bulkhead := New(1)
	for b.Loop() {
		lease, _, ok := bulkhead.TryAcquire()
		if !ok {
			b.Fatal("TryAcquire failed")
		}
		benchmarkSnapshot = lease.Release()
	}
}

func BenchmarkLeaseAcquireAndTryRelease(b *testing.B) {
	b.ReportAllocs()

	bulkhead := New(1)
	for b.Loop() {
		lease, _, ok := bulkhead.TryAcquire()
		if !ok {
			b.Fatal("TryAcquire failed")
		}
		snap, ok := lease.TryRelease()
		if !ok {
			b.Fatal("TryRelease failed")
		}
		benchmarkSnapshot = snap
	}
}

func BenchmarkLeaseTryReleaseDuplicate(b *testing.B) {
	b.ReportAllocs()

	bulkhead := New(1)
	lease, _, ok := bulkhead.TryAcquire()
	if !ok {
		b.Fatal("TryAcquire failed")
	}
	_, _ = lease.TryRelease()

	for b.Loop() {
		snap, ok := lease.TryRelease()
		if ok {
			b.Fatal("duplicate TryRelease succeeded")
		}
		benchmarkSnapshot = snap
		benchmarkOK = ok
	}
}
