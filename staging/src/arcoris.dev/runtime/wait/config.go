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

// config is the private normalized configuration used by wait primitives.
//
// The type is intentionally unexported so public callers cannot depend on its
// fields, construct partially-valid values, or observe internal representation
// changes. Public configuration must go through Option constructors.
type config struct {
	// jitterFactor is the positive one-sided jitter factor applied to fixed
	// interval sleeps between condition evaluations.
	//
	// The zero value disables jitter and preserves the exact base interval.
	// Non-zero values are validated by WithJitter before they are stored here.
	jitterFactor float64
}

// defaultOptions returns the zero-policy configuration for wait primitives.
//
// The default configuration preserves the baseline behavior of the package:
// exact fixed intervals, no jitter, no retries beyond the owning loop, no
// metrics, and no additional scheduling policy.
func defaultOptions() config {
	return config{}
}

// optionsOf normalizes a caller-supplied option list.
//
// Options are applied in order. When several options configure the same domain,
// the later option overrides the earlier one. This ordering rule makes composed
// option slices predictable without requiring each option domain to invent its
// own merge policy.
func optionsOf(opts ...Option) config {
	cfg := defaultOptions()
	cfg.apply(opts...)
	return cfg
}

// apply mutates cfg by applying opts in order.
//
// apply is a method so tests and future package code can normalize options
// without duplicating nil-option validation. It is private because callers must
// not mutate normalized wait configuration directly.
func (cfg *config) apply(opts ...Option) {
	for _, opt := range opts {
		requireOption(opt)
		opt(cfg)
	}
}

// interval returns the effective delay for one fixed-interval wait step.
//
// When jitter is disabled, interval returns base unchanged. When jitter is
// enabled, interval delegates to Jitter so the package has one implementation of
// duration spreading, validation, rounding, and saturation.
func (cfg config) interval(base time.Duration) time.Duration {
	if cfg.jitterFactor == 0 {
		return base
	}

	return Jitter(base, cfg.jitterFactor)
}
