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
// Store serializes backend operations, collection reads, watch registration,
// history mutation, and live fanout with one mutex. The zero value is not
// usable; construct Store with New. Store contains a mutex and must not be
// copied after construction.
type Store struct {
	// mu serializes every backend call that can affect the observable boundary.
	// It also protects history and watchers. This deliberately coarse lock is
	// the continuity proof for replay computation followed by live registration.
	mu sync.Mutex
	// backend is the authoritative committed-state store. It must not be
	// mutated directly after wrapping because such writes bypass history/fanout.
	backend objectstore.Store
	// history stores a bounded revision-ordered suffix of wrapper-observed
	// committed changes. It is protected by mu.
	history changeHistory
	// streamBuffer is copied from construction options and used for each stream.
	streamBuffer int
	// watchers contains live streams registered for future matching changes. It
	// is protected by mu and entries are removed on Close or terminal overflow.
	watchers map[*watcher]struct{}
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
		backend:      backend,
		history:      newChangeHistory(opts.MaxHistory),
		streamBuffer: opts.StreamBuffer,
		watchers:     make(map[*watcher]struct{}),
	}, nil
}

// removeWatcher unregisters watcher from future live fanout.
//
// It is safe to call from stream.Close. The method acquires Store.mu itself, so
// callers must not already hold Store.mu when invoking it.
func (s *Store) removeWatcher(w *watcher) {
	if s == nil {
		return
	}
	s.mu.Lock()
	delete(s.watchers, w)
	s.mu.Unlock()
}

// closeAllWithContinuityLoss terminates every live watcher after an internal
// continuity violation that prevents the wrapper from preserving its contract.
//
// The caller must hold Store.mu. Terminated watchers are removed immediately so
// later committed changes cannot enqueue onto streams that have already lost
// continuity.
func (s *Store) closeAllWithContinuityLoss(cause error) {
	err := objectwatch.ContinuityLost(cause)
	for w := range s.watchers {
		w.terminate(err)
		delete(s.watchers, w)
	}
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
