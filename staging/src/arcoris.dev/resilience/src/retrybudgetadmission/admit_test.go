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

package retrybudgetadmission

import (
	"sync"
	"sync/atomic"
	"testing"

	"arcoris.dev/resilience/retrybudget"
)

func TestAdmitterTryAdmitDelegatesOnce(t *testing.T) {
	t.Parallel()

	budget := &scriptedBudget{decisions: []retrybudget.Decision{validAllowedDecision()}}
	result := New(budget).TryAdmit(Request{})

	if !result.IsValid() || !result.Decision().IsAdmitted() {
		t.Fatalf("TryAdmit result = %+v, want valid admitted", result.Decision())
	}
	if got := budget.calls.Load(); got != 1 {
		t.Fatalf("TryAdmitRetry calls = %d, want 1", got)
	}
	metadata, ok := result.Metadata()
	if !ok || !metadata.IsAllowed() {
		t.Fatalf("metadata = (%+v,%t), want allowed decision", metadata, ok)
	}
}

func TestAdmitterTryAdmitDeniedHasNoGrant(t *testing.T) {
	t.Parallel()

	result := New(&scriptedBudget{decisions: []retrybudget.Decision{validDeniedDecision()}}).
		TryAdmit(Request{})

	if !result.IsValid() || !result.Decision().IsDenied() {
		t.Fatalf("TryAdmit result = %+v, want valid denied", result.Decision())
	}
	if result.HasGrant() {
		t.Fatal("denied retry-budget admission carried a grant")
	}
	metadata, ok := result.Metadata()
	if !ok || !metadata.IsDenied() {
		t.Fatalf("metadata = (%+v,%t), want denied decision", metadata, ok)
	}
}

func TestAdmitterConcurrentTryAdmitDoesNotOverspend(t *testing.T) {
	t.Parallel()

	const total = 100
	const allowedLimit = 25
	budget := &countingBudget{limit: allowedLimit}
	admitter := New(budget)

	var allowed atomic.Uint64
	var denied atomic.Uint64
	var wg sync.WaitGroup
	for range total {
		wg.Add(1)
		go func() {
			defer wg.Done()

			result := admitter.TryAdmit(Request{})
			if !result.IsValid() {
				t.Errorf("invalid result: %+v", result.Decision())
				return
			}
			if result.Decision().IsAdmitted() {
				allowed.Add(1)
				return
			}
			if result.Decision().IsDenied() {
				denied.Add(1)
				return
			}
			t.Errorf("unexpected result: %+v", result.Decision())
		}()
	}
	wg.Wait()

	if got := allowed.Load(); got != allowedLimit {
		t.Fatalf("allowed = %d, want %d", got, allowedLimit)
	}
	if got := denied.Load(); got != total-allowedLimit {
		t.Fatalf("denied = %d, want %d", got, total-allowedLimit)
	}
}
