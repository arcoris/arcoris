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
	"sync"
	"testing"

	storewatchapi "arcoris.dev/apimachinery/api/objectstorewatch"
	"arcoris.dev/apimachinery/api/objectwatch"
)

// fakeStream is a deterministic pull stream.
//
// It can return queued events, wait for context cancellation, or terminate with
// a configured error. That keeps reflector tests free of sleeps and background
// stream goroutines.
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

func validatorForRead(t testing.TB, read storewatchapi.CollectionRead, allowProgress bool) objectwatch.Validator {
	t.Helper()

	request, err := read.Boundary().WatchRequest(storewatchapi.WatchOptions{AllowProgress: allowProgress})
	requireNoError(t, err)
	validator, err := objectwatch.NewValidator(request)
	requireNoError(t, err)

	return validator
}

func TestConsumeStreamAppliesEventsUntilTerminalError(t *testing.T) {
	key := testKey("system", 1)
	sink := newRecordingSink(2)
	reflector := newTestReflector(t, &fakeListerWatcher{}, sink)
	reflector.lastApplied = 1
	validator := validatorForRead(t, testRead(t, 1), true)
	stream := streamWithEvents(
		changedEvent(t, updatedChange(t, key, 1, 2)),
		changedEvent(t, deletedChange(t, key, 2, 3)),
	)

	err := reflector.consumeStream(context.Background(), stream, &validator)

	requireErrorIs(t, err, objectwatch.ErrContinuityLost)
	changes := sink.recordedChanges()
	if len(changes) != 2 {
		t.Fatalf("changes = %d; want 2", len(changes))
	}
	if changes[0].Revision != 2 || changes[1].Revision != 3 {
		t.Fatalf("change revisions = %s, %s; want 2, 3", changes[0].Revision, changes[1].Revision)
	}
}

func TestConsumeStreamRejectsProgressWhenNotAllowed(t *testing.T) {
	sink := newRecordingSink(1)
	reflector := newTestReflector(t, &fakeListerWatcher{}, sink)
	reflector.lastApplied = 1
	reflector.lastProgress = 1
	validator := validatorForRead(t, testRead(t, 1), false)

	err := reflector.consumeStream(context.Background(), streamWithEvents(progressEvent(t, 2)), &validator)

	requireErrorIs(t, err, objectwatch.ErrContinuityLost)
	if sink.changeCount() != 0 {
		t.Fatalf("change count = %d; want 0", sink.changeCount())
	}
}

func TestConsumeStreamRejectsChangedEventBehindProgressBoundary(t *testing.T) {
	key := testKey("system", 1)
	sink := newRecordingSink(1)
	reflector := newTestReflector(t, &fakeListerWatcher{}, sink)
	reflector.lastApplied = 1
	reflector.lastProgress = 1
	validator := validatorForRead(t, testRead(t, 1), true)

	err := reflector.consumeStream(context.Background(), streamWithEvents(
		progressEvent(t, 10),
		changedEvent(t, updatedChange(t, key, 6, 7)),
	), &validator)

	requireErrorIs(t, err, objectwatch.ErrContinuityLost)
	if sink.changeCount() != 0 {
		t.Fatalf("change count = %d; want 0", sink.changeCount())
	}
	if reflector.lastProgress != 10 {
		t.Fatalf("lastProgress = %s; want 10", reflector.lastProgress)
	}
}

func TestConsumeStreamAcceptsProgressAtCollectionBoundary(t *testing.T) {
	sink := newRecordingSink(1)
	reflector := newTestReflector(t, &fakeListerWatcher{}, sink)
	reflector.lastApplied = 5
	reflector.lastProgress = 5
	validator := validatorForRead(t, testRead(t, 5), true)
	stream := &fakeStream{
		events:   []objectwatch.Event{progressEvent(t, 5)},
		terminal: context.Canceled,
	}

	err := reflector.consumeStream(context.Background(), stream, &validator)

	requireErrorIs(t, err, context.Canceled)
	if sink.changeCount() != 0 {
		t.Fatalf("change count = %d; want 0", sink.changeCount())
	}
	if reflector.lastProgress != 5 {
		t.Fatalf("lastProgress = %s; want 5", reflector.lastProgress)
	}
}

func TestConsumeStreamRejectsBackwardProgress(t *testing.T) {
	reflector := newTestReflector(t, &fakeListerWatcher{}, newRecordingSink(1))
	reflector.lastApplied = 1
	reflector.lastProgress = 1
	validator := validatorForRead(t, testRead(t, 1), true)

	err := reflector.consumeStream(context.Background(), streamWithEvents(
		progressEvent(t, 10),
		progressEvent(t, 9),
	), &validator)

	requireErrorIs(t, err, objectwatch.ErrContinuityLost)
	if reflector.lastProgress != 10 {
		t.Fatalf("lastProgress = %s; want 10", reflector.lastProgress)
	}
}

func TestConsumeStreamRejectsChangedEventOutsideCollection(t *testing.T) {
	reflector := newTestReflector(t, &fakeListerWatcher{}, newRecordingSink(1))
	reflector.lastApplied = 1
	validator := validatorForRead(t, testRead(t, 1), true)

	err := reflector.consumeStream(context.Background(), streamWithEvents(
		changedEvent(t, createdChange(t, otherResourceKey("system", 1), 2)),
	), &validator)

	requireErrorIs(t, err, objectwatch.ErrContinuityLost)
}

func TestConsumeStreamRejectsNonMonotonicChangedEvent(t *testing.T) {
	reflector := newTestReflector(t, &fakeListerWatcher{}, newRecordingSink(1))
	reflector.lastApplied = 1
	validator := validatorForRead(t, testRead(t, 1), true)

	err := reflector.consumeStream(context.Background(), streamWithEvents(
		changedEvent(t, createdChange(t, testKey("system", 1), 1)),
	), &validator)

	requireErrorIs(t, err, objectwatch.ErrContinuityLost)
}

func TestConsumeStreamClassifiesMalformedEvent(t *testing.T) {
	reflector := newTestReflector(t, &fakeListerWatcher{}, newRecordingSink(1))
	validator := validatorForRead(t, testRead(t, 1), true)

	err := reflector.consumeStream(context.Background(), streamWithEvents(objectwatch.Event{
		Kind:     objectwatch.EventChanged,
		Revision: 2,
	}), &validator)

	requireErrorIs(t, err, ErrInvalidEvent)
	requireErrorIs(t, err, objectwatch.ErrInvalidEvent)
}
