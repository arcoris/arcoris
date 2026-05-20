/*
  Copyright 2026 The ARCORIS Authors

  Licensed under the Apache License, Version 2.0 (the "License");
  you may not use this file except in compliance with the License.
  You may obtain a copy of the License at

      http://www.apache.org/licenses/LICENSE-2.0

  Unless required by applicable law or agreed to in writing, software
  distributed under the License is distributed on an "AS IS" BASIS,
  WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
  See the License for the specific language governing permissions and
  limitations under the License.
*/

package bulkhead

import (
	"sync"
	"sync/atomic"
	"testing"

	"arcoris.dev/admission"
)

func TestBulkheadTryAdmitGrantsLease(t *testing.T) {
	t.Parallel()

	b := New(2)
	result := b.TryAdmit(Request{Amount: 1})
	if !result.IsValid() {
		t.Fatalf("result is invalid: %+v", result.Decision())
	}
	if !result.IsAdmitted() {
		t.Fatal("result is not admitted")
	}
	if result.IsDenied() {
		t.Fatal("result is denied, want admitted")
	}
	if !result.HasGrant() {
		t.Fatal("result has no grant")
	}
	if !result.HasMetadata() {
		t.Fatal("result has no metadata")
	}
	if got, want := result.Decision(), admission.Grant(admission.ReasonAdmitted); got != want {
		t.Fatalf("decision = %+v, want %+v", got, want)
	}

	lease, ok := result.Grant()
	if !ok {
		t.Fatal("Grant returned ok=false, want true")
	}
	if lease == nil {
		t.Fatal("grant lease is nil")
	}
	defer lease.Release()

	metadata, ok := result.Metadata()
	if !ok {
		t.Fatal("Metadata returned ok=false, want true")
	}
	requireSnapshotValue(t, metadata, 2, 1, 1, 0)
}

func TestBulkheadTryAdmitDeniesWhenCapacityExhausted(t *testing.T) {
	t.Parallel()

	b := New(1)
	held := b.TryAdmit(Request{Amount: 1})
	lease, ok := held.Grant()
	if !ok {
		t.Fatal("first TryAdmit returned no lease")
	}
	defer lease.Release()

	result := b.TryAdmit(Request{Amount: 1})
	if !result.IsValid() {
		t.Fatalf("denied result is invalid: %+v", result.Decision())
	}
	if !result.IsDenied() {
		t.Fatal("result is not denied")
	}
	if result.HasGrant() {
		t.Fatal("denied result has grant")
	}
	if !result.HasMetadata() {
		t.Fatal("denied result has no metadata")
	}
	if got, want := result.Decision(), admission.Deny(admission.ReasonCapacityExhausted); got != want {
		t.Fatalf("decision = %+v, want %+v", got, want)
	}

	metadata, ok := result.Metadata()
	if !ok {
		t.Fatal("Metadata returned ok=false, want true")
	}
	requireSnapshotValue(t, metadata, 1, 1, 0, 0)
}

func TestBulkheadTryAdmitWeighted(t *testing.T) {
	t.Parallel()

	b := New(3)
	first := b.TryAdmit(Request{Amount: 2})
	lease, ok := first.Grant()
	if !ok {
		t.Fatal("first TryAdmit returned no lease")
	}

	denied := b.TryAdmit(Request{Amount: 2})
	if !denied.IsDenied() {
		t.Fatal("second TryAdmit was not denied")
	}

	lease.Release()
	third := b.TryAdmit(Request{Amount: 3})
	next, ok := third.Grant()
	if !ok {
		t.Fatal("third TryAdmit returned no lease")
	}
	defer next.Release()
	if next.Amount() != 3 {
		t.Fatalf("third lease amount = %d, want 3", next.Amount())
	}
}

func TestBulkheadTryAdmitInvalidAmountPanics(t *testing.T) {
	t.Parallel()

	b := New(1)
	requirePanic(t, "capacity: reservation amount must be positive", func() {
		_ = b.TryAdmit(Request{Amount: 0})
	})
}

func TestBulkheadTryAdmitConcurrentDoesNotOverspend(t *testing.T) {
	b := New(8)

	const workers = 64
	var admitted atomic.Uint64
	var wg sync.WaitGroup
	start := make(chan struct{})
	leases := make(chan *Lease, workers)

	for range workers {
		wg.Add(1)
		go func() {
			defer wg.Done()
			<-start

			result := b.TryAdmit(Request{Amount: 1})
			if lease, ok := result.Grant(); ok {
				admitted.Add(1)
				leases <- lease
			}
		}()
	}

	close(start)
	wg.Wait()
	close(leases)

	if got := admitted.Load(); got != 8 {
		t.Fatalf("admitted = %d, want 8", got)
	}
	requireSnapshotValue(t, b.Snapshot(), 8, 8, 0, 0)

	for lease := range leases {
		lease.Release()
	}
	requireSnapshotValue(t, b.Snapshot(), 8, 0, 8, 0)
}
