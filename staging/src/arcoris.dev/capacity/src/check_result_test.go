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

func TestCheckResultValidation(t *testing.T) {
	t.Parallel()

	if !(capacity.CheckResult{Status: capacity.ReserveStatusReserved}).IsValid() {
		t.Fatal("reserved empty diagnostic result is invalid")
	}
	if (capacity.CheckResult{Status: capacity.ReserveStatusReserved, Missing: vector(t, entry("worker_slots", 1))}).IsValid() {
		t.Fatal("reserved result with missing vector is valid")
	}
	if !(capacity.CheckResult{
		Status:  capacity.ReserveStatusInsufficient,
		Missing: vector(t, entry("worker_slots", 1)),
	}).IsValid() {
		t.Fatal("insufficient result with missing vector is invalid")
	}
	if !(capacity.CheckResult{
		Status: capacity.ReserveStatusDebt,
		Debt:   vector(t, entry("memory_bytes", 1)),
	}).IsValid() {
		t.Fatal("debt result with debt vector is invalid")
	}
}
