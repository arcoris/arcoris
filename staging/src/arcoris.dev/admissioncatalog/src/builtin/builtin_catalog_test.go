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

package builtin

import (
	"testing"

	"arcoris.dev/admission"
	"arcoris.dev/admissioncatalog"
)

func TestNewCatalog(t *testing.T) {
	t.Parallel()

	catalog := NewCatalog()

	for _, descriptor := range ReasonDescriptors() {
		if got, ok := catalog.Reason(descriptor.Reason); !ok || got != descriptor {
			t.Fatalf("admission.Reason(%q) = (%+v, %v), want standard descriptor", descriptor.Reason, got, ok)
		}
	}
	for _, descriptor := range KindDescriptors() {
		if got, ok := catalog.Kind(descriptor.Kind); !ok || got != descriptor {
			t.Fatalf("Kind(%q) = (%+v, %v), want standard descriptor", descriptor.Kind, got, ok)
		}
	}
	for _, descriptor := range ComponentDescriptors() {
		if got, ok := catalog.Component(descriptor.ID); !ok || got != descriptor {
			t.Fatalf("Component(%q) = (%+v, %v), want standard descriptor", descriptor.ID, got, ok)
		}
	}

	if catalog.LenReasons() != len(ReasonDescriptors()) {
		t.Fatalf("LenReasons = %d, want %d", catalog.LenReasons(), len(ReasonDescriptors()))
	}
	if catalog.LenKinds() != len(KindDescriptors()) {
		t.Fatalf("LenKinds = %d, want %d", catalog.LenKinds(), len(KindDescriptors()))
	}
	if catalog.LenComponents() != len(ComponentDescriptors()) {
		t.Fatalf("LenComponents = %d, want %d", catalog.LenComponents(), len(ComponentDescriptors()))
	}
}

func TestNewCatalogRegisterKindThenRegisterComponent(t *testing.T) {
	t.Parallel()

	catalog := NewCatalog()
	kind := admission.ComponentKind("custom_gate")
	componentID := admission.ComponentID("custom.gate")

	if err := catalog.RegisterKind(admissioncatalog.ComponentKindDescriptor{
		Kind: kind,
		Capabilities: admissioncatalog.NewCapabilitySet(
			admissioncatalog.CapabilityAdmit,
			admissioncatalog.CapabilityDeny,
			admissioncatalog.CapabilityEffectNone,
		),
	}); err != nil {
		t.Fatalf("RegisterKind returned error: %v", err)
	}

	component := admissioncatalog.ComponentDescriptor{
		ID:   componentID,
		Kind: kind,
		Capabilities: admissioncatalog.NewCapabilitySet(
			admissioncatalog.CapabilityAdmit,
			admissioncatalog.CapabilityDeny,
			admissioncatalog.CapabilityEffectNone,
		),
	}
	if err := catalog.RegisterComponent(component); err != nil {
		t.Fatalf("RegisterComponent returned error: %v", err)
	}
	if got, ok := catalog.Component(componentID); !ok || got != component {
		t.Fatalf("Component lookup = (%+v, %v), want registered descriptor", got, ok)
	}
}
