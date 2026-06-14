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

package objectquery

import metaidentity "arcoris.dev/apimachinery/api/meta/identity"

// NamespaceRequirement requires an exact metadata namespace value.
type NamespaceRequirement struct {
	// set distinguishes an absent requirement from an exact zero namespace match.
	set bool

	// namespace is the canonical namespace to match when set is true.
	namespace metaidentity.Namespace
}

// NameRequirement requires an exact metadata name value.
type NameRequirement struct {
	// set distinguishes an absent requirement from an exact name match.
	set bool

	// name is the canonical object name to match when set is true.
	name metaidentity.Name
}

// IsZero reports whether r carries no namespace requirement.
func (r NamespaceRequirement) IsZero() bool {
	return !r.set
}

// IsZero reports whether r carries no name requirement.
func (r NameRequirement) IsZero() bool {
	return !r.set
}
