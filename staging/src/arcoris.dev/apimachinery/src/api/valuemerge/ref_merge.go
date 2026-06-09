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

package valuemerge

import (
	"arcoris.dev/apimachinery/api/fieldpath"
	"arcoris.dev/apimachinery/api/internal/typeref"
	"arcoris.dev/apimachinery/api/types"
)

// mergeRef resolves one DescriptorRef edge without changing the semantic path.
func (m *merger) mergeRef(
	path fieldpath.Path,
	base operand,
	overlay operand,
	descriptor types.Descriptor,
	fields fieldpath.Set,
	depth int,
) (operand, error) {
	name, resolved, err := m.resolveRefDefinition(path, descriptor, depth)
	if err != nil {
		return operand{}, err
	}

	leave := m.refs.Enter(name)
	defer leave()

	return m.merge(path, base, overlay, resolved, fields, depth+1)
}

// resolveRefDefinition resolves one DescriptorRef edge and enforces recursion guards.
func (m *merger) resolveRefDefinition(
	path fieldpath.Path,
	descriptor types.Descriptor,
	depth int,
) (types.TypeName, types.Descriptor, error) {
	name, resolved, err := m.refs.Resolve(path, descriptor, depth)
	if err != nil {
		return "", types.Descriptor{}, valuemergeRefError(err)
	}

	return name, resolved, nil
}

// valuemergeRefError maps internal DescriptorRef traversal failures into this package.
func valuemergeRefError(err error) error {
	refError, ok := typeref.AsError(err)
	if !ok {
		return err
	}

	sentinel, reason := refErrorClassification(refError.Kind)
	return wrapAt(refError.Path, sentinel, reason, refError.Detail, err)
}

// refErrorClassification maps shared ref failure kinds to merge diagnostics.
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
