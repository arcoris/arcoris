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

import "arcoris.dev/apimachinery/api/objectquery"

// planIdentity narrows candidates with exact object identity requirements.
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
