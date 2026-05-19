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
	"sync"
	"testing"
)

func TestNewComponentRegistry(t *testing.T) {
	t.Parallel()

	kinds := NewBuiltinKindRegistry()
	descriptor := testComponentDescriptor("resilience.bulkhead", KindBulkhead)

	registry, err := NewComponentRegistry(kinds, descriptor)
	if err != nil {
		t.Fatalf("NewComponentRegistry returned error: %v", err)
	}
	if registry.Len() != 1 {
		t.Fatalf("Len = %d, want 1", registry.Len())
	}
}

func TestNewComponentRegistryRejectsNilKindRegistry(t *testing.T) {
	t.Parallel()

	registry, err := NewComponentRegistry(nil)
	if registry != nil {
		t.Fatal("registry should be nil on nil kind registry")
	}
	if !errors.Is(err, ErrNilKindRegistry) {
		t.Fatalf("error = %v, want ErrNilKindRegistry", err)
	}
}

func TestMustComponentRegistry(t *testing.T) {
	t.Parallel()

	registry := MustComponentRegistry(
		NewBuiltinKindRegistry(),
		testComponentDescriptor("resilience.bulkhead", KindBulkhead),
	)
	if !registry.Contains("resilience.bulkhead") {
		t.Fatal("registry should contain resilience.bulkhead")
	}
}

func TestComponentRegistryRegister(t *testing.T) {
	t.Parallel()

	registry := MustComponentRegistry(NewBuiltinKindRegistry())
	descriptor := testComponentDescriptor("resilience.bulkhead", KindBulkhead)

	if err := registry.Register(descriptor); err != nil {
		t.Fatalf("Register returned error: %v", err)
	}
	if got, ok := registry.Lookup("resilience.bulkhead"); !ok || got != descriptor {
		t.Fatalf("Lookup = (%+v, %v), want registered descriptor", got, ok)
	}
}

func TestComponentRegistryRegisterRejectsInvalidDescriptor(t *testing.T) {
	t.Parallel()

	registry := MustComponentRegistry(NewBuiltinKindRegistry())
	err := registry.Register(ComponentDescriptor{ID: "bad/id", Kind: KindBulkhead})
	if !errors.Is(err, ErrInvalidComponentDescriptor) {
		t.Fatalf("error = %v, want ErrInvalidComponentDescriptor", err)
	}
}

func TestComponentRegistryRegisterRejectsUnknownKind(t *testing.T) {
	t.Parallel()

	kinds := MustKindRegistry(testKindDescriptor(KindBulkhead))
	registry := MustComponentRegistry(kinds)
	err := registry.Register(testComponentDescriptor("resilience.retrybudget", KindRetryBudget))
	if !errors.Is(err, ErrUnknownComponentKind) {
		t.Fatalf("error = %v, want ErrUnknownComponentKind", err)
	}

	var unknown UnknownComponentKindError
	if !errors.As(err, &unknown) {
		t.Fatal("error should expose UnknownComponentKindError")
	}
	if unknown.Kind != KindRetryBudget {
		t.Fatalf("unknown kind = %q, want %q", unknown.Kind, KindRetryBudget)
	}
}

func TestComponentRegistryRegisterRejectsDuplicateID(t *testing.T) {
	t.Parallel()

	descriptor := testComponentDescriptor("resilience.bulkhead", KindBulkhead)
	registry := MustComponentRegistry(NewBuiltinKindRegistry(), descriptor)

	err := registry.Register(descriptor)
	if !errors.Is(err, ErrComponentAlreadyRegistered) {
		t.Fatalf("error = %v, want ErrComponentAlreadyRegistered", err)
	}

	var duplicate DuplicateComponentError
	if !errors.As(err, &duplicate) {
		t.Fatal("error should expose DuplicateComponentError")
	}
	if duplicate.ID != "resilience.bulkhead" {
		t.Fatalf("duplicate id = %q, want resilience.bulkhead", duplicate.ID)
	}
}

