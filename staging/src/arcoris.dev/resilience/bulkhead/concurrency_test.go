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
)

func TestLimiterConcurrentTryAcquireDoesNotExceedLimit(t *testing.T) {
	l, err := New(8)
	if err != nil {
		t.Fatalf("New returned error: %v", err)
	}

	const workers = 64
	var wg sync.WaitGroup
	var acquired atomic.Uint64

	start := make(chan struct{})
	permits := make(chan *Permit, workers)
	for range workers {
		wg.Add(1)
		go func() {
			defer wg.Done()
			<-start

			permit, dec := l.TryAcquire()
			if dec.Allowed {
				if permit == nil {
					t.Errorf("allowed decision returned nil permit")
					return
				}
				acquired.Add(1)
				permits <- permit
				return
			}
			if permit != nil {
				t.Errorf("denied decision returned non-nil permit")
			}
		}()
	}

	close(start)
	wg.Wait()
	close(permits)

	if got := acquired.Load(); got != 8 {
		t.Fatalf("acquired = %d, want 8", got)
	}

	snap := l.Snapshot()
	if got := snap.Value.Capacity.InFlight; got != 8 {
		t.Fatalf("InFlight = %d, want 8", got)
	}
	if got := snap.Value.Stats.Rejected; got != workers-8 {
		t.Fatalf("Rejected = %d, want %d", got, workers-8)
	}

	for permit := range permits {
		permit.Release()
	}

	snap = l.Snapshot()
	if got := snap.Value.Capacity.InFlight; got != 0 {
		t.Fatalf("InFlight after release = %d, want 0", got)
	}
	if got := snap.Value.Stats.Released; got != 8 {
		t.Fatalf("Released = %d, want 8", got)
	}
}

func TestPermitConcurrentReleaseIsOnce(t *testing.T) {
	l, err := New(1)
	if err != nil {
		t.Fatalf("New returned error: %v", err)
	}

	permit, dec := l.TryAcquire()
	if permit == nil || !dec.Allowed {
		t.Fatalf("TryAcquire = %v %+v, want allowed", permit, dec)
	}

	const releasers = 32
	var wg sync.WaitGroup
	start := make(chan struct{})
	for range releasers {
		wg.Add(1)
		go func() {
			defer wg.Done()
			<-start
			permit.Release()
		}()
	}

	close(start)
	wg.Wait()

	snap := l.Snapshot()
	if got := snap.Value.Stats.Released; got != 1 {
		t.Fatalf("Released = %d, want 1", got)
	}
	if got := snap.Value.Capacity.InFlight; got != 0 {
		t.Fatalf("InFlight = %d, want 0", got)
	}
}
