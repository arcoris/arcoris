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
	"testing"

	"arcoris.dev/apimachinery/api/fieldpath"
	"arcoris.dev/apimachinery/api/types"
	"arcoris.dev/apimachinery/api/value"
)

func TestMergeSelectedString(t *testing.T) {
	got, err := Merge(
		str("old"),
		str("new"),
		types.String().Descriptor(),
		pathSet(root()),
		Options{},
	)
	if err != nil {
		t.Fatalf("Merge returned error: %v", err)
	}

	requireValue(t, got, str("new"))
}

func TestMergeAtEmptyFieldsPreservesString(t *testing.T) {
	got, err := MergeAt(
		root().Field(testFieldName("name")),
		str("old"),
		str("new"),
		types.String().Descriptor(),
		fieldpath.EmptySet(),
		Options{},
	)
	if err != nil {
		t.Fatalf("MergeAt returned error: %v", err)
	}

	requireValue(t, got, str("old"))
}

func TestMergeSelectedBool(t *testing.T) {
	got, err := Merge(
		boolValue(false),
		boolValue(true),
		types.Bool().Descriptor(),
		pathSet(root()),
		Options{},
	)
	if err != nil {
		t.Fatalf("Merge returned error: %v", err)
	}

	requireValue(t, got, boolValue(true))
}

func TestMergeSelectedInteger(t *testing.T) {
	got, err := Merge(
		intValue(1),
		intValue(2),
		types.Int64().Descriptor(),
		pathSet(root()),
		Options{},
	)
	if err != nil {
		t.Fatalf("Merge returned error: %v", err)
	}

	requireValue(t, got, intValue(2))
}

func TestMergeSelectedNullReplacesValue(t *testing.T) {
	got, err := Merge(
		str("old"),
		value.NullValue(),
		types.String().Nullable().Descriptor(),
		pathSet(root()),
		Options{},
	)
	if err != nil {
		t.Fatalf("Merge returned error: %v", err)
	}
	if !got.IsNull() {
		t.Fatalf("merged value is not null")
	}
}

func TestMergeDescendantUnderScalarReturnsUnsupported(t *testing.T) {
	tests := []struct {
		name       string
		base       value.Value
		overlay    value.Value
		descriptor types.Descriptor
	}{
		{
			name:       "string",
			base:       str("old"),
			overlay:    str("new"),
			descriptor: types.String().Descriptor(),
		},
		{
			name:       "integer",
			base:       intValue(1),
			overlay:    intValue(2),
			descriptor: types.Int64().Descriptor(),
		},
		{
			name:       "bool",
			base:       boolValue(false),
			overlay:    boolValue(true),
			descriptor: types.Bool().Descriptor(),
		},
		{
			name:       "decimal",
			base:       decimalValue("1.0"),
			overlay:    decimalValue("2.0"),
			descriptor: types.Decimal().Descriptor(),
		},
		{
			name:       "null",
			base:       value.NullValue(),
			overlay:    value.NullValue(),
			descriptor: types.Null().Descriptor(),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := Merge(
				tt.base,
				tt.overlay,
				tt.descriptor,
				pathSet(root().Field(testFieldName("name"))),
				Options{},
			)

			requireErrorIs(t, err, ErrUnsupportedMerge)
		})
	}
}
