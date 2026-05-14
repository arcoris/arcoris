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

// Package jitter provides non-cryptographic delay randomization primitives.
//
// The package owns randomized wrappers over delay.Schedule, randomized delay
// schedules, and the RandomSource and RandomGenerator abstractions used to
// create per-sequence pseudo-random generators. It is intended to desynchronize
// retry, polling, reconnect, cooldown, and other owner-controlled delay streams.
// Randomness in this package is for load spreading and deterministic tests; it
// must not be used for secrets, authentication, authorization, nonce generation,
// access control, or any other security-sensitive purpose.
//
// Wrappers such as Full, Equal, Positive, and Proportional transform values
// produced by an existing delay.Schedule and preserve child sequence exhaustion.
// Uniform and Decorrelated create randomized schedules directly because they do
// not transform an existing deterministic child stream. In both cases, the
// reusable schedule stores a RandomSource, and each created sequence owns the
// RandomGenerator returned by that source unless the caller explicitly supplies
// a source that shares generator state.
//
// Package jitter does not sleep, create timers, observe contexts, execute
// operations, classify errors, retry failed work, rate limit callers, export
// metrics, log, trace, or make scheduler, admission, lifecycle, or domain
// decisions. It also does not own deterministic growth schedules or generic
// Schedule and Sequence contracts.
//
// Generic and deterministic delay schedules belong to package delay. Runtime
// waiting belongs to package wait. Retry orchestration belongs to package retry.
package jitter
