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
	"errors"
	"fmt"
	"sync"
	"testing"
	"time"

	apiidentity "arcoris.dev/apimachinery/api/identity"
	"arcoris.dev/apimachinery/api/meta"
	metaidentity "arcoris.dev/apimachinery/api/meta/identity"
	"arcoris.dev/apimachinery/api/object"
	"arcoris.dev/apimachinery/api/objectownership"
	"arcoris.dev/apimachinery/api/objectstore"
	storewatchapi "arcoris.dev/apimachinery/api/objectstorewatch"
	"arcoris.dev/apimachinery/api/objectwatch"
	"arcoris.dev/apimachinery/api/value"
)

type listResponse struct {
	read storewatchapi.CollectionRead
	err  error
}

type watchResponse struct {
	stream objectwatch.Stream
	err    error
}

type fakeListerWatcher struct {
	mu sync.Mutex

	listResponses  []listResponse
	watchResponses []watchResponse

	listRequests  []objectstore.ListRequest
	watchRequests []objectwatch.Request
}

func (f *fakeListerWatcher) ListCollection(
	_ context.Context,
	collection objectstore.ListRequest,
) (storewatchapi.CollectionRead, error) {
	f.mu.Lock()
	defer f.mu.Unlock()

	f.listRequests = append(f.listRequests, collection)
	if len(f.listResponses) == 0 {
		return storewatchapi.CollectionRead{}, errors.New("unexpected ListCollection call")
	}
	response := f.listResponses[0]
	f.listResponses = f.listResponses[1:]

	return response.read, response.err
}

func (f *fakeListerWatcher) Watch(_ context.Context, request objectwatch.Request) (objectwatch.Stream, error) {
	f.mu.Lock()
	defer f.mu.Unlock()

	f.watchRequests = append(f.watchRequests, request)
	if len(f.watchResponses) == 0 {
		return nil, errors.New("unexpected Watch call")
	}
	response := f.watchResponses[0]
	f.watchResponses = f.watchResponses[1:]

	return response.stream, response.err
}

func (f *fakeListerWatcher) addList(read storewatchapi.CollectionRead) {
	f.mu.Lock()
	defer f.mu.Unlock()

	f.listResponses = append(f.listResponses, listResponse{read: read})
}

func (f *fakeListerWatcher) addListError(err error) {
	f.mu.Lock()
	defer f.mu.Unlock()

	f.listResponses = append(f.listResponses, listResponse{err: err})
}

func (f *fakeListerWatcher) addWatch(stream objectwatch.Stream) {
	f.mu.Lock()
	defer f.mu.Unlock()

	f.watchResponses = append(f.watchResponses, watchResponse{stream: stream})
}

func (f *fakeListerWatcher) addWatchResponse(stream objectwatch.Stream, err error) {
	f.mu.Lock()
	defer f.mu.Unlock()

	f.watchResponses = append(f.watchResponses, watchResponse{stream: stream, err: err})
}

func (f *fakeListerWatcher) recordedWatchRequests() []objectwatch.Request {
	f.mu.Lock()
	defer f.mu.Unlock()

	return append([]objectwatch.Request(nil), f.watchRequests...)
}

func (f *fakeListerWatcher) recordedListRequests() []objectstore.ListRequest {
	f.mu.Lock()
	defer f.mu.Unlock()

	return append([]objectstore.ListRequest(nil), f.listRequests...)
}

type fakeStream struct {
	mu sync.Mutex

	events   []objectwatch.Event
	terminal error
	wait     bool

	closed     bool
	closeErr   error
	closeCount int

	nextStarted chan struct{}
	nextOnce    sync.Once
}

func streamWithEvents(events ...objectwatch.Event) *fakeStream {
	return &fakeStream{events: events, terminal: objectwatch.ContinuityLost(errors.New("stream exhausted"))}
}

func waitingStream() *fakeStream {
	return &fakeStream{wait: true, nextStarted: make(chan struct{})}
}

func terminalStream(err error) *fakeStream {
	return &fakeStream{terminal: err}
}

