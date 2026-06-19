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

package objectreflector

import (
	"context"
	"testing"

	"arcoris.dev/apimachinery/api/objectstore"
	storewatchapi "arcoris.dev/apimachinery/api/objectstorewatch"
	"arcoris.dev/apimachinery/runtime/objectmemorystore"
	runtimewatch "arcoris.dev/apimachinery/runtime/objectstorewatch"
)

func newTestReflector(t testing.TB, source storewatchapi.ListerWatcher, sink Sink) *Reflector {
	t.Helper()

	reflector, err := New(source, testCollection(), sink)
	requireNoError(t, err)

	return reflector
}

func TestReflectorIntegrationWithRuntimeStoreWatch(t *testing.T) {
	backend, err := objectmemorystore.New()
	requireNoError(t, err)
	source, err := runtimewatch.New(backend)
	requireNoError(t, err)

	sink := newRecordingSink(4)
	reflector, err := New(source, testCollection(), sink)
	requireNoError(t, err)

	ctx, cancel := context.WithCancel(context.Background())
	done := make(chan error, 1)
	go func() {
		done <- reflector.Run(ctx)
	}()

	waitRead(t, sink.replaceCh)

	key := testKey("system", 1)
	created, err := source.Create(context.Background(), key, testState(key, 0, "created"))
	requireNoError(t, err)
	requireChangeRevision(t, waitChange(t, sink.changeCh), objectstore.ChangeCreated, created.Revision)

	updated, err := source.Update(context.Background(), key, created.Revision, testState(key, 0, "updated"))
	requireNoError(t, err)
	requireChangeRevision(t, waitChange(t, sink.changeCh), objectstore.ChangeUpdated, updated.Revision)

	deleted, err := source.Delete(context.Background(), key, updated.Revision)
	requireNoError(t, err)
	requireChangeRevision(t, waitChange(t, sink.changeCh), objectstore.ChangeDeleted, deleted.Revision)

	cancel()
	requireErrorIs(t, <-done, context.Canceled)
}

func requireChangeRevision(
	t testing.TB,
	change objectstore.Change,
	kind objectstore.ChangeKind,
	revision objectstore.Revision,
) {
	t.Helper()

	if change.Kind != kind {
		t.Fatalf("change kind = %s; want %s", change.Kind, kind)
	}
	if change.Revision != revision {
		t.Fatalf("change revision = %s; want %s", change.Revision, revision)
	}
}
