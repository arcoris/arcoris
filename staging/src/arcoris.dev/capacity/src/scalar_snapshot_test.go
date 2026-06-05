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

package capacity_test

import (
	"testing"

	"arcoris.dev/capacity"
)

func TestNewScalarSnapshotDerivesAvailabilityAndDebt(t *testing.T) {
	t.Parallel()

	available := capacity.NewScalarSnapshot(4, 3)
	if !available.IsValid() || available.Available != 1 || available.Debt != 0 {
		t.Fatalf("available snapshot = %+v", available)
	}

	debt := capacity.NewScalarSnapshot(2, 3)
	if !debt.IsValid() || !debt.HasDebt() || debt.Available != 0 || debt.Debt != 1 {
		t.Fatalf("debt snapshot = %+v", debt)
	}
}

func TestScalarSnapshotCanReserve(t *testing.T) {
	t.Parallel()

	snap := capacity.NewScalarSnapshot(4, 2)
	if !snap.CanReserve(2) {
		t.Fatal("CanReserve(2) = false, want true")
	}
	if snap.CanReserve(3) {
		t.Fatal("CanReserve(3) = true, want false")
	}
	if capacity.NewScalarSnapshot(1, 2).CanReserve(1) {
		t.Fatal("CanReserve while in debt = true, want false")
	}
}
