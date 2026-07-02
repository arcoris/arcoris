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
	"context"
	"errors"
	"sync"
	"testing"

	"arcoris.dev/apimachinery/api/objectstore"
	"arcoris.dev/apimachinery/runtime/objectworkqueue"
)

func TestNew(t *testing.T) {
	tests := []struct {
		name    string
		queue   Enqueuer
		mapper  Mapper
		wantErr error
	}{
		{name: "nil queue", queue: nil, mapper: ChangedObject(), wantErr: ErrNilQueue},
		{name: "typed nil queue", queue: (*recordingQueue)(nil), mapper: ChangedObject(), wantErr: ErrNilQueue},
		{name: "nil mapper", queue: &recordingQueue{}, mapper: nil, wantErr: ErrNilMapper},
		{name: "typed nil mapper func", queue: &recordingQueue{}, mapper: MapperFunc(nil), wantErr: ErrNilMapper},
		{name: "typed nil mapper pointer", queue: &recordingQueue{}, mapper: (*pointerMapper)(nil), wantErr: ErrNilMapper},
		{name: "valid", queue: &recordingQueue{}, mapper: ChangedObject()},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handler, err := New(tt.queue, tt.mapper)
			if tt.wantErr != nil {
				requireErrorIs(t, err, tt.wantErr)
				if handler != nil {
					t.Fatalf("handler = %#v; want nil", handler)
				}
				return
			}

			requireNoError(t, err)
			if handler == nil {
				t.Fatalf("handler is nil")
			}
		})
	}
}

func TestHandleInvalidHandler(t *testing.T) {
	tests := []struct {
		name    string
		handler *Handler
	}{
		{name: "nil", handler: nil},
		{name: "missing queue", handler: &Handler{mapper: ChangedObject()}},
		{name: "missing mapper", handler: &Handler{queue: &recordingQueue{}}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.handler.Handle(context.Background(), createdChange(t, 1))
			requireErrorIs(t, err, ErrInvalidHandler)
		})
	}
}

func TestHandleCallsMapperOnceAndPassesChange(t *testing.T) {
	change := updatedChange(t, 1)
	mapper := &recordingMapper{}
	handler := mustHandler(t, &recordingQueue{}, mapper)

	err := handler.Handle(context.Background(), change)

	requireNoError(t, err)
	if mapper.callCount() != 1 {
		t.Fatalf("mapper calls = %d; want 1", mapper.callCount())
	}
	got := mapper.onlyChange(t)
	requireChange(t, got, change)
}

func TestHandleEnqueuesChangedObject(t *testing.T) {
	ctx := context.WithValue(context.Background(), contextKey{}, "ctx")
	change := createdChange(t, 1)
	queue := &recordingQueue{}
	handler := mustHandler(t, queue, ChangedObject())

	err := handler.Handle(ctx, change)

	requireNoError(t, err)
	queue.requireContexts(t, ctx)
	queue.requireItems(t, objectworkqueue.Item{Key: change.Key})
}

func TestHandleReturnsQueueErrorUnchanged(t *testing.T) {
	wantErr := errors.New("queue failed")
	queue := &recordingQueue{err: wantErr}
	handler := mustHandler(t, queue, ChangedObject())

	err := handler.Handle(context.Background(), createdChange(t, 1))

	if err != wantErr {
		t.Fatalf("error = %v; want %v", err, wantErr)
	}
}

func TestHandleReturnsMapperErrorUnchanged(t *testing.T) {
	wantErr := errors.New("mapper failed")
	handler := mustHandler(t, &recordingQueue{}, MapperFunc(func(objectstore.Change, EmitFunc) error {
		return wantErr
	}))

	err := handler.Handle(context.Background(), createdChange(t, 1))

	if err != wantErr {
		t.Fatalf("error = %v; want %v", err, wantErr)
	}
}

func TestHandleStopsAddingAfterFirstEmitError(t *testing.T) {
	wantErr := errors.New("first add failed")
	queue := &recordingQueue{err: wantErr}
	handler := mustHandler(t, queue, MapperFunc(func(change objectstore.Change, emit EmitFunc) error {
		_ = emit(objectworkqueue.Item{Key: change.Key})
		_ = emit(objectworkqueue.Item{Key: testKey(2)})
		return nil
	}))

	err := handler.Handle(context.Background(), createdChange(t, 1))

	if err != wantErr {
		t.Fatalf("error = %v; want %v", err, wantErr)
	}
	queue.requireItems(t, objectworkqueue.Item{Key: testKey(1)})
}

func TestHandleAllowsZeroEmits(t *testing.T) {
	queue := &recordingQueue{}
	handler := mustHandler(t, queue, MapperFunc(func(objectstore.Change, EmitFunc) error {
		return nil
	}))

	err := handler.Handle(context.Background(), createdChange(t, 1))

	requireNoError(t, err)
	queue.requireItems(t)
}

func TestHandleConcurrentCalls(t *testing.T) {
	queue := &recordingQueue{}
	handler := mustHandler(t, queue, ChangedObject())

	const calls = 128
	changes := make([]objectstore.Change, calls)
	for i := range calls {
		changes[i] = createdChange(t, i+1)
	}

	errCh := make(chan error, calls)
	var wg sync.WaitGroup
	for i := 0; i < calls; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			errCh <- handler.Handle(context.Background(), changes[id])
		}(i)
	}
	wg.Wait()
	close(errCh)

	for err := range errCh {
		requireNoError(t, err)
	}
	if queue.callCount() != calls {
		t.Fatalf("queue calls = %d; want %d", queue.callCount(), calls)
	}
}

type contextKey struct{}

type recordingQueue struct {
	mu       sync.Mutex
	err      error
	contexts []context.Context
	items    []objectworkqueue.Item
}

func (q *recordingQueue) Add(ctx context.Context, item objectworkqueue.Item) error {
	q.mu.Lock()
	defer q.mu.Unlock()

	q.contexts = append(q.contexts, ctx)
	q.items = append(q.items, item)

	return q.err
}

func (q *recordingQueue) callCount() int {
	q.mu.Lock()
	defer q.mu.Unlock()

	return len(q.items)
}

func (q *recordingQueue) requireContexts(t testing.TB, want ...context.Context) {
	t.Helper()

	q.mu.Lock()
	defer q.mu.Unlock()

	if len(q.contexts) != len(want) {
		t.Fatalf("contexts = %d; want %d", len(q.contexts), len(want))
	}
	for i := range want {
		if q.contexts[i] != want[i] {
			t.Fatalf("context[%d] = %p; want %p", i, q.contexts[i], want[i])
		}
	}
}

func (q *recordingQueue) requireItems(t testing.TB, want ...objectworkqueue.Item) {
	t.Helper()

	q.mu.Lock()
	defer q.mu.Unlock()

	if len(q.items) != len(want) {
		t.Fatalf("items = %d; want %d", len(q.items), len(want))
	}
	for i := range want {
		requireItem(t, q.items[i], want[i])
	}
}

func mustHandler(t testing.TB, queue Enqueuer, mapper Mapper) *Handler {
	t.Helper()

	handler, err := New(queue, mapper)
	requireNoError(t, err)

	return handler
}
