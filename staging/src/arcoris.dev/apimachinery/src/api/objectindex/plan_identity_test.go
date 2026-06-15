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

func TestPlanIdentityNarrowsByNamespaceAndName(t *testing.T) {
	identity := objectquery.IdentitySelector{
		Namespace: mustNamespaceEquals(t, "system"),
		Name:      mustNameEquals(t, "worker-2"),
	}

	plan := Build(testItems()).planIdentity(candidatePlan{}, identity)

	requirePlanIncludesOnly(t, plan, len(testItems()), 1)
}

func TestPlanIdentityNarrowsByNamespaceOnly(t *testing.T) {
	identity := objectquery.IdentitySelector{
		Namespace: mustNamespaceEquals(t, "system"),
	}

	plan := Build(testItems()).planIdentity(candidatePlan{}, identity)

	requirePlanIncludesOnly(t, plan, len(testItems()), 0, 1, 4)
}

func TestPlanIdentityNarrowsByNameOnly(t *testing.T) {
	identity := objectquery.IdentitySelector{
		Name: mustNameEquals(t, "worker-3"),
	}

	plan := Build(testItems()).planIdentity(candidatePlan{}, identity)

	requirePlanIncludesOnly(t, plan, len(testItems()), 2)
}

func TestPlanIdentityLeavesAbsentIdentityUnconstrained(t *testing.T) {
	plan := Build(testItems()).planIdentity(candidatePlan{}, objectquery.IdentitySelector{})

	if plan.constrained {
		t.Fatal("constrained = true; want false")
	}
}
