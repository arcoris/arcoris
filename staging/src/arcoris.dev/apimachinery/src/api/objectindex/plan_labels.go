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
	"arcoris.dev/apimachinery/api/meta/labels"
	"arcoris.dev/apimachinery/api/objectquery"
)

// planLabels intersects candidate sets for every positive label requirement.
//
// Negative label operators are residual-only in v1 because their correct
// semantics include absent-key matches. The final predicate check preserves
// those semantics after any positive narrowing.
func (idx Index) planLabels(
	plan candidatePlan,
	requirements []objectquery.LabelRequirement,
) candidatePlan {
	for _, req := range requirements {
		next, ok := idx.labelRequirementCandidates(req)
		if ok {
			plan = plan.constrain(next)
		}
	}

	return plan
}

// labelRequirementCandidates returns a positive candidate set for indexable
// label operators.
//
// The bool result is false for residual-only operators. Callers must then leave
// candidate planning unchanged and rely on the final predicate match.
func (idx Index) labelRequirementCandidates(req objectquery.LabelRequirement) (candidateSet, bool) {
	key := labels.Key(req.Key())
	switch req.Operator() {
	case objectquery.OperatorExists:
		return candidateSetFromPositions(len(idx.items), idx.byLabelKey[key]), true
	case objectquery.OperatorEquals:
		values := req.Values()
		return candidateSetFromPositions(
			len(idx.items),
			idx.byLabelValue[labelValueKey{key: key, value: labels.Value(values[0])}],
		), true
	case objectquery.OperatorIn:
		return idx.labelInCandidates(key, req.Values()), true
	default:
		return candidateSet{}, false
	}
}

// labelInCandidates unions candidates for one label In requirement.
//
// Values within an In requirement are ORed by objectquery semantics, while
// separate requirements are ANDed by candidatePlan.constrain.
func (idx Index) labelInCandidates(key labels.Key, values []string) candidateSet {
	groups := make([][]int, 0, len(values))
	for _, value := range values {
		groups = append(groups, idx.byLabelValue[labelValueKey{
			key:   key,
			value: labels.Value(value),
		}])
	}

	return unionCandidateSets(len(idx.items), groups...)
}
