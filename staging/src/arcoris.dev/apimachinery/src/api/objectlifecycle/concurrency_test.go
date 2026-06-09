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

package objectlifecycle

import (
	"context"
	"errors"
	"sync"
	"testing"

	"arcoris.dev/apimachinery/api/objectstore"
	"arcoris.dev/apimachinery/api/valuevalidation"
)

func TestConcurrentApplyDistinctObjects(t *testing.T) {
	executor := testExecutor(t)
	errs := runConcurrent(32, func(i int) error {
		_, err := executor.Apply(
			context.Background(),
			ApplyRequest{Object: testObject(i+1, "api:v1"), Owner: owner("creator")},
		)
		return err
	})

	for _, err := range errs {
		requireNoError(t, err)
	}
}

func TestConcurrentApplyAndDeleteSameObjectLeavesValidState(t *testing.T) {
	executor := testExecutor(t)
	created := createObject(t, executor, 1, "api:v1", owner("creator"))
	errs := runConcurrent(2, func(i int) error {
		if i == 0 {
			_, err := executor.Apply(
				context.Background(),
				ApplyRequest{Object: testObject(1, "api:v2"), Owner: owner("creator"), Force: true},
			)
			return err
		}

		_, err := executor.Delete(
			context.Background(),
			DeleteRequest{Resource: testGVR(), Object: testName(1), Expected: created.State.Revision},
		)
		return err
	})

	for _, err := range errs {
		if err == nil {
			continue
		}
		if !errors.Is(err, ErrConflict) &&
			!errors.Is(err, ErrNotFound) &&
			!errors.Is(err, ErrStaleRevision) &&
			!errors.Is(err, ErrAlreadyExists) {
			t.Fatalf("unexpected error: %v", err)
		}
	}

	result, err := executor.Get(context.Background(), GetRequest{Resource: testGVR(), Object: testName(1)})
	if err == nil {
		if !result.State.IsValid() {
			t.Fatalf("final state is invalid: %#v", result.State)
		}
		return
	}
	requireLifecycleError(t, err, ErrNotFound, ReasonNotFound)
}

func TestConcurrentApplySameLiveRevisionAllowsOneWinner(t *testing.T) {
	store := &oneWinnerUpdateStore{state: committedStateForFakeStore()}
	executor, err := NewExecutor(
		WithStore(store),
		WithResourceResolver(testCatalog(t)),
		WithDesiredValidator(valuevalidation.SurfaceValidator{}),
	)
	requireNoError(t, err)

	errs := runConcurrent(2, func(i int) error {
		_, err := executor.Apply(
			context.Background(),
			ApplyRequest{
				Object: testObject(1, "api:v2"),
				Owner:  owner("creator"),
				Force:  true,
			},
		)
		return err
	})

	successes := 0
	for _, err := range errs {
		switch {
		case err == nil:
			successes++
		case errors.Is(err, ErrStaleRevision), errors.Is(err, ErrConflict):
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

type oneWinnerUpdateStore struct {
	mu      sync.Mutex
	state   objectstore.State
	updated bool
}

func (s *oneWinnerUpdateStore) Get(context.Context, objectstore.Key) (objectstore.State, bool, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	return s.state, true, nil
}

func (s *oneWinnerUpdateStore) Create(context.Context, objectstore.Key, objectstore.State) (objectstore.State, error) {
	return objectstore.State{}, objectstore.ErrAlreadyExists
}

func (s *oneWinnerUpdateStore) Update(_ context.Context, _ objectstore.Key, _ objectstore.Revision, state objectstore.State) (objectstore.State, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.updated {
		return objectstore.State{}, objectstore.ErrStaleRevision
	}
	s.updated = true
	state.Revision = s.state.Revision + 1
	s.state = state
	return state, nil
}

func (s *oneWinnerUpdateStore) Delete(context.Context, objectstore.Key, objectstore.Revision) (objectstore.State, error) {
	return objectstore.State{}, objectstore.ErrConflict
}
