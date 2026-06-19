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
	"errors"
	"sync"

	"arcoris.dev/apimachinery/api/objectstore"
	storewatchapi "arcoris.dev/apimachinery/api/objectstorewatch"
	"arcoris.dev/apimachinery/api/objectwatch"
)

var (
	_ objectstore.Store              = (*Store)(nil)
	_ storewatchapi.CollectionLister = (*Store)(nil)
	_ storewatchapi.ListerWatcher    = (*Store)(nil)
	_ storewatchapi.Store            = (*Store)(nil)
	_ storewatchapi.CapableStore     = (*Store)(nil)
	_ objectwatch.Source             = (*Store)(nil)
	_ objectwatch.CapabilityReporter = (*Store)(nil)
)

// Store is a passive write-through observable wrapper around objectstore.Store.
//
// Store.mu serializes backend calls, collection reads, watch registration,
// retained history mutation, live stream registry mutation, and live fanout.
// The zero value is not usable; construct Store with New. Store contains a
// mutex and must not be copied after construction.
type Store struct {
	// mu serializes every backend call that can affect the observable boundary.
	// It also protects history and streams. This deliberately coarse lock is
	// the continuity proof for replay computation followed by live registration.
	mu sync.Mutex
	// backend is the authoritative committed-state store. It must not be
	// mutated directly after wrapping because such writes bypass history/fanout.
	backend objectstore.Store
	// history stores a bounded revision-ordered suffix of wrapper-observed
	// committed changes. It is protected by mu.
	history changeHistory
	// options is the validated construction configuration. It is immutable after
	// New returns and read only while Store.mu is held.
	options Options
	// nextStreamID is incremented under mu to give each live stream a stable
	// registry key that cannot be confused with another stream pointer.
	nextStreamID uint64
	// streams contains live streams registered for future matching changes. It
	// is protected by mu and entries are removed on Close or terminal overflow.
	streams map[uint64]*stream
}

// New constructs an observable wrapper around backend.
func New(backend objectstore.Store, options ...Option) (*Store, error) {
	if backend == nil {
		return nil, ErrNilBackend
	}
	opts, err := applyOptions(options)
	if err != nil {
		return nil, err
	}

	return &Store{
		backend: backend,
		history: newChangeHistory(opts.MaxHistory),
		options: opts,
		streams: make(map[uint64]*stream),
	}, nil
}

// unregister removes stream id from future live fanout.
//
// It is safe to call from stream.Close. The method acquires Store.mu itself, so
// callers must not already hold Store.mu when invoking it.
func (s *Store) unregister(id uint64) {
	if s == nil {
		return
	}
	s.mu.Lock()
	s.unregisterLocked(id)
	s.mu.Unlock()
}

// nextStreamIDLocked returns the next live stream registry key.
//
// Store.mu must be held.
func (s *Store) nextStreamIDLocked() uint64 {
	s.nextStreamID++
	return s.nextStreamID
}

// registerLocked registers stream for future live fanout.
//
// Store.mu must be held.
func (s *Store) registerLocked(stream *stream) {
	s.streams[stream.id] = stream
}

// unregisterLocked removes id from future live fanout.
//
// Store.mu must be held.
func (s *Store) unregisterLocked(id uint64) {
	delete(s.streams, id)
}

// loseContinuityLocked terminates every live stream after an internal
// continuity violation that prevents the wrapper from preserving its contract.
//
// Store.mu must be held. Terminated streams are removed immediately so later
// committed changes cannot enqueue onto streams that have already lost
// continuity.
func (s *Store) loseContinuityLocked(cause error) error {
	err := continuityError(cause)
	for id, stream := range s.streams {
		stream.finish(err)
		delete(s.streams, id)
	}
	return err
}

// loseCommittedContinuityLocked invalidates history after a committed mutation
// that the wrapper could not publish.
//
// Store.mu must be held. A valid revision becomes the new unsafe replay
// boundary. An invalid revision taints all historical replay because no caller
// can prove their requested start is after the missing committed mutation.
func (s *Store) loseCommittedContinuityLocked(revision objectstore.Revision, cause error) error {
	if revision.IsValid() {
		s.history.invalidateThrough(revision)
	} else {
		s.history.taint(cause)
	}

	return s.loseContinuityLocked(cause)
}

// continuityError records a wrapper/backend invariant violation.
//
// The wrapper uses objectwatch continuity errors for committed-backend states it
// cannot translate into a valid objectstore.Change. Returning such an error is
// preferable to pretending that the stream contract remains trustworthy.
func continuityError(cause error) error {
	if cause == nil {
		cause = errors.New("object store watch continuity violation")
	}
	return objectwatch.ContinuityLost(cause)
}
