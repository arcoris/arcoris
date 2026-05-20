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
	"testing"

	"arcoris.dev/admission"
	"arcoris.dev/resilience/retrybudget"
	"arcoris.dev/snapshot"
)

func TestBudgetTryAdmit(t *testing.T) {
	var budget Budget

	result := budget.TryAdmit(retrybudget.Request{})
	if !result.IsValid() {
		t.Fatalf("TryAdmit result is invalid: %+v", result.Decision())
	}
	if !result.IsAdmitted() {
		t.Fatal("TryAdmit result is not admitted")
	}
	if result.IsDenied() {
		t.Fatal("TryAdmit result is denied")
	}
	if !result.HasSideEffect() {
		t.Fatal("TryAdmit result has no committed side effect")
	}
	if result.HasGrant() {
		t.Fatal("TryAdmit result has grant")
	}
	if _, ok := result.Grant(); ok {
		t.Fatal("TryAdmit Grant() ok=true, want false")
	}
	if !result.HasMetadata() {
		t.Fatal("TryAdmit result has no metadata")
	}
	if got, want := result.Decision(), admission.Commit(admission.ReasonAdmitted); got != want {
		t.Fatalf("decision = %+v, want %+v", got, want)
	}

	metadata, ok := result.Metadata()
	if !ok {
		t.Fatal("Metadata returned ok=false, want true")
	}
	if !metadata.Value.IsValid() {
		t.Fatalf("metadata snapshot is invalid: %+v", metadata.Value)
	}
	if metadata.Value.Kind != retrybudget.KindNoop {
		t.Fatalf("metadata kind = %s, want %s", metadata.Value.Kind, retrybudget.KindNoop)
	}
}

func TestBudgetTryAdmitThroughAdmissionAdmitter(t *testing.T) {
	var admitter retrybudget.AdmissionAdmitter = Budget{}
	var generic admission.Admitter[
		retrybudget.Request,
		admission.NoGrant,
		snapshot.Snapshot[retrybudget.Snapshot],
	] = Budget{}

	if result := admitter.TryAdmit(retrybudget.Request{}); !result.IsValid() {
		t.Fatalf("AdmissionAdmitter result is invalid: %+v", result.Decision())
	}
	if result := generic.TryAdmit(retrybudget.Request{}); !result.IsValid() {
		t.Fatalf("generic admission result is invalid: %+v", result.Decision())
	}
}
