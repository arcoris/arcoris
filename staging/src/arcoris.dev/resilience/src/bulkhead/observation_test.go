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

package bulkhead

import "testing"

func TestObservationAcceptedDeniedAndValidity(t *testing.T) {
	t.Parallel()

	b := New(1)
	lease, accepted, ok := b.TryAcquire()
	if !ok {
		t.Fatal("TryAcquire failed")
	}
	defer lease.Release()

	if !accepted.Accepted() {
		t.Fatal("accepted observation reports Accepted=false")
	}
	if accepted.Denied() {
		t.Fatal("accepted observation reports Denied=true")
	}
	if !accepted.IsValid() {
		t.Fatalf("accepted observation is invalid: %+v", accepted)
	}

	_, denied, ok := b.TryAcquire()
	if ok {
		t.Fatal("denied TryAcquire returned ok=true")
	}
	if denied.Accepted() {
		t.Fatal("denied observation reports Accepted=true")
	}
	if !denied.Denied() {
		t.Fatal("denied observation reports Denied=false")
	}
	if !denied.IsValid() {
		t.Fatalf("denied observation is invalid: %+v", denied)
	}
}

func TestZeroObservationIsInvalid(t *testing.T) {
	t.Parallel()

	var observation Observation
	if observation.IsValid() {
		t.Fatal("zero observation is valid")
	}
}
