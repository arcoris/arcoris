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

// candidatePlan carries the current candidate intersection.
//
// constrained=false means no positive index constraint was useful, so selection
// must scan every indexed item and rely entirely on objectquery.Predicate.
type candidatePlan struct {
	constrained bool
	set         candidateSet
}

// constrain ANDs next into plan.
//
// The first positive constraint establishes the candidate universe. Later
// positive constraints intersect with that universe.
func (plan candidatePlan) constrain(next candidateSet) candidatePlan {
	if !plan.constrained {
		return candidatePlan{constrained: true, set: next}
	}

	return candidatePlan{
		constrained: true,
		set:         intersectCandidateSets(plan.set, next),
	}
}

// includes reports whether pos should be evaluated by the final predicate.
//
// An unconstrained plan includes every item because no positive index lookup was
// useful.
func (plan candidatePlan) includes(pos int) bool {
	return !plan.constrained || plan.set.has(pos)
}
