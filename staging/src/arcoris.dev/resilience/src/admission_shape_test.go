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
	"testing"
	"time"

	"arcoris.dev/admission"
	admissionbuiltin "arcoris.dev/admissioncatalog/builtin"
	"arcoris.dev/resilience/bulkhead"
	"arcoris.dev/resilience/deadline"
	"arcoris.dev/resilience/retrybudget"
)

func TestBulkheadAdmissionResultShape(t *testing.T) {
	t.Parallel()

	b := bulkhead.New(1)
	success := b.TryAdmit(bulkhead.Request{Amount: 1})
	if !success.IsValid() || !success.Decision().IsAdmitted() {
		t.Fatalf("bulkhead success = %+v, want valid admitted", success.Decision())
	}
	if got, want := success.Decision(), admission.GrantDecision(admission.ReasonAdmitted); got != want {
		t.Fatalf("success decision = %+v, want %+v", got, want)
	}
	if !success.HasGrant() || !success.HasMetadata() {
		t.Fatalf("success grant=%t metadata=%t, want true,true", success.HasGrant(), success.HasMetadata())
	}
	lease, ok := success.Grant()
	if !ok || lease == nil {
		t.Fatal("bulkhead success did not carry a live lease")
	}

	denied := b.TryAdmit(bulkhead.Request{Amount: 1})
	if !denied.IsValid() || !denied.Decision().IsDenied() {
		t.Fatalf("bulkhead denial = %+v, want valid denied", denied.Decision())
	}
	if got, want := denied.Decision(), admission.DenyDecision(admissionbuiltin.ReasonCapacityExhausted); got != want {
		t.Fatalf("denied decision = %+v, want %+v", got, want)
	}
	if denied.HasGrant() || !denied.HasMetadata() {
		t.Fatalf("denied grant=%t metadata=%t, want false,true", denied.HasGrant(), denied.HasMetadata())
	}
	if grant, ok := denied.Grant(); ok || grant != nil {
		t.Fatalf("denied GrantDecision() = (%#v,%t), want nil,false", grant, ok)
	}

	lease.Release()
}

func TestRetryBudgetAdmissionResultShape(t *testing.T) {
	t.Parallel()

	allowedLimiter := newRetryBudget(t, 1)
	success := allowedLimiter.TryAdmit(retrybudget.Request{})
	if !success.IsValid() || !success.Decision().IsAdmitted() {
		t.Fatalf("retrybudget success = %+v, want valid admitted", success.Decision())
	}
	if got, want := success.Decision(), admission.CommitDecision(admission.ReasonAdmitted); got != want {
		t.Fatalf("success decision = %+v, want %+v", got, want)
	}
	if success.HasGrant() || !success.HasMetadata() {
		t.Fatalf("success grant=%t metadata=%t, want false,true", success.HasGrant(), success.HasMetadata())
	}
	if metadata, ok := success.Metadata(); !ok || metadata.Value.Attempts.Retry != 1 {
		t.Fatalf("success metadata = (%+v,%t), want retry spent", metadata, ok)
	}

	deniedLimiter := newRetryBudget(t, 0)
	denied := deniedLimiter.TryAdmit(retrybudget.Request{})
	if !denied.IsValid() || !denied.Decision().IsDenied() {
		t.Fatalf("retrybudget denial = %+v, want valid denied", denied.Decision())
	}
	if got, want := denied.Decision(), admission.DenyDecision(admissionbuiltin.ReasonBudgetExhausted); got != want {
		t.Fatalf("denied decision = %+v, want %+v", got, want)
	}
	if denied.HasGrant() || !denied.HasMetadata() {
		t.Fatalf("denied grant=%t metadata=%t, want false,true", denied.HasGrant(), denied.HasMetadata())
	}
	if metadata, ok := denied.Metadata(); !ok || metadata.Value.Attempts.Retry != 0 {
		t.Fatalf("denied metadata = (%+v,%t), want no retry spent", metadata, ok)
	}
}

func TestDeadlineAdmissionResultShape(t *testing.T) {
	t.Parallel()

	allowed := deadline.TryAdmit(deadline.Request{
		Context: context.Background(),
		Now:     compositionNow,
		Min:     time.Second,
	})
	if !allowed.IsValid() || !allowed.Decision().IsAdmitted() {
		t.Fatalf("deadline allowed = %+v, want valid admitted", allowed.Decision())
	}
	if got, want := allowed.Decision(), admission.AdmitDecision(admission.ReasonAdmitted); got != want {
		t.Fatalf("allowed decision = %+v, want %+v", got, want)
	}
	if allowed.HasGrant() || !allowed.HasMetadata() {
		t.Fatalf("allowed grant=%t metadata=%t, want false,true", allowed.HasGrant(), allowed.HasMetadata())
	}

	denied := deadline.TryAdmit(deadline.Request{
		Context: contextWithDeadline(t, compositionNow.Add(-time.Second)),
		Now:     compositionNow,
		Min:     time.Second,
	})
	if !denied.IsValid() || !denied.Decision().IsDenied() {
		t.Fatalf("deadline denied = %+v, want valid denied", denied.Decision())
	}
	if got, want := denied.Decision(), admission.DenyDecision(admissionbuiltin.ReasonDeadlineExceeded); got != want {
		t.Fatalf("denied decision = %+v, want %+v", got, want)
	}
	if denied.HasGrant() || !denied.HasMetadata() {
		t.Fatalf("denied grant=%t metadata=%t, want false,true", denied.HasGrant(), denied.HasMetadata())
	}
}

func TestResilienceAdmissionSurfacesUseDistinctEffectSemantics(t *testing.T) {
	t.Parallel()

	bulkheadResult := bulkhead.New(1).TryAdmit(bulkhead.Request{Amount: 1})
	budgetResult := newRetryBudget(t, 1).TryAdmit(retrybudget.Request{})
	deadlineResult := deadline.TryAdmit(deadline.Request{
		Context: context.Background(),
		Now:     compositionNow,
	})

	if got := bulkheadResult.Decision().Effect; got != admission.EffectOwned {
		t.Fatalf("bulkhead effect = %s, want owned", got)
	}
	if got := budgetResult.Decision().Effect; got != admission.EffectCommitted {
		t.Fatalf("retrybudget effect = %s, want committed", got)
	}
	if got := deadlineResult.Decision().Effect; got != admission.EffectNone {
		t.Fatalf("deadline effect = %s, want none", got)
	}

	lease, ok := bulkheadResult.Grant()
	if ok {
		lease.Release()
	}
}
