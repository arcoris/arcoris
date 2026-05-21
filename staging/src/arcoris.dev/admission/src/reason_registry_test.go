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
	panicassert "arcoris.dev/testutil/panic"
	"errors"
	"sync"
	"testing"
)

func TestNewReasonRegistry(t *testing.T) {
	t.Parallel()

	registry, err := NewReasonRegistry(
		testReasonDescriptor("z_reason"),
		testReasonDescriptor("a_reason"),
	)
	if err != nil {
		t.Fatalf("NewReasonRegistry returned error: %v", err)
	}
	if registry.Len() != 2 {
		t.Fatalf("Len = %d, want 2", registry.Len())
	}
}

func TestNewReasonRegistryRejectsInvalidDescriptor(t *testing.T) {
	t.Parallel()

	registry, err := NewReasonRegistry(ReasonDescriptor{Reason: "bad-reason"})
	if registry != nil {
		t.Fatal("registry should be nil on invalid descriptor")
	}
	if !errors.Is(err, ErrInvalidReasonDescriptor) {
		t.Fatalf("error = %v, want ErrInvalidReasonDescriptor", err)
	}
}

func TestNewReasonRegistryRejectsDuplicateReason(t *testing.T) {
	t.Parallel()

	registry, err := NewReasonRegistry(
		testReasonDescriptor("custom_reason"),
		testReasonDescriptor("custom_reason"),
	)
	if registry != nil {
		t.Fatal("registry should be nil on duplicate reason")
	}
	if !errors.Is(err, ErrReasonAlreadyRegistered) {
		t.Fatalf("error = %v, want ErrReasonAlreadyRegistered", err)
	}
}

func TestMustReasonRegistry(t *testing.T) {
	t.Parallel()

	registry := MustReasonRegistry(testReasonDescriptor("custom_reason"))
	if !registry.Contains("custom_reason") {
		t.Fatal("registry should contain custom_reason")
	}
}

func TestMustReasonRegistryPanicsOnInvalidDescriptor(t *testing.T) {
	t.Parallel()

	defer func() {
		if got := recover(); got == nil {
			t.Fatal("expected panic from invalid descriptor")
		} else if err, ok := got.(error); !ok || !errors.Is(err, ErrInvalidReasonDescriptor) {
			t.Fatalf("panic = %v, want ErrInvalidReasonDescriptor", got)
		}
	}()

	_ = MustReasonRegistry(ReasonDescriptor{Reason: "bad-reason"})
}

func TestReasonRegistryRegister(t *testing.T) {
	t.Parallel()

	registry := MustReasonRegistry()
	descriptor := testReasonDescriptor("custom_reason")

	if err := registry.Register(descriptor); err != nil {
		t.Fatalf("Register returned error: %v", err)
	}
	if got, ok := registry.Lookup("custom_reason"); !ok || got != descriptor {
		t.Fatalf("Lookup = (%+v, %v), want registered descriptor", got, ok)
	}
}

func TestReasonRegistryRegisterRejectsInvalidDescriptor(t *testing.T) {
	t.Parallel()

	registry := MustReasonRegistry()
	err := registry.Register(ReasonDescriptor{Reason: "bad-reason"})
	if !errors.Is(err, ErrInvalidReasonDescriptor) {
		t.Fatalf("error = %v, want ErrInvalidReasonDescriptor", err)
	}

	var invalid InvalidReasonDescriptorError
	if !errors.As(err, &invalid) {
		t.Fatal("error should expose InvalidReasonDescriptorError")
	}
	if invalid.Descriptor.Reason != "bad-reason" {
		t.Fatalf("invalid reason = %q, want bad-reason", invalid.Descriptor.Reason)
	}
}

func TestReasonRegistryRegisterRejectsDuplicateReason(t *testing.T) {
	t.Parallel()

	registry := MustReasonRegistry(testReasonDescriptor("custom_reason"))
	err := registry.Register(testReasonDescriptor("custom_reason"))
	if !errors.Is(err, ErrReasonAlreadyRegistered) {
		t.Fatalf("error = %v, want ErrReasonAlreadyRegistered", err)
	}

	var duplicate DuplicateReasonError
	if !errors.As(err, &duplicate) {
		t.Fatal("error should expose DuplicateReasonError")
	}
	if duplicate.Reason != "custom_reason" {
		t.Fatalf("duplicate reason = %q, want custom_reason", duplicate.Reason)
	}
}

