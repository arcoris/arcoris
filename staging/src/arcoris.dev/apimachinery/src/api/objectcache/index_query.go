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
	"arcoris.dev/apimachinery/api/meta/labels"
	"arcoris.dev/apimachinery/api/objectquery"
	"arcoris.dev/apimachinery/api/objectstore"
)

type candidatePlan struct {
	constrained bool
	keys        keySet
}

func (idx indexes) plan(predicate objectquery.Predicate) candidatePlan {
	query := predicate.Query()

	var plan candidatePlan
	plan = idx.planIdentity(plan, query.Identity)
	plan = idx.planLabels(plan, query.Labels.Requirements())
	plan = idx.planAnnotations(plan, query.Annotations.Requirements())

	return plan
}

func (plan candidatePlan) constrain(next keySet) candidatePlan {
	if !plan.constrained {
		return candidatePlan{constrained: true, keys: next.clone()}
	}

	return candidatePlan{
		constrained: true,
		keys:        intersectKeySets(plan.keys, next),
	}
}

func (plan candidatePlan) includes(key objectstore.Key) bool {
	return !plan.constrained || plan.keys.has(key)
}

func (idx indexes) planIdentity(
	plan candidatePlan,
	identity objectquery.IdentitySelector,
) candidatePlan {
	namespace, hasNamespace := identity.Namespace.Namespace()
	name, hasName := identity.Name.Name()

	switch {
	case hasNamespace && hasName:
		return plan.constrain(idx.byObject[objectNameKey{namespace: namespace, name: name}])
	case hasNamespace:
		return plan.constrain(idx.byNamespace[namespace])
	case hasName:
		return plan.constrain(idx.byName[name])
	default:
		return plan
	}
}

func (idx indexes) planLabels(
	plan candidatePlan,
	requirements []objectquery.LabelRequirement,
) candidatePlan {
	for _, req := range requirements {
		if next, ok := idx.labelCandidates(req); ok {
			plan = plan.constrain(next)
		}
	}

	return plan
}

func (idx indexes) labelCandidates(req objectquery.LabelRequirement) (keySet, bool) {
	key := labels.Key(req.Key())
	switch req.Operator() {
	case objectquery.OperatorExists:
		return idx.byLabelKey[key], true
	case objectquery.OperatorEquals:
		values := req.Values()
		return idx.byLabelValue[labelValueKey{key: key, value: labels.Value(values[0])}], true
	case objectquery.OperatorIn:
		sets := make([]keySet, 0, len(req.Values()))
		for _, value := range req.Values() {
			sets = append(sets, idx.byLabelValue[labelValueKey{
				key:   key,
				value: labels.Value(value),
			}])
		}
		return unionKeySets(sets...), true
	default:
		return nil, false
	}
}

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
