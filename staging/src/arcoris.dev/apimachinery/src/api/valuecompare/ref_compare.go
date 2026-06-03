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
	"arcoris.dev/apimachinery/api/internal/typeref"
	"arcoris.dev/apimachinery/api/types"
)

// compareRef resolves one TypeRef edge without changing the semantic path.
//
// Result paths describe payload locations, not descriptor definition paths.
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

	leave := c.refs.Enter(name)
	defer leave()

	return c.compare(path, oldOperand, newOperand, resolved, depth+1)
}

// resolveRefDefinition resolves one TypeRef edge and enforces recursion guards.
//
// The caller owns marking the returned name in c.refs while it descends
// into the resolved descriptor.
func (c *comparer) resolveRefDefinition(
	path fieldpath.Path,
	descriptor types.Type,
	depth int,
) (types.TypeName, types.Type, error) {
	name, resolved, err := c.refs.Resolve(path, descriptor, depth)
	if err != nil {
		return "", types.Type{}, valuecompareRefError(err)
	}

	return name, resolved, nil
}

// valuecompareRefError maps internal TypeRef traversal failures into this
// package's public sentinel and reason model.
func valuecompareRefError(err error) error {
	refError, ok := typeref.AsError(err)
	if !ok {
		return err
	}

	sentinel, reason := refErrorClassification(refError.Kind)
	return wrapAt(refError.Path, sentinel, reason, refError.Detail, err)
}

func refErrorClassification(kind typeref.FailureKind) (error, ErrorReason) {
	switch kind {
	case typeref.FailureInvalidDescriptor:
		return ErrInvalidDescriptor, ErrorReasonInvalidDescriptor
	case typeref.FailureUnresolvedRef:
		return ErrUnresolvedRef, ErrorReasonUnresolvedRef
	case typeref.FailureReferenceCycle:
		return ErrReferenceCycle, ErrorReasonReferenceCycle
	default:
		return ErrInvalidDescriptor, ErrorReasonInvalidDescriptor
	}
}
