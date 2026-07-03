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

	"arcoris.dev/apimachinery/api/objectquery"
	"arcoris.dev/apimachinery/api/objectstore"
	storewatchapi "arcoris.dev/apimachinery/api/objectstorewatch"
	"arcoris.dev/apimachinery/runtime/objectreflector"
	"arcoris.dev/apimachinery/runtime/objectworkqueue"
)

var _ objectreflector.Sink = (*ReflectorSink)(nil)

func TestNewReflectorSink(t *testing.T) {
	queue := newSinkQueue()
	tests := []struct {
		name    string
		config  ReflectorSinkConfig
		wantErr error
	}{
		{
			name: "nil queue",
			config: ReflectorSinkConfig{
				Listed:  ListedObject(),
				Changed: ChangedObject(),
			},
			wantErr: ErrNilQueue,
		},
		{
			name: "typed nil queue",
			config: ReflectorSinkConfig{
				Queue:   (*recordingQueue)(nil),
				Listed:  ListedObject(),
				Changed: ChangedObject(),
			},
			wantErr: ErrNilQueue,
		},
		{
			name: "nil listed mapper",
			config: ReflectorSinkConfig{
				Queue:   queue,
				Changed: ChangedObject(),
			},
			wantErr: ErrNilListItemMapper,
		},
		{
			name: "typed nil listed mapper func",
			config: ReflectorSinkConfig{
				Queue:   queue,
				Listed:  ListItemMapperFunc(nil),
				Changed: ChangedObject(),
			},
			wantErr: ErrNilListItemMapper,
		},
		{
			name: "typed nil listed mapper pointer",
			config: ReflectorSinkConfig{
				Queue:   queue,
				Listed:  (*pointerListItemMapper)(nil),
				Changed: ChangedObject(),
			},
			wantErr: ErrNilListItemMapper,
		},
		{
			name: "nil changed mapper",
			config: ReflectorSinkConfig{
				Queue:  queue,
				Listed: ListedObject(),
			},
			wantErr: ErrNilMapper,
		},
		{
			name: "typed nil changed mapper func",
			config: ReflectorSinkConfig{
				Queue:   queue,
				Listed:  ListedObject(),
				Changed: MapperFunc(nil),
			},
			wantErr: ErrNilMapper,
		},
		{
			name: "typed nil changed mapper pointer",
			config: ReflectorSinkConfig{
				Queue:   queue,
				Listed:  ListedObject(),
				Changed: (*pointerMapper)(nil),
			},
			wantErr: ErrNilMapper,
		},
		{
			name: "zero predicate",
			config: ReflectorSinkConfig{
				Queue:   queue,
				Listed:  ListedObject(),
				Changed: ChangedObject(),
			},
		},
		{
			name: "valid",
			config: ReflectorSinkConfig{
				Queue:     queue,
				Predicate: mustPredicate(t, objectquery.All()),
				Listed:    ListedObject(),
				Changed:   ChangedObject(),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sink, err := NewReflectorSink(tt.config)
			if tt.wantErr != nil {
				requireErrorIs(t, err, tt.wantErr)
				if sink != nil {
					t.Fatalf("sink = %#v; want nil", sink)
				}
				return
			}

			requireNoError(t, err)
			if sink == nil {
				t.Fatalf("sink is nil")
			}
			requireKnownKeys(t, sink)
		})
	}
}

func TestReflectorSinkInvalidReceiver(t *testing.T) {
	tests := []struct {
		name string
		sink *ReflectorSink
	}{
		{name: "nil", sink: nil},
		{name: "missing queue", sink: &ReflectorSink{listed: ListedObject(), changed: ChangedObject(), known: map[objectstore.Key]objectstore.ListItem{}}},
		{name: "typed nil queue", sink: &ReflectorSink{queue: (*recordingQueue)(nil), listed: ListedObject(), changed: ChangedObject(), known: map[objectstore.Key]objectstore.ListItem{}}},
		{name: "missing listed", sink: &ReflectorSink{queue: newSinkQueue(), changed: ChangedObject(), known: map[objectstore.Key]objectstore.ListItem{}}},
		{name: "typed nil listed", sink: &ReflectorSink{queue: newSinkQueue(), listed: (*pointerListItemMapper)(nil), changed: ChangedObject(), known: map[objectstore.Key]objectstore.ListItem{}}},
		{name: "missing changed", sink: &ReflectorSink{queue: newSinkQueue(), listed: ListedObject(), known: map[objectstore.Key]objectstore.ListItem{}}},
		{name: "typed nil changed", sink: &ReflectorSink{queue: newSinkQueue(), listed: ListedObject(), changed: (*pointerMapper)(nil), known: map[objectstore.Key]objectstore.ListItem{}}},
		{name: "missing known", sink: &ReflectorSink{queue: newSinkQueue(), listed: ListedObject(), changed: ChangedObject()}},
	}

	for _, tt := range tests {
		t.Run(tt.name+"/replace", func(t *testing.T) {
			err := tt.sink.Replace(context.Background(), testRead(t, 1))
			requireErrorIs(t, err, ErrInvalidReflectorSink)
		})
		t.Run(tt.name+"/apply", func(t *testing.T) {
			err := tt.sink.ApplyChange(context.Background(), createdChange(t, 1))
			requireErrorIs(t, err, ErrInvalidReflectorSink)
		})
	}
}

