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

package objectstorewatch

import (
	"context"
	"fmt"
	"sync"
	"testing"
	"time"

	"arcoris.dev/apimachinery/api/objectwatch"
)

func TestConcurrentWatchAndCreate(t *testing.T) {
	store := testRuntimeStore(t, WithStreamBuffer(16))
	request := watchRequestAfter(t, testCollection(), 0)
	errs := make(chan error, 16)
	var wg sync.WaitGroup

	for i := 0; i < 8; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			stream, err := store.Watch(context.Background(), request)
			if err != nil {
				errs <- fmt.Errorf("watch %d: %w", i, err)
				return
			}
			if err := stream.Close(); err != nil {
				errs <- fmt.Errorf("close %d: %w", i, err)
			}
		}(i)
	}
	for i := 0; i < 8; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			key := testKey("system", i)
			_, err := store.Create(context.Background(), key, stateForKey(key, "created"))
			if err != nil {
				errs <- fmt.Errorf("create %d: %w", i, err)
			}
		}(i)
	}

	wg.Wait()
	close(errs)
	for err := range errs {
		requireNoError(t, err)
	}
}

func TestConcurrentCloseAndFanout(t *testing.T) {
	store := testRuntimeStore(t, WithStreamBuffer(16))
	stream := watchAfter(t, store, testCollection(), 0)
	errs := make(chan error, 2)
	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		defer wg.Done()
		if err := stream.Close(); err != nil {
			errs <- err
		}
	}()
	wg.Add(1)
	go func() {
		defer wg.Done()
		key := testKey("system", 1)
		_, err := store.Create(context.Background(), key, stateForKey(key, "created"))
		if err != nil {
			errs <- err
		}
	}()

	wg.Wait()
	close(errs)
	for err := range errs {
		requireNoError(t, err)
	}
	_, err := stream.Next(context.Background())
	requireErrorIs(t, err, objectwatch.ErrClosed)
}

func TestConcurrentNextCallsAreRaceSafe(t *testing.T) {
	store := testRuntimeStore(t, WithStreamBuffer(16))
	createObject(t, store, testKey("system", 1), "one")
	createObject(t, store, testKey("system", 2), "two")
	stream := watchAfter(t, store, testCollection(), 0)
	defer func() { requireNoError(t, stream.Close()) }()

	results := make(chan objectwatch.Event, 2)
	errs := make(chan error, 2)
	var wg sync.WaitGroup
	for i := 0; i < 2; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			ctx, cancel := context.WithTimeout(context.Background(), time.Second)
			defer cancel()
			event, err := stream.Next(ctx)
			if err != nil {
				errs <- err
				return
			}
			results <- event
		}()
	}

	wg.Wait()
	close(results)
	close(errs)
	for err := range errs {
		requireNoError(t, err)
	}
	if len(results) != 2 {
		t.Fatalf("results = %d; want 2", len(results))
	}
}
