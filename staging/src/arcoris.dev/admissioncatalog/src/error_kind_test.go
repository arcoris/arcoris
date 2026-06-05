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

func TestKindErrorsSupportIsAndAs(t *testing.T) {
	invalid := InvalidComponentKindDescriptorError{Descriptor: kindDescriptor(testKind)}
	if !errors.Is(invalid, ErrInvalidComponentKindDescriptor) {
		t.Fatal("invalid kind error does not match sentinel")
	}
	var invalidTyped InvalidComponentKindDescriptorError
	if !errors.As(invalid, &invalidTyped) {
		t.Fatal("invalid kind error does not expose typed value")
	}

	duplicate := DuplicateComponentKindDeclarationError{Kind: testKind}
	if !errors.Is(duplicate, ErrDuplicateComponentKindDeclaration) {
		t.Fatal("duplicate kind error does not match sentinel")
	}
	var duplicateTyped DuplicateComponentKindDeclarationError
	if !errors.As(duplicate, &duplicateTyped) {
		t.Fatal("duplicate kind error does not expose typed value")
	}
}

func TestKindErrorsExposeDetails(t *testing.T) {
	invalid := InvalidComponentKindDescriptorError{Descriptor: kindDescriptor(testKind)}
	if invalid.Descriptor.Kind != testKind {
		t.Fatalf("Descriptor.Kind = %s, want %s", invalid.Descriptor.Kind, testKind)
	}

	duplicate := DuplicateComponentKindDeclarationError{Kind: testKind}
	if duplicate.Kind != testKind {
		t.Fatalf("Kind = %s, want %s", duplicate.Kind, testKind)
	}
}
