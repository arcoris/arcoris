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

package valuevalidation

import (
	"arcoris.dev/apimachinery/api/fieldpath"
	"arcoris.dev/apimachinery/api/types"
	"arcoris.dev/apimachinery/api/value"
)

// validateRef resolves a TypeRef and validates the value at the same semantic path.
func (v *validator) validateRef(
	path fieldpath.Path,
	val value.Value,
	descriptor types.Type,
	depth int,
) {
	if depth >= v.maxDepth {
		v.add(path, ErrReferenceCycle, ErrorReasonReferenceCycle, "maximum TypeRef validation depth reached")
		return
	}

	view, ok := descriptor.Ref()
	if !ok {
		v.add(path, ErrInvalidDescriptor, ErrorReasonInvalidDescriptor, "descriptor is not a reference")
		return
	}

	name := view.Name()
	if v.resolver == nil {
		v.addf(path, ErrUnresolvedRef, ErrorReasonUnresolvedRef, "reference %q has no resolver", name)
		return
	}
	if v.resolving[name] {
		v.addf(path, ErrReferenceCycle, ErrorReasonReferenceCycle, "reference %q is recursive", name)
		return
	}

	definition, ok := v.resolver.ResolveType(name)
	if !ok {
		v.addf(path, ErrUnresolvedRef, ErrorReasonUnresolvedRef, "reference %q was not found", name)
		return
	}

	v.resolving[name] = true
	defer delete(v.resolving, name)

	v.validate(path, val, definition.Type(), depth+1)
}
