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

	"arcoris.dev/apimachinery/api/objectstore"
	"arcoris.dev/apimachinery/runtime/objectworkqueue"
)

func TestPredicateFuncNil(t *testing.T) {
	var predicate PredicateFunc
	ok, err := predicate.Match(createdChange(t, 1))

	if ok {
		t.Fatalf("ok = true; want false")
	}
	requireErrorIs(t, err, ErrNilPredicate)
}

func TestPredicateFuncDelegates(t *testing.T) {
	change := createdChange(t, 1)
	wantErr := errors.New("predicate failed")
	called := false
	predicate := PredicateFunc(func(got objectstore.Change) (bool, error) {
		called = true
		requireChange(t, got, change)
		return true, wantErr
	})

	ok, err := predicate.Match(change)

	if !ok {
		t.Fatalf("ok = false; want true")
	}
	if err != wantErr {
		t.Fatalf("error = %v; want %v", err, wantErr)
	}
	if !called {
		t.Fatalf("predicate function was not called")
	}
}

func TestFilterRejectsNilDependencies(t *testing.T) {
	tests := []struct {
		name      string
		predicate Predicate
		mapper    Mapper
		wantErr   error
	}{
		{name: "nil predicate", predicate: nil, mapper: ChangedObject(), wantErr: ErrNilPredicate},
		{name: "typed nil predicate", predicate: PredicateFunc(nil), mapper: ChangedObject(), wantErr: ErrNilPredicate},
		{name: "nil mapper", predicate: PredicateFunc(func(objectstore.Change) (bool, error) { return true, nil }), mapper: nil, wantErr: ErrNilMapper},
		{name: "typed nil mapper", predicate: PredicateFunc(func(objectstore.Change) (bool, error) { return true, nil }), mapper: MapperFunc(nil), wantErr: ErrNilMapper},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := Filter(tt.predicate, tt.mapper).Map(createdChange(t, 1), func(objectworkqueue.Item) error {
				return nil
			})
			requireErrorIs(t, err, tt.wantErr)
		})
	}
}

func TestFilterReturnsPredicateErrorUnchanged(t *testing.T) {
	wantErr := errors.New("predicate failed")
	mapperCalled := false
	mapper := MapperFunc(func(objectstore.Change, EmitFunc) error {
		mapperCalled = true
		return nil
	})
	filtered := Filter(PredicateFunc(func(objectstore.Change) (bool, error) {
		return false, wantErr
	}), mapper)

	err := filtered.Map(createdChange(t, 1), func(objectworkqueue.Item) error { return nil })

	if err != wantErr {
		t.Fatalf("error = %v; want %v", err, wantErr)
	}
	if mapperCalled {
		t.Fatalf("mapper was called")
	}
}

func TestFilterSkipsMapperWhenPredicateFalse(t *testing.T) {
	mapperCalled := false
	filtered := Filter(
		PredicateFunc(func(objectstore.Change) (bool, error) { return false, nil }),
		MapperFunc(func(objectstore.Change, EmitFunc) error {
			mapperCalled = true
			return nil
		}),
	)

	err := filtered.Map(createdChange(t, 1), func(objectworkqueue.Item) error { return nil })

	requireNoError(t, err)
	if mapperCalled {
		t.Fatalf("mapper was called")
	}
}

func TestFilterCallsMapperWhenPredicateTrue(t *testing.T) {
	mapperCalled := false
	filtered := Filter(
		PredicateFunc(func(objectstore.Change) (bool, error) { return true, nil }),
		MapperFunc(func(change objectstore.Change, emit EmitFunc) error {
			mapperCalled = true
			return emit(objectworkqueue.Item{Key: change.Key})
		}),
	)

	var got []objectworkqueue.Item
	err := filtered.Map(createdChange(t, 1), func(item objectworkqueue.Item) error {
		got = append(got, item)
		return nil
	})

	requireNoError(t, err)
	if !mapperCalled {
		t.Fatalf("mapper was not called")
	}
	if len(got) != 1 {
		t.Fatalf("items = %d; want 1", len(got))
	}
}

func TestFilterPropagatesMapperErrorUnchanged(t *testing.T) {
	wantErr := errors.New("mapper failed")
	filtered := Filter(
		PredicateFunc(func(objectstore.Change) (bool, error) { return true, nil }),
		MapperFunc(func(objectstore.Change, EmitFunc) error { return wantErr }),
	)

	err := filtered.Map(createdChange(t, 1), func(objectworkqueue.Item) error { return nil })

	if err != wantErr {
		t.Fatalf("error = %v; want %v", err, wantErr)
	}
}

func TestFilterPropagatesEmitErrorUnchanged(t *testing.T) {
	wantErr := errors.New("emit failed")
	filtered := Filter(
		PredicateFunc(func(objectstore.Change) (bool, error) { return true, nil }),
		ChangedObject(),
	)

	err := filtered.Map(createdChange(t, 1), func(objectworkqueue.Item) error { return wantErr })

	if err != wantErr {
		t.Fatalf("error = %v; want %v", err, wantErr)
	}
}
