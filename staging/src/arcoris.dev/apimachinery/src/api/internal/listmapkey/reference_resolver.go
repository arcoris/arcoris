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

package listmapkey

import (
	"fmt"

	"arcoris.dev/apimachinery/api/fieldpath"
	"arcoris.dev/apimachinery/api/types"
)

// referenceResolver resolves descriptor references for one ListMap key extraction.
type referenceResolver struct {
	typeResolver types.Resolver
	maxDepth     int
	activeRefs   map[types.TypeName]bool
}

// newReferenceResolver creates bounded resolver state for one extraction call.
func newReferenceResolver(opts Options) referenceResolver {
	return referenceResolver{
		typeResolver: opts.Resolver,
		maxDepth:     opts.normalizedMaxDepth(),
		activeRefs:   make(map[types.TypeName]bool),
	}
}

// resolve resolves a descriptor until it reaches a non-ref descriptor.
func (r referenceResolver) resolve(
	path fieldpath.Path,
	descriptor types.Type,
	depth int,
) (types.Type, error) {
	if descriptor.Code() != types.TypeRef {
		return descriptor, nil
	}

	if depth >= r.maxDepth {
		return types.Type{}, failure(
			path,
			FailureReferenceCycle,
			"maximum TypeRef ListMap key extraction depth reached",
		)
	}

	referenceView, ok := descriptor.Ref()
	if !ok {
		return types.Type{}, failure(
			path,
			FailureInvalidDescriptor,
			"descriptor is not a reference",
		)
	}

	referenceName := referenceView.Name()
	if r.typeResolver == nil {
		return types.Type{}, failure(
			path,
			FailureUnresolvedRef,
			fmt.Sprintf("reference %q has no resolver", referenceName),
		)
	}

	if r.activeRefs[referenceName] {
		return types.Type{}, failure(
			path,
			FailureReferenceCycle,
			fmt.Sprintf("reference %q is recursive", referenceName),
		)
	}

	typeDefinition, ok := r.typeResolver.ResolveType(referenceName)
	if !ok {
		return types.Type{}, failure(
			path,
			FailureUnresolvedRef,
			fmt.Sprintf("reference %q was not found", referenceName),
		)
	}

	r.activeRefs[referenceName] = true
	resolvedDescriptor, err := r.resolve(path, typeDefinition.Type(), depth+1)
	delete(r.activeRefs, referenceName)

	return resolvedDescriptor, err
}
