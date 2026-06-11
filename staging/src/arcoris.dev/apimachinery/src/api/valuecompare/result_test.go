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
	"testing"
)

func TestEmptyResult(t *testing.T) {
	got := EmptyResult()

	requireResult(t, got, nil, nil, nil)
}

func TestZeroResultIsValid(t *testing.T) {
	var got Result

	requireNoError(t, got.ValidateStructure())
	requireResult(t, got, nil, nil, nil)
}

func TestResultIsEmpty(t *testing.T) {
	empty := EmptyResult()
	if !empty.IsEmpty() {
		t.Fatalf("EmptyResult().IsEmpty() = false")
	}

	changed, err := empty.withModified(rootField("name"))
	requireNoError(t, err)
	if changed.IsEmpty() {
		t.Fatalf("changed result IsEmpty() = true")
	}
}

func TestResultChanged(t *testing.T) {
	name := rootField("name")
	result, err := EmptyResult().withModified(name)
	requireNoError(t, err)

	requireSet(t, "changed", result.Changed(), name)
}

func TestResultChangedIsUnion(t *testing.T) {
	added := fieldpath.MustSet(rootField("new"))
	removed := fieldpath.MustSet(rootField("old"))
	modified := fieldpath.MustSet(rootField("same"))
	result := MustResult(added, removed, modified)

	requireDisjointResult(t, result)
	requireSet(t, "changed", result.Changed(), rootField("new"), rootField("old"), rootField("same"))
}

func TestResultAccessorsReturnExpectedSets(t *testing.T) {
	added := fieldpath.MustSet(rootField("new"))
	removed := fieldpath.MustSet(rootField("old"))
	modified := fieldpath.MustSet(rootField("same"))
	result := MustResult(added, removed, modified)

	requireSet(t, "added", result.Added(), rootField("new"))
	requireSet(t, "removed", result.Removed(), rootField("old"))
	requireSet(t, "modified", result.Modified(), rootField("same"))
}

func TestUnionSetsReturnsNonEmptySide(t *testing.T) {
	set := fieldpath.MustSet(rootField("name"))

	requireSet(t, "left empty", unionSets(fieldpath.EmptySet(), set), rootField("name"))
	requireSet(t, "right empty", unionSets(set, fieldpath.EmptySet()), rootField("name"))
}
