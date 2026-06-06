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

func TestLeaseAmountReflectsWeightedLease(t *testing.T) {
	t.Parallel()

	b := New(2)
	lease, _, ok := b.TryAcquireAmount(2)
	if !ok {
		t.Fatal("TryAcquireAmount failed")
	}
	if lease.Amount() != 2 {
		t.Fatalf("Amount() = %d, want 2", lease.Amount())
	}

	lease.Release()
	if lease.Amount() != 2 {
		t.Fatalf("Amount() after release = %d, want 2", lease.Amount())
	}
}

func TestLeaseReleasedReflectsLeaseState(t *testing.T) {
	t.Parallel()

	b := New(1)
	lease, _, ok := b.TryAcquire()
	if !ok {
		t.Fatal("TryAcquire failed")
	}
	if lease.Released() {
		t.Fatal("Released() before release = true, want false")
	}

	lease.Release()
	if !lease.Released() {
		t.Fatal("Released() after release = false, want true")
	}
}
