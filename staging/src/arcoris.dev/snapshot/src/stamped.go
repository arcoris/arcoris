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

package snapshot

import "time"

// Stamped is a point-in-time value with local update time metadata.
//
// R is the authoritative revision type for the source. Updated is local
// publication metadata. It is not a distributed timestamp and does not imply
// cluster-wide ordering or domain event time.
type Stamped[R comparable, T any] struct {
	// Revision is the source or domain revision at which Value was observed.
	Revision R

	// Updated is the local commit or publication time of Value.
	Updated time.Time

	// Value is the typed read-model value observed at Revision.
	Value T
}

// ChangedSince reports whether the stamped snapshot revision differs from prev.
//
// The check is equality-based for the same reason as Snapshot.ChangedSince.
func (s Stamped[R, T]) ChangedSince(prev R) bool {
	return s.Revision != prev
}

// Age returns the duration from Updated to now.
//
// Age accepts now explicitly so callers can use their own clock policy and avoid
// hidden reads from the runtime clock.
func (s Stamped[R, T]) Age(now time.Time) time.Duration {
	return now.Sub(s.Updated)
}

// Snapshot drops Updated and returns the lightweight snapshot representation.
func (s Stamped[R, T]) Snapshot() Snapshot[R, T] {
	return Snapshot[R, T]{
		Revision: s.Revision,
		Value:    s.Value,
	}
}

// WithValue returns a stamped snapshot with the same Revision and Updated time
// and a different Value.
//
// WithValue does not clone either value.
func (s Stamped[R, T]) WithValue(value T) Stamped[R, T] {
	return Stamped[R, T]{
		Revision: s.Revision,
		Updated:  s.Updated,
		Value:    value,
	}
}
