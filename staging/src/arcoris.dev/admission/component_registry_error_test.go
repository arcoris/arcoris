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
	"testing"
)

func TestComponentRegistryErrorsSupportIsAndAs(t *testing.T) {
	t.Parallel()

	if !errors.Is(ErrNilKindRegistry, ErrNilKindRegistry) {
		t.Fatal("nil kind registry sentinel should match itself")
	}

	invalid := InvalidComponentDescriptorError{
		Descriptor: ComponentDescriptor{ID: "bad/id"},
	}
	if !errors.Is(invalid, ErrInvalidComponentDescriptor) {
		t.Fatal("invalid descriptor error should match sentinel")
	}

	unknown := UnknownComponentKindError{Kind: "custom_kind"}
	if !errors.Is(unknown, ErrUnknownComponentKind) {
		t.Fatal("unknown kind error should match sentinel")
	}
	var unknownKind UnknownComponentKindError
	if !errors.As(unknown, &unknownKind) {
		t.Fatal("unknown kind error should support errors.As")
	}
	if unknownKind.Kind != "custom_kind" {
		t.Fatalf("unknown kind = %q, want custom_kind", unknownKind.Kind)
	}

	duplicate := DuplicateComponentError{ID: "resilience.bulkhead"}
	if !errors.Is(duplicate, ErrComponentAlreadyRegistered) {
		t.Fatal("duplicate component error should match sentinel")
	}
	var duplicateComponent DuplicateComponentError
	if !errors.As(duplicate, &duplicateComponent) {
		t.Fatal("duplicate component error should support errors.As")
	}
	if duplicateComponent.ID != "resilience.bulkhead" {
		t.Fatalf("duplicate id = %q, want resilience.bulkhead", duplicateComponent.ID)
	}
}
