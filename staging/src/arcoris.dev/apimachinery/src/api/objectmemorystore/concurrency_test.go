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

package objectmemorystore

import (
	"context"
	"errors"
	"sync"
	"testing"

	"arcoris.dev/apimachinery/api/objectstore"
)

func TestConcurrentCreateSameKeyAllowsOneWinner(t *testing.T) {
	store := testStore(t)
	key := testKey(1)
	errs := runConcurrent(32, func(i int) error {
		_, err := store.Create(context.Background(), key, testState("created"))
		return err
	})

	successes := 0
	for _, err := range errs {
		switch {
		case err == nil:
			successes++
		case errors.Is(err, objectstore.ErrAlreadyExists):
		default:
			t.Fatalf("unexpected error: %v", err)
		}
	}
	if successes != 1 {
		t.Fatalf("successes = %d; want 1", successes)
	}
}

func TestConcurrentCreateDifferentKeysAllSucceed(t *testing.T) {
	store := testStore(t)
	errs := runConcurrent(64, func(i int) error {
		_, err := store.Create(context.Background(), testKey(i), testState("created"))
		return err
	})

	for _, err := range errs {
		requireNoError(t, err)
	}
}

func TestConcurrentUpdateSameKeySameRevisionAllowsOneWinner(t *testing.T) {
	store := testStore(t)
	key := testKey(1)
	created := createState(t, store, key, "created")
	errs := runConcurrent(32, func(i int) error {
		_, err := store.Update(context.Background(), key, created.Revision, testState("updated"))
		return err
	})

	successes := 0
	for _, err := range errs {
		switch {
		case err == nil:
			successes++
		case errors.Is(err, objectstore.ErrConflict), errors.Is(err, objectstore.ErrStaleRevision):
		default:
			t.Fatalf("unexpected error: %v", err)
		}
	}
	if successes != 1 {
		t.Fatalf("successes = %d; want 1", successes)
	}
}

func TestConcurrentUpdateDifferentKeysAllSucceed(t *testing.T) {
	store := testStore(t)
	const workers = 64
	keys := make([]objectstore.Key, workers)
	states := make([]objectstore.State, workers)
	for i := 0; i < workers; i++ {
		keys[i] = testKey(i)
		states[i] = createState(t, store, keys[i], "created")
	}

	errs := runConcurrent(workers, func(i int) error {
		updated, err := store.Update(context.Background(), keys[i], states[i].Revision, testState("updated"))
		if err == nil {
			states[i] = updated
		}
		return err
	})

	for _, err := range errs {
		requireNoError(t, err)
	}
	for i := 0; i < workers; i++ {
		got, ok, err := store.Get(context.Background(), keys[i])
		requireNoError(t, err)
		if !ok {
			t.Fatalf("key %d missing after update", i)
		}
		requireDesiredString(t, got, "updated")
	}
}

func TestConcurrentGetDuringUpdateIsRaceFree(t *testing.T) {
	store := testStore(t)
	key := testKey(1)
	current := createState(t, store, key, "created")

	var mu sync.Mutex
	for i := 0; i < 32; i++ {
		next, err := store.Update(context.Background(), key, current.Revision, testState("updated"))
		requireNoError(t, err)
		current = next
	}

	errs := runConcurrent(64, func(i int) error {
		if i%2 == 0 {
			mu.Lock()
			expected := current.Revision
			mu.Unlock()
			next, err := store.Update(context.Background(), key, expected, testState("again"))
			if err == nil {
				mu.Lock()
				current = next
				mu.Unlock()
			}
			return err
		}
		_, _, err := store.Get(context.Background(), key)
		return err
	})

	for _, err := range errs {
		if err != nil && !errors.Is(err, objectstore.ErrConflict) && !errors.Is(err, objectstore.ErrStaleRevision) {
			t.Fatalf("unexpected error: %v", err)
		}
	}
}

func TestRevisionMonotonicUnderConcurrentWrites(t *testing.T) {
	store := testStore(t)
	const workers = 128
	created := make([]objectstore.State, workers)
	errs := runConcurrent(workers, func(i int) error {
		state, err := store.Create(context.Background(), testKey(i), testState("created"))
		if err == nil {
			created[i] = state
		}
		return err
	})

	seen := make(map[objectstore.Revision]struct{}, workers)
	for i, err := range errs {
		requireNoError(t, err)
		if !created[i].Revision.IsValid() {
			t.Fatalf("revision %d is invalid", i)
		}
		if _, exists := seen[created[i].Revision]; exists {
			t.Fatalf("duplicate revision %v", created[i].Revision)
		}
		seen[created[i].Revision] = struct{}{}
	}
}

func TestConcurrentDeleteAndUpdateSameKeyHasOneTerminalWinner(t *testing.T) {
	store := testStore(t)
	key := testKey(1)
	created := createState(t, store, key, "created")
	errs := runConcurrent(2, func(i int) error {
		if i == 0 {
			_, err := store.Delete(context.Background(), key, created.Revision)
			return err
		}
		_, err := store.Update(context.Background(), key, created.Revision, testState("updated"))
		return err
	})

	successes := 0
	for _, err := range errs {
		switch {
		case err == nil:
			successes++
		case errors.Is(err, objectstore.ErrConflict), errors.Is(err, objectstore.ErrNotFound), errors.Is(err, objectstore.ErrStaleRevision):
		default:
			t.Fatalf("unexpected error: %v", err)
		}
	}
	if successes != 1 {
		t.Fatalf("successes = %d; want 1", successes)
	}

	deleteErr := errs[0]
	updateErr := errs[1]
	got, ok, err := store.Get(context.Background(), key)
	requireNoError(t, err)
	switch {
	case deleteErr == nil:
		if ok {
			t.Fatalf("delete won but object is still visible")
		}
		_, err := store.Update(context.Background(), key, created.Revision, testState("after-delete"))
		requireErrorIs(t, err, objectstore.ErrNotFound)
	case updateErr == nil:
		if !ok {
			t.Fatalf("update won but object is tombstoned")
		}
		requireDesiredString(t, got, "updated")
		_, err := store.Delete(context.Background(), key, created.Revision)
		requireErrorIs(t, err, objectstore.ErrStaleRevision)
	default:
		t.Fatalf("no terminal winner: delete=%v update=%v", deleteErr, updateErr)
	}
}

func TestConcurrentDeleteSameKeyAllowsOneWinner(t *testing.T) {
	store := testStore(t)
	key := testKey(1)
	created := createState(t, store, key, "created")
	errs := runConcurrent(32, func(i int) error {
		_, err := store.Delete(context.Background(), key, created.Revision)
		return err
	})

	successes := 0
	for _, err := range errs {
		switch {
		case err == nil:
			successes++
		case errors.Is(err, objectstore.ErrConflict), errors.Is(err, objectstore.ErrNotFound), errors.Is(err, objectstore.ErrStaleRevision):
		default:
			t.Fatalf("unexpected error: %v", err)
		}
	}
	if successes != 1 {
		t.Fatalf("successes = %d; want 1", successes)
	}
}

func runConcurrent(n int, fn func(int) error) []error {
	errs := make([]error, n)
	var wg sync.WaitGroup
	wg.Add(n)
	for i := 0; i < n; i++ {
		i := i
		go func() {
			defer wg.Done()
			errs[i] = fn(i)
		}()
	}
	wg.Wait()

	return errs
}
