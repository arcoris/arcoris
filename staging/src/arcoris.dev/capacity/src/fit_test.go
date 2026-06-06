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

func TestFitValidityAndPredicates(t *testing.T) {
	fits := capacity.Fit{Refusal: capacity.RefusalNone}
	if !fits.Fits() || fits.Refused() || !fits.IsValid() {
		t.Fatalf("fit success predicates invalid: %#v", fits)
	}

	missing := capacity.Fit{
		Refusal: capacity.RefusalInsufficient,
		Missing: vector(t, entry("worker_slots", 1)),
	}
	if missing.Fits() || !missing.Refused() || !missing.IsValid() {
		t.Fatalf("missing fit predicates invalid: %#v", missing)
	}

	debt := capacity.Fit{
		Refusal: capacity.RefusalDebt,
		Debt:    vector(t, entry("memory_bytes", 2)),
	}
	if debt.Fits() || !debt.Refused() || !debt.IsValid() {
		t.Fatalf("debt fit predicates invalid: %#v", debt)
	}
}

func TestFitValidityAllowsMixedDiagnostics(t *testing.T) {
	debtWithMissing := capacity.Fit{
		Refusal: capacity.RefusalDebt,
		Missing: vector(t, entry("worker_slots", 1)),
		Debt:    vector(t, entry("memory_bytes", 2)),
	}
	if !debtWithMissing.IsValid() {
		t.Fatalf("debt with missing diagnostics is invalid: %#v", debtWithMissing)
	}

	unknownWithDebt := capacity.Fit{
		Refusal: capacity.RefusalUnknownResource,
		Missing: vector(t, entry("queue_slots", 1)),
		Debt:    vector(t, entry("memory_bytes", 2)),
	}
	if !unknownWithDebt.IsValid() {
		t.Fatalf("unknown with debt diagnostics is invalid: %#v", unknownWithDebt)
	}
}