func (s *fakeStream) Next(ctx context.Context) (objectwatch.Event, error) {
	if s.nextStarted != nil {
		s.nextOnce.Do(func() { close(s.nextStarted) })
	}

	s.mu.Lock()
	if s.closed {
		s.mu.Unlock()
		return objectwatch.Event{}, objectwatch.Closed(nil)
	}
	if len(s.events) > 0 {
		event := s.events[0]
		copy(s.events, s.events[1:])
		s.events[len(s.events)-1] = objectwatch.Event{}
		s.events = s.events[:len(s.events)-1]
		s.mu.Unlock()
		return event, nil
	}
	terminal := s.terminal
	wait := s.wait
	s.mu.Unlock()

	if terminal != nil {
		return objectwatch.Event{}, terminal
	}
	if wait {
		<-ctx.Done()
		return objectwatch.Event{}, ctx.Err()
	}

	return objectwatch.Event{}, objectwatch.ContinuityLost(errors.New("stream exhausted"))
}

func (s *fakeStream) Close() error {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.closed = true
	s.closeCount++

	return s.closeErr
}

func (s *fakeStream) closes() int {
	s.mu.Lock()
	defer s.mu.Unlock()

	return s.closeCount
}

type recordingSink struct {
	mu sync.Mutex

	replaceErr error
	applyErr   error

	reads   []storewatchapi.CollectionRead
	changes []objectstore.Change

	replaceCh chan storewatchapi.CollectionRead
	changeCh  chan objectstore.Change
}

func newRecordingSink(buffer int) *recordingSink {
	return &recordingSink{
		replaceCh: make(chan storewatchapi.CollectionRead, buffer),
		changeCh:  make(chan objectstore.Change, buffer),
	}
}

func (s *recordingSink) Replace(_ context.Context, read storewatchapi.CollectionRead) error {
	if s.replaceErr != nil {
		return s.replaceErr
	}

	read = read.Clone()
	s.mu.Lock()
	s.reads = append(s.reads, read)
	s.mu.Unlock()

	select {
	case s.replaceCh <- read:
	default:
	}

	return nil
}

func (s *recordingSink) ApplyChange(_ context.Context, change objectstore.Change) error {
	if s.applyErr != nil {
		return s.applyErr
	}

	change = change.Clone()
	s.mu.Lock()
	s.changes = append(s.changes, change)
	s.mu.Unlock()

	select {
	case s.changeCh <- change:
	default:
	}

	return nil
}

func (s *recordingSink) replaceCount() int {
	s.mu.Lock()
	defer s.mu.Unlock()

	return len(s.reads)
}

func (s *recordingSink) changeCount() int {
	s.mu.Lock()
	defer s.mu.Unlock()

	return len(s.changes)
}

func (s *recordingSink) recordedChanges() []objectstore.Change {
	s.mu.Lock()
	defer s.mu.Unlock()

	out := make([]objectstore.Change, len(s.changes))
	for i, change := range s.changes {
		out[i] = change.Clone()
	}

	return out
}

func testResource() apiidentity.GroupVersionResource {
	return apiidentity.GroupVersionResource{
		Group:    "control.arcoris.dev",
		Version:  "v1",
		Resource: "workers",
	}
}

func otherResource() apiidentity.GroupVersionResource {
	return apiidentity.GroupVersionResource{
		Group:    "control.arcoris.dev",
		Version:  "v1",
		Resource: "tasks",
	}
}

func testCollection() objectstore.ListRequest {
	return objectstore.ListRequest{Resource: testResource(), Scope: objectstore.AllNamespaces()}
}

func namespaceCollection(namespace metaidentity.Namespace) objectstore.ListRequest {
	return objectstore.ListRequest{Resource: testResource(), Scope: objectstore.MustNamespace(namespace)}
}

func testKey(namespace metaidentity.Namespace, index int) objectstore.Key {
	return objectstore.MustKey(testResource(), metaidentity.ObjectName{
		Namespace: namespace,
		Name:      metaidentity.Name(fmt.Sprintf("worker-%d", index)),
	})
}

