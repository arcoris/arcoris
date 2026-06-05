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

package admissioncatalog

import (
	"errors"
	"testing"

	"arcoris.dev/admission"
)

const (
	testReason         admission.Reason        = "custom_reason"
	testKind           admission.ComponentKind = "custom_kind"
	testComponent      admission.ComponentID   = "custom.component"
	testOtherReason    admission.Reason        = "other_reason"
	testOtherKind      admission.ComponentKind = "other_kind"
	testOtherComponent admission.ComponentID   = "other.component"
)

func requirePanicMessage(t *testing.T, want string, fn func()) {
	t.Helper()

	defer func() {
		got := recover()
		if got == nil {
			t.Fatalf("function did not panic; want %q", want)
		}
		if got != want {
			t.Fatalf("panic = %v, want %q", got, want)
		}
	}()
	fn()
}

func requireErrorIs[T error](t *testing.T, err error, target error) T {
	t.Helper()

	if err == nil {
		t.Fatal("error is nil")
	}
	if !errors.Is(err, target) {
		t.Fatalf("errors.Is(%v, %v) = false", err, target)
	}
	var typed T
	if !errors.As(err, &typed) {
		t.Fatalf("errors.As(%T) = false for %v", typed, err)
	}
	return typed
}

func reasonDescriptor(reason admission.Reason) ReasonDescriptor {
	return ReasonDescriptor{
		Reason:  reason,
		Summary: reason.String() + " summary",
		DeclaredCapabilities: NewCapabilitySet(
			NewOutcomeSet(OutcomeCapabilityAdmit),
			NewEffectSet(EffectCapabilityNone),
		),
	}
}

func kindDescriptor(kind admission.ComponentKind) ComponentKindDescriptor {
	return ComponentKindDescriptor{
		Kind:    kind,
		Summary: kind.String() + " summary",
		DeclaredCapabilities: NewCapabilitySet(
			NewOutcomeSet(OutcomeCapabilityAdmit),
			NewEffectSet(EffectCapabilityNone),
		),
	}
}

func componentDescriptor(id admission.ComponentID, kind admission.ComponentKind) ComponentDescriptor {
	return ComponentDescriptor{
		ID:      id,
		Kind:    kind,
		Summary: id.String() + " summary",
		DeclaredCapabilities: NewCapabilitySet(
			NewOutcomeSet(OutcomeCapabilityDeny),
			NewEffectSet(EffectCapabilityNone),
		),
	}
}

func validInput() Input {
	return Input{
		Reasons: []ReasonDescriptor{
			reasonDescriptor(testReason),
			reasonDescriptor(testOtherReason),
		},
		Kinds: []ComponentKindDescriptor{
			kindDescriptor(testKind),
			kindDescriptor(testOtherKind),
		},
		Components: []ComponentDescriptor{
			componentDescriptor(testComponent, testKind),
			componentDescriptor(testOtherComponent, testOtherKind),
		},
	}
}

func mustCatalog(t *testing.T, input Input) *Catalog {
	t.Helper()

	catalog, err := Build(input)
	if err != nil {
		t.Fatalf("Build returned error: %v", err)
	}
	return catalog
}
