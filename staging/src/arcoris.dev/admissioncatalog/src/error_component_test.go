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

func TestComponentErrorsSupportIsAndAs(t *testing.T) {
	invalid := InvalidComponentDescriptorError{Descriptor: componentDescriptor(testComponent, testKind)}
	if !errors.Is(invalid, ErrInvalidComponentDescriptor) {
		t.Fatal("invalid component error does not match sentinel")
	}
	var invalidTyped InvalidComponentDescriptorError
	if !errors.As(invalid, &invalidTyped) {
		t.Fatal("invalid component error does not expose typed value")
	}

	duplicate := DuplicateComponentDeclarationError{ID: testComponent}
	if !errors.Is(duplicate, ErrDuplicateComponentDeclaration) {
		t.Fatal("duplicate component error does not match sentinel")
	}
	var duplicateTyped DuplicateComponentDeclarationError
	if !errors.As(duplicate, &duplicateTyped) {
		t.Fatal("duplicate component error does not expose typed value")
	}

	unknown := UnknownComponentKindError{ComponentID: testComponent, Kind: testKind}
	if !errors.Is(unknown, ErrUnknownComponentKind) {
		t.Fatal("unknown component kind error does not match sentinel")
	}
	var unknownTyped UnknownComponentKindError
	if !errors.As(unknown, &unknownTyped) {
		t.Fatal("unknown component kind error does not expose typed value")
	}
}

func TestComponentErrorsExposeDetails(t *testing.T) {
	invalid := InvalidComponentDescriptorError{Descriptor: componentDescriptor(testComponent, testKind)}
	if invalid.Descriptor.ID != testComponent {
		t.Fatalf("Descriptor.ID = %s, want %s", invalid.Descriptor.ID, testComponent)
	}

	duplicate := DuplicateComponentDeclarationError{ID: testComponent}
	if duplicate.ID != testComponent {
		t.Fatalf("ID = %s, want %s", duplicate.ID, testComponent)
	}

	unknown := UnknownComponentKindError{ComponentID: testComponent, Kind: testKind}
	if unknown.ComponentID != testComponent || unknown.Kind != testKind {
		t.Fatalf("UnknownComponentKindError = %+v", unknown)
	}
}
