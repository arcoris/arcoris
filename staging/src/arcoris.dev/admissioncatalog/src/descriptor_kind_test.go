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

func TestComponentKindDescriptorIsValid(t *testing.T) {
	descriptor := kindDescriptor(testKind)
	if !descriptor.IsValid() {
		t.Fatal("descriptor is invalid")
	}
}

func TestComponentKindDescriptorRejectsInvalidFields(t *testing.T) {
	tests := []struct {
		name       string
		descriptor ComponentKindDescriptor
	}{
		{name: "kind", descriptor: ComponentKindDescriptor{Kind: admission.ComponentKind("bad-kind")}},
		{name: "summary", descriptor: ComponentKindDescriptor{Kind: testKind, Summary: "bad\nsummary"}},
		{
			name: "capabilities",
			descriptor: ComponentKindDescriptor{
				Kind:                 testKind,
				DeclaredCapabilities: NewCapabilitySet(NewOutcomeSet(), EffectSet(1<<7)),
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

func TestComponentKindDescriptorSummary(t *testing.T) {
	empty := kindDescriptor(testKind)
	empty.Summary = ""
	if !empty.IsValid() {
		t.Fatal("empty summary should be valid")
	}

	nonEmpty := kindDescriptor(testKind)
	if got, want := nonEmpty.Summary, testKind.String()+" summary"; got != want {
		t.Fatalf("Summary = %q, want %q", got, want)
	}
}
