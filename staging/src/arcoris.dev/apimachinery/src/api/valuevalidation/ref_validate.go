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
	"arcoris.dev/apimachinery/api/internal/typeref"
	"arcoris.dev/apimachinery/api/types"
	"arcoris.dev/apimachinery/api/value"
)

// validateRef resolves a DescriptorRef and validates the value at the same semantic path.
func (v *validator) validateRef(
	path fieldpath.Path,
	val value.Value,
	descriptor types.Descriptor,
	depth int,
) {
	name, resolved, err := v.refs.Resolve(path, descriptor, depth)
	if err != nil {
		v.addRefError(err)
		return
	}

	leave := v.refs.Enter(name)
	defer leave()

	v.validate(path, val, resolved, depth+1)
}

// addRefError maps shared DescriptorRef traversal failures into validation diagnostics.
func (v *validator) addRefError(err error) {
	refError, ok := typeref.AsError(err)
	if !ok {
		v.wrap(
			fieldpath.Root(),
			ErrInvalidDescriptor,
			ErrorReasonInvalidDescriptor,
			"DescriptorRef traversal failed",
			err,
		)
		return
	}

	sentinel, reason := refErrorClassification(refError.Kind)
	v.wrap(refError.Path, sentinel, reason, refError.Detail, err)
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
