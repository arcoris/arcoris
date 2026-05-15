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

package fixedwindow

import (
	"sync"
	"sync/atomic"
	"testing"

	"arcoris.dev/resilience/retrybudget"
)

func TestLimiterConcurrentUseIsRaceFree(t *testing.T) {
	limiter, _ := newTestLimiter(t, WithRatio(1), WithMinRetries(100))

	var wg sync.WaitGroup
	for range 16 {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for range 1000 {
				limiter.RecordOriginal()
				_ = limiter.TryAdmitRetry()
				_ = limiter.Snapshot()
				_ = limiter.Revision()
			}
		}()
	}
	wg.Wait()

	snap := limiter.Snapshot()
	requireValidSnapshot(t, snap)
}

func TestLimiterConcurrentRetryAdmissionDoesNotOverspend(t *testing.T) {
	limiter, _ := newTestLimiter(t, WithRatio(0), WithMinRetries(25))

	var allowed atomic.Uint64
	var denied atomic.Uint64
	var wg sync.WaitGroup
	for range 100 {
		wg.Add(1)
		go func() {
			defer wg.Done()
			decision := limiter.TryAdmitRetry()
			if decision.Reason == retrybudget.ReasonAllowed {
				allowed.Add(1)
			} else if decision.Reason == retrybudget.ReasonExhausted {
				denied.Add(1)
			} else {
				t.Errorf("unexpected reason: %s", decision.Reason)
			}
		}()
	}
	wg.Wait()

	if got := allowed.Load(); got != 25 {
		t.Fatalf("allowed retries = %d, want 25", got)
	}
	if got := denied.Load(); got != 75 {
		t.Fatalf("denied retries = %d, want 75", got)
	}

	snap := limiter.Snapshot()
	if snap.Value.Attempts.Retry != 25 {
		t.Fatalf("snapshot retry attempts = %d, want 25", snap.Value.Attempts.Retry)
	}
	if !snap.Value.Capacity.Exhausted {
		t.Fatalf("capacity = %+v, want exhausted", snap.Value.Capacity)
	}
}
