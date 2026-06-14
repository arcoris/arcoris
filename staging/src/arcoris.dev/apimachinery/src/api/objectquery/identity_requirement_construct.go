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

// NamespaceEquals constructs an exact namespace requirement.
//
// A zero namespace is valid and matches objects whose storage identity has no
// namespace. Invalid non-zero namespace syntax is rejected by api/meta/identity.
func NamespaceEquals(namespace metaidentity.Namespace) (NamespaceRequirement, error) {
	if err := namespace.ValidateLexical(); err != nil {
		return NamespaceRequirement{}, wrapf(
			"query.identity.namespace",
			ErrInvalidQuery,
			ErrorReasonInvalidIdentity,
			err,
			"namespace requirement is invalid",
		)
	}

	return NamespaceRequirement{set: true, namespace: namespace}, nil
}

// NameEquals constructs an exact object name requirement.
func NameEquals(name metaidentity.Name) (NameRequirement, error) {
	if err := name.ValidateLexical(); err != nil {
		return NameRequirement{}, wrapf(
			"query.identity.name",
			ErrInvalidQuery,
			ErrorReasonInvalidIdentity,
			err,
			"name requirement is invalid",
		)
	}

	return NameRequirement{set: true, name: name}, nil
}
