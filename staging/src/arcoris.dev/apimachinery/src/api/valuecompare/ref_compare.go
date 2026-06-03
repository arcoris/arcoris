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

package valuecompare

import (
	"arcoris.dev/apimachinery/api/fieldpath"
	"arcoris.dev/apimachinery/api/types"
)

// compareRef resolves a TypeRef and compares at the same semantic path.
func (c *comparer) compareRef(
	path fieldpath.Path,
	oldOperand operand,
	newOperand operand,
	descriptor types.Type,
	depth int,
) (Result, error) {
	name, resolved, err := c.resolveRefDefinition(path, descriptor, depth)
	if err != nil {
		return Result{}, err
	}

	c.resolving[name] = true
	defer delete(c.resolving, name)

	return c.compare(path, oldOperand, newOperand, resolved, depth+1)
}

// resolveRefDefinition resolves one TypeRef edge and checks recursion guards.
func (c *comparer) resolveRefDefinition(
	path fieldpath.Path,
	descriptor types.Type,
	depth int,
) (types.TypeName, types.Type, error) {
	if depth >= c.maxDepth {
		return "", types.Type{}, errorAt(
			path,
			ErrReferenceCycle,
			ErrorReasonReferenceCycle,
			"maximum TypeRef comparison depth reached",
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
	if c.resolver == nil {
		return "", types.Type{}, errorfAt(
			path,
			ErrUnresolvedRef,
			ErrorReasonUnresolvedRef,
			"reference %q has no resolver",
			name,
		)
	}
	if c.resolving[name] {
		return "", types.Type{}, errorfAt(
			path,
			ErrReferenceCycle,
			ErrorReasonReferenceCycle,
			"reference %q is recursive",
			name,
		)
	}

	definition, ok := c.resolver.ResolveType(name)
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