func otherResourceKey(namespace metaidentity.Namespace, index int) objectstore.Key {
	return objectstore.MustKey(otherResource(), metaidentity.ObjectName{
		Namespace: namespace,
		Name:      metaidentity.Name(fmt.Sprintf("task-%d", index)),
	})
}

func testState(key objectstore.Key, revision objectstore.Revision, desired string) objectstore.State {
	return objectstore.State{
		Object: object.NewObserved[value.Value, value.Value](
			meta.FromGroupVersionKind(apiidentity.GroupVersionKind{
				Group:   key.Resource.Group,
				Version: key.Resource.Version,
				Kind:    "Worker",
			}),
			meta.ObjectMeta{Name: key.Object.Name, Namespace: key.Object.Namespace},
			value.StringValue(desired),
			value.StringValue("observed-"+desired),
		),
		Ownership: objectownership.EmptyState(),
		Revision:  revision,
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

func testReadForCollection(
	t testing.TB,
	collection objectstore.ListRequest,
	revision objectstore.Revision,
	items ...objectstore.ListItem,
) storewatchapi.CollectionRead {
	t.Helper()

	read, err := storewatchapi.NewCollectionRead(collection, objectstore.ListResult{
		Items:    items,
		Revision: revision,
	})
	requireNoError(t, err)

	return read
}

func listItem(key objectstore.Key, revision objectstore.Revision, desired string) objectstore.ListItem {
	return objectstore.ListItem{Key: key, State: testState(key, revision, desired)}
}

func createdChange(t testing.TB, key objectstore.Key, revision objectstore.Revision) objectstore.Change {
	t.Helper()

	change, err := objectstore.NewCreatedChange(key, testState(key, revision, "created"))
	requireNoError(t, err)

	return change
}

func updatedChange(t testing.TB, key objectstore.Key, beforeRevision, afterRevision objectstore.Revision) objectstore.Change {
	t.Helper()

	change, err := objectstore.NewUpdatedChange(
		key,
		testState(key, beforeRevision, "before"),
		testState(key, afterRevision, "after"),
	)
	requireNoError(t, err)

	return change
}

func deletedChange(t testing.TB, key objectstore.Key, beforeRevision, deleteRevision objectstore.Revision) objectstore.Change {
	t.Helper()

	change, err := objectstore.NewDeletedChange(key, testState(key, beforeRevision, "deleted"), deleteRevision)
	requireNoError(t, err)

	return change
}

func changedEvent(t testing.TB, change objectstore.Change) objectwatch.Event {
	t.Helper()

	event, err := objectwatch.Changed(change)
	requireNoError(t, err)

	return event
}

func progressEvent(t testing.TB, revision objectstore.Revision) objectwatch.Event {
	t.Helper()

	event, err := objectwatch.Progress(revision)
	requireNoError(t, err)

	return event
}

func restartEvent(t testing.TB) objectwatch.Event {
	t.Helper()

	event, err := objectwatch.RestartRequired(objectwatch.RestartContinuityLost, 0)
	requireNoError(t, err)

	return event
}

func newTestReflector(t testing.TB, source storewatchapi.ListerWatcher, sink Sink) *Reflector {
	t.Helper()

	reflector, err := New(source, testCollection(), sink)
	requireNoError(t, err)

	return reflector
}

func waitRead(t testing.TB, ch <-chan storewatchapi.CollectionRead) storewatchapi.CollectionRead {
	t.Helper()

	select {
	case read := <-ch:
		return read
	case <-time.After(time.Second):
		t.Fatalf("timed out waiting for collection read")
		return storewatchapi.CollectionRead{}
	}
}

func waitChange(t testing.TB, ch <-chan objectstore.Change) objectstore.Change {
	t.Helper()

	select {
	case change := <-ch:
		return change
	case <-time.After(time.Second):
		t.Fatalf("timed out waiting for change")
		return objectstore.Change{}
	}
}

func requireNoChange(t testing.TB, ch <-chan objectstore.Change) {
	t.Helper()

	select {
	case change := <-ch:
		t.Fatalf("unexpected change: %#v", change)
	default:
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
		t.Fatalf("errors.Is(%v, %v) = false", err, target)
	}
}
