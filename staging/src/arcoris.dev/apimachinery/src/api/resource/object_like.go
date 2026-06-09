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

package resource

import (
	"fmt"

	"arcoris.dev/apimachinery/api/types"
)

// requireObjectLike converts object-like probing into a VersionDefinition
// validation error.
//
// The caller provides the reason because Desired and Observed use distinct
// diagnostics while sharing the same structural root rule.
func requireObjectLike(desc types.Descriptor, resolver types.Resolver, path string, reason ErrorReason, label string) error {
	ok, detail := objectLike(desc, resolver, make(map[types.TypeName]bool), label)
	if ok {
		return nil
	}

	return versionError(path, reason, detail)
}

// objectLike reports whether desc is a direct object or resolver-proven object
// reference.
//
// Nil resolvers intentionally do not prove references. api/types may accept a
// syntactically valid unresolved DescriptorRef during local validation, but resource
// surfaces need object-like proof because resource definitions define API
// object roots.
func objectLike(desc types.Descriptor, resolver types.Resolver, resolving map[types.TypeName]bool, label string) (bool, string) {
	switch desc.Code() {
	case types.DescriptorObject:
		return true, ""

	case types.DescriptorRef:
		view, _ := desc.AsRef()
		name := view.Name()

		if resolver == nil {
			return false, fmt.Sprintf(
				"%s root reference %q requires a resolver so the resource surface can be proven object-like",
				label,
				name,
			)
		}

		if resolving[name] {
			return false, fmt.Sprintf("%s root reference %q is recursive", label, name)
		}

		def, ok := resolver.Resolve(name)
		if !ok {
			return false, fmt.Sprintf(
				"%s root reference %q was not found in resolver",
				label,
				name,
			)
		}

		next := make(map[types.TypeName]bool, len(resolving)+1)
		for candidate, active := range resolving {
			next[candidate] = active
		}
		next[name] = true

		return objectLike(def.Descriptor(), resolver, next, label)

	default:
		return false, fmt.Sprintf(
			"%s root must be object or reference to object, got %s",
			label,
			desc.Code(),
		)
	}
}
