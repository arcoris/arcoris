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
type Option func(*config)

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

	return func(cfg *config) {
		cfg.jitterFactor = factor
	}
}
