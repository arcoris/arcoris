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

// Option configures a Holder at construction time.
//
// Options are applied in order by New. Passing nil as an explicit option or
// option argument panics because it is a programming error; omitting an option
// entirely selects the package default.
type Option[T any] func(*config[T])

var (
	// ErrNilClock reports an attempt to configure a nil clock.
	ErrNilClock = errors.New("liveconfig: nil clock")

	// ErrNilClone reports an attempt to configure a nil clone function.
	ErrNilClone = errors.New("liveconfig: nil clone function")

	// ErrNilNormalizer reports an attempt to configure a nil normalizer.
	ErrNilNormalizer = errors.New("liveconfig: nil normalizer")

	// ErrNilValidator reports an attempt to configure a nil validator.
	ErrNilValidator = errors.New("liveconfig: nil validator")

	// ErrNilEqual reports an attempt to configure a nil equality function.
	ErrNilEqual = errors.New("liveconfig: nil equal function")
)

// WithClock configures the local timestamp source used for published snapshots.
//
// The clock is used only when a value is published and Stamped.Updated is set.
// It does not control retries, sleeps, timers, reload loops, or source polling.
func WithClock[T any](clk clock.PassiveClock) Option[T] {
	if clk == nil {
		panic(ErrNilClock)
	}

	return func(cfg *config[T]) {
		cfg.clock = clk
	}
}

// WithClone configures the ownership-boundary clone function.
//
// The clone function is called before normalization, validation, equality, and
// publication. Use it to detach maps, slices, pointers, or other mutable nested
// state from caller-owned input.
func WithClone[T any](clone CloneFunc[T]) Option[T] {
	if clone == nil {
		panic(ErrNilClone)
	}

	return func(cfg *config[T]) {
		cfg.clone = clone
	}
}

// WithNormalizer configures canonicalization for candidate values.
//
// The normalizer may fill defaults or convert equivalent representations into a
// stable form. It runs after cloning and before validation for both New and
// Apply.
func WithNormalizer[T any](normalize Normalizer[T]) Option[T] {
	if normalize == nil {
		panic(ErrNilNormalizer)
	}

	return func(cfg *config[T]) {
		cfg.normalize = normalize
	}
}

// WithValidator configures candidate validation.
//
// The validator checks the final normalized candidate. A validation error
// rejects the candidate and preserves the previous last-good value.
func WithValidator[T any](validate Validator[T]) Option[T] {
	if validate == nil {
		panic(ErrNilValidator)
	}

	return func(cfg *config[T]) {
		cfg.validate = validate
	}
}

// WithEqual configures logical equality for suppressing no-op publications.
//
// The equality function compares the current published value with the prepared
// candidate. When it returns true, Apply succeeds without publishing a new
// revision.
func WithEqual[T any](equal EqualFunc[T]) Option[T] {
	if equal == nil {
		panic(ErrNilEqual)
	}

	return func(cfg *config[T]) {
		cfg.equal = equal
	}
}
