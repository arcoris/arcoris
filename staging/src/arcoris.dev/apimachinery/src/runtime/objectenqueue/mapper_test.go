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
	"fmt"
	"sync"
	"testing"

	apiidentity "arcoris.dev/apimachinery/api/identity"
	"arcoris.dev/apimachinery/api/meta"
	metaidentity "arcoris.dev/apimachinery/api/meta/identity"
	"arcoris.dev/apimachinery/api/meta/labels"
	"arcoris.dev/apimachinery/api/object"
	"arcoris.dev/apimachinery/api/objectownership"
	"arcoris.dev/apimachinery/api/objectstore"
	"arcoris.dev/apimachinery/api/value"
	"arcoris.dev/apimachinery/runtime/objectworkqueue"
)

func TestMapperFuncNil(t *testing.T) {
	var mapper MapperFunc
	err := mapper.Map(createdChange(t, 1), func(objectworkqueue.Item) error {
		return nil
	})

	requireErrorIs(t, err, ErrNilMapper)
}

func TestMapperFuncDelegates(t *testing.T) {
	change := createdChange(t, 1)
	wantErr := errors.New("delegate failed")
	called := false
	mapper := MapperFunc(func(got objectstore.Change, emit EmitFunc) error {
		called = true
		requireChange(t, got, change)
		if emit == nil {
			t.Fatalf("emit is nil")
		}
		return wantErr
	})

	err := mapper.Map(change, func(objectworkqueue.Item) error { return nil })

	if err != wantErr {
		t.Fatalf("error = %v; want %v", err, wantErr)
	}
	if !called {
		t.Fatalf("mapper function was not called")
	}
}

func requireNoError(t testing.TB, err error) {
	t.Helper()

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func requireErrorIs(t testing.TB, err error, target error) {
	t.Helper()

	if !errors.Is(err, target) {
		t.Fatalf("error = %v; want errors.Is(%v)", err, target)
	}
}

func requireItem(t testing.TB, got objectworkqueue.Item, want objectworkqueue.Item) {
	t.Helper()

	if !got.Key.Equal(want.Key) {
		t.Fatalf("item = %s; want %s", got.Key, want.Key)
	}
}

func requireChange(t testing.TB, got objectstore.Change, want objectstore.Change) {
	t.Helper()

	if got.Kind != want.Kind || got.Revision != want.Revision || !got.Key.Equal(want.Key) {
		t.Fatalf("change identity = %#v; want %#v", got, want)
	}
}

type pointerMapper struct{}

func (*pointerMapper) Map(objectstore.Change, EmitFunc) error {
	return nil
}

type recordingMapper struct {
	mu      sync.Mutex
	changes []objectstore.Change
}

func (m *recordingMapper) Map(change objectstore.Change, _ EmitFunc) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.changes = append(m.changes, change)

	return nil
}

func (m *recordingMapper) callCount() int {
	m.mu.Lock()
	defer m.mu.Unlock()

	return len(m.changes)
}

func (m *recordingMapper) onlyChange(t testing.TB) objectstore.Change {
	t.Helper()

	m.mu.Lock()
	defer m.mu.Unlock()

	if len(m.changes) != 1 {
		t.Fatalf("changes = %d; want 1", len(m.changes))
	}

	return m.changes[0]
}

func createdChange(t testing.TB, id int) objectstore.Change {
	t.Helper()

	key := testKey(id)
	change, err := objectstore.NewCreatedChange(key, testState(key, 1, "created"))
	requireNoError(t, err)

	return change
}

func updatedChange(t testing.TB, id int) objectstore.Change {
	t.Helper()

	key := testKey(id)
	change, err := objectstore.NewUpdatedChange(
		key,
		testState(key, 1, "before"),
		testState(key, 2, "after"),
	)
	requireNoError(t, err)

	return change
}

func deletedChange(t testing.TB, id int) objectstore.Change {
	t.Helper()

	key := testKey(id)
	change, err := objectstore.NewDeletedChange(key, testState(key, 1, "deleted"), 2)
	requireNoError(t, err)

	return change
}

func testKey(id int) objectstore.Key {
	return objectstore.MustKey(testResource(), metaidentity.ObjectName{
		Namespace: metaidentity.Namespace("default"),
		Name:      metaidentity.Name(fmt.Sprintf("unit-%d", id)),
	})
}

func testResource() apiidentity.GroupVersionResource {
	return apiidentity.GroupVersionResource{
		Group:    apiidentity.Group("control.arcoris.dev"),
		Version:  apiidentity.Version("v1"),
		Resource: apiidentity.Resource("units"),
	}
}

func testState(key objectstore.Key, revision objectstore.Revision, desired string) objectstore.State {
	labelSet, err := labels.FromStrings(map[string]string{"env": desired})
	if err != nil {
		panic(err)
	}

	return objectstore.State{
		Object: object.NewObserved[value.Value, value.Value](
			meta.FromGroupVersionKind(apiidentity.GroupVersionKind{
				Group:   key.Resource.Group,
				Version: key.Resource.Version,
				Kind:    apiidentity.Kind("Unit"),
			}),
			meta.ObjectMeta{
				Name:      key.Object.Name,
				Namespace: key.Object.Namespace,
				Labels:    labelSet,
			},
			value.StringValue(desired),
			value.StringValue("observed-"+desired),
		),
		Ownership: objectownership.EmptyState(),
		Revision:  revision,
	}
}
