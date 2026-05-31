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

package objectvalidation

import (
	apiidentity "arcoris.dev/apimachinery/api/identity"
	"arcoris.dev/apimachinery/api/resource"
)

// resolveVersion selects the contract version named by object type metadata.
//
// Exposed and canonical markers are not checked here; they are serving and
// storage/conversion concerns, not baseline resource contract conformance.
func resolveVersion(
	gvk apiidentity.GroupVersionKind,
	def resource.Definition,
) (resource.VersionDefinition, error) {
	version, ok := def.Version(gvk.Version)
	if !ok {
		return resource.VersionDefinition{}, errorf(
			pathResourceVersions,
			ErrVersionNotDefined,
			ErrorReasonVersionNotDefined,
			"resource %s does not define object API version %q",
			def.GroupKind(),
			gvk.Version,
		)
	}

	return version, nil
}
