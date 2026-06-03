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
)

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
