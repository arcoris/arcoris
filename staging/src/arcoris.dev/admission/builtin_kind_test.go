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

func TestBuiltinKindDescriptors(t *testing.T) {
	t.Parallel()

	descriptors := BuiltinKindDescriptors()
	wantKinds := map[ComponentKind]bool{
		KindBulkhead:        false,
		KindRetryBudget:     false,
		KindDeadline:        false,
		KindRateLimiter:     false,
		KindQueue:           false,
		KindScheduler:       false,
		KindWorkerPool:      false,
		KindOverloadGate:    false,
		KindTenantIsolation: false,
	}

	for _, descriptor := range descriptors {
		if !descriptor.IsValid() {
			t.Fatalf("built-in descriptor should be valid: %+v", descriptor)
		}
		if descriptor.Capabilities.IsZero() {
			t.Fatalf("built-in descriptor should declare capabilities: %+v", descriptor)
		}
		if _, known := wantKinds[descriptor.Kind]; !known {
			t.Fatalf("unexpected built-in kind %q", descriptor.Kind)
		}
		wantKinds[descriptor.Kind] = true
	}
	for kind, found := range wantKinds {
		if !found {
			t.Fatalf("missing built-in kind %q", kind)
		}
	}
}

func TestBuiltinKindDescriptorsReturnsCopy(t *testing.T) {
	t.Parallel()

	descriptors := BuiltinKindDescriptors()
	descriptors[0].Kind = "mutated_kind"

	fresh := BuiltinKindDescriptors()
	if fresh[0].Kind == "mutated_kind" {
		t.Fatal("mutating returned descriptors should not mutate built-in catalog")
	}
}

func TestNewBuiltinKindRegistry(t *testing.T) {
	t.Parallel()

	registry := NewBuiltinKindRegistry()
	for _, descriptor := range BuiltinKindDescriptors() {
		if got, ok := registry.Lookup(descriptor.Kind); !ok || got != descriptor {
			t.Fatalf("Lookup(%q) = (%+v, %v), want built-in descriptor", descriptor.Kind, got, ok)
		}
	}
}
