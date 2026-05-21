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
	"errors"
	"fmt"
	"sync"
	"testing"
)

func TestNewCatalog(t *testing.T) {
	t.Parallel()

	kinds := NewBuiltinKindRegistry()
	catalog, err := NewCatalog(
		NewBuiltinReasonRegistry(),
		kinds,
		NewBuiltinComponentRegistry(kinds),
	)
	if err != nil {
		t.Fatalf("NewCatalog returned error: %v", err)
	}
	if catalog.LenReasons() == 0 {
		t.Fatal("catalog should expose reason registry")
	}
	if catalog.LenKinds() == 0 {
		t.Fatal("catalog should expose kind registry")
	}
	if catalog.LenComponents() == 0 {
		t.Fatal("catalog should expose component registry")
	}
}

func TestCatalogLookupMethods(t *testing.T) {
	t.Parallel()

	catalog := testCatalog()

	if got, ok := catalog.Reason(ReasonDenied); !ok || got.Reason != ReasonDenied {
		t.Fatalf("Reason = (%+v, %v), want descriptor,true", got, ok)
	}
	if got, ok := catalog.Kind(KindBulkhead); !ok || got.Kind != KindBulkhead {
		t.Fatalf("Kind = (%+v, %v), want descriptor,true", got, ok)
	}
	if got, ok := catalog.Component("resilience.bulkhead"); !ok || got.ID != "resilience.bulkhead" {
		t.Fatalf("Component = (%+v, %v), want descriptor,true", got, ok)
	}
}

func TestCatalogListMethodsReturnSortedCopies(t *testing.T) {
	t.Parallel()

	catalog := testCatalog()

	reasons := catalog.Reasons()
	if len(reasons) == 0 {
		t.Fatal("Reasons should not be empty")
	}
	for i := 1; i < len(reasons); i++ {
		if reasons[i-1].Reason.String() > reasons[i].Reason.String() {
			t.Fatalf("Reasons order[%d:%d] = %q,%q, want sorted",
				i-1,
				i,
				reasons[i-1].Reason,
				reasons[i].Reason,
			)
		}
	}
	reasons[0].Reason = "mutated_reason"
	if catalog.LenReasons() == 0 || catalog.Reasons()[0].Reason == "mutated_reason" {
		t.Fatal("mutating Reasons result should not mutate catalog")
	}

	kinds := catalog.Kinds()
	if len(kinds) == 0 {
		t.Fatal("Kinds should not be empty")
	}
	for i := 1; i < len(kinds); i++ {
		if kinds[i-1].Kind.String() > kinds[i].Kind.String() {
			t.Fatalf("Kinds order[%d:%d] = %q,%q, want sorted",
				i-1,
				i,
				kinds[i-1].Kind,
				kinds[i].Kind,
			)
		}
	}
	kinds[0].Kind = "mutated_kind"
	if catalog.LenKinds() == 0 || catalog.Kinds()[0].Kind == "mutated_kind" {
		t.Fatal("mutating Kinds result should not mutate catalog")
	}

	components := catalog.Components()
	if len(components) == 0 {
		t.Fatal("Components should not be empty")
	}
	for i := 1; i < len(components); i++ {
		if components[i-1].ID.String() > components[i].ID.String() {
			t.Fatalf("Components order[%d:%d] = %q,%q, want sorted",
				i-1,
				i,
				components[i-1].ID,
				components[i].ID,
			)
		}
	}
	components[0].ID = "resilience.mutated"
	if catalog.LenComponents() == 0 || catalog.Components()[0].ID == "resilience.mutated" {
		t.Fatal("mutating Components result should not mutate catalog")
	}
}

func TestCatalogRegisterMethodsDelegateToRegistries(t *testing.T) {
	t.Parallel()

	reasons := MustReasonRegistry()
	kinds := MustKindRegistry(testKindDescriptor("custom_kind"))
	components := MustComponentRegistry(kinds)

	catalog, err := NewCatalog(reasons, kinds, components)
	if err != nil {
		t.Fatalf("NewCatalog returned error: %v", err)
	}

	reason := testReasonDescriptor("custom_reason")
	if err := catalog.RegisterReason(reason); err != nil {
		t.Fatalf("RegisterReason returned error: %v", err)
	}
	if got, ok := catalog.Reason("custom_reason"); !ok || got != reason {
		t.Fatalf("Reason lookup = (%+v, %v), want registered descriptor", got, ok)
	}

	kind := testKindDescriptor("second_kind")
	if err := catalog.RegisterKind(kind); err != nil {
		t.Fatalf("RegisterKind returned error: %v", err)
	}
	if got, ok := catalog.Kind("second_kind"); !ok || got != kind {
		t.Fatalf("Kind lookup = (%+v, %v), want registered descriptor", got, ok)
	}

	component := testComponentDescriptor("custom.component", "second_kind")
	if err := catalog.RegisterComponent(component); err != nil {
		t.Fatalf("RegisterComponent returned error: %v", err)
	}
	if got, ok := catalog.Component("custom.component"); !ok || got != component {
		t.Fatalf("Component lookup = (%+v, %v), want registered descriptor", got, ok)
	}
}

