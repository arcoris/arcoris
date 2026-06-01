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

// validator carries one validation run's resolver, recursion, and diagnostic state.
type validator struct {
	resolver  types.Resolver
	maxDepth  int
	maxErrors int

	resolving map[types.TypeName]bool
	errors    ErrorList
}

// newValidator normalizes options into an executable validation state.
func newValidator(opts Options) *validator {
	maxDepth := opts.MaxDepth
	if maxDepth <= 0 {
		maxDepth = defaultMaxDepth
	}

	maxErrors := opts.MaxErrors
	if maxErrors <= 0 {
		maxErrors = defaultMaxErrors
	}

	return &validator{
		resolver:  opts.Resolver,
		maxDepth:  maxDepth,
		maxErrors: maxErrors,
		resolving: make(map[types.TypeName]bool),
	}
}

// validate dispatches concrete value validation by descriptor type code.
func (v *validator) validate(
	path fieldpath.Path,
	val value.Value,
	descriptor types.Type,
	depth int,
) {
	if v.shouldStop() {
		return
	}

	if val.IsZero() {
		v.add(path, ErrInvalidValue, ErrorReasonInvalidZero, "value is the invalid zero Value")
		return
	}

	if !descriptor.IsValid() {
		v.add(
			path,
			ErrInvalidDescriptor,
			ErrorReasonInvalidDescriptor,
			"descriptor has no valid type code",
		)
		return
	}

	if val.IsNull() {
		// A TypeRef can resolve to a nullable target even when the reference
		// descriptor itself is not marked nullable. Resolve before applying
		// nullability so reusable semantic types keep their own value contract.
		if descriptor.Code() == types.TypeRef && !descriptor.Nullable() {
			v.validateRef(path, val, descriptor, depth)
			return
		}

		v.validateNull(path, descriptor)
		return
	}

	switch descriptor.Code() {
	case types.TypeNull:
		v.addKindMismatch(path, val.Kind(), value.KindNull, descriptor.Code())
	case types.TypeBool:
		v.validateBool(path, val, descriptor)
	case types.TypeString:
		v.validateString(path, val, descriptor)
	case types.TypeBytes:
		v.validateBytes(path, val, descriptor)
	case types.TypeInt8,
		types.TypeInt16,
		types.TypeInt32,
		types.TypeInt64:
		v.validateSignedInteger(path, val, descriptor)
	case types.TypeUint8,
		types.TypeUint16,
		types.TypeUint32,
		types.TypeUint64:
		v.validateUnsignedInteger(path, val, descriptor)
	case types.TypeFloat32:
		v.validateFloat32(path, val, descriptor)
	case types.TypeFloat64:
		v.validateFloat64(path, val, descriptor)
	case types.TypeDecimal:
		v.validateDecimal(path, val, descriptor)
	case types.TypeTimestamp,
		types.TypeDate,
		types.TypeTime,
		types.TypeDuration:
		v.validateTemporal(path, val, descriptor)
	case types.TypeObject:
		v.validateObject(path, val, descriptor, depth)
	case types.TypeMap:
		v.validateMap(path, val, descriptor, depth)
	case types.TypeList:
		v.validateList(path, val, descriptor, depth)
	case types.TypeRef:
		v.validateRef(path, val, descriptor, depth)
	default:
		v.add(
			path,
			ErrInvalidDescriptor,
			ErrorReasonInvalidDescriptor,
			"descriptor has an unsupported type code",
		)
	}
}

// add stores one direct diagnostic when collection has capacity.
func (v *validator) add(path fieldpath.Path, err error, reason ErrorReason, detail string) {
	if v.shouldStop() {
		return
	}

	v.errors = append(v.errors, errorAt(path, err, reason, detail))
}

// addf stores one formatted direct diagnostic when collection has capacity.
func (v *validator) addf(
	path fieldpath.Path,
	err error,
	reason ErrorReason,
	format string,
	args ...any,
) {
	if v.shouldStop() {
		return
	}

	v.errors = append(v.errors, errorfAt(path, err, reason, format, args...))
}

// wrap stores one diagnostic that preserves a lower-layer cause.
func (v *validator) wrap(
	path fieldpath.Path,
	err error,
	reason ErrorReason,
	detail string,
	cause error,
) {
	if v.shouldStop() {
		return
	}

	v.errors = append(v.errors, wrapAt(path, err, reason, detail, cause))
}

// shouldStop reports whether the configured diagnostic budget has been reached.
func (v *validator) shouldStop() bool {
	return len(v.errors) >= v.maxErrors
}

// result returns nil or the collected validation diagnostics.
func (v *validator) result() error {
	if len(v.errors) == 0 {
		return nil
	}

	return v.errors
}
