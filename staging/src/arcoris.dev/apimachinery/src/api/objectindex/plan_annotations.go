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
	"arcoris.dev/apimachinery/api/meta/annotations"
	"arcoris.dev/apimachinery/api/objectquery"
)

// planAnnotations intersects candidate sets for every positive annotation
// requirement.
//
// Negative annotation operators are residual-only in v1 for the same reason as
// labels: absent-key matching is semantic query behavior, not index behavior.
func (idx Index) planAnnotations(
	plan candidatePlan,
	requirements []objectquery.AnnotationRequirement,
) candidatePlan {
	for _, req := range requirements {
		next, ok := idx.annotationRequirementCandidates(req)
		if ok {
			plan = plan.constrain(next)
		}
	}

	return plan
}

// annotationRequirementCandidates returns a positive candidate set for
// indexable annotation operators.
//
// The bool result is false when the requirement must remain a residual
// predicate check.
func (idx Index) annotationRequirementCandidates(
	req objectquery.AnnotationRequirement,
) (candidateSet, bool) {
	key := annotations.Key(req.Key())
	switch req.Operator() {
	case objectquery.OperatorExists:
		return candidateSetFromPositions(len(idx.items), idx.byAnnotationKey[key]), true
	case objectquery.OperatorEquals:
		values := req.Values()
		return candidateSetFromPositions(
			len(idx.items),
			idx.byAnnotationValue[annotationValueKey{
				key:   key,
				value: annotations.Value(values[0]),
			}],
		), true
	case objectquery.OperatorIn:
		return idx.annotationInCandidates(key, req.Values()), true
	default:
		return candidateSet{}, false
	}
}

// annotationInCandidates unions candidates for one annotation In requirement.
//
// The union implements the OR semantics inside the value set; later
// candidatePlan intersections provide AND semantics across requirements.
func (idx Index) annotationInCandidates(key annotations.Key, values []string) candidateSet {
	groups := make([][]int, 0, len(values))
	for _, value := range values {
		groups = append(groups, idx.byAnnotationValue[annotationValueKey{
			key:   key,
			value: annotations.Value(value),
		}])
	}

	return unionCandidateSets(len(idx.items), groups...)
}
