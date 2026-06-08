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

package resilience_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	"arcoris.dev/admission"
	"arcoris.dev/resilience/bulkhead"
	"arcoris.dev/resilience/bulkheadadmission"
	"arcoris.dev/resilience/deadline"
	"arcoris.dev/resilience/deadlineadmission"
	"arcoris.dev/resilience/retrybudget"
	"arcoris.dev/resilience/retrybudget/fixedwindow"
	"arcoris.dev/resilience/retrybudgetadmission"
	"arcoris.dev/snapshot"
)

func TestManualCompositionReleasesBulkheadLeaseWhenRetryBudgetDenies(t *testing.T) {
	t.Parallel()

	b := bulkhead.New(1)
	budget := newRetryBudget(t, 0)

	composed := manuallyCompose(t, b, budget, allowingDeadlineRequest(), true)

	if composed.Lease != nil && !composed.Lease.Released() {
		t.Fatal("bulkhead lease was not released after retrybudget denial")
	}
	requireBulkheadSnapshot(t, b.Snapshot(), 1, 0, 1, 0)
	if !composed.BudgetResult.IsValid() {
		t.Fatalf("retrybudget result is invalid: %+v", composed.BudgetResult.Decision())
	}
	if !composed.BudgetResult.Decision().IsDenied() {
		t.Fatalf("retrybudget result should be denied: %+v", composed.BudgetResult.Decision())
	}
	if composed.BudgetResult.HasGrant() {
		t.Fatal("retrybudget denial carried a grant")
	}
}

func TestManualCompositionReleasesBulkheadLeaseWhenDeadlineDenies(t *testing.T) {
	t.Parallel()

	b := bulkhead.New(1)
	budget := newRetryBudget(t, 1)
	req := deadlineadmission.Request{
		Context: contextWithDeadline(t, compositionNow.Add(-time.Second)),
		Now:     compositionNow,
	}

	composed := manuallyCompose(t, b, budget, req, false)

	if !composed.DeadlineResult.IsValid() {
		t.Fatalf("deadline result is invalid: %+v", composed.DeadlineResult.Decision())
	}
	if !composed.DeadlineResult.Decision().IsDenied() {
		t.Fatalf("deadline result should be denied: %+v", composed.DeadlineResult.Decision())
	}
	if composed.DeadlineResult.HasGrant() {
		t.Fatal("deadline denial carried a grant")
	}
	if composed.Lease != nil && !composed.Lease.Released() {
		t.Fatal("bulkhead lease was not released after deadline denial")
	}
	requireBulkheadSnapshot(t, b.Snapshot(), 1, 0, 1, 0)
}

func TestManualCompositionDeniedBeforeBulkheadDoesNotAcquireLease(t *testing.T) {
	t.Parallel()

	b := bulkhead.New(1)
	req := deadlineadmission.Request{
		Context: contextWithDeadline(t, compositionNow.Add(-time.Second)),
		Now:     compositionNow,
	}

	deadlineResult := deadlineadmission.TryAdmit(req)
	if !deadlineResult.IsValid() || !deadlineResult.Decision().IsDenied() {
		t.Fatalf("deadline result = %+v, want valid denial", deadlineResult.Decision())
	}

	var lease *bulkhead.Lease
	if deadlineResult.Decision().IsAdmitted() {
		bulkheadResult := bulkheadadmission.New(b).
			TryAdmit(bulkheadadmission.Request{Amount: 1})
		lease, _ = bulkheadResult.Grant()
	}

	if lease != nil {
		t.Fatal("lease acquired after earlier deadline denial")
	}
	requireBulkheadSnapshot(t, b.Snapshot(), 1, 0, 1, 0)
}

func TestManualCompositionKeepsLeaseLiveWhenAllStagesAdmit(t *testing.T) {
	t.Parallel()

	b := bulkhead.New(1)
	budget := newRetryBudget(t, 1)

	composed := manuallyCompose(t, b, budget, allowingDeadlineRequest(), true)

	if composed.Lease == nil {
		t.Fatal("successful composition returned nil lease")
	}
	if composed.Lease.Released() {
		t.Fatal("successful composition released lease before returning it")
	}
	requireBulkheadSnapshot(t, b.Snapshot(), 1, 1, 0, 0)

	composed.Lease.Release()
	requireBulkheadSnapshot(t, b.Snapshot(), 1, 0, 1, 0)
}

func TestManualCompositionReleasesLeaseWhenLaterStageReturnsInvalidResult(t *testing.T) {
	t.Parallel()

	b := bulkhead.New(1)
	bulkheadResult := bulkheadadmission.New(b).
		TryAdmit(bulkheadadmission.Request{Amount: 1})
	lease, ok := bulkheadResult.Grant()
	if !ok || lease == nil {
		t.Fatal("bulkhead did not return a live lease")
	}

	releaseOnFailure := true
	defer func() {
		if releaseOnFailure {
			_, _ = lease.TryRelease()
		}
	}()

	invalid := admission.Result[admission.NoGrant, string]{}
	if invalid.IsValid() {
		t.Fatal("test setup invalid result is valid")
	}
	if !invalid.IsValid() {
		_, _ = lease.TryRelease()
		releaseOnFailure = false
	}

	if !lease.Released() {
		t.Fatal("lease was not released after invalid later result")
	}
	requireBulkheadSnapshot(t, b.Snapshot(), 1, 0, 1, 0)
}

