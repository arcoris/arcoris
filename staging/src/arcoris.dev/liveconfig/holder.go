/*
  Copyright 2026 The ARCORIS Authors

  Licensed under the Apache License, Version 2.0 (the "License");
  you may not use this file except in compliance with the License.
  You may obtain a copy of the License at

      http://www.apache.org/licenses/LICENSE-2.0

  Unless required by applicable law or agreed to in writing, software
  distributed under the License is distributed on an "AS IS" BASIS,
  WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
  See the License for the specific language governing permissions and
  limitations under the License.
*/

package liveconfig

import (
	"errors"
	"sync"

	"arcoris.dev/snapshot"
)

// Holder owns the current last-good live configuration for one component or
// policy domain.
//
// A Holder serializes write-side Apply calls with mu, publishes accepted values
// through a snapshot.Publisher, and exposes Snapshot, Stamped, and Revision as
// read-side methods. Reads are delegated to the publisher and do not mutate
// holder state or execute validation, normalization, or source reload logic.
// This keeps read paths cheap and makes the holder suitable for read-mostly
// runtime policy such as limits, thresholds, retry knobs, and schedules.
//
// The holder does not own any input source. File watchers, environment readers,
// remote control-plane clients, and subscriber notification loops should call
// Apply from outside the package after they have built a candidate value.
//
// Holder is safe for concurrent use. Holder must be constructed with New and
// must not be copied after first use.
type Holder[T any] struct {
	// noCopy lets go vet report accidental Holder copies after first use.
	noCopy noCopy

	// mu serializes write-side Apply calls and protects lastErr.
	//
	// The snapshot.Publisher is safe for concurrent publishing, but Holder still
	// serializes Apply so the candidate preparation, equality check, publication,
	// and LastError update are observed as one coherent write-side operation.
	mu sync.Mutex

	// cfg contains the immutable construction policy used by New and Apply.
	//
	// Keeping these functions together makes the candidate pipeline explicit:
	// clone, then normalize, then validate, then optionally compare.
	cfg config[T]

	// pub publishes accepted immutable configuration values for read-side
	// Snapshot, Stamped, and Revision calls.
	pub *snapshot.Publisher[T]

	// lastErr is the last rejected Apply error. Successful changed and no-op
	// applies clear it.
	lastErr error
}

// ErrNilHolder reports a method call on a nil Holder receiver.
//
// Holder methods panic with this value instead of returning zero snapshots from
// a nil receiver. A nil Holder has no publisher, no last-good value, and no
// meaningful revision, so treating it as a programming error keeps failures
// explicit.
var ErrNilHolder = errors.New("liveconfig: nil holder")

// New creates a Holder containing initial as its first last-good value.
//
// New applies the same candidate pipeline used by Apply: clone, normalize,
// validate, and publish. If the initial value is rejected, New returns an error
// and no Holder. A successful New always publishes initial at a non-zero
// revision so readers never observe an uninitialized holder.
func New[T any](initial T, opts ...Option[T]) (*Holder[T], error) {
	cfg := newConfig(opts...)
	h := &Holder[T]{
		cfg: cfg,
		pub: snapshot.NewPublisher[T](
			snapshot.WithClock(cfg.clock),
		),
	}

	cur, err := h.prepare(initial)
	if err != nil {
		return nil, err
	}

	h.pub.Publish(cur)
	return h, nil
}

// Snapshot returns the current lightweight live configuration snapshot.
//
// Snapshot delegates to the internal snapshot.Publisher. It does not take the
// Holder write mutex, does not clone the value, and does not update LastError.
// The returned Value must be treated as immutable.
func (h *Holder[T]) Snapshot() snapshot.Snapshot[T] {
	requireHolder(h)
	return h.pub.Snapshot()
}

// Stamped returns the current stamped live configuration snapshot.
//
// Stamped includes the local publication time assigned when the current value
// was accepted. It has the same immutability and read-side behavior as Snapshot.
func (h *Holder[T]) Stamped() snapshot.Stamped[T] {
	requireHolder(h)
	return h.pub.Stamped()
}

// Revision returns the current source-local configuration revision.
//
// Revision is a cheap read-side change check for consumers that do not need the
// value itself. The revision is local to this holder.
func (h *Holder[T]) Revision() snapshot.Revision {
	requireHolder(h)
	return h.pub.Revision()
}

// LastError returns the most recent rejected Apply error.
//
// LastError is diagnostic state for reload loops and operators. It is set only
// when Apply rejects a candidate after clone, normalization, or validation. It
// is cleared after any successful Apply attempt, including an EqualFunc no-op,
// because the most recent candidate was accepted.
func (h *Holder[T]) LastError() error {
	requireHolder(h)

	h.mu.Lock()
	defer h.mu.Unlock()

	return h.lastErr
}

// requireHolder panics when h is nil.
func requireHolder[T any](h *Holder[T]) {
	if h == nil {
		panic(ErrNilHolder)
	}
}
