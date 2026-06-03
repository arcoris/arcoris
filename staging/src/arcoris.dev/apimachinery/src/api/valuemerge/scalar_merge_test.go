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

	"arcoris.dev/apimachinery/api/types"
	"arcoris.dev/apimachinery/api/value"
)

func TestMergeSelectedString(t *testing.T) {
	got, err := Merge(
		str("old"),
		str("new"),
		types.String().Type(),
		pathSet(root()),
		Options{},
	)
	if err != nil {
		t.Fatalf("Merge returned error: %v", err)
	}

	requireValue(t, got, str("new"))
}

func TestMergeUnselectedStringPreserved(t *testing.T) {
	got, err := MergeAt(
		root().Field("name"),
		str("old"),
		str("new"),
		types.String().Type(),
		pathSet(root().Field("other")),
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
		types.Bool().Type(),
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
		types.Int64().Type(),
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
		types.String().Nullable().Type(),
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
	_, err := Merge(
		str("old"),
		str("new"),
		types.String().Type(),
		pathSet(root().Field("name")),
		Options{},
	)

	requireErrorIs(t, err, ErrUnsupportedMerge)
}
