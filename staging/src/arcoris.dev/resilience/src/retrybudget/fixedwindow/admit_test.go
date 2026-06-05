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

	"arcoris.dev/admission"
	admissionbuiltin "arcoris.dev/admissioncatalog/builtin"
	"arcoris.dev/resilience/retrybudget"
	"arcoris.dev/snapshot"
)

func TestLimiterTryAdmitAllowedMapsToCommittedAdmissionResult(t *testing.T) {
	limiter, _ := newTestLimiter(t, WithRatio(0), WithMinRetries(2))

	result := limiter.TryAdmit(retrybudget.Request{})
	if !result.IsValid() {
		t.Fatalf("TryAdmit result is invalid: %+v", result.Decision())
	}
	if !result.Decision().IsAdmitted() {
		t.Fatal("TryAdmit result is not admitted")
	}
	if result.Decision().IsDenied() {
		t.Fatal("TryAdmit result is denied, want admitted")
	}
	if !result.Decision().HasSideEffect() {
		t.Fatal("TryAdmit result has no committed side effect")
	}
	if result.HasGrant() {
		t.Fatal("TryAdmit result has grant, want none")
	}
	if _, ok := result.Grant(); ok {
		t.Fatal("TryAdmit GrantDecision() ok=true, want false")
	}
	if !result.HasMetadata() {
		t.Fatal("TryAdmit result has no metadata")
	}
	if got, want := result.Decision(), admission.CommitDecision(admission.ReasonAdmitted); got != want {
		t.Fatalf("decision = %+v, want %+v", got, want)
	}

	metadata, ok := result.Metadata()
	if !ok {
		t.Fatal("Metadata returned ok=false, want true")
	}
	requireValidSnapshot(t, metadata)
	if metadata.Value.Attempts.Retry != 1 {
		t.Fatalf("Retry attempts = %d, want 1", metadata.Value.Attempts.Retry)
	}
	if metadata.Value.Capacity.Available != 1 {
		t.Fatalf("Available = %d, want 1", metadata.Value.Capacity.Available)
	}
}

func TestLimiterTryAdmitDeniedMapsToBudgetExhaustedAdmissionResult(t *testing.T) {
	limiter, _ := newTestLimiter(t, WithRatio(0), WithMinRetries(1))

	allowed := limiter.TryAdmit(retrybudget.Request{})
	if !allowed.IsValid() || !allowed.Decision().IsAdmitted() {
		t.Fatalf("first TryAdmit result = %+v, want valid admitted", allowed.Decision())
	}
	prev := limiter.Revision()

	result := limiter.TryAdmit(retrybudget.Request{})
	if !result.IsValid() {
		t.Fatalf("denied TryAdmit result is invalid: %+v", result.Decision())
	}
	if !result.Decision().IsDenied() {
		t.Fatal("TryAdmit result is not denied")
	}
	if result.Decision().IsAdmitted() {
		t.Fatal("TryAdmit result is admitted, want denied")
	}
	if result.Decision().HasSideEffect() {
		t.Fatal("TryAdmit denied result has side effect")
	}
	if result.HasGrant() {
		t.Fatal("TryAdmit denied result has grant")
	}
	if _, ok := result.Grant(); ok {
		t.Fatal("TryAdmit denied GrantDecision() ok=true, want false")
	}
	if !result.HasMetadata() {
		t.Fatal("TryAdmit denied result has no metadata")
	}
	if got, want := result.Decision(), admission.DenyDecision(admissionbuiltin.ReasonBudgetExhausted); got != want {
		t.Fatalf("decision = %+v, want %+v", got, want)
	}

	metadata, ok := result.Metadata()
	if !ok {
		t.Fatal("Metadata returned ok=false, want true")
	}
	requireValidSnapshot(t, metadata)
	if metadata.Revision != prev {
		t.Fatalf("denied revision = %d, want stable %d", metadata.Revision, prev)
	}
	if metadata.Value.Attempts.Retry != 1 {
		t.Fatalf("Retry attempts = %d, want 1", metadata.Value.Attempts.Retry)
	}
}

func TestFixedWindowTryAdmitAdmissionAllowedMapsToCommittedNoGrant(t *testing.T) {
	t.Parallel()

	limiter, _ := newTestLimiter(t, WithRatio(0), WithMinRetries(1))

	result := limiter.TryAdmit(retrybudget.Request{})
	if !result.IsValid() {
		t.Fatalf("result is invalid: %+v", result.Decision())
	}
	if got, want := result.Decision(), admission.CommitDecision(admission.ReasonAdmitted); got != want {
		t.Fatalf("decision = %+v, want %+v", got, want)
	}
	if result.HasGrant() {
		t.Fatal("committed retrybudget admission carried a grant")
	}
}

func TestFixedWindowTryAdmitAdmissionDeniedMapsToDeniedNoGrant(t *testing.T) {
	t.Parallel()

	limiter, _ := newTestLimiter(t, WithRatio(0), WithMinRetries(0))

	result := limiter.TryAdmit(retrybudget.Request{})
	if !result.IsValid() {
		t.Fatalf("result is invalid: %+v", result.Decision())
	}
	if got, want := result.Decision(), admission.DenyDecision(admissionbuiltin.ReasonBudgetExhausted); got != want {
		t.Fatalf("decision = %+v, want %+v", got, want)
	}
	if result.HasGrant() {
		t.Fatal("denied retrybudget admission carried a grant")
	}
	if metadata, ok := result.Metadata(); !ok || metadata.Value.Attempts.Retry != 0 {
		t.Fatalf("metadata = (%+v,%t), want no retry spent", metadata, ok)
	}
}

func TestLimiterTryAdmitThroughAdmissionAdmitter(t *testing.T) {
	limiter, _ := newTestLimiter(t, WithRatio(0), WithMinRetries(1))
	var admitter retrybudget.AdmissionAdmitter = limiter
	var generic admission.Admitter[
		retrybudget.Request,
		admission.NoGrant,
		snapshot.Snapshot[retrybudget.Snapshot],
	] = limiter

	if result := admitter.TryAdmit(retrybudget.Request{}); !result.IsValid() {
		t.Fatalf("AdmissionAdmitter result is invalid: %+v", result.Decision())
	}

	limiter, _ = newTestLimiter(t, WithRatio(0), WithMinRetries(1))
	generic = limiter
	if result := generic.TryAdmit(retrybudget.Request{}); !result.IsValid() {
		t.Fatalf("generic admission result is invalid: %+v", result.Decision())
	}
}
