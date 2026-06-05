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
	"testing"

	"arcoris.dev/admission"
)

func TestComponentDescriptorIsValid(t *testing.T) {
	descriptor := componentDescriptor(testComponent, testKind)
	if !descriptor.IsValid() {
		t.Fatal("descriptor is invalid")
	}
}

func TestComponentDescriptorRejectsInvalidFields(t *testing.T) {
	tests := []struct {
		name       string
		descriptor ComponentDescriptor
	}{
		{name: "id", descriptor: ComponentDescriptor{ID: admission.ComponentID("bad id"), Kind: testKind}},
		{name: "kind", descriptor: ComponentDescriptor{ID: testComponent, Kind: admission.ComponentKind("bad-kind")}},
		{name: "summary", descriptor: ComponentDescriptor{ID: testComponent, Kind: testKind, Summary: "bad\nsummary"}},
		{
			name: "capabilities",
			descriptor: ComponentDescriptor{
				ID:                   testComponent,
				Kind:                 testKind,
				DeclaredCapabilities: NewCapabilitySet(OutcomeSet(1<<7), NewEffectSet()),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.descriptor.IsValid() {
				t.Fatal("descriptor is valid")
			}
		})
	}
}

func TestComponentDescriptorDoesNotRequireKindMembership(t *testing.T) {
	descriptor := componentDescriptor(testComponent, testKind)
	if !descriptor.IsValid() {
		t.Fatal("descriptor validity should not require catalog membership")
	}
}

func TestComponentDescriptorSummary(t *testing.T) {
	empty := componentDescriptor(testComponent, testKind)
	empty.Summary = ""
	if !empty.IsValid() {
		t.Fatal("empty summary should be valid")
	}

	nonEmpty := componentDescriptor(testComponent, testKind)
	if got, want := nonEmpty.Summary, testComponent.String()+" summary"; got != want {
		t.Fatalf("Summary = %q, want %q", got, want)
	}
}
