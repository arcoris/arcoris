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

package typeref

import (
	"arcoris.dev/apimachinery/api/fieldpath"
	"arcoris.dev/apimachinery/api/types"
)

// Resolve resolves one DescriptorRef edge and leaves stack ownership to the caller.
//
// Call Enter with the returned name before descending into the resolved
// descriptor. Keeping stack entry separate lets callers preserve their own
// control flow while sharing the failure checks.
func (r *Resolver) Resolve(
	path fieldpath.Path,
	descriptor types.Descriptor,
	depth int,
) (types.TypeName, types.Descriptor, error) {
	if depth >= r.maxDepth {
		return "", types.Descriptor{}, failure(
			path,
			FailureReferenceCycle,
			"maximum DescriptorRef traversal depth reached",
		)
	}

	view, ok := descriptor.AsRef()
	if !ok {
		return "", types.Descriptor{}, failure(
			path,
			FailureInvalidDescriptor,
			"descriptor is not a reference",
		)
	}

	name := view.Name()
	if r.resolver == nil {
		return "", types.Descriptor{}, failuref(
			path,
			FailureUnresolvedRef,
			"reference %q has no resolver",
			name,
		)
	}
	if r.active[name] {
		return "", types.Descriptor{}, failuref(
			path,
			FailureReferenceCycle,
			"reference %q is recursive",
			name,
		)
	}

	definition, ok := r.resolver.Resolve(name)
	if !ok {
		return "", types.Descriptor{}, failuref(
			path,
			FailureUnresolvedRef,
			"reference %q was not found",
			name,
		)
	}

	return name, definition.Descriptor(), nil
}
