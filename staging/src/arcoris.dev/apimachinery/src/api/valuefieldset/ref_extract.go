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
	"arcoris.dev/apimachinery/api/internal/typeref"
	"arcoris.dev/apimachinery/api/types"
	"arcoris.dev/apimachinery/api/value"
)

// extractRef resolves a DescriptorRef and extracts paths at the same semantic location.
func (e *extractor) extractRef(
	path fieldpath.Path,
	val value.Value,
	descriptor types.Descriptor,
	depth int,
) (fieldpath.Set, error) {
	name, resolved, err := e.resolveRefDefinition(path, descriptor, depth, e.refs)
	if err != nil {
		return fieldpath.Set{}, err
	}

	leave := e.refs.Enter(name)
	set, err := e.extract(path, val, resolved, depth+1)
	leave()

	return set, err
}

// resolveRefDescriptor resolves references for selector-descriptor inspection.
//
// It uses an isolated stack because selector extraction may inspect references
// before the actual list item is recursively traversed.
func (e *extractor) resolveRefDescriptor(
	path fieldpath.Path,
	descriptor types.Descriptor,
	depth int,
) (types.Descriptor, error) {
	resolved, err := typeref.New(e.resolver, e.maxDepth).ResolveFinal(path, descriptor, depth)
	if err != nil {
		return types.Descriptor{}, valuefieldsetRefError(err)
	}

	return resolved, nil
}

// resolveRefDefinition resolves the immediate DescriptorRef target.
func (e *extractor) resolveRefDefinition(
	path fieldpath.Path,
	descriptor types.Descriptor,
	depth int,
	references *typeref.Resolver,
) (types.TypeName, types.Descriptor, error) {
	name, resolved, err := references.Resolve(path, descriptor, depth)
	if err != nil {
		return "", types.Descriptor{}, valuefieldsetRefError(err)
	}

	return name, resolved, nil
}

// valuefieldsetRefError maps shared DescriptorRef traversal failures into this
// package's structured error model.
func valuefieldsetRefError(err error) error {
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
