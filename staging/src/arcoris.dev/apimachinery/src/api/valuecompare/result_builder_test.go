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

func TestResultBuilderBuildsCanonicalBuckets(t *testing.T) {
	var builder resultBuilder
	builder.AddAddedSet(fieldpath.MustSet(rootField("b"), rootField("a")))
	builder.AddRemovedSet(fieldpath.MustSet(rootField("old")))
	builder.AddModified(rootField("same"))

	got, err := builder.Build()
	requireNoError(t, err)

	requireResult(
		t,
		got,
		paths(rootField("a"), rootField("b")),
		paths(rootField("old")),
		paths(rootField("same")),
	)
}

func TestResultBuilderAddsResult(t *testing.T) {
	child := MustResult(
		fieldpath.MustSet(rootField("new")),
		fieldpath.MustSet(rootField("old")),
		fieldpath.MustSet(rootField("same")),
	)
	var builder resultBuilder
	builder.AddResult(child)

	got, err := builder.Build()
	requireNoError(t, err)

	requireResult(t, got, paths(rootField("new")), paths(rootField("old")), paths(rootField("same")))
}
