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

package objectreconcilewrite

import (
	"context"
	"testing"

	"arcoris.dev/apimachinery/api/objectlifecycle"
	"arcoris.dev/apimachinery/api/value"
	"arcoris.dev/apimachinery/runtime/objectreconciler"
)

func TestUpdateObservedBuildsRequestFromCurrent(t *testing.T) {
	key := testKey("task-1")
	owner := testOwner()
	current := newCurrent(t, key, 7, "desired")

	req, err := current.UpdateObserved(value.StringValue("healthy"), owner)

	requireNoError(t, err)
	if req.Resource != key.Resource {
		t.Fatalf("resource = %#v; want %#v", req.Resource, key.Resource)
	}
	if req.Object != key.Object {
		t.Fatalf("object = %#v; want %#v", req.Object, key.Object)
	}
	if req.Expected != 7 {
		t.Fatalf("expected = %s; want 7", req.Expected)
	}
	if req.Owner != owner {
		t.Fatalf("owner = %q; want %q", req.Owner, owner)
	}
	got, ok := req.Observed.AsString()
	if !ok || got != "healthy" {
		t.Fatalf("observed = %q/%v; want healthy/true", got, ok)
	}
}

func TestUpdateObservedDetachesObservedValue(t *testing.T) {
	current := newCurrent(t, testKey("task-1"), 7, "desired")
	observed := value.BytesValue([]byte{1, 2, 3})

	req, err := current.UpdateObserved(observed, testOwner())
	requireNoError(t, err)
	bytes, ok := req.Observed.AsBytes()
	if !ok {
		t.Fatal("observed is not bytes")
	}
	bytes[0] = 9
	again, _ := req.Observed.AsBytes()
	if again[0] != 1 {
		t.Fatalf("observed was not detached: %#v", again)
	}
}

func TestUpdateObservedRejectsZeroCurrent(t *testing.T) {
	var current Current

	_, err := current.UpdateObserved(value.StringValue("healthy"), testOwner())

	requireErrorIs(t, err, ErrInvalidCurrent)
}

func TestCurrentUpdateObservedRequestCanBeInjectedIntoObservedUpdater(t *testing.T) {
	key := testKey("task-1")
	snapshot := testSnapshot(t, 100, testItem(key, 7, "desired"))
	writer := &recordingObservedUpdater{}
	current, found, err := FromSnapshot(objectreconciler.Request{Key: key}, snapshot)
	requireNoError(t, err)
	if !found {
		t.Fatal("found = false; want true")
	}
	req, err := current.UpdateObserved(value.StringValue("healthy"), testOwner())
	requireNoError(t, err)

	_, err = writer.UpdateObserved(context.Background(), req)

	requireNoError(t, err)
	if !writer.called {
		t.Fatal("writer was not called")
	}
	if writer.request.Expected != 7 {
		t.Fatalf("writer request expected = %s; want 7", writer.request.Expected)
	}
	if writer.request.Resource != key.Resource || writer.request.Object != key.Object {
		t.Fatalf("writer request identity = %#v/%#v; want %#v/%#v",
			writer.request.Resource,
			writer.request.Object,
			key.Resource,
			key.Object,
		)
	}
}

type recordingObservedUpdater struct {
	called  bool
	request objectlifecycle.UpdateObservedRequest
}

var _ objectlifecycle.ObservedUpdater = (*recordingObservedUpdater)(nil)

func (r *recordingObservedUpdater) UpdateObserved(
	_ context.Context,
	req objectlifecycle.UpdateObservedRequest,
) (objectlifecycle.Result, error) {
	r.called = true
	r.request = req

	return objectlifecycle.Result{}, nil
}
