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

// validateResourceMatch checks the GVK family identity against the resource.
//
// The resource collection name is intentionally not checked here. Object type
// metadata carries GVK, while resource names belong to GVR routing and serving
// layers outside this package.
func validateResourceMatch(
	gvk apiidentity.GroupVersionKind,
	def resource.Definition,
) error {
	if gvk.Group != def.Group() {
		return errorf(
			pathObjectTypeMeta,
			ErrResourceMismatch,
			ErrorReasonResourceMismatch,
			"object group %q does not match resource group %q",
			gvk.Group,
			def.Group(),
		)
	}

	if gvk.Kind != def.Kind() {
		return errorf(
			pathObjectTypeMeta,
			ErrResourceMismatch,
			ErrorReasonResourceMismatch,
			"object kind %q does not match resource kind %q",
			gvk.Kind,
			def.Kind(),
		)
	}

	return nil
}
