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
