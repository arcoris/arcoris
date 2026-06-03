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
	"arcoris.dev/apimachinery/api/internal/valuepresence"
	"arcoris.dev/apimachinery/api/types"
	"arcoris.dev/apimachinery/api/value"
)

// Compare reports semantic changes between two present payload values at "$".
//
// The descriptor is expected to have been validated before comparison. Compare
// performs only local traversal checks for blockers such as invalid zero values,
// unusable descriptor views, kind mismatches, unresolved TypeRef values, and
// invalid ListMap keys.
func Compare(
	oldValue value.Value,
	newValue value.Value,
	descriptor types.Type,
	opts Options,
) (Result, error) {
	return CompareAt(fieldpath.RootPath(), oldValue, newValue, descriptor, opts)
}

// CompareAt reports semantic changes between two present payload values at path.
//
// The supplied base path is preserved in every returned path and diagnostic.
// This lets callers compare nested payload surfaces without rewriting root-based
// results. Invalid base paths are reported as ErrInvalidPath.
func CompareAt(
	path fieldpath.Path,
	oldValue value.Value,
	newValue value.Value,
	descriptor types.Type,
	opts Options,
) (Result, error) {
	if err := path.Validate(); err != nil {
		return Result{}, wrapAt(
			path,
			ErrInvalidPath,
			ErrorReasonInvalidPath,
			"base field path is invalid",
			err,
		)
	}

	run := newComparer(opts)
	return run.compare(
		path,
		valuepresence.Present(oldValue),
		valuepresence.Present(newValue),
		descriptor,
		0,
	)
}

// compare is the recursive descriptor dispatcher.
//
// Presence is handled before inspecting value kind so absent and explicit null
// remain distinct. Non-absent values then dispatch by descriptor code at the
// same semantic path.
func (c *comparer) compare(
	path fieldpath.Path,
	oldOperand valuepresence.Operand,
	newOperand valuepresence.Operand,
	descriptor types.Type,
	depth int,
) (Result, error) {
	if result, done, err := c.comparePresence(path, oldOperand, newOperand, descriptor); done {
		return result, err
	}

	oldValue := oldOperand.Value()
	newValue := newOperand.Value()
	if err := requireComparableInputs(path, oldValue, newValue, descriptor); err != nil {
		return Result{}, err
	}
	if oldValue.IsNull() || newValue.IsNull() {
		return c.compareNull(path, oldValue, newValue)
	}

	switch descriptor.Code() {
	case types.TypeNull:
		return c.compareNullDescriptor(path, oldValue, newValue, descriptor)
	case types.TypeBool,
		types.TypeString,
		types.TypeBytes,
		types.TypeInt8,
		types.TypeInt16,
		types.TypeInt32,
		types.TypeInt64,
		types.TypeUint8,
		types.TypeUint16,
		types.TypeUint32,
		types.TypeUint64,
		types.TypeFloat32,
		types.TypeFloat64,
		types.TypeDecimal,
		types.TypeTimestamp,
		types.TypeDate,
		types.TypeTime,
		types.TypeDuration:
		return c.compareScalar(path, oldValue, newValue, descriptor)
	case types.TypeObject:
		return c.compareObject(path, oldValue, newValue, descriptor, depth)
	case types.TypeMap:
		return c.compareMap(path, oldValue, newValue, descriptor, depth)
	case types.TypeList:
		return c.compareList(path, oldValue, newValue, descriptor, depth)
	case types.TypeRef:
		return c.compareRef(path, oldOperand, newOperand, descriptor, depth)
	default:
		return Result{}, errorAt(
			path,
			ErrInvalidDescriptor,
			ErrorReasonInvalidDescriptor,
			"descriptor has an unsupported type code",
		)
	}
}

// comparePresence resolves absent/present transitions before kind checks.
//
// Added and removed present values are expanded with valuefieldset so the result
// contains the same semantic subtree paths that ownership and apply layers will
// later use. Present null is not absence and flows through normal comparison.
func (c *comparer) comparePresence(
	path fieldpath.Path,
	oldOperand valuepresence.Operand,
	newOperand valuepresence.Operand,
	descriptor types.Type,
) (Result, bool, error) {
	switch {
	case oldOperand.Absent() && newOperand.Absent():
		return EmptyResult(), true, nil
	case oldOperand.Absent():
		result, err := c.addSubtree(path, newOperand.Value(), descriptor, EmptyResult())
		return result, true, err
	case newOperand.Absent():
		result, err := c.removeSubtree(path, oldOperand.Value(), descriptor, EmptyResult())
		return result, true, err
	default:
		return Result{}, false, nil
	}
}

// requireComparableInputs rejects blockers that prevent descriptor dispatch.
func requireComparableInputs(
	path fieldpath.Path,
	oldValue value.Value,
	newValue value.Value,
	descriptor types.Type,
) error {
	if oldValue.IsZero() || newValue.IsZero() {
		return errorAt(
			path,
			ErrInvalidValue,
			ErrorReasonInvalidZero,
			"value is the invalid zero Value",
		)
	}
	if !descriptor.IsValid() {
		return errorAt(
			path,
			ErrInvalidDescriptor,
			ErrorReasonInvalidDescriptor,
			"descriptor has no valid type code",
		)
	}

	return nil
}

// requireKind reports when a concrete payload kind cannot satisfy a descriptor.
func requireKind(path fieldpath.Path, val value.Value, expected value.Kind, code types.TypeCode) error {
	if val.Kind() == expected {
		return nil
	}

	return errorfAt(
		path,
		ErrKindMismatch,
		ErrorReasonKindMismatch,
		"value kind %s does not match descriptor %s; expected %s",
		val.Kind(),
		code,
		expected,
	)
}

// compareNull compares explicit null as present leaf data.
//
// This is used only after both sides are known present. Absent/null transitions
// are handled earlier by comparePresence.
func (c *comparer) compareNull(
	path fieldpath.Path,
	oldValue value.Value,
	newValue value.Value,
) (Result, error) {
	if oldValue.IsNull() && newValue.IsNull() {
		return EmptyResult(), nil
	}

	return EmptyResult().withModified(path)
}

// compareNullDescriptor verifies both present values satisfy TypeNull.
func (c *comparer) compareNullDescriptor(
	path fieldpath.Path,
	oldValue value.Value,
	newValue value.Value,
	descriptor types.Type,
) (Result, error) {
	if err := requireKind(path, oldValue, value.KindNull, descriptor.Code()); err != nil {
		return Result{}, err
	}
	if err := requireKind(path, newValue, value.KindNull, descriptor.Code()); err != nil {
		return Result{}, err
	}

	return EmptyResult(), nil
}
