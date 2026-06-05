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

package admission

import "testing"

func TestDeferDecision(t *testing.T) {
	t.Parallel()

	requireDecision(t, DeferDecision(ReasonDeferred), Decision{
		Outcome: OutcomeDeferred,
		Reason:  ReasonDeferred,
		Effect:  EffectNone,
	})
}

func TestDeferredResult(t *testing.T) {
	t.Parallel()

	result := DeferredResult(ReasonDeferred, "metadata")
	requireResultShape(t, result, DeferDecision(ReasonDeferred), false, true)
}

func TestDeferredForResult(t *testing.T) {
	t.Parallel()

	result := DeferredForResult[string](ReasonDeferred, "metadata")
	requireResultShape(t, result, DeferDecision(ReasonDeferred), false, true)
}

func TestDeferredNoMetadataResult(t *testing.T) {
	t.Parallel()

	result := DeferredNoMetadataResult(ReasonDeferred)
	requireResultShape(t, result, DeferDecision(ReasonDeferred), false, false)
}

func TestDeferredConstructorsWithInvalidReasonReturnInvalidValues(t *testing.T) {
	t.Parallel()

	invalid := Reason("bad-reason")
	if DeferDecision(invalid).IsValid() {
		t.Fatal("DeferDecision with invalid reason is valid")
	}
	if DeferredResult(invalid, "metadata").IsValid() {
		t.Fatal("DeferredResult with invalid reason is valid")
	}
	if DeferredForResult[string](invalid, "metadata").IsValid() {
		t.Fatal("DeferredForResult with invalid reason is valid")
	}
	if DeferredNoMetadataResult(invalid).IsValid() {
		t.Fatal("DeferredNoMetadataResult with invalid reason is valid")
	}
}

func TestDeferredForResultDoesNotRetainGrantReferences(t *testing.T) {
	t.Parallel()

	type grant struct{ value string }
	result := DeferredForResult[*grant](ReasonDeferred, "metadata")
	if got, ok := result.Grant(); ok || got != nil {
		t.Fatalf("Grant() = (%v, %t), want nil,false", got, ok)
	}
}
