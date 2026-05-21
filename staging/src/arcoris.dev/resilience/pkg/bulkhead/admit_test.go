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
	"fmt"
	"sync"
	"testing"

	"arcoris.dev/admission"
	"arcoris.dev/snapshot"
	panicassert "arcoris.dev/testutil/panic"
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
	if lease.Amount() != 1 {
		t.Fatalf("lease amount = %d, want 1", lease.Amount())
	}

	metadata, ok := result.Metadata()
	if !ok {
		t.Fatal("Metadata returned ok=false, want true")
	}
	requireSnapshotValue(t, metadata, 2, 1, 1, 0)

	lease.Release()
	requireSnapshotValue(t, b.Snapshot(), 2, 0, 2, 0)
}

func TestBulkheadTryAdmitDeniesWhenCapacityExhausted(t *testing.T) {
	t.Parallel()

	b := New(1)
	held := b.TryAdmit(Request{Amount: 1})
	if !held.IsValid() {
		t.Fatalf("first TryAdmit result is invalid: %+v", held.Decision())
	}
	lease, ok := held.Grant()
	if !ok {
		t.Fatal("first TryAdmit returned no lease")
	}
	if lease == nil {
		t.Fatal("first TryAdmit returned nil lease")
	}
	defer lease.Release()

	result := b.TryAdmit(Request{Amount: 1})
	if !result.IsValid() {
		t.Fatalf("denied result is invalid: %+v", result.Decision())
	}
	if !result.IsDenied() {
		t.Fatal("result is not denied")
	}
	if result.IsAdmitted() {
		t.Fatal("result is admitted, want denied")
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
	if grant, ok := result.Grant(); ok || grant != nil {
		t.Fatalf("denied grant = (%#v, %t), want (nil, false)", grant, ok)
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
	if !first.IsValid() {
		t.Fatalf("first result is invalid: %+v", first.Decision())
	}
	if !first.IsAdmitted() {
		t.Fatal("first result is not admitted")
	}
	lease, ok := first.Grant()
	if !ok {
		t.Fatal("first TryAdmit returned no lease")
	}
	if lease == nil {
		t.Fatal("first TryAdmit returned nil lease")
	}
	if lease.Amount() != 2 {
		t.Fatalf("first lease amount = %d, want 2", lease.Amount())
	}

	denied := b.TryAdmit(Request{Amount: 2})
	if !denied.IsValid() {
		t.Fatalf("denied weighted result is invalid: %+v", denied.Decision())
	}
	if !denied.IsDenied() {
		t.Fatal("second TryAdmit was not denied")
	}

	lease.Release()
	third := b.TryAdmit(Request{Amount: 3})
	if !third.IsValid() {
		t.Fatalf("third result is invalid: %+v", third.Decision())
	}
	if !third.IsAdmitted() {
		t.Fatal("third result is not admitted")
	}
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
	panicassert.RequireMessage(t, "capacity: reservation amount must be positive", func() {
		_ = b.TryAdmit(Request{Amount: 0})
	})
}

func TestBulkheadTryAdmitThroughAdmitterInterface(t *testing.T) {
	t.Parallel()

	b := New(1)
	var admitter admission.Admitter[Request, *Lease, snapshot.Snapshot[Snapshot]] = b

	result := admitter.TryAdmit(Request{Amount: 1})
	if !result.IsValid() {
		t.Fatalf("interface result is invalid: %+v", result.Decision())
	}
	lease, ok := result.Grant()
	if !ok {
		t.Fatal("interface result returned no lease")
	}
	if lease == nil {
		t.Fatal("interface result returned nil lease")
	}
	defer lease.Release()
}

func TestBulkheadTryAdmitConcurrentDoesNotOverspend(t *testing.T) {
	b := New(8)

	const workers = 64
	var wg sync.WaitGroup
	start := make(chan struct{})
	errCh := make(chan error, workers)
	leases := make(chan *Lease, workers)

	for range workers {
		wg.Add(1)
		go func() {
			defer wg.Done()
			<-start

			result := b.TryAdmit(Request{Amount: 1})
			if !result.IsValid() {
				errCh <- fmt.Errorf("invalid result: %+v", result.Decision())
				return
			}

			switch {
			case result.IsAdmitted():
				lease, ok := result.Grant()
				if !ok {
					errCh <- fmt.Errorf("admitted result has no grant: %+v", result.Decision())
					return
				}
				if lease == nil {
					errCh <- fmt.Errorf("admitted result has nil grant: %+v", result.Decision())
					return
				}
				leases <- lease

			case result.IsDenied():
				if result.HasGrant() {
					errCh <- fmt.Errorf("denied result has grant: %+v", result.Decision())
					return
				}
				if grant, ok := result.Grant(); ok || grant != nil {
					errCh <- fmt.Errorf("denied Grant() = (%#v, %t), want (nil, false)", grant, ok)
					return
				}

			default:
				errCh <- fmt.Errorf("unexpected result outcome: %+v", result.Decision())
			}
		}()
	}

	close(start)
	wg.Wait()
	close(errCh)
	close(leases)

	for err := range errCh {
		if err != nil {
			t.Fatalf("unexpected concurrent TryAdmit error: %v", err)
		}
	}

	admittedLeases := make([]*Lease, 0, 8)
	for lease := range leases {
		admittedLeases = append(admittedLeases, lease)
	}
	if got := len(admittedLeases); got != 8 {
		t.Fatalf("admitted = %d, want 8", got)
	}
	requireSnapshotValue(t, b.Snapshot(), 8, 8, 0, 0)

	for _, lease := range admittedLeases {
		lease.Release()
	}
	requireSnapshotValue(t, b.Snapshot(), 8, 0, 8, 0)
}
