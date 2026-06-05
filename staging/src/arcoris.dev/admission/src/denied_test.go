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

func TestDenyDecision(t *testing.T) {
	t.Parallel()

	requireDecision(t, DenyDecision(ReasonDenied), Decision{
		Outcome: OutcomeDenied,
		Reason:  ReasonDenied,
		Effect:  EffectNone,
	})
}

func TestDeniedResult(t *testing.T) {
	t.Parallel()

	result := DeniedResult(ReasonDenied, "metadata")
	requireResultShape(t, result, DenyDecision(ReasonDenied), false, true)
}

func TestDeniedForResult(t *testing.T) {
	t.Parallel()

	result := DeniedForResult[string](ReasonDenied, "metadata")
	requireResultShape(t, result, DenyDecision(ReasonDenied), false, true)
}

func TestDeniedNoMetadataResult(t *testing.T) {
	t.Parallel()

	result := DeniedNoMetadataResult(ReasonDenied)
	requireResultShape(t, result, DenyDecision(ReasonDenied), false, false)
}

func TestDeniedConstructorsWithInvalidReasonReturnInvalidValues(t *testing.T) {
	t.Parallel()

	invalid := Reason("bad-reason")
	if DenyDecision(invalid).IsValid() {
		t.Fatal("DenyDecision with invalid reason is valid")
	}
	if DeniedResult(invalid, "metadata").IsValid() {
		t.Fatal("DeniedResult with invalid reason is valid")
	}
	if DeniedForResult[string](invalid, "metadata").IsValid() {
		t.Fatal("DeniedForResult with invalid reason is valid")
	}
	if DeniedNoMetadataResult(invalid).IsValid() {
		t.Fatal("DeniedNoMetadataResult with invalid reason is valid")
	}
}

func TestDeniedForResultDoesNotRetainGrantReferences(t *testing.T) {
	t.Parallel()

	type grant struct{ value string }
	result := DeniedForResult[*grant](ReasonDenied, "metadata")
	if got, ok := result.Grant(); ok || got != nil {
		t.Fatalf("Grant() = (%v, %t), want nil,false", got, ok)
	}
}
