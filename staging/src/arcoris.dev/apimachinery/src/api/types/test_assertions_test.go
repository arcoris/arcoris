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

// requireCode fails t when typ does not use the expected TypeCode.
func requireCode(t *testing.T, typ Type, want TypeCode) {
	t.Helper()
	requireEqual(t, typ.Code(), want)
}

// requireNullable fails t when typ has unexpected nullability.
func requireNullable(t *testing.T, typ Type, want bool) {
	t.Helper()
	requireEqual(t, typ.Nullable(), want)
}

// requireValidType fails t when typ is not structurally valid.
func requireValidType(t *testing.T, typ Type, resolver Resolver) {
	t.Helper()
	requireNoError(t, ValidateType(typ, resolver))
}

// requireInvalidType fails t when typ is valid or has the wrong broad error.
func requireInvalidType(t *testing.T, typ Type, resolver Resolver, target error) {
	t.Helper()
	requireErrorIs(t, ValidateType(typ, resolver), target)
}

// requireTypeError returns the structured TypeError diagnostics or fails t.
func requireTypeError(t *testing.T, err error, target error, path string, reason TypeErrorReason, detailContains string) *TypeError {
	t.Helper()
	requireErrorIs(t, err, target)
	var typeErr *TypeError
	if !errors.As(err, &typeErr) {
		t.Fatalf("expected TypeError, got %T", err)
	}
	requireEqual(t, typeErr.Path, path)
	requireEqual(t, typeErr.Reason, reason)
	if detailContains != "" && !strings.Contains(typeErr.Detail, detailContains) {
		t.Fatalf("expected detail containing %q, got %q", detailContains, typeErr.Detail)
	}
	return typeErr
}

// requireStringView returns the exact TypeString view or fails t.
func requireStringView(t *testing.T, typ Type) StringView {
	t.Helper()
	view, ok := typ.String()
	requireEqual(t, ok, true)
	return view
}

// requireBytesView returns the exact TypeBytes view or fails t.
func requireBytesView(t *testing.T, typ Type) BytesView {
	t.Helper()
	view, ok := typ.Bytes()
	requireEqual(t, ok, true)
	return view
}

// requireInt8View returns the exact TypeInt8 view or fails t.
func requireInt8View(t *testing.T, typ Type) Int8View {
	t.Helper()
	view, ok := typ.Int8()
	requireEqual(t, ok, true)
	return view
}

// requireInt16View returns the exact TypeInt16 view or fails t.
func requireInt16View(t *testing.T, typ Type) Int16View {
	t.Helper()
	view, ok := typ.Int16()
	requireEqual(t, ok, true)
	return view
}

// requireInt32View returns the exact TypeInt32 view or fails t.
func requireInt32View(t *testing.T, typ Type) Int32View {
	t.Helper()
	view, ok := typ.Int32()
	requireEqual(t, ok, true)
	return view
}

// requireInt64View returns the exact TypeInt64 view or fails t.
func requireInt64View(t *testing.T, typ Type) Int64View {
	t.Helper()
	view, ok := typ.Int64()
	requireEqual(t, ok, true)
	return view
}

// requireUint8View returns the exact TypeUint8 view or fails t.
func requireUint8View(t *testing.T, typ Type) Uint8View {
	t.Helper()
	view, ok := typ.Uint8()
	requireEqual(t, ok, true)
	return view
}

// requireUint16View returns the exact TypeUint16 view or fails t.
func requireUint16View(t *testing.T, typ Type) Uint16View {
	t.Helper()
	view, ok := typ.Uint16()
	requireEqual(t, ok, true)
	return view
}

// requireUint32View returns the exact TypeUint32 view or fails t.
func requireUint32View(t *testing.T, typ Type) Uint32View {
	t.Helper()
	view, ok := typ.Uint32()
	requireEqual(t, ok, true)
	return view
}

// requireUint64View returns the exact TypeUint64 view or fails t.
func requireUint64View(t *testing.T, typ Type) Uint64View {
	t.Helper()
	view, ok := typ.Uint64()
	requireEqual(t, ok, true)
	return view
}

// requireFloat32View returns the exact TypeFloat32 view or fails t.
func requireFloat32View(t *testing.T, typ Type) Float32View {
	t.Helper()
	view, ok := typ.Float32()
	requireEqual(t, ok, true)
	return view
}

// requireFloat64View returns the exact TypeFloat64 view or fails t.
func requireFloat64View(t *testing.T, typ Type) Float64View {
	t.Helper()
	view, ok := typ.Float64()
	requireEqual(t, ok, true)
	return view
}

// requireDecimalView returns the exact TypeDecimal view or fails t.
func requireDecimalView(t *testing.T, typ Type) DecimalView {
	t.Helper()
	view, ok := typ.Decimal()
	requireEqual(t, ok, true)
	return view
}

// requireTimestampView returns the exact TypeTimestamp view or fails t.
func requireTimestampView(t *testing.T, typ Type) TimestampView {
	t.Helper()
	view, ok := typ.Timestamp()
	requireEqual(t, ok, true)
	return view
}

// requireDateView returns the exact TypeDate view or fails t.
func requireDateView(t *testing.T, typ Type) DateView {
	t.Helper()
	view, ok := typ.Date()
	requireEqual(t, ok, true)
	return view
}

// requireTimeView returns the exact TypeTime view or fails t.
func requireTimeView(t *testing.T, typ Type) TimeView {
	t.Helper()
	view, ok := typ.Time()
	requireEqual(t, ok, true)
	return view
}

// requireDurationView returns the exact TypeDuration view or fails t.
func requireDurationView(t *testing.T, typ Type) DurationView {
	t.Helper()
	view, ok := typ.Duration()
	requireEqual(t, ok, true)
	return view
}

// requireObjectView returns the exact TypeObject view or fails t.
func requireObjectView(t *testing.T, typ Type) ObjectView {
	t.Helper()
	view, ok := typ.Object()
	requireEqual(t, ok, true)
	return view
}

// requireListView returns the exact TypeList view or fails t.
func requireListView(t *testing.T, typ Type) ListView {
	t.Helper()
	view, ok := typ.List()
	requireEqual(t, ok, true)
	return view
}

// requireMapView returns the exact TypeMap view or fails t.
func requireMapView(t *testing.T, typ Type) MapView {
	t.Helper()
	view, ok := typ.Map()
	requireEqual(t, ok, true)
	return view
}

// requireRefView returns the exact TypeRef view or fails t.
func requireRefView(t *testing.T, typ Type) RefView {
	t.Helper()
	view, ok := typ.Ref()
	requireEqual(t, ok, true)
	return view
}

// objectTypeForField creates a TypeObject around a finalized field descriptor.
func objectTypeForField(field FieldDescriptor) Type {
	typ := Type{code: TypeObject}
	typ.object.unknown = UnknownReject
	typ.object.fields = []FieldDescriptor{field}
	return typ
}
