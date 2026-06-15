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

package objectindex

import (
	"testing"

	"arcoris.dev/apimachinery/api/objectquery"
)

func TestPlanCandidatesUsesCanonicalPredicateQuery(t *testing.T) {
	query := objectquery.Query{
		Labels: mustLabelSelector(
			t,
			mustLabelEquals(t, "tier", "backend"),
			mustLabelIn(t, "env", "qa", "prod", "prod"),
		),
	}
	predicate, err := objectquery.Compile(query)
	requireNoError(t, err)

	plan := Build(testItems()).planCandidates(predicate)

	requirePlanIncludesOnly(t, plan, len(testItems()), 0, 1)
}

func TestPlanCandidatesLeavesResidualOnlyQueryUnconstrained(t *testing.T) {
	query := objectquery.Query{
		Labels: mustLabelSelector(t, mustLabelNotEquals(t, "env", "prod")),
	}
	predicate, err := objectquery.Compile(query)
	requireNoError(t, err)

	plan := Build(testItems()).planCandidates(predicate)

	if plan.constrained {
		t.Fatal("constrained = true; want false")
	}
}
