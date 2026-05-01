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

// Package backoff provides reusable recipes for generating delay sequences.
//
// The package is intentionally narrow. It computes time.Duration values and
// leaves all runtime meaning to callers. It does not sleep, create timers,
// observe context cancellation, execute operations, classify errors, retry
// failed work, log, trace, export metrics, rate limit callers, schedule queue
// items, or make domain decisions.
//
// A typical owner wires this package like this:
//
//	schedule := backoff.FullJitter(
//		backoff.Cap(
//			backoff.Exponential(100*time.Millisecond, 2.0),
//			5*time.Second,
//		),
//	)
//	sequence := schedule.NewSequence()
//
//	delay, ok := sequence.Next()
//	// The owner decides whether to sleep, retry, stop, record metrics,
//	// check context cancellation, or apply any other runtime policy.
//	_ = delay
//	_ = ok
//
// # Schedule and Sequence
//
// Schedule is the immutable reusable recipe. It is suitable for storage in a
// policy object, controller, client configuration, or component option. A
// Schedule returned by this package does not contain per-owner cursor state,
// attempt state, previous-delay state, or mutable pseudo-random generator state.
//
// Sequence is the mutable single-owner stream created from a Schedule. A
// Sequence may hold an index, a remaining count, previous delay state, or a
// pseudo-random generator. For that reason, Sequence values are not required to
// be safe for concurrent calls to Next unless a concrete implementation
// explicitly documents stronger guarantees.
//
// The contract is:
//
//	Schedule.NewSequence() -> non-nil Sequence
//	Sequence.Next()       -> non-negative delay with ok=true, or ok=false
//
// Finite exhaustion is represented by ok=false. Exhaustion is not an error in
// this package. A retry package may interpret exhaustion as retry exhaustion; a
// polling owner may interpret it as shutdown; a test may interpret it as an
// expected boundary.
//
// # Constructor validation
//
// Invalid constructor input is a programming error. Constructors and package
// adapters panic with stable package-local diagnostic strings instead of
// returning errors. This keeps runtime retry, polling, reconnect, and cooldown
// loops simple and makes invalid policy construction fail at the boundary where
// the schedule is created.
//
// Examples of invalid input include negative durations, nil child schedules,
// nil random options, nil random sources, non-finite multipliers, and closed
// ranges whose upper bound is smaller than their lower bound.
//
// # Duration arithmetic
//
// Duration arithmetic is centralized in duration.go. Algorithms saturate at the
// largest representable time.Duration instead of wrapping into negative values.
// This is important because negative available delays violate the Sequence
// contract and have no meaningful runtime interpretation for timers, clocks,
// waits, retries, polling loops, reconnect loops, or controller cooldowns.
//
// For example, a long-running exponential schedule eventually returns
// maxDuration repeatedly instead of overflowing:
//
//	sequence := backoff.Exponential(time.Second, 2).NewSequence()
//	// ... after enough calls, later values saturate instead of wrapping.
//
// # Randomness
//
// Random and jitter schedules use non-cryptographic pseudo-randomness only for
// desynchronization and deterministic tests. Randomness in this package MUST NOT
// be used for secrets, authentication, authorization, nonce generation, access
// control, or any other security-sensitive purpose.
//
// RandomSource separates schedule reuse from per-sequence random state. A
// reusable schedule stores a RandomSource. Each Sequence receives a generator
// from that source. WithSeed intentionally creates a fresh rand.Rand for every
// sequence so independent sequences from one schedule are deterministic without
// sharing mutable generator state.
//
// Example:
//
//	schedule := backoff.Random(time.Second, 5*time.Second, backoff.WithSeed(42))
//	left := schedule.NewSequence()
//	right := schedule.NewSequence()
//	// left and right produce the same deterministic stream, but do not share
//	// one mutable random generator.
//
// # Composition
//
// Schedules are designed to compose mechanically. Limit makes an infinite
// schedule finite. Cap bounds available child delays. Jitter wrappers transform
// available child delays while preserving child exhaustion.
//
// Composition order is significant. For example, Cap(FullJitter(s), max) caps
// the final jittered value, while FullJitter(Cap(s, max)) randomizes after the
// child has already been capped. Prefer an outer Cap when maxDelay must be a
// hard final upper bound.
//
// # Non-goals
//
// This package deliberately does not define retry policy, retryable error
// classification, context handling, sleeps, fake clocks, wait loops, logging,
// tracing, metrics, admission control, scheduling, or overload policy. Those
// belong to higher-level packages that can interpret delay values in their own
// domain.
package backoff
