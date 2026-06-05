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
)

func TestReasonDescriptors(t *testing.T) {
	descriptors := ReasonDescriptors()
	if len(descriptors) == 0 {
		t.Fatal("ReasonDescriptors returned empty slice")
	}
	requireReason(t, descriptors, admission.ReasonAdmitted)
	requireReason(t, descriptors, ReasonCapacityExhausted)
	requireReason(t, descriptors, ReasonBudgetExhausted)
	requireReason(t, descriptors, ReasonDeadlineExceeded)
	requireReason(t, descriptors, ReasonCanceled)
}

func TestReasonDescriptorsAreFresh(t *testing.T) {
	first := ReasonDescriptors()
	second := ReasonDescriptors()
	first[0].Reason = "mutated_reason"

	if second[0].Reason == "mutated_reason" {
		t.Fatal("ReasonDescriptors shared mutable slice storage")
	}
}

func TestReasonDescriptorsAreValid(t *testing.T) {
	for _, descriptor := range ReasonDescriptors() {
		if !descriptor.IsValid() {
			t.Fatalf("descriptor is invalid: %+v", descriptor)
		}
		if descriptor.Summary == "" {
			t.Fatalf("descriptor has empty summary: %+v", descriptor)
		}
	}
}
