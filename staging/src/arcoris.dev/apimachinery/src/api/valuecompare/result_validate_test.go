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
	"testing"

	"arcoris.dev/apimachinery/api/fieldpath"
)

func TestNewResultValidatesBuckets(t *testing.T) {
	added := fieldpath.MustSet(rootField("new"))
	removed := fieldpath.MustSet(rootField("old"))
	modified := fieldpath.MustSet(rootField("same"))

	got, err := NewResult(added, removed, modified)
	requireNoError(t, err)
	requireResult(t, got, paths(rootField("new")), paths(rootField("old")), paths(rootField("same")))
}

func TestNewResultRejectsOverlappingBuckets(t *testing.T) {
	_, err := NewResult(
		fieldpath.MustSet(rootField("shared")),
		fieldpath.MustSet(rootField("shared")),
		fieldpath.EmptySet(),
	)

	requireErrorIs(t, err, ErrInvalidResult)
	requireErrorReason(t, err, ErrorReasonOverlappingResultPath)
	requireErrorPath(t, err, "$.shared")
}

func TestMustResultPanicsOnInvalidResult(t *testing.T) {
	defer func() {
		if recover() == nil {
			t.Fatalf("MustResult() did not panic")
		}
	}()

	_ = MustResult(
		fieldpath.MustSet(rootField("shared")),
		fieldpath.MustSet(rootField("shared")),
		fieldpath.EmptySet(),
	)
}

func TestResultBucketsAreDisjoint(t *testing.T) {
	got, err := Compare(
		valueRecord(
			"same", "old",
			"removed", "gone",
		),
		valueRecord(
			"same", "new",
			"added", "here",
		),
		typesObject("same", "removed", "added"),
		Options{},
	)
	requireNoError(t, err)

	requireDisjointResult(t, got)
}
