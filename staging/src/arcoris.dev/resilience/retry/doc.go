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

// Package retry provides bounded retry orchestration for ARCORIS component
// internals.
//
// The package executes caller-owned operations and decides whether another
// attempt may be scheduled after an operation-owned error. It combines operation
// execution, retryability classification, retry-owned limits, delay sequence
// consumption, clock-backed waiting, retry-owned context interruption, terminal
// outcome metadata, and observer events.
//
// Package retry is an orchestration layer. It does not define delay formulas,
// protocol retry rules, storage retry rules, controller reconciliation policy,
// retry budgets, circuit breakers, hedged requests, rate limiters, queue
// scheduling, logging backends, metrics exporters, tracing exporters, goroutine
// supervision, or component lifecycle transitions.
//
// # Execution model
//
// A retry execution is started by Do or DoValue. Do executes an Operation that
// returns only an error. DoValue executes a ValueOperation that returns a value
// and an error.
//
// A retry execution may call the operation zero, one, or many times:
//
//   - zero times when the retry-owned context is already stopped before the first
//     attempt starts;
//   - one time when the first attempt succeeds, fails with a non-retryable error,
//     or the configured attempt limit allows only the initial attempt;
//   - many times when operation errors are classified as retryable and
//     retry-owned boundaries still allow another attempt.
//
// Attempt numbering is one-based. Attempt 1 is the initial operation call.
// Attempts 2 and later are retry operation calls. Attempt metadata is represented
// by Attempt and is intentionally separate from operation errors, retry delays,
// observer events, and terminal outcomes.
//
// The retry package does not infer idempotency, replay safety, transaction
// safety, protocol semantics, or external side-effect safety. Callers MUST
// configure retries only for operations whose retry behavior is valid for their
// domain.
//
// # Delay integration
//
// Package retry uses package delay for delay streams.
//
// Retry configuration stores a delay.Schedule. A fresh delay.Sequence is
// created for each Do or DoValue execution. The sequence is owned by that single
// retry execution and is not shared across executions.
//
// The retry loop consumes a delay sequence only after an operation attempt
// fails with an error that the configured Classifier accepts as retryable and
// after retry-owned attempt limits still allow another attempt.
//
// When Sequence.Next returns ok=false, the delay package has not failed; it has
// reported finite sequence exhaustion. At the retry layer, that means a
// retryable operation failure cannot be followed by another scheduled attempt.
// Retry reports this as retry-owned exhaustion with
// StopReasonDelayExhausted and ErrExhausted.
//
// A zero delay is valid and means retry may proceed to the next attempt
// immediately after retry-owned context checks. A negative delay with ok=true is
// a programming error because it violates the delay.Sequence contract.
//
// # Clock integration
//
// Package retry uses package clock for runtime time.
//
// The configured clock is used for attempt timestamps, terminal outcome
// timestamps, elapsed-time checks, and retry delay timers. This keeps retry
// execution deterministic in tests when callers use a fake clock and avoids
// direct dependency on time.Sleep or time.Timer in the retry loop.
//
// Package retry does not use package wait for delay execution. Package wait owns
// lower-level waiting helpers; retry owns higher-level operation retry
// orchestration and requires clock-backed delay behavior for deterministic tests.
//
// # Deadline integration
//
// Package retry uses package deadline to avoid scheduling retry delays that
// would consume or exceed the owning context deadline budget.
//
// After a retryable operation failure, retry checks the selected delay against
// the remaining context deadline budget. A delay equal to or greater than the
// remaining budget stops execution before retry emits a delay event or sleeps,
// because it would leave no usable budget for the next attempt.
//
// This boundary is reported as ErrExhausted with StopReasonDeadline. It is
// distinct from ErrInterrupted, which is returned only when retry observes that
// the context has already stopped or stops while retry is waiting at a
// retry-owned boundary.
//
// Deadline checks use the configured retry clock as the observation time.
// Callers that use fake clocks in tests should create context deadlines from
// the same time base as the configured retry clock. Production code should
// normally use the standard runtime time domain for both context deadlines and
// the retry clock.
//
// # Classification model
//
// Classifier decides whether an operation-owned error may be retried. The default
// classifier is NeverRetry, making the default configuration conservative.
//
// A nil error is never retryable. Nil means the operation succeeded, not that a
// retry should be scheduled.
//
// Package retry provides generic classifiers only. Protocol-specific and
// domain-specific retry rules, such as HTTP status code handling, gRPC status
// handling, database serialization failures, storage conflicts, or controller
// reconciliation conflicts, belong in adapter packages or caller-owned
// Classifier implementations.
//
// # Limits and exhaustion
//
// Retry execution is bounded by retry-owned limits:
//
//   - max attempts;
//   - max elapsed runtime;
//   - owning context deadline budget before retry delay scheduling;
//   - finite delay sequence exhaustion.
//
// Max attempts includes the initial operation call. WithMaxAttempts(1) allows
// only the initial attempt and no retry attempts. WithMaxAttempts(0) is invalid.
//
// A zero max elapsed duration disables elapsed-time limiting. A positive max
// elapsed duration bounds one retry execution according to the configured clock.
// Retry should stop before sleeping for a delay that would consume or exceed the
// remaining elapsed-time budget.
//
// Retry also checks the owning context deadline budget before sleeping for a
// retry delay. A selected delay equal to or greater than the remaining deadline
// budget stops execution with ErrExhausted and StopReasonDeadline.
//
// Retry-owned exhaustion is classified with ErrExhausted. Exhausted errors carry
// an Outcome and unwrap to the last operation-owned error when one is available.
//
// Non-retryable operation errors are not retry-owned exhaustion. They remain
// operation-owned results and are normally returned unchanged.
//
// # Context ownership
//
// Retry observes the context passed to Do or DoValue at retry-owned boundaries,
// such as before an attempt and while waiting between attempts.
//
// When retry observes that its owning context is stopped at a retry boundary, it
// returns ErrInterrupted and preserves the underlying context sentinel or custom
// context cause.
//
// Raw context errors returned by operations remain operation-owned. Context
// stops observed by retry are ErrInterrupted. Context deadline budget rejection
// before a retry delay is ErrExhausted with StopReasonDeadline.
//
// Raw context.Canceled or context.DeadlineExceeded values returned by an
// operation are not automatically classified as retry-owned interruption.
// Operations can return context errors for operation-owned work, and retry cannot
// infer ownership unless retry itself observed the context stop.
//
// # Outcome model
//
// Outcome records terminal retry metadata: number of operation calls, start time,
// finish time, last operation-owned error, and StopReason.
//
// Outcome is metadata, not an error wrapper and not a retry decision. It is used
// by terminal observer events and retry-owned errors such as ErrExhausted.
//
// The zero Outcome is invalid. This prevents empty completion metadata from
// accidentally looking like a successful retry execution.
//
// # Observer model
//
// Observers receive Event values describing retry execution progress. Events are
// observer-facing metadata for attempt start, attempt failure, retry delay, and
// retry stop.
//
// Observers are notification boundaries only. They do not decide retryability,
// choose delays, change limits, mutate retry configuration, or affect operation
// execution.
//
// Observers are called synchronously in registration order. Observer failures are
// not represented in retry's error model because Observer does not return error.
// The retry package does not recover observer panics.
//
// Observer implementations that are shared across concurrent retry executions
// must provide their own synchronization.
//
// # Default behavior
//
// The default configuration is intentionally conservative:
//
//   - clock.RealClock{} for runtime time;
//   - delay.Immediate() as the default delay schedule;
//   - NeverRetry as the default classifier;
//   - one max attempt;
//   - no max elapsed limit;
//   - no observers.
//
// Therefore Do(ctx, op) and DoValue(ctx, op) execute the operation at most once
// unless callers explicitly configure retryability and limits.
//
// # Error ownership
//
// Package retry defines two retry-owned error classifications:
//
//   - ErrExhausted for retry-owned exhaustion;
//   - ErrInterrupted for retry-owned context interruption.
//
// ErrExhausted is returned when retry-owned boundaries prevent another attempt
// after a retryable operation failure, including context deadline budget
// rejection before a retry delay. ErrInterrupted is returned when retry
// observes its owning context stop at a retry boundary.
//
// Operation-owned errors that are not retried are returned unchanged. The retry
// package does not wrap all operation failures and does not create a
// "non-retryable" error type.
//
// # Panic policy
//
// Package retry panics on programming errors such as nil context, nil operation,
// nil option, nil clock, nil delay schedule, nil classifier, nil observer,
// zero max attempts, negative max elapsed duration, or a negative delay
// returned with ok=true.
//
// Package retry does not recover panics raised by operations, observers, clocks,
// delay schedules, or classifiers. Panic recovery, if required, belongs to the
// operation owner, observer implementation, runtime supervisor, or an explicit
// wrapper outside this package.
//
// # File ownership
//
// The package is intentionally split by responsibility:
//
//   - operation.go defines executable operation contracts;
//   - attempt.go defines per-attempt metadata;
//   - stop_reason.go defines terminal stop reasons;
//   - outcome.go defines terminal completion metadata;
//   - classifier.go and classifier_func.go define retryability classification;
//   - event_kind.go and event.go define observer event metadata;
//   - observer.go and observer_func.go define observer contracts;
//   - error_kind.go, error_exhausted.go, error_interrupted.go, and
//     error_context_stop.go define retry-owned error classification;
//   - validate.go defines package-local programming-error validation;
//   - option.go and option_*.go define normalized retry configuration;
//   - delay.go defines clock-backed retry-owned waiting;
//   - loop_deadline.go owns context deadline budget checks for retry delay
//     scheduling;
//   - loop.go defines the private execution engine;
//   - do.go and do_value.go define public entry points.
//
// # Non-goals
//
// Package retry deliberately does not provide:
//
//   - HTTP retry adapters;
//   - gRPC retry adapters;
//   - storage retry adapters;
//   - controller-runtime retry helpers;
//   - retry budgets;
//   - circuit breakers;
//   - hedged or speculative requests;
//   - token buckets or rate limiters;
//   - worker pools;
//   - goroutine supervision;
//   - lifecycle management;
//   - logging, metrics, or tracing backends;
//   - queue or scheduler policy.
package retry
