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

func TestNewBuiltinCatalog(t *testing.T) {
	t.Parallel()

	catalog := NewBuiltinCatalog()

	for _, descriptor := range BuiltinReasonDescriptors() {
		if got, ok := catalog.Reason(descriptor.Reason); !ok || got != descriptor {
			t.Fatalf("Reason(%q) = (%+v, %v), want built-in descriptor", descriptor.Reason, got, ok)
		}
	}
	for _, descriptor := range BuiltinKindDescriptors() {
		if got, ok := catalog.Kind(descriptor.Kind); !ok || got != descriptor {
			t.Fatalf("Kind(%q) = (%+v, %v), want built-in descriptor", descriptor.Kind, got, ok)
		}
	}
	for _, descriptor := range BuiltinComponentDescriptors() {
		if got, ok := catalog.Component(descriptor.ID); !ok || got != descriptor {
			t.Fatalf("Component(%q) = (%+v, %v), want built-in descriptor", descriptor.ID, got, ok)
		}
	}

	if catalog.LenReasons() != len(BuiltinReasonDescriptors()) {
		t.Fatalf("LenReasons = %d, want %d", catalog.LenReasons(), len(BuiltinReasonDescriptors()))
	}
	if catalog.LenKinds() != len(BuiltinKindDescriptors()) {
		t.Fatalf("LenKinds = %d, want %d", catalog.LenKinds(), len(BuiltinKindDescriptors()))
	}
	if catalog.LenComponents() != len(BuiltinComponentDescriptors()) {
		t.Fatalf("LenComponents = %d, want %d", catalog.LenComponents(), len(BuiltinComponentDescriptors()))
	}
}

func TestNewBuiltinCatalogRegisterKindThenRegisterComponent(t *testing.T) {
	t.Parallel()

	catalog := NewBuiltinCatalog()
	kind := ComponentKind("custom_gate")
	componentID := ComponentID("custom.gate")

	if err := catalog.RegisterKind(ComponentKindDescriptor{
		Kind: kind,
		Capabilities: NewCapabilitySet(
			CapabilityAdmit,
			CapabilityDeny,
			CapabilityEffectNone,
		),
	}); err != nil {
		t.Fatalf("RegisterKind returned error: %v", err)
	}

	component := ComponentDescriptor{
		ID:   componentID,
		Kind: kind,
		Capabilities: NewCapabilitySet(
			CapabilityAdmit,
			CapabilityDeny,
			CapabilityEffectNone,
		),
	}
	if err := catalog.RegisterComponent(component); err != nil {
		t.Fatalf("RegisterComponent returned error: %v", err)
	}
	if got, ok := catalog.Component(componentID); !ok || got != component {
		t.Fatalf("Component lookup = (%+v, %v), want registered descriptor", got, ok)
	}
}
