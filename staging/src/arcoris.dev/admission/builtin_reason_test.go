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

import "testing"

func TestBuiltinReasonDescriptors(t *testing.T) {
	t.Parallel()

	descriptors := BuiltinReasonDescriptors()
	wantReasons := map[Reason]bool{
		ReasonAdmitted:          false,
		ReasonDenied:            false,
		ReasonQueued:            false,
		ReasonDeferred:          false,
		ReasonCapacityExhausted: false,
		ReasonBudgetExhausted:   false,
		ReasonRateLimited:       false,
		ReasonOverloaded:        false,
		ReasonBackpressured:     false,
		ReasonClosed:            false,
		ReasonDraining:          false,
		ReasonDeadlineExceeded:  false,
		ReasonCanceled:          false,
		ReasonPolicyDenied:      false,
	}

	for _, descriptor := range descriptors {
		if !descriptor.IsValid() {
			t.Fatalf("built-in descriptor should be valid: %+v", descriptor)
		}
		if descriptor.Capabilities.IsZero() {
			t.Fatalf("built-in descriptor should declare capabilities: %+v", descriptor)
		}
		if found, known := wantReasons[descriptor.Reason]; !known {
			t.Fatalf("unexpected built-in reason %q", descriptor.Reason)
		} else if found {
			t.Fatalf("duplicate built-in reason %q", descriptor.Reason)
		}
		wantReasons[descriptor.Reason] = true
	}
	for reason, found := range wantReasons {
		if !found {
			t.Fatalf("missing built-in reason %q", reason)
		}
	}
}

func TestBuiltinReasonDescriptorsReturnsCopy(t *testing.T) {
	t.Parallel()

	descriptors := BuiltinReasonDescriptors()
	descriptors[0].Reason = "mutated_reason"

	fresh := BuiltinReasonDescriptors()
	if fresh[0].Reason == "mutated_reason" {
		t.Fatal("mutating returned descriptors should not mutate built-in catalog")
	}
}

func TestNewBuiltinReasonRegistry(t *testing.T) {
	t.Parallel()

	registry := NewBuiltinReasonRegistry()
	for _, descriptor := range BuiltinReasonDescriptors() {
		if got, ok := registry.Lookup(descriptor.Reason); !ok || got != descriptor {
			t.Fatalf("Lookup(%q) = (%+v, %v), want built-in descriptor", descriptor.Reason, got, ok)
		}
	}
}
