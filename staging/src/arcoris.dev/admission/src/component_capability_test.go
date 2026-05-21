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

func TestCapabilitySet(t *testing.T) {
	t.Parallel()

	set := NewCapabilitySet(
		CapabilityAdmit,
		CapabilityDeny,
		CapabilityEffectOwned,
	)

	if !set.IsValid() {
		t.Fatal("set should be valid")
	}
	if !set.Has(CapabilityAdmit) {
		t.Fatal("set should contain admit capability")
	}
	if !set.Has(CapabilityDeny) {
		t.Fatal("set should contain deny capability")
	}
	if !set.Has(CapabilityEffectOwned) {
		t.Fatal("set should contain owned-effect capability")
	}
	if set.Has(CapabilityQueue) {
		t.Fatal("set should not contain queue capability")
	}
	if set.Has(0) {
		t.Fatal("zero capability should never be reported as present")
	}
}

func TestCapabilitySetZeroAndUnknownBits(t *testing.T) {
	t.Parallel()

	var zero CapabilitySet
	if !zero.IsValid() {
		t.Fatal("zero set should be valid")
	}
	if !zero.IsZero() {
		t.Fatal("zero set should report unspecified capabilities")
	}

	unknown := CapabilitySet(1 << 15)
	if unknown.IsValid() {
		t.Fatal("unknown capability bit should be invalid")
	}
	if unknown.IsZero() {
		t.Fatal("unknown capability bit should not be treated as zero")
	}
}
