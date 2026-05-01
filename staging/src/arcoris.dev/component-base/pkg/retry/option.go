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

package retry

import (
	"time"

	"arcoris.dev/component-base/pkg/backoff"
	"arcoris.dev/component-base/pkg/clock"
)

// Option configures one retry execution.
//
// Options are applied to a private config before Do or DoValue starts executing
// an operation. They do not mutate a retry execution after it has started and
// they are not retained as option values after configuration normalization.
//
// Option is intentionally function-based. This keeps the public construction API
// stable while allowing narrowly-scoped configuration domains to evolve without
// exposing a mutable public configuration struct.
//
// A nil Option is a programming error. Retry rejects nil options instead of
// silently ignoring them so invalid conditional option composition is visible at
// the configuration boundary.
type Option func(*config)

// config contains normalized retry execution settings.
//
// config is intentionally package-local. Public callers configure retry through
// Option constructors, while the runtime loop receives a complete normalized
// configuration.
//
// config must remain limited to retry-owned mechanics:
//
//   - clock dependency;
//   - backoff schedule;
//   - retryability classifier;
//   - retry-owned limits;
//   - observer list.
//
// It must not contain operation business logic, protocol-specific retry policy,
// storage-specific retry policy, controller reconciliation behavior, circuit
// breaker state, retry budget state, metrics exporters, tracing exporters, or
// logging backends directly.
type config struct {
	// clock provides retry execution time.
	//
	// Retry uses the clock for attempt timestamps, terminal outcome timestamps,
	// elapsed-time checks, and delay timers. The default is clock.RealClock{}.
	clock clock.Clock

	// backoff is the reusable delay schedule used by retry execution.
	//
	// A fresh Sequence is created from this Schedule for each Do or DoValue call.
	// The retry package must not store or share a mutable backoff.Sequence in
	// config.
	backoff backoff.Schedule

	// classifier decides whether an operation-owned error may be retried.
	//
	// The default is NeverRetry. This keeps generic retry conservative: callers
	// must explicitly opt into retrying operation failures.
	classifier Classifier

	// maxAttempts is the maximum number of operation calls allowed.
	//
	// The value includes the initial operation call. A value of one means no retry
	// attempts beyond the initial call. A value of zero is invalid.
	maxAttempts uint

	// maxElapsed limits the total elapsed runtime of one retry execution.
	//
	// A zero value disables elapsed-time limiting. Negative values are invalid.
	maxElapsed time.Duration

	// observers are notified synchronously about retry events.
	//
	// Observers are called in registration order. The slice is owned by config
	// normalization and must not alias caller-owned slices.
	observers []Observer
}

// defaultConfig returns the conservative retry configuration.
//
// The default configuration performs exactly one operation attempt, does not
// retry operation errors, uses an immediate backoff schedule that is normally not
// consumed under the default attempt limit, uses real runtime time, has no
// elapsed-time limit, and registers no observers.
func defaultConfig() config {
	return config{
		clock:       clock.RealClock{},
		backoff:     backoff.Immediate(),
		classifier:  NeverRetry(),
		maxAttempts: 1,
		maxElapsed:  0,
		observers:   nil,
	}
}

// configOf returns normalized retry configuration for opts.
//
// Options are applied in order. When several options configure a single-value
// domain, the later option wins. Observer options append observers in order.
// This mirrors ordinary functional-option behavior and makes composed option
// lists deterministic.
func configOf(opts ...Option) config {
	config := defaultConfig()
	config.apply(opts...)
	return config
}

// apply mutates c by applying opts in order.
//
// apply is private because callers must not mutate normalized retry
// configuration directly. Nil options are rejected immediately through
// requireOption.
func (c *config) apply(opts ...Option) {
	for _, opt := range opts {
		requireOption(opt)
		opt(c)
	}
}
