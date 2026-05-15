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

package noop

import (
	"math"
	"sync"
	"testing"

	"arcoris.dev/resilience/retrybudget"
	"arcoris.dev/snapshot"
)

func TestNewReturnsBudget(t *testing.T) {
	budget := New()

	var _ retrybudget.Budget = budget
	var _ snapshot.Source[retrybudget.Snapshot] = budget
	var _ snapshot.RevisionSource = budget
}

func TestBudgetSnapshot(t *testing.T) {
	budget := New()
	snap := budget.Snapshot()

	if snap.Revision != snapshot.ZeroRevision.Next() {
		t.Fatalf("Snapshot().Revision = %d, want %d", snap.Revision, snapshot.ZeroRevision.Next())
	}
	if !snap.Value.IsValid() {
		t.Fatalf("Snapshot().Value is invalid: %+v", snap.Value)
	}
	if snap.Value.Kind != retrybudget.KindNoop {
		t.Fatalf("Snapshot().Value.Kind = %s, want %s", snap.Value.Kind, retrybudget.KindNoop)
	}
	if snap.Value.Attempts.HasTraffic() {
		t.Fatalf("Snapshot().Value.Attempts.HasTraffic() = true, want false")
	}
	if snap.Value.Capacity.Allowed != math.MaxUint64 {
		t.Fatalf("Snapshot().Value.Capacity.Allowed = %d, want %d", snap.Value.Capacity.Allowed, uint64(math.MaxUint64))
	}
	if snap.Value.Capacity.Available != math.MaxUint64 {
		t.Fatalf("Snapshot().Value.Capacity.Available = %d, want %d", snap.Value.Capacity.Available, uint64(math.MaxUint64))
	}
	if snap.Value.Capacity.Exhausted {
		t.Fatalf("Snapshot().Value.Capacity.Exhausted = true, want false")
	}
	if snap.Value.Window.Bounded {
		t.Fatalf("Snapshot().Value.Window.Bounded = true, want false")
	}
	if snap.Value.Policy.Bounded {
		t.Fatalf("Snapshot().Value.Policy.Bounded = true, want false")
	}
}

func TestBudgetRecordOriginalDoesNotChangeSnapshot(t *testing.T) {
	budget := New()
	before := budget.Snapshot()

	budget.RecordOriginal()
	budget.RecordOriginal()
	after := budget.Snapshot()

	if before != after {
		t.Fatalf("Snapshot changed after RecordOriginal: before=%+v after=%+v", before, after)
	}
}

func TestBudgetTryAdmitRetry(t *testing.T) {
	budget := New()
	decision := budget.TryAdmitRetry()

	if !decision.IsAllowed() {
		t.Fatal("TryAdmitRetry denied retry")
	}
	if decision.Reason != retrybudget.ReasonAllowed {
		t.Fatalf("Decision.Reason = %s, want %s", decision.Reason, retrybudget.ReasonAllowed)
	}
	if !decision.IsValid() {
		t.Fatalf("Decision is invalid: %+v", decision)
	}
	if decision.Snapshot != budget.Snapshot() {
		t.Fatalf("Decision.Snapshot = %+v, want Snapshot() %+v", decision.Snapshot, budget.Snapshot())
	}
}

func TestBudgetTryAdmitRetryDoesNotChangeSnapshot(t *testing.T) {
	budget := New()
	before := budget.Snapshot()

	for i := 0; i < 16; i++ {
		if decision := budget.TryAdmitRetry(); !decision.IsAllowed() {
			t.Fatalf("TryAdmitRetry iteration %d denied retry", i)
		}
	}
	after := budget.Snapshot()

	if before != after {
		t.Fatalf("Snapshot changed after TryAdmitRetry: before=%+v after=%+v", before, after)
	}
}

func TestBudgetRevision(t *testing.T) {
	budget := New()

	if got, want := budget.Revision(), snapshot.ZeroRevision.Next(); got != want {
		t.Fatalf("Revision() = %d, want %d", got, want)
	}
}

func TestZeroValueBudget(t *testing.T) {
	var budget Budget

	if !budget.TryAdmitRetry().IsAllowed() {
		t.Fatal("zero-value Budget denied retry")
	}
	if !budget.Snapshot().Value.IsValid() {
		t.Fatalf("zero-value Budget snapshot is invalid: %+v", budget.Snapshot())
	}
}

func TestBudgetConcurrentUse(t *testing.T) {
	budget := New()
	const goroutines = 16
	const iterations = 128

	var wg sync.WaitGroup
	wg.Add(goroutines)
	for i := 0; i < goroutines; i++ {
		go func() {
			defer wg.Done()
			for j := 0; j < iterations; j++ {
				budget.RecordOriginal()
				if decision := budget.TryAdmitRetry(); !decision.IsAllowed() {
					t.Errorf("TryAdmitRetry denied retry: %+v", decision)
				}
				if snap := budget.Snapshot(); !snap.Value.IsValid() {
					t.Errorf("Snapshot invalid: %+v", snap)
				}
				if rev := budget.Revision(); rev.IsZero() {
					t.Error("Revision returned zero")
				}
			}
		}()
	}
	wg.Wait()
}