func TestComponentRegistryLookupAndContains(t *testing.T) {
	t.Parallel()

	descriptor := testComponentDescriptor("resilience.bulkhead", KindBulkhead)
	registry := MustComponentRegistry(NewBuiltinKindRegistry(), descriptor)

	if got, ok := registry.Lookup("resilience.bulkhead"); !ok || got != descriptor {
		t.Fatalf("Lookup existing = (%+v, %v), want descriptor,true", got, ok)
	}
	if got, ok := registry.Lookup("resilience.missing"); ok || got != (ComponentDescriptor{}) {
		t.Fatalf("Lookup missing = (%+v, %v), want zero,false", got, ok)
	}
	if got, ok := registry.Lookup("bad/id"); ok || got != (ComponentDescriptor{}) {
		t.Fatalf("Lookup invalid = (%+v, %v), want zero,false", got, ok)
	}
	if !registry.Contains("resilience.bulkhead") {
		t.Fatal("Contains should report registered component")
	}
	if registry.Contains("bad/id") {
		t.Fatal("Contains should reject invalid ID")
	}
}

func TestComponentRegistryListIsSortedCopy(t *testing.T) {
	t.Parallel()

	registry := MustComponentRegistry(
		NewBuiltinKindRegistry(),
		testComponentDescriptor("resilience.retrybudget", KindRetryBudget),
		testComponentDescriptor("resilience.bulkhead", KindBulkhead),
		testComponentDescriptor("resilience.deadline", KindDeadline),
	)

	list := registry.List()
	gotOrder := []ComponentID{list[0].ID, list[1].ID, list[2].ID}
	wantOrder := []ComponentID{
		"resilience.bulkhead",
		"resilience.deadline",
		"resilience.retrybudget",
	}
	for i := range wantOrder {
		if gotOrder[i] != wantOrder[i] {
			t.Fatalf("order[%d] = %q, want %q", i, gotOrder[i], wantOrder[i])
		}
	}

	list[0].ID = "resilience.mutated"
	if registry.Contains("resilience.mutated") {
		t.Fatal("mutating List result should not mutate registry")
	}
}

func TestComponentRegistryNilReceiverPanics(t *testing.T) {
	t.Parallel()

	var registry *ComponentRegistry
	assertPanicString(t, nilComponentRegistryPanic, func() {
		_ = registry.Register(testComponentDescriptor("resilience.bulkhead", KindBulkhead))
	})
	assertPanicString(t, nilComponentRegistryPanic, func() {
		_, _ = registry.Lookup("resilience.bulkhead")
	})
	assertPanicString(t, nilComponentRegistryPanic, func() {
		_ = registry.Contains("resilience.bulkhead")
	})
	assertPanicString(t, nilComponentRegistryPanic, func() {
		_ = registry.List()
	})
	assertPanicString(t, nilComponentRegistryPanic, func() {
		_ = registry.Len()
	})
}

func TestComponentRegistryConcurrentAccess(t *testing.T) {
	registry := MustComponentRegistry(NewBuiltinKindRegistry())

	var wg sync.WaitGroup
	for i := 0; i < 32; i++ {
		i := i
		wg.Add(1)
		go func() {
			defer wg.Done()

			id := ComponentID("resilience.bulkhead_" +
				string(rune('a'+i/26)) +
				string(rune('a'+i%26)))
			_ = registry.Register(testComponentDescriptor(id, KindBulkhead))
			_, _ = registry.Lookup(id)
			_ = registry.Contains(id)
			_ = registry.List()
			_ = registry.Len()
		}()
	}
	wg.Wait()
}

// testComponentDescriptor builds a syntactically valid component descriptor for
// registry tests.
//
// The helper intentionally does not register the kind. Individual tests decide
// whether the descriptor should pass catalog-level kind validation or fail with
// UnknownComponentKindError.
func testComponentDescriptor(
	id ComponentID,
	kind ComponentKind,
) ComponentDescriptor {
	return ComponentDescriptor{
		ID:   id,
		Kind: kind,
		Capabilities: NewCapabilitySet(
			CapabilityAdmit,
			CapabilityDeny,
			CapabilityEffectNone,
		),
	}
}
