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

package wait

import "time"

// Option configures optional mechanical behavior for wait primitives.
//
// Option is intentionally small and function-based. It keeps public wait APIs
// stable while allowing narrowly-scoped behavior, such as interval jitter, to be
// added without exposing a mutable configuration struct.
//
// Options are evaluated when a wait primitive starts. Option constructors should
// validate their arguments immediately when that produces clearer failures. A
// nil Option is a programming error and causes the receiving wait primitive to
// panic.
//
// Option does not represent retry policy, backoff policy, scheduler policy,
// observability configuration, or condition semantics. Those concerns belong to
// higher-level packages that build on top of wait.
type Option func(*options)

// options is the private normalized configuration used by wait primitives.
//
// The type is intentionally unexported so public callers cannot depend on its
// fields, construct partially-valid values, or observe internal representation
// changes. Public configuration must go through Option constructors.
type options struct {
	// jitterFactor is the positive one-sided jitter factor applied to fixed
	// interval sleeps between condition evaluations.
	//
	// The zero value disables jitter and preserves the exact base interval.
	// Non-zero values are validated by WithJitter before they are stored here.
	jitterFactor float64
}

// WithJitter returns an option that applies positive one-sided jitter to fixed
// wait intervals.
//
// A factor of 0 disables jitter. A factor greater than 0 allows the effective
// delay to grow by up to factor*base. For example, WithJitter(0.2) allows a
// one-second interval to be delayed by a random value in [0, 200ms], producing
// an effective interval in [1s, 1.2s].
//
// Jitter is applied independently to each interval delay. It is intended only to
// desynchronize otherwise identical loops. It does not implement exponential
// backoff, retry budgeting, fairness, scheduling, or overload-control policy.
//
// If multiple jitter options are supplied to the same wait primitive, the last
// one wins. This mirrors common functional-option behavior and keeps option
// composition deterministic.
//
// WithJitter panics when factor is negative, NaN, or infinite.
func WithJitter(factor float64) Option {
	requireJitterFactor(factor)

	return func(o *options) {
		o.jitterFactor = factor
	}
}

// defaultOptions returns the zero-policy configuration for wait primitives.
//
// The default configuration preserves the baseline behavior of the package:
// exact fixed intervals, no jitter, no retries beyond the owning loop, no
// metrics, and no additional scheduling policy.
func defaultOptions() options {
	return options{}
}

// optionsOf normalizes a caller-supplied option list.
//
// Options are applied in order. When several options configure the same domain,
// the later option overrides the earlier one. This ordering rule makes composed
// option slices predictable without requiring each option domain to invent its
// own merge policy.
func optionsOf(opts ...Option) options {
	config := defaultOptions()
	config.apply(opts...)
	return config
}

// apply mutates o by applying opts in order.
//
// apply is a method so tests and future package code can normalize options
// without duplicating nil-option validation. It is private because callers must
// not mutate normalized wait configuration directly.
func (o *options) apply(opts ...Option) {
	for _, opt := range opts {
		requireOption(opt)
		opt(o)
	}
}

// interval returns the effective delay for one fixed-interval wait step.
//
// When jitter is disabled, interval returns base unchanged. When jitter is
// enabled, interval delegates to Jitter so the package has one implementation of
// duration spreading, validation, rounding, and saturation.
func (o options) interval(base time.Duration) time.Duration {
	if o.jitterFactor == 0 {
		return base
	}

	return Jitter(base, o.jitterFactor)
}
