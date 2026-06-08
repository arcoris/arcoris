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

package fixedwindow

import (
	"testing"
	"time"

	"arcoris.dev/resilience/retrybudget"
)

func TestLimiterTryAdmitRetryAllowed(t *testing.T) {
	l, _ := newTestLimiter(t, WithRatio(RatioZero), WithMinRetries(2))

	first := l.TryAdmitRetry()
	requireDecision(t, first, true, retrybudget.ReasonAllowed)
	if first.Snapshot.Value.Attempts.Retry != 1 {
		t.Fatalf("Retry attempts = %d, want 1", first.Snapshot.Value.Attempts.Retry)
	}
	if first.Snapshot.Value.Capacity.Available != 1 {
		t.Fatalf("Available = %d, want 1", first.Snapshot.Value.Capacity.Available)
	}

	second := l.TryAdmitRetry()
	requireDecision(t, second, true, retrybudget.ReasonAllowed)
	if second.Snapshot.Value.Attempts.Retry != 2 {
		t.Fatalf("Retry attempts = %d, want 2", second.Snapshot.Value.Attempts.Retry)
	}
	if !second.Snapshot.Value.Capacity.Exhausted {
		t.Fatalf("Capacity = %+v, want exhausted", second.Snapshot.Value.Capacity)
	}
}

func TestLimiterTryAdmitRetryMinimumAvailableBeforeOriginalTraffic(t *testing.T) {
	l, _ := newTestLimiter(t, WithRatio(RatioZero), WithMinRetries(1))

	decision := l.TryAdmitRetry()
	requireDecision(t, decision, true, retrybudget.ReasonAllowed)
	if got := decision.Snapshot.Value.Attempts.Original; got != 0 {
		t.Fatalf("Original attempts = %d, want 0", got)
	}
	if got := decision.Snapshot.Value.Attempts.Retry; got != 1 {
		t.Fatalf("Retry attempts = %d, want 1", got)
	}
}

func TestLimiterTryAdmitRetryZeroMinimumRequiresTrafficAllowance(t *testing.T) {
	l, _ := newTestLimiter(t, WithRatio(RatioOne), WithMinRetries(0))

	denied := l.TryAdmitRetry()
	requireDecision(t, denied, false, retrybudget.ReasonExhausted)

	l.RecordOriginal()
	allowed := l.TryAdmitRetry()
	requireDecision(t, allowed, true, retrybudget.ReasonAllowed)
}

func TestLimiterTryAdmitRetryDenied(t *testing.T) {
	l, _ := newTestLimiter(t, WithRatio(RatioZero), WithMinRetries(1))

	allowed := l.TryAdmitRetry()
	requireDecision(t, allowed, true, retrybudget.ReasonAllowed)
	prev := l.Revision()

	denied := l.TryAdmitRetry()
	requireDecision(t, denied, false, retrybudget.ReasonExhausted)
	if denied.Snapshot.Revision != prev {
		t.Fatalf("denied revision = %d, want stable %d", denied.Snapshot.Revision, prev)
	}
	if l.Revision() != prev {
		t.Fatalf("limiter revision = %d, want stable %d", l.Revision(), prev)
	}
	if denied.Snapshot.Value.Attempts.Retry != 1 {
		t.Fatalf("Retry attempts = %d, want 1", denied.Snapshot.Value.Attempts.Retry)
	}
}

func TestFixedWindowTryAdmitRetryDeniedDoesNotSpendRetry(t *testing.T) {
	l, _ := newTestLimiter(t, WithRatio(RatioZero), WithMinRetries(1))

	allowed := l.TryAdmitRetry()
	requireDecision(t, allowed, true, retrybudget.ReasonAllowed)

	denied := l.TryAdmitRetry()
	requireDecision(t, denied, false, retrybudget.ReasonExhausted)
	if got, want := denied.Snapshot.Value.Attempts.Retry, allowed.Snapshot.Value.Attempts.Retry; got != want {
		t.Fatalf("denied retry attempts = %d, want unchanged %d", got, want)
	}
	if got, want := l.Snapshot().Value.Attempts.Retry, allowed.Snapshot.Value.Attempts.Retry; got != want {
		t.Fatalf("published retry attempts = %d, want unchanged %d", got, want)
	}
}

func TestFixedWindowTryAdmitRetryAllowedSpendsBeforeReturning(t *testing.T) {
	l, _ := newTestLimiter(t, WithRatio(RatioZero), WithMinRetries(1))

	decision := l.TryAdmitRetry()
	requireDecision(t, decision, true, retrybudget.ReasonAllowed)
	if got := decision.Snapshot.Value.Attempts.Retry; got != 1 {
		t.Fatalf("decision retry attempts = %d, want 1", got)
	}
	if got := l.Snapshot().Value.Attempts.Retry; got != 1 {
		t.Fatalf("published retry attempts = %d, want 1", got)
	}
}

func TestLimiterTryAdmitRetryDeniedAfterRotationPublishesRotatedSnapshot(t *testing.T) {
	l, clk := newTestLimiter(t, WithWindow(time.Second), WithRatio(RatioZero), WithMinRetries(1))

	allowed := l.TryAdmitRetry()
	requireDecision(t, allowed, true, retrybudget.ReasonAllowed)
	prev := l.Revision()

	clk.Add(time.Second)
	rotatedAllowed := l.TryAdmitRetry()
	requireDecision(t, rotatedAllowed, true, retrybudget.ReasonAllowed)
	if rotatedAllowed.Snapshot.Revision == prev {
		t.Fatalf("rotation did not advance revision: %d", rotatedAllowed.Snapshot.Revision)
	}
	if rotatedAllowed.Snapshot.Value.Attempts.Retry != 1 || rotatedAllowed.Snapshot.Value.Attempts.Original != 0 {
		t.Fatalf("rotated attempts = %+v, want retry=1 original=0", rotatedAllowed.Snapshot.Value.Attempts)
	}

	prev = l.Revision()
	denied := l.TryAdmitRetry()
	requireDecision(t, denied, false, retrybudget.ReasonExhausted)
	if denied.Snapshot.Revision != prev {
		t.Fatalf("denied revision = %d, want stable %d", denied.Snapshot.Revision, prev)
	}
}

func TestLimiterTryAdmitRetryDeniedAfterRotationWithZeroMinimum(t *testing.T) {
	l, clk := newTestLimiter(t, WithWindow(time.Second), WithRatio(RatioZero), WithMinRetries(0))
	prev := l.Revision()

	clk.Add(time.Second)
	denied := l.TryAdmitRetry()
	requireDecision(t, denied, false, retrybudget.ReasonExhausted)
	if denied.Snapshot.Revision == prev {
		t.Fatalf("rotation denial did not publish new revision")
	}
	if !denied.Snapshot.Value.Window.StartedAt.Equal(fixedWindowTestNow.Add(time.Second)) {
		t.Fatalf("Window.StartedAt = %s, want rotated start", denied.Snapshot.Value.Window.StartedAt)
	}
	if got := denied.Snapshot.Value.Attempts.Retry; got != 0 {
		t.Fatalf("denied rotated retry attempts = %d, want 0", got)
	}
}
