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

package objectquery

import "testing"

// TestBuildPlanExtractsKeyResourceAndAndConstraints verifies AND constraints
// are accumulated while OR/NOT remain residual-only.
func TestBuildPlanExtractsKeyResourceAndAndConstraints(t *testing.T) {
	query := mustAnd(t,
		mustQ(ResourceEquals(testResource)),
		mustQ(ObjectWithName("worker-1")),
	)

	plan := buildPlan(mustPredicate(t, query).expr)

	if len(plan.constraints) != 2 {
		t.Fatalf("constraints = %d; want 2", len(plan.constraints))
	}
}

// TestBuildPlanLeavesOrAndNotResidual verifies unsafe boolean shapes do not
// produce misleading narrowing hints.
func TestBuildPlanLeavesOrAndNotResidual(t *testing.T) {
	orQuery := mustOr(t, mustQ(LabelEquals("env", "prod")), mustQ(LabelEquals("env", "qa")))
	notQuery := mustNot(t, mustQ(LabelEquals("env", "prod")))

	if got := buildPlan(mustPredicate(t, orQuery).expr); len(got.constraints) != 0 {
		t.Fatalf("OR constraints = %d; want 0", len(got.constraints))
	}
	if got := buildPlan(mustPredicate(t, notQuery).expr); len(got.constraints) != 0 {
		t.Fatalf("NOT constraints = %d; want 0", len(got.constraints))
	}
}
