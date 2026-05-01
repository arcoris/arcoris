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

// Package clock provides runtime time abstractions for ARCORIS component
// internals.
//
// The package exists to make time-dependent runtime code explicit, testable, and
// independent from direct calls to package time. ARCORIS components use time for
// admission checks, queue wake-ups, retry delays, dispatch timeouts, worker
// heartbeats, lease monitoring, controller loops, sampling cadence, cooldowns,
// and deterministic tests.
//
// clock is part of component-base. It is a low-level runtime utility, not API
// machinery and not a scheduling policy package.
//
// # Scope
//
// The package defines a small set of runtime-time contracts and implementations:
//
//   - PassiveClock for read-only time access;
//   - Clock for time access plus waiting primitives;
//   - Timer for one-shot waits;
//   - Ticker for periodic wake-ups;
//   - RealClock for production code backed by the Go standard library;
//   - FakeClock for deterministic tests.
//
// The package owns only the mechanics of reading time and waiting for time to
// pass. Higher-level packages own the meaning of time in their own domains.
//
// For example:
//
//   - admission packages decide whether a request deadline has expired;
//   - retry packages decide how long a backoff should be;
//   - queue packages decide when delayed work becomes eligible;
//   - worker packages decide how heartbeats are interpreted;
//   - lease packages decide ownership and fencing rules;
//   - controller packages decide sampling cadence and cooldown policy.
//
// clock supplies the time source used by those packages. It does not encode
// those policies itself.
//
// # PassiveClock versus Clock
//
// Components should depend on the smallest clock contract they need.
//
// PassiveClock is for components that only need to read time:
//
//   - metadata writers;
//   - admission evaluators;
//   - queue eligibility checks;
//   - lazy rate/window calculations;
//   - diagnostics and sampling code.
//
// Clock is for components that own time-based runtime behavior:
//
//   - scheduler loops;
//   - adaptive controller loops;
//   - worker heartbeat loops;
//   - retry executors;
//   - dispatch timeout monitors;
//   - lease reapers;
//   - background reconciliation loops.
//
// A component that only needs Now or Since should accept PassiveClock instead of
// Clock. This keeps read-only code independent from timers, tickers, sleeps, and
// runtime loop ownership.
//
// # Wall time and elapsed time
//
// time.Time values can serve different purposes depending on where they are
// used.
//
// Wall-clock timestamps are suitable for externally visible event time:
//
//   - created-at timestamps;
//   - updated-at timestamps;
//   - admitted-at timestamps;
//   - started-at timestamps;
//   - completed-at timestamps;
//   - audit and diagnostic event timestamps.
//
// Elapsed-time measurements are suitable for process-local runtime decisions:
//
//   - latency;
//   - timeout age;
//   - retry delay;
//   - queue wait duration;
//   - heartbeat freshness;
//   - controller cooldown;
//   - sampling interval.
//
// In Go, time values returned by time.Now may contain process-local monotonic
// information. RealClock preserves the standard library behavior by delegating
// Now to time.Now and Since to time.Since. That monotonic information is useful
// for local elapsed-time measurements while the value remains in the same
// process.
//
// Callers must not rely on monotonic elapsed-time semantics after a time.Time has
// been serialized, persisted, transmitted over the network, loaded from an API
// object, or reconstructed from text. Serialized timestamps are wall-clock
// values, not process-local monotonic instants.
//
// # RealClock
//
// RealClock is the production implementation of Clock.
//
// RealClock is stateless, zero-value usable, and backed by the Go standard
// library. It adapts package time to the clock interfaces without adding
// scheduling, admission, retry, lease, controller, metrics, or distributed-time
// policy.
//
// RealClock reads the local process clock. It does not provide cluster-wide time
// synchronization and must not be treated as a distributed ordering source.
//
// # FakeClock
//
// FakeClock is a deterministic Clock implementation for tests.
//
// FakeClock does not observe real runtime time. Its current time changes only
// when test code explicitly advances it with Set or Step. Timers, tickers,
// After, and Sleep created by a FakeClock are released only by fake-time
// advancement.
//
// FakeClock is intended for tests of time-dependent runtime components without
// real sleeps or timing races. It is useful for testing:
//
//   - delayed queue eligibility;
//   - retry and backoff waits;
//   - controller ticks;
//   - worker heartbeat loops;
//   - timeout behavior;
//   - cooldown behavior;
//   - lease reaper loops;
//   - scheduler wake-up behavior.
//
// FakeClock is monotonic. Moving fake time backwards is rejected because backward
// movement makes already-fired waiters, timers, and tickers ambiguous.
//
// FakeClock is safe for concurrent use by tests and by the component under test,
// but it must not be copied after first use.
//
// # Timers and tickers
//
// Timer is a one-shot waiting primitive. Use Timer when a component needs
// explicit lifecycle ownership over a wait, such as Stop or Reset.
//
// Ticker is a periodic wake-up primitive. Use Ticker for repeated runtime loops
// such as scheduler ticks, controller sampling loops, heartbeat loops, and lease
// reaper loops.
//
// Timers and tickers belong to the Clock that created them. A Timer or Ticker
// created by RealClock is driven by real runtime time. A Timer or Ticker created
// by FakeClock is driven only by explicit fake-time advancement.
//
// Timer and Ticker channels are not closed by Stop. Callers must not use channel
// close as a lifecycle signal.
//
// Fake timer and ticker delivery is non-blocking. If a fake timer or ticker
// channel already contains an unread value, a later fake delivery may be
// dropped. Tests that reuse timers or tickers and need to observe every firing
// must coordinate receives before Reset or further fake-time advancement.
//
// # Fake ticker delivery
//
// Fake tickers intentionally deliver at most one tick per fake-clock advancement
// operation.
//
// If Set or Step moves fake time across several ticker intervals, the ticker
// delivers one tick and is rescheduled relative to the new fake time. This avoids
// unbounded bursts of missed ticks and gives controller-loop tests deterministic
// one-wakeup-per-advance behavior.
//
// This is a fake-clock testing policy. It is not a distributed scheduling
// guarantee and must not be interpreted as production runtime behavior.
//
// # Distributed systems boundaries
//
// clock provides physical/runtime time abstractions. It does not provide
// distributed causality or replicated-state ordering.
//
// The following concepts do not belong in this package:
//
//   - resource versions;
//   - generations;
//   - observed generations;
//   - lease epochs;
//   - fencing tokens;
//   - Lamport clocks;
//   - vector clocks;
//   - hybrid logical clocks;
//   - cluster time oracle;
//   - clock-skew monitoring.
//
// Those concepts should be modeled by API machinery, runtime ownership
// protocols, distributed-state packages, or future logical-time packages when
// the system requires them.
//
// # File ownership
//
// The package intentionally separates files by responsibility.
//
// Public contracts live in:
//
//   - clock.go;
//   - timer.go;
//   - ticker.go.
//
// Production standard-library adapters live in:
//
//   - real_clock.go;
//   - real_timer.go;
//   - real_ticker.go.
//
// Deterministic fake-clock implementation lives in:
//
//   - fake_clock.go;
//   - fake_time.go;
//   - fake_waiter.go;
//   - fake_timer.go;
//   - fake_ticker.go.
//
// The ownership rules are:
//
//   - clock.go defines PassiveClock and Clock only;
//   - timer.go defines the Timer contract only;
//   - ticker.go defines the Ticker contract only;
//   - real_clock.go defines RealClock and factory methods only;
//   - real_timer.go defines the time.Timer adapter only;
//   - real_ticker.go defines the time.Ticker adapter only;
//   - fake_clock.go defines FakeClock state and construction only;
//   - fake_time.go owns fake time reads and advancement;
//   - fake_waiter.go owns After, Sleep, HasWaiters, and waiter lifecycle;
//   - fake_timer.go owns fake Timer lifecycle;
//   - fake_ticker.go owns fake Ticker lifecycle.
//
// Files should not absorb neighboring responsibilities merely for convenience.
// In particular, fake timer and fake ticker lifecycle logic should not be folded
// into fake_clock.go, and real timer or ticker adapter methods should not be
// folded into real_clock.go.
//
// # Non-goals
//
// This package is not a general timeutils package.
//
// It must not define:
//
//   - retry policy;
//   - backoff algorithms;
//   - rate limiting;
//   - token buckets;
//   - window calculations;
//   - EWMA calculations;
//   - scheduler policy;
//   - queue policy;
//   - admission policies;
//   - lease ownership protocols;
//   - API timestamp wrappers;
//   - JSON serialization helpers;
//   - metrics/exporters;
//   - observability instruments.
//
// Higher-level packages may build those concepts using clock, but clock itself
// must remain a small runtime-time abstraction layer.
//
// # Dependency policy
//
// Production code in this package should depend only on the Go standard library.
//
// clock must not depend on apimachinery, scheduler packages, queue packages,
// worker packages, observability exporters, metrics packages, or external test
// assertion libraries.
package clock
