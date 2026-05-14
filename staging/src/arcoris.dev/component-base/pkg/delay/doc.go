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

// Package delay provides reusable delay sequence contracts and deterministic
// delay schedules for ARCORIS component internals.
//
// The package owns Schedule and Sequence. A Schedule is an immutable reusable
// recipe for creating per-owner delay streams. A Sequence is a single-owner
// stream of concrete time.Duration values. Package constructors return schedule
// values that can be stored in component configuration and reused across many
// executions, while each NewSequence call creates the mutable stream owned by one
// retry loop, polling loop, reconnect loop, cooldown path, or test.
//
// Sequence exhaustion is finite and explicit. Next returns ok=false when no
// further delay value is available, and callers must ignore the accompanying
// duration in that case. A result of delay=0, ok=true is different: it means an
// available delay whose value is immediate continuation. A negative delay with
// ok=true violates the Sequence contract. Package adapters and wrappers panic at
// that boundary instead of silently passing an invalid runtime duration to
// clocks, timers, retry orchestration, or waiting code.
//
// The package includes deterministic schedules such as Immediate, Fixed,
// Delays, Linear, Exponential, and Fibonacci, plus deterministic wrappers such
// as Cap and Limit. Arithmetic used by deterministic growth schedules saturates
// at the largest representable time.Duration instead of wrapping into negative
// values.
//
// Package delay does not sleep, create timers or tickers, observe contexts,
// execute operations, classify errors, retry failed work, randomize delays,
// export metrics, log, trace, or make scheduler, admission, retry, lifecycle,
// health, or domain decisions. It only describes and transforms streams of
// non-negative duration values.
//
// Randomized delay transforms and randomized schedules belong to package jitter.
// Runtime waiting belongs to package wait. Retry orchestration belongs to
// package retry. Clock and fake-time behavior belong to package clock.
package delay
