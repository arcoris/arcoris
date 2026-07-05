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

package objectindex

import (
	"errors"
	"testing"

	"arcoris.dev/apimachinery/api/objectstore"
)

func TestExtractorFuncNilReturnsErrNilExtractor(t *testing.T) {
	item := testItem(testKey("task-1"), 1, "worker-a")

	err := ExtractorFunc(nil).Extract(item, func(Value) error { return nil })

	requireErrorIs(t, err, ErrNilExtractor)
}

func TestExtractorFuncDelegates(t *testing.T) {
	item := testItem(testKey("task-1"), 1, "worker-a")
	called := false
	fn := ExtractorFunc(func(got objectstore.ListItem, emit EmitFunc) error {
		called = true
		if !got.Key.Equal(item.Key) {
			t.Fatalf("item key = %#v; want %#v", got.Key, item.Key)
		}

		return emit("worker-a")
	})

	requireNoError(t, fn.Extract(item, func(value Value) error {
		if value != "worker-a" {
			t.Fatalf("value = %q; want worker-a", value)
		}

		return nil
	}))
	if !called {
		t.Fatal("function was not called")
	}
}

func TestExtractorFuncReturnsDelegatedErrorUnchanged(t *testing.T) {
	item := testItem(testKey("task-1"), 1, "worker-a")
	extractorErr := errors.New("extract failed")
	fn := ExtractorFunc(func(objectstore.ListItem, EmitFunc) error {
		return extractorErr
	})

	err := fn.Extract(item, func(Value) error { return nil })

	requireErrorIs(t, err, extractorErr)
}

func TestExtractValuesRejectsEmptyValue(t *testing.T) {
	item := testItem(testKey("task-1"), 1, "worker-a")
	definitions := map[Name]Extractor{
		"worker": fixedExtractor(""),
	}

	_, err := extractValues([]Name{"worker"}, definitions, item)

	requireErrorIs(t, err, ErrInvalidIndex)
}
