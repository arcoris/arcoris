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
	"arcoris.dev/admission"
	"errors"
	"testing"
)

func TestReasonRegistryErrorsSupportIsAndAs(t *testing.T) {
	t.Parallel()

	invalid := InvalidReasonDescriptorError{
		Descriptor: ReasonDescriptor{Reason: "bad-reason"},
	}
	if !errors.Is(invalid, ErrInvalidReasonDescriptor) {
		t.Fatal("invalid descriptor error should match sentinel")
	}
	var invalidDescriptor InvalidReasonDescriptorError
	if !errors.As(invalid, &invalidDescriptor) {
		t.Fatal("invalid descriptor error should support errors.As")
	}
	if invalidDescriptor.Descriptor.Reason != "bad-reason" {
		t.Fatalf("invalid reason = %q, want bad-reason", invalidDescriptor.Descriptor.Reason)
	}

	duplicate := DuplicateReasonError{Reason: admission.ReasonDenied}
	if !errors.Is(duplicate, ErrReasonAlreadyRegistered) {
		t.Fatal("duplicate reason error should match sentinel")
	}
	var duplicateReason DuplicateReasonError
	if !errors.As(duplicate, &duplicateReason) {
		t.Fatal("duplicate reason error should support errors.As")
	}
	if duplicateReason.Reason != admission.ReasonDenied {
		t.Fatalf("duplicate reason = %q, want %q", duplicateReason.Reason, admission.ReasonDenied)
	}
}