func TestReflectorSinkConcurrentReplaceAndApplyChange(t *testing.T) {
	queue := newSinkQueue()
	sink := newTestReflectorSink(t, queue)

	const calls = 64
	reads := make([]storewatchapi.CollectionRead, calls)
	changes := make([]objectstore.Change, calls)
	for i := range calls {
		id := i + 1
		reads[i] = testRead(t, objectstore.Revision(id), testListItem(id, objectstore.Revision(id), "listed"))
		changes[i] = createdChange(t, id)
	}

	errCh := make(chan error, calls)
	var wg sync.WaitGroup
	for i := range calls {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			if id%2 == 0 {
				errCh <- sink.Replace(context.Background(), reads[id])
				return
			}
			errCh <- sink.ApplyChange(context.Background(), changes[id])
		}(i)
	}
	wg.Wait()
	close(errCh)

	for err := range errCh {
		requireNoError(t, err)
	}
}

func newTestReflectorSink(t testing.TB, queue Enqueuer, opts ...func(*ReflectorSinkConfig)) *ReflectorSink {
	t.Helper()

	config := ReflectorSinkConfig{
		Queue:   queue,
		Listed:  ListedObject(),
		Changed: ChangedObject(),
	}
	for _, opt := range opts {
		opt(&config)
	}

	sink, err := NewReflectorSink(config)
	requireNoError(t, err)

	return sink
}

func withPredicate(predicate objectquery.Predicate) func(*ReflectorSinkConfig) {
	return func(config *ReflectorSinkConfig) {
		config.Predicate = predicate
	}
}

func withListed(mapper ListItemMapper) func(*ReflectorSinkConfig) {
	return func(config *ReflectorSinkConfig) {
		config.Listed = mapper
	}
}

func withChanged(mapper Mapper) func(*ReflectorSinkConfig) {
	return func(config *ReflectorSinkConfig) {
		config.Changed = mapper
	}
}

type sinkQueue struct {
	mu sync.Mutex

	err    error
	failAt int

	contexts []context.Context
	items    []objectworkqueue.Item
}

func newSinkQueue() *sinkQueue {
	return &sinkQueue{}
}

func (q *sinkQueue) Add(ctx context.Context, item objectworkqueue.Item) error {
	q.mu.Lock()
	defer q.mu.Unlock()

	q.contexts = append(q.contexts, ctx)
	q.items = append(q.items, item)
	if q.err != nil && (q.failAt == 0 || len(q.items) == q.failAt) {
		return q.err
	}

	return nil
}

func (q *sinkQueue) reset() {
	q.mu.Lock()
	defer q.mu.Unlock()

	q.contexts = nil
	q.items = nil
}

func (q *sinkQueue) requireContexts(t testing.TB, want ...context.Context) {
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

func (q *sinkQueue) requireItems(t testing.TB, want ...objectworkqueue.Item) {
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

func testRead(t testing.TB, revision objectstore.Revision, items ...objectstore.ListItem) storewatchapi.CollectionRead {
	t.Helper()

	read, err := storewatchapi.NewCollectionRead(testCollection(), objectstore.ListResult{
		Items:    items,
		Revision: revision,
	})
	requireNoError(t, err)

	return read
}

func testCollection() objectstore.ListRequest {
	return objectstore.ListRequest{Resource: testResource(), Scope: objectstore.AllNamespaces()}
}

func testListItem(id int, revision objectstore.Revision, desired string) objectstore.ListItem {
	key := testKey(id)

	return objectstore.ListItem{Key: key, State: testState(key, revision, desired)}
}

func requireListItem(t testing.TB, got objectstore.ListItem, want objectstore.ListItem) {
	t.Helper()

	if !got.Key.Equal(want.Key) || got.State.Revision != want.State.Revision {
		t.Fatalf("list item = %#v; want %#v", got, want)
	}
}

func requireKnownKeys(t testing.TB, sink *ReflectorSink, want ...objectstore.Key) {
	t.Helper()

	sink.mu.Lock()
	defer sink.mu.Unlock()

	if len(sink.order) != len(want) {
		t.Fatalf("known order = %v; want %v", sink.order, want)
	}
	if len(sink.known) != len(want) {
		t.Fatalf("known count = %d; want %d", len(sink.known), len(want))
	}
	for i := range want {
		if !sink.order[i].Equal(want[i]) {
			t.Fatalf("known order = %v; want %v", sink.order, want)
		}
		if _, ok := sink.known[want[i]]; !ok {
			t.Fatalf("known missing key %s", want[i])
		}
	}
}

func requireKnownRevision(t testing.TB, sink *ReflectorSink, key objectstore.Key, want objectstore.Revision) {
	t.Helper()

	sink.mu.Lock()
	defer sink.mu.Unlock()

	item, ok := sink.known[key]
	if !ok {
		t.Fatalf("known missing key %s", key)
	}
	if item.State.Revision != want {
		t.Fatalf("known revision = %s; want %s", item.State.Revision, want)
	}
}

func requireUnknown(t testing.TB, sink *ReflectorSink, key objectstore.Key) {
	t.Helper()

	sink.mu.Lock()
	defer sink.mu.Unlock()

	if _, ok := sink.known[key]; ok {
		t.Fatalf("known contains key %s", key)
	}
}

func requireErrorSame(t testing.TB, got error, want error) {
	t.Helper()

	if got != want {
		t.Fatalf("error = %v; want %v", got, want)
	}
}

var errMapperFailed = errors.New("mapper failed")
var errQueueFailed = errors.New("queue failed")
