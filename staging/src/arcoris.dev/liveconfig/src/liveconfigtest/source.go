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

package liveconfigtest

import (
	"sync"

	"arcoris.dev/snapshot"
)

// ControlledSource is a deterministic test source for revisioned snapshots.
//
// ControlledSource serializes publications and exposes the standard snapshot
// read contracts. It is useful for tests of components that consume a live
// configuration source but should not depend on a full live configuration holder.
//
// ControlledSource does not clone values. Published values must be treated as
// immutable by tests after Publish returns. Use CloneConfig before publishing the
// Config fixture when mutation isolation is part of the test.
//
// The zero value is usable. It starts with no published value and lazily creates
// its publisher on the first publication.
type ControlledSource[T any] struct {
	// mu serializes calls to Publish and PublishStamped so revision order matches
	// publication order even when tests publish from several goroutines. It also
	// protects lazy zero-value publisher creation.
	mu sync.Mutex

	// publisher owns the latest immutable published value.
	publisher *snapshot.Publisher[T]
}

// NewControlledSource creates a source and publishes initial as revision 1.
func NewControlledSource[T any](initial T, opts ...snapshot.Option) *ControlledSource[T] {
	src := NewEmptyControlledSource[T](opts...)
	src.Publish(initial)
	return src
}

// NewEmptyControlledSource creates a source with no published value.
//
// Snapshot returns the zero snapshot until Publish or PublishStamped is called.
func NewEmptyControlledSource[T any](opts ...snapshot.Option) *ControlledSource[T] {
	return &ControlledSource[T]{
		publisher: snapshot.NewPublisher[T](opts...),
	}
}

// Publish publishes val and returns the resulting lightweight snapshot.
func (s *ControlledSource[T]) Publish(val T) snapshot.Snapshot[T] {
	return s.PublishStamped(val).Snapshot()
}

// PublishStamped publishes val and returns the resulting stamped snapshot.
func (s *ControlledSource[T]) PublishStamped(val T) snapshot.Stamped[T] {
	s.mu.Lock()
	defer s.mu.Unlock()

	return s.publisherLocked().PublishStamped(val)
}

// PublishMany publishes vals in order and returns the final snapshot.
//
// When vals is empty, PublishMany returns the current snapshot without
// publishing a new revision.
func (s *ControlledSource[T]) PublishMany(vals ...T) snapshot.Snapshot[T] {
	var snap snapshot.Snapshot[T]
	if len(vals) == 0 {
		return s.Snapshot()
	}
	for _, val := range vals {
		snap = s.Publish(val)
	}
	return snap
}

// Snapshot returns the latest lightweight snapshot.
func (s *ControlledSource[T]) Snapshot() snapshot.Snapshot[T] {
	pub := s.publisherOrNil()
	if pub == nil {
		return snapshot.Snapshot[T]{}
	}
	return pub.Snapshot()
}

// Stamped returns the latest stamped snapshot.
func (s *ControlledSource[T]) Stamped() snapshot.Stamped[T] {
	pub := s.publisherOrNil()
	if pub == nil {
		return snapshot.Stamped[T]{}
	}
	return pub.Stamped()
}

// Revision returns the latest source-local revision.
func (s *ControlledSource[T]) Revision() snapshot.Revision {
	pub := s.publisherOrNil()
	if pub == nil {
		return snapshot.ZeroRevision
	}
	return pub.Revision()
}

// Current returns the latest snapshot value.
//
// Current is a convenience for tests that only care about the value. It has the
// same zero-value behavior as Snapshot().Value.
func (s *ControlledSource[T]) Current() T {
	return s.Snapshot().Value
}

// publisherLocked returns the source publisher, creating one for a zero-value
// ControlledSource. s.mu must be held by the caller.
func (s *ControlledSource[T]) publisherLocked() *snapshot.Publisher[T] {
	if s.publisher == nil {
		s.publisher = snapshot.NewPublisher[T]()
	}
	return s.publisher
}

// publisherOrNil returns the current publisher without creating one.
func (s *ControlledSource[T]) publisherOrNil() *snapshot.Publisher[T] {
	s.mu.Lock()
	defer s.mu.Unlock()

	return s.publisher
}
