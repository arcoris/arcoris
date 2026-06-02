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

package valuefieldset

import (
	"arcoris.dev/apimachinery/api/fieldpath"
	"arcoris.dev/apimachinery/api/types"
	"arcoris.dev/apimachinery/api/value"
)

// extractRef resolves a TypeRef and extracts paths at the same semantic location.
func (e *extractor) extractRef(
	path fieldpath.Path,
	val value.Value,
	descriptor types.Type,
	depth int,
) (fieldpath.Set, error) {
	name, resolved, err := e.resolveRefDefinition(path, descriptor, depth, e.resolving)
	if err != nil {
		return fieldpath.Set{}, err
	}

	e.resolving[name] = true
	set, err := e.extract(path, val, resolved, depth+1)
	delete(e.resolving, name)

	return set, err
}

// resolveRefDescriptor resolves references for selector-descriptor inspection.
//
// It uses an isolated stack because selector extraction may inspect references
// before the actual list item is recursively traversed.
func (e *extractor) resolveRefDescriptor(
	path fieldpath.Path,
	descriptor types.Type,
	depth int,
) (types.Type, error) {
	return e.resolveRefDescriptorWithStack(
		path,
		descriptor,
		depth,
		make(map[types.TypeName]bool),
	)
}

// resolveRefDescriptorWithStack resolves a TypeRef chain into a non-ref descriptor.
func (e *extractor) resolveRefDescriptorWithStack(
	path fieldpath.Path,
	descriptor types.Type,
	depth int,
	resolving map[types.TypeName]bool,
) (types.Type, error) {
	name, resolved, err := e.resolveRefDefinition(path, descriptor, depth, resolving)
	if err != nil {
		return types.Type{}, err
	}

	if resolved.Code() != types.TypeRef {
		return resolved, nil
	}

	resolving[name] = true
	target, err := e.resolveRefDescriptorWithStack(path, resolved, depth+1, resolving)
	delete(resolving, name)

	return target, err
}

// resolveRefDefinition resolves the immediate TypeRef target.
func (e *extractor) resolveRefDefinition(
	path fieldpath.Path,
	descriptor types.Type,
	depth int,
	resolving map[types.TypeName]bool,
) (types.TypeName, types.Type, error) {
	if depth >= e.maxDepth {
		return "", types.Type{}, errorAt(
			path,
			ErrReferenceCycle,
			ErrorReasonReferenceCycle,
			"maximum TypeRef extraction depth reached",
		)
	}

	view, ok := descriptor.Ref()
	if !ok {
		return "", types.Type{}, errorAt(
			path,
			ErrInvalidDescriptor,
			ErrorReasonInvalidDescriptor,
			"descriptor is not a reference",
		)
	}

	name := view.Name()
	if e.resolver == nil {
		return "", types.Type{}, errorfAt(
			path,
			ErrUnresolvedRef,
			ErrorReasonUnresolvedRef,
			"reference %q has no resolver",
			name,
		)
	}
	if resolving[name] {
		return "", types.Type{}, errorfAt(
			path,
			ErrReferenceCycle,
			ErrorReasonReferenceCycle,
			"reference %q is recursive",
			name,
		)
	}

	definition, ok := e.resolver.ResolveType(name)
	if !ok {
		return "", types.Type{}, errorfAt(
			path,
			ErrUnresolvedRef,
			ErrorReasonUnresolvedRef,
			"reference %q was not found",
			name,
		)
	}

	return name, definition.Type(), nil
}
