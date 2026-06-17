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
	"errors"
	"sync"

	"arcoris.dev/apimachinery/api/objectwatch"
)

// stream is one bounded pull-based watch stream.
//
// stream has no goroutine. Writers enqueue directly with a non-blocking send
// while holding Store.mu. Readers pull through Next. Close wakes blocked
// readers through done and unregisters the watcher from the parent Store.
type stream struct {
	// parent owns the live watcher registry used by Close.
	parent *Store
	// watcher is the registry key associated with this stream.
	watcher *watcher
	// events is a bounded queue of cloned, validated watch events.
	events chan objectwatch.Event
	// done is closed exactly once when the stream reaches a terminal state.
	done chan struct{}

	// mu protects terminal and closed. It is intentionally independent from
	// Store.mu so Next can observe terminal state without touching the registry.
	mu sync.Mutex
	// terminal is the error returned after done is closed.
	terminal error
	// closed records whether terminal has been set.
	closed bool
	// once makes Close and continuity termination idempotent.
	once sync.Once
}

// newStream constructs a stream with a bounded event queue.
func newStream(parent *Store, watcher *watcher, buffer int) *stream {
	return &stream{
		parent:  parent,
		watcher: watcher,
		events:  make(chan objectwatch.Event, buffer),
		done:    make(chan struct{}),
	}
}

// Next returns the next event or a terminal/context error.
//
// The method never returns a normal EOF. If the stream has reached a terminal
// state, the terminal objectwatch-compatible error is returned instead.
func (s *stream) Next(ctx context.Context) (objectwatch.Event, error) {
	if ctx == nil {
		ctx = context.Background()
	}
	if err := s.terminalError(); err != nil {
		return objectwatch.Event{}, err
	}

	select {
	case <-ctx.Done():
		return objectwatch.Event{}, ctx.Err()
	case <-s.done:
		return objectwatch.Event{}, s.terminalError()
	case event := <-s.events:
		if err := s.terminalError(); err != nil {
			return objectwatch.Event{}, err
		}
		return event.Clone(), nil
	}
}

// Close unregisters and closes the stream. It is idempotent.
func (s *stream) Close() error {
	if s == nil {
		return nil
	}
	s.parent.removeWatcher(s.watcher)
	s.closeWithError(objectwatch.Closed(errors.New("stream closed")))
	return nil
}

// enqueue attempts to queue event without blocking.
//
// A false return means the caller must close the stream with continuity loss.
// Valid events are cloned before queuing so callers cannot mutate queued data.
func (s *stream) enqueue(event objectwatch.Event) bool {
	if err := event.Validate(); err != nil {
		s.closeWithError(continuityError(err))
		return true
	}
	if s.terminalError() != nil {
		return true
	}

	select {
	case s.events <- event.Clone():
		return true
	default:
		return false
	}
}

// closeWithError records terminal err and wakes blocked readers.
//
// It does not unregister from Store. Callers that initiate user-visible Close
// must unregister first; callers already holding Store.mu remove the watcher
// themselves after invoking this method.
func (s *stream) closeWithError(err error) {
	if err == nil {
		err = objectwatch.Closed(errors.New("stream closed"))
	}
	s.once.Do(func() {
		s.mu.Lock()
		s.terminal = err
		s.closed = true
		s.mu.Unlock()
		close(s.done)
	})
}

// terminalError returns the stream's terminal error if it is closed.
func (s *stream) terminalError() error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if !s.closed {
		return nil
	}
	if s.terminal != nil {
		return s.terminal
	}

	return objectwatch.Closed(errors.New("stream closed"))
}
