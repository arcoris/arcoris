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

import "testing"

func TestEffectSet(t *testing.T) {
	set := NewEffectSet(EffectCapabilityNone, EffectCapabilityOwned)
	if !set.Has(EffectCapabilityNone) {
		t.Fatal("set does not contain none")
	}
	if !set.Has(EffectCapabilityOwned) {
		t.Fatal("set does not contain owned")
	}
	if set.Has(EffectCapabilityQueued) {
		t.Fatal("set unexpectedly contains queued")
	}
	if !set.IsValid() {
		t.Fatal("set is invalid")
	}
	if set.IsZero() {
		t.Fatal("set is zero")
	}
}

func TestEffectSetZeroIsValidAndUnspecified(t *testing.T) {
	var set EffectSet
	if !set.IsValid() {
		t.Fatal("zero effect set is invalid")
	}
	if !set.IsZero() {
		t.Fatal("zero effect set is not zero")
	}
	if set.Has(0) {
		t.Fatal("zero capability reported present")
	}
	if set.Has(EffectCapability(1 << 7)) {
		t.Fatal("unknown capability reported present in zero set")
	}
}

func TestEffectSetRejectsUnknownBits(t *testing.T) {
	set := NewEffectSet(EffectCapability(1 << 7))
	if set.IsValid() {
		t.Fatal("unknown effect bits were accepted")
	}
	if EffectSet(1<<7)&set == 0 {
		t.Fatal("constructor dropped unknown effect bits")
	}
}

func TestEffectSetWithPreservesExistingBits(t *testing.T) {
	set := NewEffectSet(EffectCapabilityNone).With(EffectCapabilityOwned)
	if !set.Has(EffectCapabilityNone) || !set.Has(EffectCapabilityOwned) {
		t.Fatal("With did not preserve existing effect bits")
	}
}
