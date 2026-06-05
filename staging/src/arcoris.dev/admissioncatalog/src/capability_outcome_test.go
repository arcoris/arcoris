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

func TestOutcomeSet(t *testing.T) {
	set := NewOutcomeSet(OutcomeCapabilityAdmit, OutcomeCapabilityDeny)
	if !set.Has(OutcomeCapabilityAdmit) {
		t.Fatal("set does not contain admit")
	}
	if !set.Has(OutcomeCapabilityDeny) {
		t.Fatal("set does not contain deny")
	}
	if set.Has(OutcomeCapabilityQueue) {
		t.Fatal("set unexpectedly contains queue")
	}
	if !set.IsValid() {
		t.Fatal("set is invalid")
	}
	if set.IsZero() {
		t.Fatal("set is zero")
	}
}

func TestOutcomeSetZeroIsValidAndUnspecified(t *testing.T) {
	var set OutcomeSet
	if !set.IsValid() {
		t.Fatal("zero outcome set is invalid")
	}
	if !set.IsZero() {
		t.Fatal("zero outcome set is not zero")
	}
	if set.Has(0) {
		t.Fatal("zero capability reported present")
	}
}

func TestOutcomeSetRejectsUnknownBits(t *testing.T) {
	set := OutcomeSet(1 << 7)
	if set.IsValid() {
		t.Fatal("unknown outcome bits were accepted")
	}
}
