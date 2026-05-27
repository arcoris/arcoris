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

import (
	"testing"

	panicassert "arcoris.dev/testutil/panic"
)

func TestBuiltinCatalogConstructorsReturnIndependentRegistries(t *testing.T) {
	t.Parallel()

	firstReasons := NewBuiltinReasonRegistry()
	secondReasons := NewBuiltinReasonRegistry()
	if err := firstReasons.Register(testReasonDescriptor("custom_reason")); err != nil {
		t.Fatalf("Register reason returned error: %v", err)
	}
	if secondReasons.Contains("custom_reason") {
		t.Fatal("custom reason registered in one built-in registry leaked into another")
	}

	firstKinds := NewBuiltinKindRegistry()
	secondKinds := NewBuiltinKindRegistry()
	if err := firstKinds.Register(testKindDescriptor("custom_kind")); err != nil {
		t.Fatalf("Register kind returned error: %v", err)
	}
	if secondKinds.Contains("custom_kind") {
		t.Fatal("custom kind registered in one built-in registry leaked into another")
	}

	firstComponents := NewBuiltinComponentRegistry(firstKinds)
	secondComponents := NewBuiltinComponentRegistry(secondKinds)
	component := testComponentDescriptor("custom.component", "custom_kind")
	if err := firstComponents.Register(component); err != nil {
		t.Fatalf("Register component returned error: %v", err)
	}
	if secondComponents.Contains("custom.component") {
		t.Fatal("custom component registered in one built-in registry leaked into another")
	}
}

func TestOrdinaryResultValidationDoesNotRequireRegistryMembership(t *testing.T) {
	t.Parallel()

	reason := Reason("owner_defined_reason")
	decision := Decision{
		Outcome: OutcomeAdmitted,
		Reason:  reason,
		Effect:  EffectOwned,
	}
	result := resultWith(decision, someString("grant"), noneString())

	if !decision.IsValid() {
		t.Fatalf("custom syntactically valid reason should not require registry membership: %+v", decision)
	}
	if !result.IsValid() {
		t.Fatalf("custom syntactically valid reason should not require registry membership: %+v", result.Decision())
	}

	registry := NewBuiltinReasonRegistry()
	if registry.Contains(reason) {
		t.Fatal("test setup custom reason unexpectedly exists in built-in registry")
	}
}

func TestCatalogDoesNotBehaveAsGlobalSingleton(t *testing.T) {
	t.Parallel()

	first := NewBuiltinCatalog()
	second := NewBuiltinCatalog()

	if err := first.RegisterReason(testReasonDescriptor("custom_reason")); err != nil {
		t.Fatalf("RegisterReason returned error: %v", err)
	}
	if err := first.RegisterKind(testKindDescriptor("custom_kind")); err != nil {
		t.Fatalf("RegisterKind returned error: %v", err)
	}
	component := testComponentDescriptor("custom.component", "custom_kind")
	if err := first.RegisterComponent(component); err != nil {
		t.Fatalf("RegisterComponent returned error: %v", err)
	}

	if _, ok := second.Reason("custom_reason"); ok {
		t.Fatal("custom reason registered in one catalog leaked into another")
	}
	if _, ok := second.Kind("custom_kind"); ok {
		t.Fatal("custom kind registered in one catalog leaked into another")
	}
	if _, ok := second.Component("custom.component"); ok {
		t.Fatal("custom component registered in one catalog leaked into another")
	}
}

func TestRegistryZeroValueIsUsableButNilReceiverPanics(t *testing.T) {
	t.Parallel()

	var reasons ReasonRegistry
	if err := reasons.Register(testReasonDescriptor("custom_reason")); err != nil {
		t.Fatalf("zero-value ReasonRegistry Register returned error: %v", err)
	}
	if !reasons.Contains("custom_reason") {
		t.Fatal("zero-value ReasonRegistry did not contain registered reason")
	}

	var kinds KindRegistry
	if err := kinds.Register(testKindDescriptor("custom_kind")); err != nil {
		t.Fatalf("zero-value KindRegistry Register returned error: %v", err)
	}
	if !kinds.Contains("custom_kind") {
		t.Fatal("zero-value KindRegistry did not contain registered kind")
	}

	var components ComponentRegistry
	if got := components.Len(); got != 0 {
		t.Fatalf("zero-value ComponentRegistry Len = %d, want 0", got)
	}
	if _, ok := components.Lookup("custom.component"); ok {
		t.Fatal("zero-value ComponentRegistry lookup returned ok=true")
	}

	var nilReasons *ReasonRegistry
	panicassert.RequireMessage(t, nilReasonRegistryPanic, func() {
		_ = nilReasons.Len()
	})
	var nilKinds *KindRegistry
	panicassert.RequireMessage(t, nilKindRegistryPanic, func() {
		_ = nilKinds.Len()
	})
	var nilComponents *ComponentRegistry
	panicassert.RequireMessage(t, nilComponentRegistryPanic, func() {
		_ = nilComponents.Len()
	})
}
