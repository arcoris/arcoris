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

	"arcoris.dev/apimachinery/api/objectstore"
	"arcoris.dev/apimachinery/api/objectwatch"
)

// stream is one bounded pull-based watch stream.
//
// stream has no goroutine. Historical replay is stored separately from the
// bounded live queue so StreamBuffer limits only live backlog, never retained
// replay. Writers enqueue live events directly with a non-blocking send while
// holding Store.mu. Readers pull through Next. Close wakes blocked readers
// through done and unregisters the stream from the parent Store.
type stream struct {
	// parent owns the live stream registry used by Close.
	parent *Store
	// id is the Store registry key associated with this stream.
	id uint64
	// request is immutable after registration and defines live fanout matching.
	request objectwatch.Request
	// live is a bounded queue of cloned, validated live watch events.
	live chan objectwatch.Event
	// done is closed exactly once when the stream reaches a terminal state.
	done chan struct{}

	// mu protects replay, terminal, and finished. It is intentionally independent
	// from Store.mu so Next can observe terminal state without touching the
	// registry.
	mu sync.Mutex
	// replay is drained before live. It is never constrained by StreamBuffer.
	replay []objectwatch.Event
	// terminal is the error returned after done is closed.
	terminal error
	// finished records whether terminal has been set.
	finished bool
	// once makes Close and continuity termination idempotent.
	once sync.Once
}

// newStream constructs a stream with detached replay and a bounded live queue.
func newStream(
	parent *Store,
	id uint64,
	request objectwatch.Request,
	replay []objectstore.Change,
	buffer int,
) (*stream, error) {
	events, err := replayEvents(replay)
	if err != nil {
		return nil, err
	}

	return &stream{
		parent:  parent,
		id:      id,
		request: request,
		replay:  events,
		live:    make(chan objectwatch.Event, buffer),
		done:    make(chan struct{}),
	}, nil
}

// Next returns the next event or a terminal/context error.
//
// The method never returns a normal EOF. If the stream has reached a terminal
// state, the terminal objectwatch-compatible error is returned instead.
// Concurrent Next calls are safe, but callers should not rely on which goroutine
// receives a particular event.
func (s *stream) Next(ctx context.Context) (objectwatch.Event, error) {
	if ctx == nil {
		ctx = context.Background()
	}
	if event, ok, err := s.nextReplay(); err != nil || ok {
		return event, err
	}

	select {
	case <-ctx.Done():
		return objectwatch.Event{}, ctx.Err()
	case <-s.done:
		return objectwatch.Event{}, s.terminalError()
	case event := <-s.live:
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
	s.parent.unregister(s.id)
	s.finish(objectwatch.Closed(errors.New("stream closed")))
	return nil
}

// enqueue attempts to queue event without blocking.
//
// A false return means the caller must unregister the stream and finish it with
// continuity loss if enqueue did not already do so. Valid events are cloned
// before queuing so callers cannot mutate queued data.
func (s *stream) enqueue(event objectwatch.Event) bool {
	if err := event.Validate(); err != nil {
		s.finish(continuityError(err))
		return false
	}
	if s.terminalError() != nil {
		return true
	}

	select {
	case s.live <- event.Clone():
		return true
	default:
		return false
	}
}

// finish records terminal err and wakes blocked readers.
//
// It does not unregister from Store. Callers that initiate user-visible Close
// must unregister first; callers already holding Store.mu remove the stream
// themselves before invoking this method.
func (s *stream) finish(err error) {
	if err == nil {
		err = objectwatch.Closed(errors.New("stream closed"))
	}
	s.once.Do(func() {
		s.mu.Lock()
		s.terminal = err
		s.finished = true
		s.mu.Unlock()
		close(s.done)
	})
}

// nextReplay returns the next replay event before the stream reads live events.
func (s *stream) nextReplay() (objectwatch.Event, bool, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.finished {
		return objectwatch.Event{}, false, s.terminal
	}
	if len(s.replay) == 0 {
		return objectwatch.Event{}, false, nil
	}

	event := s.replay[0].Clone()
	s.replay[0] = objectwatch.Event{}
	s.replay = s.replay[1:]
	return event, true, nil
}

// terminalError returns the stream's terminal error if it is closed.
func (s *stream) terminalError() error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if !s.finished {
		return nil
	}
	if s.terminal != nil {
		return s.terminal
	}

	return objectwatch.Closed(errors.New("stream closed"))
}
