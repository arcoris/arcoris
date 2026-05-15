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
	"time"
)

// Stamped is a point-in-time value with local update time metadata.
//
// Updated is the local time at which the source committed or published Value. It
// is not a distributed timestamp and does not imply cluster-wide ordering. It is
// also not necessarily the time at which a domain event was originally observed;
// domain packages must model observation time explicitly when they need it.
type Stamped[T any] struct {
	// Revision is the source-local revision of Value.
	Revision Revision

	// Updated is the local commit or publication time of Value.
	Updated time.Time

	// Value is the typed value observed at Revision.
	Value T
}

// IsZeroRevision reports whether the stamped snapshot has ZeroRevision.
func (s Stamped[T]) IsZeroRevision() bool {
	return s.Revision.IsZero()
}

// ChangedSince reports whether the stamped snapshot revision differs from
// revision.
func (s Stamped[T]) ChangedSince(rev Revision) bool {
	return s.Revision.ChangedSince(rev)
}

// Age returns the duration from Updated to now.
//
// Age accepts now explicitly so callers can use their own clock policy and avoid
// hidden reads from the runtime clock.
func (s Stamped[T]) Age(now time.Time) time.Duration {
	return now.Sub(s.Updated)
}

// WithValue returns a stamped snapshot with the same Revision and Updated time
// and a different Value.
//
// The method does not clone either value.
func (s Stamped[T]) WithValue(val T) Stamped[T] {
	return Stamped[T]{
		Revision: s.Revision,
		Updated:  s.Updated,
		Value:    val,
	}
}

// Snapshot drops the Updated timestamp and returns the lightweight snapshot
// representation of s.
func (s Stamped[T]) Snapshot() Snapshot[T] {
	return Snapshot[T]{
		Revision: s.Revision,
		Value:    s.Value,
	}
}
