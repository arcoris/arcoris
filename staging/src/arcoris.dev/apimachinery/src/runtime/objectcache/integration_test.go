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

package objectcache

import (
	"context"
	"testing"
	"time"

	"arcoris.dev/apimachinery/runtime/objectmemorystore"
	"arcoris.dev/apimachinery/runtime/objectreflector"
	runtimewatch "arcoris.dev/apimachinery/runtime/objectstorewatch"
)

func TestCacheIntegrationWithReflectorAndObservableStore(t *testing.T) {
	backend, err := objectmemorystore.New()
	requireNoError(t, err)
	source, err := runtimewatch.New(backend)
	requireNoError(t, err)
	cache, err := New(testCollection(), WithHistory(HistoryPolicy{RetainedVersionsPerObject: 3}))
	requireNoError(t, err)
	reflector, err := objectreflector.New(source, testCollection(), cache)
	requireNoError(t, err)

	ctx, cancel := context.WithCancel(context.Background())
	done := make(chan error, 1)
	go func() {
		done <- reflector.Run(ctx)
	}()

	waitCtx, waitCancel := context.WithTimeout(context.Background(), time.Second)
	defer waitCancel()
	waitUntil(t, waitCtx, cache.Ready)

	key := testKey("system", 1)
	created, err := source.Create(context.Background(), key, testState(key, 0, "created"))
	requireNoError(t, err)
	waitUntil(t, waitCtx, func() bool {
		result, err := cache.Get(key)
		return err == nil && result.Found && result.State.Revision == created.Revision && desiredString(t, result.State) == "created"
	})

	updated, err := source.Update(context.Background(), key, created.Revision, testState(key, 0, "updated"))
	requireNoError(t, err)
	waitUntil(t, waitCtx, func() bool {
		result, err := cache.Get(key)
		return err == nil && result.Found && result.State.Revision == updated.Revision && desiredString(t, result.State) == "updated"
	})

	previous, err := cache.PreviousLive(key, updated.Revision)
	requireNoError(t, err)
	if !previous.Found || desiredString(t, previous.State) != "created" {
		t.Fatalf("PreviousLive() = %#v; want created version", previous)
	}

	deleted, err := source.Delete(context.Background(), key, updated.Revision)
	requireNoError(t, err)
	waitUntil(t, waitCtx, func() bool {
		result, err := cache.Get(key)
		return err == nil && result.Revision == deleted.Revision && !result.Found
	})
	historical, err := cache.GetAt(key, updated.Revision)
	requireNoError(t, err)
	if !historical.Found || desiredString(t, historical.State) != "updated" {
		t.Fatalf("GetAt(updated) = %#v; want updated version", historical)
	}
	previous, err = cache.PreviousLive(key, deleted.Revision)
	requireNoError(t, err)
	if !previous.Found || desiredString(t, previous.State) != "updated" {
		t.Fatalf("PreviousLive(delete) = %#v; want updated version", previous)
	}

	cancel()
	select {
	case err := <-done:
		requireErrorIs(t, err, context.Canceled)
	case <-time.After(time.Second):
		t.Fatalf("reflector did not stop after cancellation")
	}
}
