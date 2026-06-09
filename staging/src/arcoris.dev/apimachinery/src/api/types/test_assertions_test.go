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

package types

import (
	"errors"
	"strings"
	"testing"
)

// requireNoError fails t when err is non-nil.
func requireNoError(t *testing.T, err error) {
	t.Helper()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

// requireErrorIs fails t when err is nil or not classified as target.
func requireErrorIs(t *testing.T, err, target error) {
	t.Helper()
	if !errors.Is(err, target) {
		t.Fatalf("expected error matching %v, got %v", target, err)
	}
}

// requireEqual fails t when got and want are not equal.
func requireEqual[T comparable](t *testing.T, got, want T) {
	t.Helper()
	if got != want {
		t.Fatalf("got %v, want %v", got, want)
	}
}

// requireCode fails t when desc does not use the expected DescriptorKind.
func requireCode(t *testing.T, desc Descriptor, want DescriptorKind) {
	t.Helper()
	requireEqual(t, desc.Code(), want)
}

// requireNullable fails t when desc has unexpected nullability.
func requireNullable(t *testing.T, desc Descriptor, want bool) {
	t.Helper()
	requireEqual(t, desc.Nullable(), want)
}

// requireValidDescriptor fails t when desc is not structurally valid.
func requireValidDescriptor(t *testing.T, desc Descriptor, resolver Resolver) {
	t.Helper()
	if resolver == nil {
		requireNoError(t, ValidateLocal(desc))
		return
	}
	requireNoError(t, ValidateResolved(desc, resolver))
}

// requireInvalidDescriptor fails t when desc is valid or has the wrong broad error.
func requireInvalidDescriptor(t *testing.T, desc Descriptor, resolver Resolver, target error) {
	t.Helper()
	if resolver == nil {
		requireErrorIs(t, ValidateLocal(desc), target)
		return
	}
	requireErrorIs(t, ValidateResolved(desc, resolver), target)
}

// validateTestDescriptor keeps tests terse while the public API stays explicit.
func validateTestDescriptor(desc Descriptor, resolver Resolver) error {
	if resolver == nil {
		return ValidateLocal(desc)
	}

	return ValidateResolved(desc, resolver)
}

// requireDescriptorError returns the structured DescriptorError diagnostics or fails t.
func requireDescriptorError(
	t *testing.T,
	err error,
	target error,
	path string,
	reason DescriptorErrorReason,
	detailContains string,
) *DescriptorError {
	t.Helper()
	requireErrorIs(t, err, target)

	var descriptorErr *DescriptorError
	if !errors.As(err, &descriptorErr) {
		t.Fatalf("expected DescriptorError, got %T", err)
	}

	requireEqual(t, descriptorErr.Path, path)
	requireEqual(t, descriptorErr.Reason, reason)

	if detailContains != "" && !strings.Contains(descriptorErr.Detail, detailContains) {
		t.Fatalf("expected detail containing %q, got %q", detailContains, descriptorErr.Detail)
	}

	return descriptorErr
}

// requireStringView returns the exact DescriptorString view or fails t.
func requireStringView(t *testing.T, desc Descriptor) StringView {
	t.Helper()
	view, ok := desc.AsString()
	requireEqual(t, ok, true)
	return view
}

// requireBytesView returns the exact DescriptorBytes view or fails t.
func requireBytesView(t *testing.T, desc Descriptor) BytesView {
	t.Helper()
	view, ok := desc.AsBytes()
	requireEqual(t, ok, true)
	return view
}

// requireInt8View returns the exact DescriptorInt8 view or fails t.
func requireInt8View(t *testing.T, desc Descriptor) Int8View {
	t.Helper()
	view, ok := desc.AsInt8()
	requireEqual(t, ok, true)
	return view
}

// requireInt16View returns the exact DescriptorInt16 view or fails t.
func requireInt16View(t *testing.T, desc Descriptor) Int16View {
	t.Helper()
	view, ok := desc.AsInt16()
	requireEqual(t, ok, true)
	return view
}

// requireInt32View returns the exact DescriptorInt32 view or fails t.
func requireInt32View(t *testing.T, desc Descriptor) Int32View {
	t.Helper()
	view, ok := desc.AsInt32()
	requireEqual(t, ok, true)
	return view
}

// requireInt64View returns the exact DescriptorInt64 view or fails t.
func requireInt64View(t *testing.T, desc Descriptor) Int64View {
	t.Helper()
	view, ok := desc.AsInt64()
	requireEqual(t, ok, true)
	return view
}

// requireUint8View returns the exact DescriptorUint8 view or fails t.
func requireUint8View(t *testing.T, desc Descriptor) Uint8View {
	t.Helper()
	view, ok := desc.AsUint8()
	requireEqual(t, ok, true)
	return view
}

// requireUint16View returns the exact DescriptorUint16 view or fails t.
func requireUint16View(t *testing.T, desc Descriptor) Uint16View {
	t.Helper()
	view, ok := desc.AsUint16()
	requireEqual(t, ok, true)
	return view
}

// requireUint32View returns the exact DescriptorUint32 view or fails t.
func requireUint32View(t *testing.T, desc Descriptor) Uint32View {
	t.Helper()
	view, ok := desc.AsUint32()
	requireEqual(t, ok, true)
	return view
}

// requireUint64View returns the exact DescriptorUint64 view or fails t.
func requireUint64View(t *testing.T, desc Descriptor) Uint64View {
	t.Helper()
	view, ok := desc.AsUint64()
	requireEqual(t, ok, true)
	return view
}

// requireFloat32View returns the exact DescriptorFloat32 view or fails t.
func requireFloat32View(t *testing.T, desc Descriptor) Float32View {
	t.Helper()
	view, ok := desc.AsFloat32()
	requireEqual(t, ok, true)
	return view
}

// requireFloat64View returns the exact DescriptorFloat64 view or fails t.
func requireFloat64View(t *testing.T, desc Descriptor) Float64View {
	t.Helper()
	view, ok := desc.AsFloat64()
	requireEqual(t, ok, true)
	return view
}

// requireDecimalView returns the exact DescriptorDecimal view or fails t.
func requireDecimalView(t *testing.T, desc Descriptor) DecimalView {
	t.Helper()
	view, ok := desc.AsDecimal()
	requireEqual(t, ok, true)
	return view
}

// requireTimestampView returns the exact DescriptorTimestamp view or fails t.
func requireTimestampView(t *testing.T, desc Descriptor) TimestampView {
	t.Helper()
	view, ok := desc.AsTimestamp()
	requireEqual(t, ok, true)
	return view
}

// requireDateView returns the exact DescriptorDate view or fails t.
func requireDateView(t *testing.T, desc Descriptor) DateView {
	t.Helper()
	view, ok := desc.AsDate()
	requireEqual(t, ok, true)
	return view
}

// requireTimeView returns the exact DescriptorTime view or fails t.
func requireTimeView(t *testing.T, desc Descriptor) TimeView {
	t.Helper()
	view, ok := desc.AsTime()
	requireEqual(t, ok, true)
	return view
}

// requireDurationView returns the exact DescriptorDuration view or fails t.
func requireDurationView(t *testing.T, desc Descriptor) DurationView {
	t.Helper()
	view, ok := desc.AsDuration()
	requireEqual(t, ok, true)
	return view
}

// requireObjectView returns the exact DescriptorObject view or fails t.
func requireObjectView(t *testing.T, desc Descriptor) ObjectView {
	t.Helper()
	view, ok := desc.AsObject()
	requireEqual(t, ok, true)
	return view
}

// requireListView returns the exact DescriptorList view or fails t.
func requireListView(t *testing.T, desc Descriptor) ListView {
	t.Helper()
	view, ok := desc.AsList()
	requireEqual(t, ok, true)
	return view
}

// requireMapView returns the exact DescriptorMap view or fails t.
func requireMapView(t *testing.T, desc Descriptor) MapView {
	t.Helper()
	view, ok := desc.AsMap()
	requireEqual(t, ok, true)
	return view
}

// requireRefView returns the exact DescriptorRef view or fails t.
func requireRefView(t *testing.T, desc Descriptor) RefView {
	t.Helper()
	view, ok := desc.AsRef()
	requireEqual(t, ok, true)
	return view
}

// objectTypeForField creates a DescriptorObject around a finalized field descriptor.
func objectTypeForField(field FieldDescriptor) Descriptor {
	desc := Descriptor{code: DescriptorObject}
	desc.object.unknown = UnknownReject
	desc.object.fields = []FieldDescriptor{field}
	return desc
}
