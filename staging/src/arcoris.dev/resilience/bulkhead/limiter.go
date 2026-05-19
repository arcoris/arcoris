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

package bulkhead

import (
	"errors"
	"sync"

	"arcoris.dev/snapshot"
)

// Limiter is a local non-blocking concurrency bulkhead.
//
// Limiter owns coherent mutable admission state under mu and publishes immutable
// read models through snapshot.Publisher. TryAcquire and permit release are the
// only operations that mutate limiter state. Snapshot, Stamped, and Revision are
// delegated to the publisher and do not acquire mu.
//
// Limiter is safe for concurrent use. A Limiter must be constructed with New and
// must not be copied after first use.
type Limiter struct {
	// noCopy lets go vet report accidental Limiter copies after first use.
	noCopy noCopy

	// mu protects all mutable limiter state below.
	mu sync.Mutex

	// cfg contains construction-time limiter policy.
	cfg config

	// inFlight is the number of currently held permits.
	inFlight uint64

	// acquired is the lifetime number of successful acquisitions.
	acquired uint64

	// rejected is the lifetime number of rejected acquisitions.
	rejected uint64

	// released is the lifetime number of releases.
	released uint64

	// published exposes the latest immutable read model.
	published *snapshot.Publisher[Snapshot]
}

var (
	// ErrNilLimiter reports a method call on a nil Limiter receiver.
	ErrNilLimiter = errors.New("bulkhead: nil limiter")

	// ErrUninitializedLimiter reports use of a zero-value Limiter.
	//
	// Limiter requires construction through New so the initial snapshot and
	// publisher are installed before readers can observe the limiter.
	ErrUninitializedLimiter = errors.New("bulkhead: uninitialized limiter")
)

// New creates a Limiter with the provided concurrency limit.
//
// The limit is required because capacity is the defining bulkhead policy. A
// successful New publishes an initial snapshot at a non-zero revision with zero
// in-flight permits and zero lifetime counters.
func New(limit uint64, opts ...Option) (*Limiter, error) {
	cfg := newConfig(limit, opts...)
	if err := validateConfig(cfg); err != nil {
		return nil, err
	}

	l := &Limiter{
		cfg: cfg,
		published: snapshot.NewPublisher[Snapshot](
			snapshot.WithClock(cfg.clock),
		),
	}
	l.publishLocked()

	return l, nil
}

// TryAcquire attempts to acquire one permit without waiting.
//
// When capacity is available, TryAcquire returns a release-once Permit and an
// allowed Decision. When the limiter is full, TryAcquire returns nil and a denied
// Decision with ReasonFull. Both successful and rejected attempts publish a new
// snapshot because lifetime counters change in either branch.
func (l *Limiter) TryAcquire() (*Permit, Decision) {
	requireLimiter(l)

	l.mu.Lock()
	defer l.mu.Unlock()

	if l.inFlight >= l.cfg.limit {
		l.rejected++
		snap := l.publishLocked()
		return nil, Decision{
			Allowed:  false,
			Reason:   ReasonFull,
			Snapshot: snap,
		}
	}

	l.inFlight++
	l.acquired++
	snap := l.publishLocked()

	return &Permit{limiter: l}, Decision{
		Allowed:  true,
		Reason:   ReasonAllowed,
		Snapshot: snap,
	}
}

// Snapshot returns the current lightweight limiter snapshot.
//
// Snapshot delegates to the internal snapshot.Publisher. It does not take the
// limiter mutex and does not mutate limiter state.
func (l *Limiter) Snapshot() snapshot.Snapshot[Snapshot] {
	requireLimiter(l)
	return l.published.Snapshot()
}

// Stamped returns the current stamped limiter snapshot.
//
// Stamped includes the local publication time assigned when the current read
// model was accepted by the limiter.
func (l *Limiter) Stamped() snapshot.Stamped[Snapshot] {
	requireLimiter(l)
	return l.published.Stamped()
}

// Revision returns the current source-local limiter revision.
func (l *Limiter) Revision() snapshot.Revision {
	requireLimiter(l)
	return l.published.Revision()
}

// release returns one permit to the limiter.
//
// release is called only by Permit.Release after the permit has won its
// release-once guard.
func (l *Limiter) release() {
	l.mu.Lock()
	defer l.mu.Unlock()

	if l.inFlight == 0 {
		panic(ErrReleaseUnderflow)
	}

	l.inFlight--
	l.released++
	l.publishLocked()
}

// snapshotValueLocked builds the current immutable read model.
//
// The caller must hold l.mu, except during New before l is published to other
// goroutines.
func (l *Limiter) snapshotValueLocked() Snapshot {
	return Snapshot{
		Capacity: newCapacitySnapshot(l.cfg.limit, l.inFlight),
		Stats: StatsSnapshot{
			Acquired: l.acquired,
			Rejected: l.rejected,
			Released: l.released,
		},
	}
}

// publishLocked publishes the current immutable read model.
//
// The caller must hold l.mu, except during New before l is published to other
// goroutines.
func (l *Limiter) publishLocked() snapshot.Snapshot[Snapshot] {
	return l.published.Publish(l.snapshotValueLocked())
}

// requireLimiter panics when l is nil or was not constructed with New.
func requireLimiter(l *Limiter) {
	if l == nil {
		panic(ErrNilLimiter)
	}
	if l.published == nil {
		panic(ErrUninitializedLimiter)
	}
}
