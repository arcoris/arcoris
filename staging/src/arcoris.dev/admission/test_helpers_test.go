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

// someString builds a present string Maybe through the admission-local helper.
//
// Tests use this instead of importing arcoris.dev/value/maybe directly so the
// package-local alias remains the only optional-value vocabulary inside
// admission tests.
func someString(value string) Maybe[string] {
	return some(value)
}

// noneString builds an absent string Maybe through the admission-local helper.
func noneString() Maybe[string] {
	return none[string]()
}

// noneMetadata builds an absent NoMetadata Maybe for invalid-shape tests.
func noneMetadata() Maybe[NoMetadata] {
	return none[NoMetadata]()
}

// assertPanicString verifies the exact panic string used for stable nil
// receiver contracts.
func assertPanicString(
	t *testing.T,
	want string,
	call func(),
) {
	t.Helper()

	defer func() {
		got := recover()
		if got == nil {
			t.Fatalf("expected panic %q", want)
		}
		if got != want {
			t.Fatalf("panic = %q, want %q", got, want)
		}
	}()

	call()
}

// requireCapability verifies that set advertises capability.
//
// Built-in descriptor tests use this helper to make semantic capability
// expectations read as catalog requirements instead of low-level bit checks.
func requireCapability(
	t *testing.T,
	set CapabilitySet,
	capability Capability,
) {
	t.Helper()

	if !set.Has(capability) {
		t.Fatalf("capabilities %08b should contain %08b", set, capability)
	}
}
