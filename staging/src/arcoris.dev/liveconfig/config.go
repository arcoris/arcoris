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

	"arcoris.dev/chrono/clock"
)

// config contains Holder construction policy.
//
// The configuration is internal so the package can preserve a small public API
// while still keeping clone, normalization, validation, equality, and clock
// policy explicit and testable.
//
// config is built once by New and then treated as immutable. Apply reads it
// under the Holder write mutex, but option functions are never retained or
// called after construction.
type config[T any] struct {
	// clock provides local publication timestamps for snapshot.Stamped values.
	clock clock.PassiveClock

	// clone copies values across the Holder ownership boundary.
	clone CloneFunc[T]

	// normalize converts a cloned candidate into canonical form before validation.
	normalize Normalizer[T]

	// validate checks a normalized candidate before publication.
	validate Validator[T]

	// equal suppresses no-op publications when configured.
	equal EqualFunc[T]
}

// defaultConfig returns Holder defaults.
//
// The default clone is identity, which is correct only for immutable or
// copy-safe configuration values. The default normalizer, validator, and equal
// functions are absent. Without an equal function, each valid Apply publishes a
// new revision.
func defaultConfig[T any]() config[T] {
	return config[T]{
		clock: clock.RealClock{},
		clone: identityClone[T],
	}
}

// ErrNilOption reports a nil Option in a constructor option list.
//
// New treats a nil option as a programmer error because there is no sensible
// default behavior for executing an absent configuration function.
var ErrNilOption = errors.New("liveconfig: nil option")

// newConfig builds Holder configuration from opts.
func newConfig[T any](opts ...Option[T]) config[T] {
	cfg := defaultConfig[T]()
	for _, opt := range opts {
		if opt == nil {
			panic(ErrNilOption)
		}
		opt(&cfg)
	}
	return cfg
}
