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

func TestNewKindRegistry(t *testing.T) {
	t.Parallel()

	registry, err := NewKindRegistry(
		testKindDescriptor("z_kind"),
		testKindDescriptor("a_kind"),
	)
	if err != nil {
		t.Fatalf("NewKindRegistry returned error: %v", err)
	}
	if registry.Len() != 2 {
		t.Fatalf("Len = %d, want 2", registry.Len())
	}
}

func TestMustKindRegistry(t *testing.T) {
	t.Parallel()

	registry := MustKindRegistry(testKindDescriptor("custom_kind"))
	if !registry.Contains("custom_kind") {
		t.Fatal("registry should contain custom_kind")
	}
}

func TestKindRegistryRegister(t *testing.T) {
	t.Parallel()

	registry := MustKindRegistry()
	descriptor := testKindDescriptor("custom_kind")

	if err := registry.Register(descriptor); err != nil {
		t.Fatalf("Register returned error: %v", err)
	}
	if got, ok := registry.Lookup("custom_kind"); !ok || got != descriptor {
		t.Fatalf("Lookup = (%+v, %v), want registered descriptor", got, ok)
	}
}

func TestKindRegistryRegisterRejectsInvalidDescriptor(t *testing.T) {
	t.Parallel()

	registry := MustKindRegistry()
	err := registry.Register(ComponentKindDescriptor{Kind: "bad-kind"})
	if !errors.Is(err, ErrInvalidComponentKindDescriptor) {
		t.Fatalf("error = %v, want ErrInvalidComponentKindDescriptor", err)
	}
}

func TestKindRegistryRegisterRejectsDuplicateKind(t *testing.T) {
	t.Parallel()

	registry := MustKindRegistry(testKindDescriptor("custom_kind"))
	err := registry.Register(testKindDescriptor("custom_kind"))
	if !errors.Is(err, ErrComponentKindAlreadyRegistered) {
		t.Fatalf("error = %v, want ErrComponentKindAlreadyRegistered", err)
	}

	var duplicate DuplicateComponentKindError
	if !errors.As(err, &duplicate) {
		t.Fatal("error should expose DuplicateComponentKindError")
	}
	if duplicate.Kind != "custom_kind" {
		t.Fatalf("duplicate kind = %q, want custom_kind", duplicate.Kind)
	}
}

func TestKindRegistryLookupAndContains(t *testing.T) {
	t.Parallel()

	descriptor := testKindDescriptor("custom_kind")
	registry := MustKindRegistry(descriptor)

	if got, ok := registry.Lookup("custom_kind"); !ok || got != descriptor {
		t.Fatalf("Lookup existing = (%+v, %v), want descriptor,true", got, ok)
	}
	if got, ok := registry.Lookup("missing_kind"); ok || got != (ComponentKindDescriptor{}) {
		t.Fatalf("Lookup missing = (%+v, %v), want zero,false", got, ok)
	}
	if got, ok := registry.Lookup("bad-kind"); ok || got != (ComponentKindDescriptor{}) {
		t.Fatalf("Lookup invalid = (%+v, %v), want zero,false", got, ok)
	}
	if !registry.Contains("custom_kind") {
		t.Fatal("Contains should report registered kind")
	}
	if registry.Contains("bad-kind") {
		t.Fatal("Contains should reject invalid kind")
	}
}

func TestKindRegistryListIsSortedCopy(t *testing.T) {
	t.Parallel()

	registry := MustKindRegistry(
		testKindDescriptor("z_kind"),
		testKindDescriptor("a_kind"),
		testKindDescriptor("m_kind"),
	)

	list := registry.List()
	gotOrder := []ComponentKind{list[0].Kind, list[1].Kind, list[2].Kind}
	wantOrder := []ComponentKind{"a_kind", "m_kind", "z_kind"}
	for i := range wantOrder {
		if gotOrder[i] != wantOrder[i] {
			t.Fatalf("order[%d] = %q, want %q", i, gotOrder[i], wantOrder[i])
		}
	}

	list[0].Kind = "mutated_kind"
	if registry.Contains("mutated_kind") {
		t.Fatal("mutating List result should not mutate registry")
	}
}

func TestKindRegistryNilReceiverPanics(t *testing.T) {
	t.Parallel()

	var registry *KindRegistry
	assertPanicString(t, nilKindRegistryPanic, func() {
		_ = registry.Register(testKindDescriptor("custom_kind"))
	})
	assertPanicString(t, nilKindRegistryPanic, func() {
		_, _ = registry.Lookup("custom_kind")
	})
	assertPanicString(t, nilKindRegistryPanic, func() {
		_ = registry.Contains("custom_kind")
	})
	assertPanicString(t, nilKindRegistryPanic, func() {
		_ = registry.List()
	})
	assertPanicString(t, nilKindRegistryPanic, func() {
		_ = registry.Len()
	})
}

func TestKindRegistryConcurrentAccess(t *testing.T) {
	registry := MustKindRegistry()

	var wg sync.WaitGroup
	for i := 0; i < 32; i++ {
		i := i
		wg.Add(1)
		go func() {
			defer wg.Done()

			kind := ComponentKind("custom_kind_" +
				string(rune('a'+i/26)) +
				string(rune('a'+i%26)))
			_ = registry.Register(testKindDescriptor(kind))
			_, _ = registry.Lookup(kind)
			_ = registry.Contains(kind)
			_ = registry.List()
			_ = registry.Len()
		}()
	}
	wg.Wait()
}

func testKindDescriptor(kind ComponentKind) ComponentKindDescriptor {
	return ComponentKindDescriptor{
		Kind: kind,
		Capabilities: NewCapabilitySet(
			CapabilityAdmit,
			CapabilityDeny,
			CapabilityEffectNone,
		),
	}
}
