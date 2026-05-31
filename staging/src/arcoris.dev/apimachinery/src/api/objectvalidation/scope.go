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
	apiobject "arcoris.dev/apimachinery/api/object"
	"arcoris.dev/apimachinery/api/resource"
)

// validateScope applies only minimal resource scope compatibility.
//
// Global resources must not carry a namespace. Namespaced resources may carry
// one or omit it because namespace defaulting and persistence requirements
// belong to higher layers.
func validateScope[D any, O any](
	obj apiobject.Object[D, O],
	def resource.Definition,
) error {
	switch def.Scope() {
	case resource.ScopeGlobal:
		if !obj.ObjectMeta.Namespace.IsZero() {
			return errorf(
				pathObjectNamespace,
				ErrInvalidScope,
				ErrorReasonInvalidScope,
				"global resource %s must not carry namespace %q",
				def.GroupKind(),
				obj.ObjectMeta.Namespace,
			)
		}

	case resource.ScopeNamespaced:
		return nil

	default:
		return errorf(
			pathResourceScope,
			ErrInvalidScope,
			ErrorReasonInvalidScope,
			"resource %s has invalid scope %q",
			def.GroupKind(),
			def.Scope(),
		)
	}

	return nil
}