func Example_manualAdmissionComposition_releaseOnLaterDeny() {
	b := bulkhead.New(1)
	budget, _ := fixedwindow.New(fixedwindow.WithRatio(fixedwindow.RatioZero), fixedwindow.WithMinRetries(0))

	result := bulkheadadmission.New(b).TryAdmit(bulkheadadmission.Request{Amount: 1})
	lease, ok := result.Grant()
	if !ok {
		fmt.Println("bulkhead denied")
		return
	}

	releaseOnFailure := true
	defer func() {
		if releaseOnFailure {
			_, _ = lease.TryRelease()
		}
	}()

	budgetResult := retrybudgetadmission.New(budget).
		TryAdmit(retrybudgetadmission.Request{})
	if !budgetResult.Decision().IsAdmitted() {
		fmt.Println("released after later denial")
		return
	}

	releaseOnFailure = false
	fmt.Println("admitted")

	// Output:
	// released after later denial
}

func Example_manualAdmissionComposition_returnOwnedLease() {
	b := bulkhead.New(1)
	budget, _ := fixedwindow.New(fixedwindow.WithRatio(fixedwindow.RatioZero), fixedwindow.WithMinRetries(1))

	result := bulkheadadmission.New(b).TryAdmit(bulkheadadmission.Request{Amount: 1})
	lease, ok := result.Grant()
	if !ok {
		fmt.Println("bulkhead denied")
		return
	}

	releaseOnFailure := true
	defer func() {
		if releaseOnFailure {
			_, _ = lease.TryRelease()
		}
	}()

	budgetResult := retrybudgetadmission.New(budget).
		TryAdmit(retrybudgetadmission.Request{})
	if !budgetResult.Decision().IsAdmitted() {
		fmt.Println("budget denied")
		return
	}

	releaseOnFailure = false
	defer lease.Release()

	fmt.Println("caller owns lease")

	// Output:
	// caller owns lease
}

type composedAdmission struct {
	BulkheadResult admission.Result[*bulkhead.Lease, bulkhead.Observation]
	BudgetResult   admission.Result[admission.NoGrant, retrybudget.Decision]
	DeadlineResult admission.Result[admission.NoGrant, deadline.Decision]
	Lease          *bulkhead.Lease
}

var compositionNow = time.Date(2026, 5, 27, 12, 0, 0, 0, time.UTC)

func manuallyCompose(
	t *testing.T,
	b *bulkhead.Bulkhead,
	budget *fixedwindow.Limiter,
	deadlineReq deadlineadmission.Request,
	budgetAfterBulkhead bool,
) composedAdmission {
	t.Helper()

	var out composedAdmission

	out.DeadlineResult = deadlineadmission.TryAdmit(deadlineReq)
	if !out.DeadlineResult.IsValid() || out.DeadlineResult.Decision().IsDenied() {
		return out
	}

	out.BulkheadResult = bulkheadadmission.New(b).
		TryAdmit(bulkheadadmission.Request{Amount: 1})
	if !out.BulkheadResult.IsValid() || out.BulkheadResult.Decision().IsDenied() {
		return out
	}

	lease, ok := out.BulkheadResult.Grant()
	if !ok || lease == nil {
		t.Fatalf("bulkhead admitted without a live lease: %+v", out.BulkheadResult.Decision())
	}
	out.Lease = lease

	releaseOnFailure := true
	defer func() {
		if releaseOnFailure {
			_, _ = lease.TryRelease()
		}
	}()

	if budgetAfterBulkhead {
		out.BudgetResult = retrybudgetadmission.New(budget).
			TryAdmit(retrybudgetadmission.Request{})
		if !out.BudgetResult.IsValid() || out.BudgetResult.Decision().IsDenied() {
			return out
		}
	}

	releaseOnFailure = false
	return out
}

func allowingDeadlineRequest() deadlineadmission.Request {
	return deadlineadmission.Request{
		Context: context.Background(),
		Now:     compositionNow,
		Min:     time.Second,
	}
}

func contextWithDeadline(t *testing.T, at time.Time) context.Context {
	ctx, cancel := context.WithDeadline(context.Background(), at)
	t.Cleanup(cancel)
	return ctx
}

func newRetryBudget(t *testing.T, minRetries uint64) *fixedwindow.Limiter {
	t.Helper()

	limiter, err := fixedwindow.New(
		fixedwindow.WithRatio(fixedwindow.RatioZero),
		fixedwindow.WithMinRetries(minRetries),
	)
	if err != nil {
		t.Fatalf("fixedwindow.New() error = %v", err)
	}
	return limiter
}

func requireBulkheadSnapshot(
	t *testing.T,
	snap snapshot.Snapshot[bulkhead.Snapshot],
	limit bulkhead.Amount,
	reserved bulkhead.Amount,
	available bulkhead.Amount,
	debt bulkhead.Amount,
) {
	t.Helper()

	if snap.Value.Limit != limit ||
		snap.Value.Reserved != reserved ||
		snap.Value.Available != available ||
		snap.Value.Debt != debt {
		t.Fatalf("snapshot = %+v, want limit=%d reserved=%d available=%d debt=%d",
			snap.Value,
			limit,
			reserved,
			available,
			debt,
		)
	}
}
