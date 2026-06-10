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

package listmapkey

import (
	"errors"
	"testing"

	"arcoris.dev/apimachinery/api/fieldpath"
	"arcoris.dev/apimachinery/api/types"
	"arcoris.dev/apimachinery/api/value"
)

type resolverFunc func(types.TypeName) (types.Definition, bool)

func (f resolverFunc) Resolve(name types.TypeName) (types.Definition, bool) {
	return f(name)
}

func requireNoError(t *testing.T, err error) {
	t.Helper()

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func requireErrorKind(t *testing.T, err error, want FailureKind) {
	t.Helper()

	var extractionError *Error
	if !errors.As(err, &extractionError) {
		t.Fatalf("error type = %T, want *Error", err)
	}

	if extractionError.Kind != want {
		t.Fatalf("failure kind = %q, want %q", extractionError.Kind, want)
	}
}

func requireEqual[T comparable](t *testing.T, got, want T) {
	t.Helper()

	if got != want {
		t.Fatalf("got %v, want %v", got, want)
	}
}

func objectElement(fields ...types.FieldExpr) types.Descriptor {
	return types.Object(fields...).Descriptor()
}

func objectWith(name string, val value.Value) value.Value {
	return value.MustRecordValue(value.MustRecordMember(name, val))
}

func objectWithMembers(members ...value.RecordMember) value.Value {
	return value.MustRecordValue(members...)
}

func conditionPath(index int) fieldpath.Path {
	return fieldpath.Root().Field(testFieldName("conditions")).Index(index)
}

func testFieldName(name string) fieldpath.FieldName {
	return fieldpath.MustFieldName(name)
}
