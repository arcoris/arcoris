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

package objectcache

import (
	"arcoris.dev/apimachinery/api/objectquery"
	"arcoris.dev/apimachinery/api/objectstore"
)

// candidatePlan describes the current indexed narrowing state for one query.
//
// constrained=false means the query has no positive indexed requirements yet,
// so all ordered keys remain candidates.
type candidatePlan struct {
	// constrained marks whether keys is an actual candidate set or an implicit
	// all-keys universe.
	constrained bool

	// keys is the current intersection of positive indexed requirements.
	keys keySet
}

// plan inspects the canonical predicate query and builds an indexed candidate
// plan. Negative requirements intentionally remain residual and are handled by
// Predicate.Match during the final scan.
func (idx indexes) plan(predicate objectquery.Predicate) candidatePlan {
	query := predicate.Query()

	var plan candidatePlan
	plan = idx.planIdentity(plan, query.Identity)
	plan = idx.planLabels(plan, query.Labels.Requirements())
	plan = idx.planAnnotations(plan, query.Annotations.Requirements())

	return plan
}

// constrain ANDs one positive indexed requirement into the current plan.
func (plan candidatePlan) constrain(next keySet) candidatePlan {
	if !plan.constrained {
		return candidatePlan{constrained: true, keys: next.clone()}
	}

	return candidatePlan{
		constrained: true,
		keys:        intersectKeySets(plan.keys, next),
	}
}

// includes reports whether key survives the indexed candidate plan.
func (plan candidatePlan) includes(key objectstore.Key) bool {
	return !plan.constrained || plan.keys.has(key)
}
