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
)

func TestMergeAtomicListExactSelectedReplacesWholeList(t *testing.T) {
	descriptor := types.ListOf(types.String()).Atomic().Type()
	base := list(str("old"))
	overlay := list(str("new"), str("next"))

	got, err := Merge(base, overlay, descriptor, pathSet(root()), Options{})
	if err != nil {
		t.Fatalf("Merge returned error: %v", err)
	}

	requireValue(t, got, overlay)
}

func TestMergeAtomicListDescendantSelectionUnsupported(t *testing.T) {
	descriptor := types.ListOf(types.String()).Atomic().Type()

	_, err := Merge(
		list(str("old")),
		list(str("new")),
		descriptor,
		pathSet(root().Index(0)),
		Options{},
	)

	requireErrorIs(t, err, ErrUnsupportedMerge)
}

func TestMergeSetListExactSelectedReplacesWholeList(t *testing.T) {
	descriptor := types.ListOf(types.String()).Set().Type()
	base := list(str("old"))
	overlay := list(str("new"), str("next"))

	got, err := Merge(base, overlay, descriptor, pathSet(root()), Options{})
	if err != nil {
		t.Fatalf("Merge returned error: %v", err)
	}

	requireValue(t, got, overlay)
}

func TestMergeSetListDescendantSelectionUnsupported(t *testing.T) {
	descriptor := types.ListOf(types.String()).Set().Type()

	_, err := Merge(
		list(str("old")),
		list(str("new")),
		descriptor,
		pathSet(root().Index(0)),
		Options{},
	)

	requireErrorIs(t, err, ErrUnsupportedMerge)
}
