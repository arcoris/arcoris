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

package valuefieldset

import (
	"testing"

	"arcoris.dev/apimachinery/api/types"
	"arcoris.dev/apimachinery/api/value"
)

func TestExtractOrderedListUsesIndexPaths(t *testing.T) {
	path := rootField("args")
	val := value.MustListValue(
		value.StringValue("a"),
		value.StringValue("b"),
	)
	listView, _ := val.List()
	element := types.String().Type()
	run := newExtractor(Options{})

	got, err := run.extractIndexedList(path, listView, element, 0)
	requireNoError(t, err)

	requireFieldSet(t, got, path.Index(0), path.Index(1))
}

func TestExtractOrderedListEmptyIncludesListPath(t *testing.T) {
	path := rootField("args")
	val := value.MustListValue()

	got, err := ExtractAt(path, val, types.ListOf(types.String()).Type(), Options{})
	requireNoError(t, err)

	requireFieldSet(t, got, path)
}

func TestExtractAtomicListIncludesOnlyListPath(t *testing.T) {
	path := rootField("args")
	val := value.MustListValue(
		value.StringValue("a"),
		value.StringValue("b"),
	)

	got, err := ExtractAt(
		path,
		val,
		types.ListOf(types.String()).Atomic().Type(),
		Options{},
	)
	requireNoError(t, err)

	requireFieldSet(t, got, path)
}

func TestExtractSetListIncludesOnlyListPath(t *testing.T) {
	path := rootField("args")
	val := value.MustListValue(
		value.StringValue("a"),
		value.StringValue("b"),
	)

	got, err := ExtractAt(
		path,
		val,
		types.ListOf(types.String()).Set().Type(),
		Options{},
	)
	requireNoError(t, err)

	requireFieldSet(t, got, path)
}
