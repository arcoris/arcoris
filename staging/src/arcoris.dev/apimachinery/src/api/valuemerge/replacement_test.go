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

	"arcoris.dev/apimachinery/api/internal/valuepresence"
	"arcoris.dev/apimachinery/api/types"
)

func TestReplaceSubtreeCopiesOverlay(t *testing.T) {
	got, err := newMerger(Options{}).replaceSubtree(
		root(),
		valuepresence.Present(str("new")),
		types.String().Descriptor(),
		0,
	)
	if err != nil {
		t.Fatalf("replaceSubtree returned error: %v", err)
	}

	requireValue(t, got.Value(), str("new"))
}

func TestReplaceSubtreeAbsentRemoves(t *testing.T) {
	got, err := newMerger(Options{}).replaceSubtree(
		root(),
		valuepresence.Absent(),
		types.String().Descriptor(),
		0,
	)
	if err != nil {
		t.Fatalf("replaceSubtree returned error: %v", err)
	}
	if got.Present() {
		t.Fatalf("replacement is present")
	}
}

func TestMergeExactSelectedReplacementChecksKind(t *testing.T) {
	_, err := Merge(
		obj(),
		str("invalid"),
		types.Object().Descriptor(),
		pathSet(root()),
		Options{},
	)

	requireErrorIs(t, err, ErrKindMismatch)
}
