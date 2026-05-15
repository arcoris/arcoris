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

package retrybudget

import "time"

// WindowSnapshot describes the time window represented by a retry-budget
// snapshot.
type WindowSnapshot struct {
	// StartedAt is the inclusive start of a bounded window.
	StartedAt time.Time

	// EndsAt is the exclusive end of a bounded window.
	EndsAt time.Time

	// Duration is the bounded window duration.
	Duration time.Duration

	// Bounded reports whether StartedAt, EndsAt, and Duration are meaningful.
	Bounded bool
}

// IsValid reports whether s is internally consistent.
func (s WindowSnapshot) IsValid() bool {
	if !s.Bounded {
		return s.StartedAt.IsZero() && s.EndsAt.IsZero() && s.Duration == 0
	}
	if s.Duration <= 0 {
		return false
	}
	if s.StartedAt.IsZero() || s.EndsAt.IsZero() {
		return false
	}
	if !s.EndsAt.Equal(s.StartedAt.Add(s.Duration)) {
		return false
	}
	return s.EndsAt.After(s.StartedAt)
}

// IsBounded reports whether s describes a finite retry-budget window.
func (s WindowSnapshot) IsBounded() bool {
	return s.Bounded
}

// Contains reports whether t is inside the bounded window.
//
// Bounded windows are start-inclusive and end-exclusive. Unbounded windows never
// contain a specific time instant because they do not publish a finite interval.
func (s WindowSnapshot) Contains(t time.Time) bool {
	if !s.Bounded || !s.IsValid() {
		return false
	}
	return !t.Before(s.StartedAt) && t.Before(s.EndsAt)
}
