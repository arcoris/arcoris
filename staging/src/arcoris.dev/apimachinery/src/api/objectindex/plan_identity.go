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

import "arcoris.dev/apimachinery/api/objectquery"

// planIdentity narrows candidates with exact storage-identity requirements.
//
// Namespace/name together use the combined object-name index. A single
// namespace or name requirement uses its single-dimension index. An absent
// identity selector leaves the existing candidate plan unchanged.
func (idx Index) planIdentity(
	plan candidatePlan,
	identity objectquery.IdentitySelector,
) candidatePlan {
	namespace, hasNamespace := identity.Namespace.Namespace()
	name, hasName := identity.Name.Name()

	switch {
	case hasNamespace && hasName:
		return plan.constrain(candidateSetFromPositions(
			len(idx.items),
			idx.byObjectName[objectNameKey{namespace: namespace, name: name}],
		))
	case hasNamespace:
		return plan.constrain(candidateSetFromPositions(len(idx.items), idx.byNamespace[namespace]))
	case hasName:
		return plan.constrain(candidateSetFromPositions(len(idx.items), idx.byName[name]))
	default:
		return plan
	}
}
