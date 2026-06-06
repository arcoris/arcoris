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

import (
	"sync"
	"sync/atomic"
	"testing"
)

func TestConcurrentAcquireDoesNotOverspend(t *testing.T) {
	b := New(8)

	const workers = 64
	var acquired atomic.Uint64
	var wg sync.WaitGroup
	start := make(chan struct{})
	leases := make(chan *Lease, workers)

	for range workers {
		wg.Add(1)
		go func() {
			defer wg.Done()
			<-start

			lease, _, ok := b.TryAcquire()
			if ok {
				acquired.Add(1)
				leases <- lease
			}
		}()
	}

	close(start)
	wg.Wait()
	close(leases)

	if got := acquired.Load(); got != 8 {
		t.Fatalf("acquired = %d, want 8", got)
	}
	requireSnapshotValue(t, b.Snapshot(), 8, 8, 0, 0)

	for lease := range leases {
		lease.Release()
	}
	requireSnapshotValue(t, b.Snapshot(), 8, 0, 8, 0)
}

func TestConcurrentReleaseIsRaceSafe(t *testing.T) {
	b := New(32)
	leases := make([]*Lease, 0, 32)
	for range 32 {
		lease, _, ok := b.TryAcquire()
		if !ok {
			t.Fatal("TryAcquire failed")
		}
		leases = append(leases, lease)
	}

	var wg sync.WaitGroup
	for _, lease := range leases {
		wg.Add(1)
		go func(l *Lease) {
			defer wg.Done()
			l.Release()
		}(lease)
	}
	wg.Wait()

	requireSnapshotValue(t, b.Snapshot(), 32, 0, 32, 0)
}

func TestBulkheadLeaseTryReleaseConcurrentWithReleasedIsRaceFree(t *testing.T) {
	b := New(1)
	lease, _, ok := b.TryAcquire()
	if !ok {
		t.Fatal("TryAcquire failed")
	}
	beforeRelease := b.Revision()

	var released atomic.Uint64
	var wg sync.WaitGroup
	for range 64 {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for range 100 {
				_ = lease.Released()
				if _, ok := lease.TryRelease(); ok {
					released.Add(1)
				}
			}
		}()
	}
	wg.Wait()

	if got := released.Load(); got != 1 {
		t.Fatalf("successful releases = %d, want 1", got)
	}
	if !lease.Released() {
		t.Fatal("lease is not released")
	}
	requireSnapshotValue(t, b.Snapshot(), 1, 0, 1, 0)
	afterRelease := b.Revision()
	if afterRelease == beforeRelease {
		t.Fatal("successful concurrent TryRelease did not advance revision")
	}

	if _, ok := lease.TryRelease(); ok {
		t.Fatal("TryRelease after concurrent release returned ok=true")
	}
	if got := b.Revision(); got != afterRelease {
		t.Fatalf("duplicate TryRelease advanced revision: got %d, want %d", got, afterRelease)
	}
}