func TestCatalogRegisterKindThenRegisterComponent(t *testing.T) {
	t.Parallel()

	reasons := MustReasonRegistry()
	kinds := MustKindRegistry()
	components := MustComponentRegistry(kinds)
	catalog, err := NewCatalog(reasons, kinds, components)
	if err != nil {
		t.Fatalf("NewCatalog returned error: %v", err)
	}

	kind := testKindDescriptor("custom_gate")
	component := testComponentDescriptor("custom.gate", "custom_gate")

	if err := catalog.RegisterKind(kind); err != nil {
		t.Fatalf("RegisterKind returned error: %v", err)
	}
	if err := catalog.RegisterComponent(component); err != nil {
		t.Fatalf("RegisterComponent returned error: %v", err)
	}
	if got, ok := catalog.Component("custom.gate"); !ok || got != component {
		t.Fatalf("Component lookup = (%+v, %v), want registered descriptor", got, ok)
	}
}

func TestCatalogRegisterComponentRejectsUnknownKind(t *testing.T) {
	t.Parallel()

	kinds := MustKindRegistry()
	components := MustComponentRegistry(kinds)
	catalog, err := NewCatalog(MustReasonRegistry(), kinds, components)
	if err != nil {
		t.Fatalf("NewCatalog returned error: %v", err)
	}

	err = catalog.RegisterComponent(testComponentDescriptor("custom.gate", "custom_gate"))
	if !errors.Is(err, ErrUnknownComponentKind) {
		t.Fatalf("error = %v, want ErrUnknownComponentKind", err)
	}

	var unknown UnknownComponentKindError
	if !errors.As(err, &unknown) {
		t.Fatal("error should expose UnknownComponentKindError")
	}
	if unknown.Kind != "custom_gate" {
		t.Fatalf("unknown kind = %q, want custom_gate", unknown.Kind)
	}
}

func TestCatalogConcurrentAccess(t *testing.T) {
	catalog := testCatalog()

	var wg sync.WaitGroup
	errCh := make(chan error, 32*3)
	for i := 0; i < 32; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()

			suffix := string(rune('a'+i/26)) + string(rune('a'+i%26))
			reason := Reason("custom_reason_" + suffix)
			kind := ComponentKind("custom_kind_" + suffix)
			componentID := ComponentID("custom.component_" + suffix)

			if err := catalog.RegisterReason(testReasonDescriptor(reason)); err != nil {
				errCh <- fmt.Errorf("register reason %q: %w", reason, err)
			}
			if err := catalog.RegisterKind(testKindDescriptor(kind)); err != nil {
				errCh <- fmt.Errorf("register kind %q: %w", kind, err)
			}
			if err := catalog.RegisterComponent(testComponentDescriptor(componentID, kind)); err != nil {
				errCh <- fmt.Errorf("register component %q: %w", componentID, err)
			}

			_, _ = catalog.Reason(reason)
			_, _ = catalog.Kind(kind)
			_, _ = catalog.Component(componentID)
			_ = catalog.Reasons()
			_ = catalog.Kinds()
			_ = catalog.Components()
			_ = catalog.LenReasons()
			_ = catalog.LenKinds()
			_ = catalog.LenComponents()
		}(i)
	}
	wg.Wait()
	close(errCh)

	for err := range errCh {
		if err != nil {
			t.Fatalf("unexpected concurrent catalog error: %v", err)
		}
	}
}

func TestCatalogNilReceiverPanics(t *testing.T) {
	t.Parallel()

	var catalog *Catalog
	assertPanicString(t, nilCatalogPanic, func() {
		_, _ = catalog.Reason(ReasonDenied)
	})
	assertPanicString(t, nilCatalogPanic, func() {
		_, _ = catalog.Kind(KindBulkhead)
	})
	assertPanicString(t, nilCatalogPanic, func() {
		_, _ = catalog.Component("resilience.bulkhead")
	})
	assertPanicString(t, nilCatalogPanic, func() {
		_ = catalog.Reasons()
	})
	assertPanicString(t, nilCatalogPanic, func() {
		_ = catalog.Kinds()
	})
	assertPanicString(t, nilCatalogPanic, func() {
		_ = catalog.Components()
	})
	assertPanicString(t, nilCatalogPanic, func() {
		_ = catalog.RegisterReason(testReasonDescriptor("custom_reason"))
	})
	assertPanicString(t, nilCatalogPanic, func() {
		_ = catalog.RegisterKind(testKindDescriptor("custom_kind"))
	})
	assertPanicString(t, nilCatalogPanic, func() {
		_ = catalog.RegisterComponent(testComponentDescriptor("custom.component", KindBulkhead))
	})
	assertPanicString(t, nilCatalogPanic, func() {
		_ = catalog.LenReasons()
	})
	assertPanicString(t, nilCatalogPanic, func() {
		_ = catalog.LenKinds()
	})
	assertPanicString(t, nilCatalogPanic, func() {
		_ = catalog.LenComponents()
	})
}

func testCatalog() *Catalog {
	kinds := NewBuiltinKindRegistry()
	catalog, err := NewCatalog(
		NewBuiltinReasonRegistry(),
		kinds,
		NewBuiltinComponentRegistry(kinds),
	)
	if err != nil {
		panic(err)
	}
	return catalog
}
