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

func TestReasonErrorsSupportIsAndAs(t *testing.T) {
	invalid := InvalidReasonDescriptorError{Descriptor: reasonDescriptor(testReason)}
	if !errors.Is(invalid, ErrInvalidReasonDescriptor) {
		t.Fatal("invalid reason error does not match sentinel")
	}
	var invalidTyped InvalidReasonDescriptorError
	if !errors.As(invalid, &invalidTyped) {
		t.Fatal("invalid reason error does not expose typed value")
	}

	duplicate := DuplicateReasonDeclarationError{Reason: testReason}
	if !errors.Is(duplicate, ErrDuplicateReasonDeclaration) {
		t.Fatal("duplicate reason error does not match sentinel")
	}
	var duplicateTyped DuplicateReasonDeclarationError
	if !errors.As(duplicate, &duplicateTyped) {
		t.Fatal("duplicate reason error does not expose typed value")
	}
}

func TestReasonErrorsExposeDetails(t *testing.T) {
	invalid := InvalidReasonDescriptorError{
		Descriptor: reasonDescriptor(testReason),
		Path:       "input.reasons[0]",
	}
	if invalid.Descriptor.Reason != testReason {
		t.Fatalf("Descriptor.Reason = %s, want %s", invalid.Descriptor.Reason, testReason)
	}
	if invalid.Path == "" {
		t.Fatal("Path is empty")
	}

	duplicate := DuplicateReasonDeclarationError{Reason: testReason}
	if duplicate.Reason != testReason {
		t.Fatalf("Reason = %s, want %s", duplicate.Reason, testReason)
	}
}
