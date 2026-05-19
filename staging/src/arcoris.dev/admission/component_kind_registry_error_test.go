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

func TestKindRegistryErrorsSupportIsAndAs(t *testing.T) {
	t.Parallel()

	invalid := InvalidComponentKindDescriptorError{
		Descriptor: ComponentKindDescriptor{Kind: "bad-kind"},
	}
	if !errors.Is(invalid, ErrInvalidComponentKindDescriptor) {
		t.Fatal("invalid descriptor error should match sentinel")
	}
	var invalidKind InvalidComponentKindDescriptorError
	if !errors.As(invalid, &invalidKind) {
		t.Fatal("invalid descriptor error should support errors.As")
	}
	if invalidKind.Descriptor.Kind != "bad-kind" {
		t.Fatalf("invalid kind = %q, want bad-kind", invalidKind.Descriptor.Kind)
	}

	duplicate := DuplicateComponentKindError{Kind: KindBulkhead}
	if !errors.Is(duplicate, ErrComponentKindAlreadyRegistered) {
		t.Fatal("duplicate kind error should match sentinel")
	}

	var duplicateKind DuplicateComponentKindError
	if !errors.As(duplicate, &duplicateKind) {
		t.Fatal("duplicate kind error should support errors.As")
	}
	if duplicateKind.Kind != KindBulkhead {
		t.Fatalf("duplicate kind = %q, want %q", duplicateKind.Kind, KindBulkhead)
	}
}
