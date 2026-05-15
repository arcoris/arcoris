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

package snapshot

import (
	"sync"
	"time"

	"arcoris.dev/atomicx"
	"arcoris.dev/chrono/clock"
)

// Publisher atomically publishes immutable copy-on-write values.
//
// Publisher is the fast baseline for read-mostly values such as observer lists,
// handler registries, immutable routing tables, and precomputed dispatch plans.
// Reads use an atomic pointer load and do not clone the value.
//
// Writes are serialized so each commit advances the source-local revision and
// stores the corresponding record as one ordered publication. Reads remain
// lock-free by loading the latest published record pointer.
//
// Values passed to Publish or PublishStamped must be treated as immutable after
// publication. If T contains slices, maps, pointers, or other mutable state, the
// caller must build a fresh copy before publishing it and must not mutate that
// copy after Publish returns.
//
// Publisher is safe for concurrent use and zero-value usable. A zero-value
// Publisher has no published record and returns a zero snapshot until the first
// Publish. Publisher must not be copied after first use.
type Publisher[T any] struct {
	// noCopy prevents accidental copies after first use.
	noCopy noCopy

	// mu serializes Publish and PublishStamped commits.
	//
	// The mutex preserves monotonic visible revisions under concurrent
	// publishers. Without it, one goroutine could reserve a newer revision and
	// store it before another goroutine stores an older reserved revision.
	mu sync.Mutex

	// nextRevision is the latest source-local revision assigned to a published
	// record.
	//
	// nextRevision is protected by mu. Readers do not read this field; they load
	// the revision from the latest immutable record through ptr.
	nextRevision Revision

	// clock provides local publication timestamps for Stamped values.
	//
	// A nil clock means the zero-value Publisher should lazily use RealClock.
	clock clock.PassiveClock

	// ptr is a padded atomic pointer to the latest immutable published record.
	//
	// The padding isolates the hot pointer cell from neighboring fields. It does
	// not provide publication ordering; write-side ordering is provided by mu.
	ptr atomicx.PaddedPointer[record[T]]
}

// record is an immutable published value.
//
// Once a record is stored in Publisher.ptr, it must never be mutated. Readers may
// load and use records without locks.
type record[T any] struct {
	// revision is the source-local revision assigned at publication time.
	revision Revision

	// updated is the local publication time.
	updated time.Time

	// value is the immutable published value.
	value T
}

// NewPublisher creates a Publisher configured with opts.
//
// The returned Publisher has no published value. Snapshot returns a zero snapshot
// until Publish is called.
func NewPublisher[T any](opts ...Option) *Publisher[T] {
	cfg := newConfig(opts...)
	return &Publisher[T]{
		clock: cfg.clock,
	}
}

// Snapshot returns the latest lightweight snapshot.
//
// Snapshot is lock-free. If no value has been published, Snapshot returns the
// zero Snapshot[T]. The returned value is not cloned; it is the immutable value
// stored in the latest published record.
func (p *Publisher[T]) Snapshot() Snapshot[T] {
	stamped := p.Stamped()
	return stamped.Snapshot()
}

// Stamped returns the latest stamped snapshot.
//
// Stamped is lock-free. If no value has been published, Stamped returns the zero
// Stamped[T]. The returned value is not cloned.
func (p *Publisher[T]) Stamped() Stamped[T] {
	rec := p.ptr.Load()
	if rec == nil {
		return Stamped[T]{}
	}

	return Stamped[T]{
		Revision: rec.revision,
		Updated:  rec.updated,
		Value:    rec.value,
	}
}

// Revision returns the revision of the latest visible published record.
//
// Revision returns ZeroRevision before the first Publish.
func (p *Publisher[T]) Revision() Revision {
	rec := p.ptr.Load()
	if rec == nil {
		return ZeroRevision
	}

	return rec.revision
}

// Publish publishes next and returns the resulting lightweight snapshot.
//
// Publish does not clone next. Callers must not mutate next after publication.
func (p *Publisher[T]) Publish(next T) Snapshot[T] {
	return p.PublishStamped(next).Snapshot()
}

// PublishStamped publishes next and returns the resulting stamped snapshot.
//
// PublishStamped assigns a fresh source-local revision, records the local
// publication time using the configured PassiveClock, stores a new immutable
// record, and returns the published stamped snapshot. Concurrent publish calls
// are serialized so readers cannot observe revision rollback. If revision
// overflow is detected, PublishStamped panics before storing a new record.
func (p *Publisher[T]) PublishStamped(next T) Stamped[T] {
	p.mu.Lock()
	defer p.mu.Unlock()

	rev := p.nextRevision.Next()
	updated := p.passiveClock().Now()
	rec := &record[T]{
		revision: rev,
		updated:  updated,
		value:    next,
	}

	p.nextRevision = rev
	p.ptr.Store(rec)

	return Stamped[T]{
		Revision: rev,
		Updated:  updated,
		Value:    next,
	}
}

// passiveClock returns the Publisher clock.
//
// A zero-value Publisher has no configured clock, so it lazily falls back to
// clock.RealClock. NewPublisher should be used when deterministic timestamps are
// required in tests.
func (p *Publisher[T]) passiveClock() clock.PassiveClock {
	if p.clock != nil {
		return p.clock
	}

	return clock.RealClock{}
}