func TestReasonRegistryLookupAndContains(t *testing.T) {
	t.Parallel()

	descriptor := testReasonDescriptor("custom_reason")
	registry := MustReasonRegistry(descriptor)

	if got, ok := registry.Lookup("custom_reason"); !ok || got != descriptor {
		t.Fatalf("Lookup existing = (%+v, %v), want descriptor,true", got, ok)
	}
	if got, ok := registry.Lookup("missing_reason"); ok || got != (ReasonDescriptor{}) {
		t.Fatalf("Lookup missing = (%+v, %v), want zero,false", got, ok)
	}
	if got, ok := registry.Lookup("bad-reason"); ok || got != (ReasonDescriptor{}) {
		t.Fatalf("Lookup invalid = (%+v, %v), want zero,false", got, ok)
	}
	if !registry.Contains("custom_reason") {
		t.Fatal("Contains should report registered reason")
	}
	if registry.Contains("bad-reason") {
		t.Fatal("Contains should reject invalid reason")
	}
}

func TestReasonRegistryZeroValue(t *testing.T) {
	t.Parallel()

	var registry ReasonRegistry
	if got := registry.Len(); got != 0 {
		t.Fatalf("Len = %d, want 0", got)
	}
	if got, ok := registry.Lookup(ReasonDenied); ok || got != (ReasonDescriptor{}) {
		t.Fatalf("Lookup = (%+v, %v), want zero,false", got, ok)
	}
	if list := registry.List(); len(list) != 0 {
		t.Fatalf("List length = %d, want 0", len(list))
	}

	descriptor := testReasonDescriptor("custom_reason")
	if err := registry.Register(descriptor); err != nil {
		t.Fatalf("Register returned error: %v", err)
	}
	if got, ok := registry.Lookup("custom_reason"); !ok || got != descriptor {
		t.Fatalf("Lookup registered = (%+v, %v), want descriptor,true", got, ok)
	}
}

func TestReasonRegistryListIsSortedCopy(t *testing.T) {
	t.Parallel()

	registry := MustReasonRegistry(
		testReasonDescriptor("z_reason"),
		testReasonDescriptor("a_reason"),
		testReasonDescriptor("m_reason"),
	)

	list := registry.List()
	gotOrder := []Reason{list[0].Reason, list[1].Reason, list[2].Reason}
	wantOrder := []Reason{"a_reason", "m_reason", "z_reason"}
	for i := range wantOrder {
		if gotOrder[i] != wantOrder[i] {
			t.Fatalf("order[%d] = %q, want %q", i, gotOrder[i], wantOrder[i])
		}
	}

	list[0].Reason = "mutated_reason"
	if registry.Contains("mutated_reason") {
		t.Fatal("mutating List result should not mutate registry")
	}
}

func TestReasonRegistryNilReceiverPanics(t *testing.T) {
	t.Parallel()

	var registry *ReasonRegistry
	panicassert.RequireMessage(t, nilReasonRegistryPanic, func() {
		_ = registry.Register(testReasonDescriptor("custom_reason"))
	})
	panicassert.RequireMessage(t, nilReasonRegistryPanic, func() {
		_, _ = registry.Lookup("custom_reason")
	})
	panicassert.RequireMessage(t, nilReasonRegistryPanic, func() {
		_ = registry.Contains("custom_reason")
	})
	panicassert.RequireMessage(t, nilReasonRegistryPanic, func() {
		_ = registry.List()
	})
	panicassert.RequireMessage(t, nilReasonRegistryPanic, func() {
		_ = registry.Len()
	})
}

func TestReasonRegistryConcurrentAccess(t *testing.T) {
	registry := MustReasonRegistry()

	var wg sync.WaitGroup
	for i := 0; i < 32; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()

			reason := Reason("custom_reason_" +
				string(rune('a'+i/26)) +
				string(rune('a'+i%26)))
			_ = registry.Register(testReasonDescriptor(reason))
			_, _ = registry.Lookup(reason)
			_ = registry.Contains(reason)
			_ = registry.List()
			_ = registry.Len()
		}(i)
	}
	wg.Wait()
}

// testReasonDescriptor builds a valid reason descriptor with a small common
// capability surface.
//
// Registry tests use this helper so the individual test cases can focus on
// registry behavior rather than repeating descriptor boilerplate.
func testReasonDescriptor(reason Reason) ReasonDescriptor {
	return ReasonDescriptor{
		Reason: reason,
		Capabilities: NewCapabilitySet(
			CapabilityAdmit,
			CapabilityDeny,
			CapabilityEffectNone,
		),
	}
}
