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
	"arcoris.dev/apimachinery/api/meta/annotations"
	"arcoris.dev/apimachinery/api/objectquery"
)

// planAnnotations narrows candidates for positive annotation requirements.
func (idx indexes) planAnnotations(
	plan candidatePlan,
	requirements []objectquery.AnnotationRequirement,
) candidatePlan {
	for _, req := range requirements {
		if next, ok := idx.annotationCandidates(req); ok {
			plan = plan.constrain(next)
		}
	}

	return plan
}

// annotationCandidates returns candidates for one indexable annotation
// requirement. Non-indexable negative operators are returned as residual-only.
func (idx indexes) annotationCandidates(req objectquery.AnnotationRequirement) (keySet, bool) {
	key := annotations.Key(req.Key())
	switch req.Operator() {
	case objectquery.OperatorExists:
		return idx.byAnnotationKey[key], true
	case objectquery.OperatorEquals:
		values := req.Values()
		return idx.byAnnotationValue[annotationValueKey{
			key:   key,
			value: annotations.Value(values[0]),
		}], true
	case objectquery.OperatorIn:
		sets := make([]keySet, 0, len(req.Values()))
		for _, value := range req.Values() {
			sets = append(sets, idx.byAnnotationValue[annotationValueKey{
				key:   key,
				value: annotations.Value(value),
			}])
		}
		return unionKeySets(sets...), true
	default:
		return nil, false
	}
}
