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

package objectenqueue

import (
	"errors"
	"testing"

	"arcoris.dev/apimachinery/api/objectquery"
	"arcoris.dev/apimachinery/api/objectstore"
	"arcoris.dev/apimachinery/runtime/objectworkqueue"
)

func TestFilterQueryZeroPredicateMatchesValidChanges(t *testing.T) {
	var got []objectworkqueue.Item
	err := FilterQuery(objectquery.Predicate{}, ChangedObject()).Map(createdChange(t, 1), func(item objectworkqueue.Item) error {
		got = append(got, item)
		return nil
	})

	requireNoError(t, err)
	if len(got) != 1 {
		t.Fatalf("items = %d; want 1", len(got))
	}
	requireItem(t, got[0], objectworkqueue.Item{Key: testKey(1)})
}

func TestFilterQueryProjectionSemantics(t *testing.T) {
	tests := []struct {
		name   string
		query  objectquery.Query
		change objectstore.Change
		want   bool
	}{
		{name: "ignored", query: mustLabelEquals(t, "env", "prod"), change: createdChange(t, 1), want: false},
		{name: "entered", query: mustLabelEquals(t, "env", "after"), change: updatedChange(t, 2), want: true},
		{name: "updated", query: objectquery.All(), change: updatedChange(t, 3), want: true},
		{name: "left", query: mustLabelEquals(t, "env", "before"), change: updatedChange(t, 4), want: true},
		{name: "deleted left", query: mustLabelEquals(t, "env", "deleted"), change: deletedChange(t, 5), want: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			predicate := mustPredicate(t, tt.query)
			mapperCalled := false
			mapper := MapperFunc(func(change objectstore.Change, emit EmitFunc) error {
				mapperCalled = true
				return emit(objectworkqueue.Item{Key: change.Key})
			})

			var got []objectworkqueue.Item
			err := FilterQuery(predicate, mapper).Map(tt.change, func(item objectworkqueue.Item) error {
				got = append(got, item)
				return nil
			})

			requireNoError(t, err)
			if mapperCalled != tt.want {
				t.Fatalf("mapperCalled = %v; want %v", mapperCalled, tt.want)
			}
			if len(got) != boolInt(tt.want) {
				t.Fatalf("items = %d; want %d", len(got), boolInt(tt.want))
			}
		})
	}
}

func TestFilterQueryReturnsProjectChangeError(t *testing.T) {
	mapperCalled := false
	mapper := MapperFunc(func(objectstore.Change, EmitFunc) error {
		mapperCalled = true
		return nil
	})

	err := FilterQuery(objectquery.Predicate{}, mapper).Map(objectstore.Change{}, func(objectworkqueue.Item) error {
		return nil
	})

	requireErrorIs(t, err, objectquery.ErrInvalidChange)
	if mapperCalled {
		t.Fatalf("mapper was called")
	}
}

func TestFilterQueryReturnsMapperErrorUnchanged(t *testing.T) {
	wantErr := errors.New("mapper failed")
	mapper := MapperFunc(func(objectstore.Change, EmitFunc) error {
		return wantErr
	})

	err := FilterQuery(objectquery.Predicate{}, mapper).Map(createdChange(t, 1), func(objectworkqueue.Item) error {
		return nil
	})

	if err != wantErr {
		t.Fatalf("error = %v; want %v", err, wantErr)
	}
}

func TestFilterQueryReturnsEmitErrorUnchanged(t *testing.T) {
	wantErr := errors.New("emit failed")

	err := FilterQuery(objectquery.Predicate{}, ChangedObject()).Map(createdChange(t, 1), func(objectworkqueue.Item) error {
		return wantErr
	})

	if err != wantErr {
		t.Fatalf("error = %v; want %v", err, wantErr)
	}
}

func TestFilterQueryRejectsNilMapper(t *testing.T) {
	tests := []struct {
		name   string
		mapper Mapper
	}{
		{name: "nil", mapper: nil},
		{name: "typed nil func", mapper: MapperFunc(nil)},
		{name: "typed nil pointer", mapper: (*pointerMapper)(nil)},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := FilterQuery(objectquery.Predicate{}, tt.mapper).Map(createdChange(t, 1), func(objectworkqueue.Item) error {
				return nil
			})
			requireErrorIs(t, err, ErrNilMapper)
		})
	}
}

func mustLabelEquals(t testing.TB, key string, value string) objectquery.Query {
	t.Helper()

	query, err := objectquery.LabelEquals(key, value)
	requireNoError(t, err)

	return query
}

func mustPredicate(t testing.TB, query objectquery.Query) objectquery.Predicate {
	t.Helper()

	predicate, err := objectquery.Compile(query)
	requireNoError(t, err)

	return predicate
}

func boolInt(v bool) int {
	if v {
		return 1
	}

	return 0
}
