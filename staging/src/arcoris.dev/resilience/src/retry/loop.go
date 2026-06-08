// Copyright 2026 The ARCORIS Authors.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package retry

import "context"

// run executes op with retry orchestration.
//
// run is the private retry engine behind Do and DoValue. It owns the runtime
// control flow that combines operation execution, retryability classification,
// retry-owned limits, delay sequence consumption, clock-backed delays, context
// interruption, and observer events.
//
// run does not define operation semantics. The caller remains responsible for
// idempotency, replay safety, transactional safety, and external side effects.
// run also does not define protocol-specific retry rules, storage retry rules,
// retry budgets, circuit breakers, hedging, metrics exporters, tracing exporters,
// or logging backends.
//
// The returned value is meaningful only when err is nil. Every non-panic
// terminal path returns a valid Outcome that matches the terminal stop event.
// On failure, run returns the zero value of T, the terminal Outcome, and an
// error describing the terminal retry decision.
func run[T any](
	ctx context.Context,
	op ValueOperation[T],
	cfg config,
) (T, Outcome, error) {
	requireContext(ctx)
	requireValueOperation(op)

	var zero T

	execution := newRetryExecution(cfg)

	for {
		if err := execution.contextStop(ctx); err != nil {
			outcome, err := execution.interrupted(ctx, err)
			return zero, outcome, err
		}

		attempt := execution.nextAttempt(ctx)

		val, err := op(ctx)
		if err == nil {
			outcome := execution.succeeded(ctx)
			return val, outcome, nil
		}

		execution.recordFailure(ctx, attempt, err)

		if !execution.retryable(err) {
			outcome, err := execution.nonRetryable(ctx, err)
			return zero, outcome, err
		}

		if execution.maxAttemptsReached() {
			outcome, err := execution.exhausted(ctx, StopReasonMaxAttempts)
			return zero, outcome, err
		}

		if err := execution.contextStop(ctx); err != nil {
			outcome, err := execution.interrupted(ctx, err)
			return zero, outcome, err
		}

		delay, ok := execution.nextDelay()
		if !ok {
			outcome, err := execution.exhausted(ctx, StopReasonDelayExhausted)
			return zero, outcome, err
		}

		if execution.contextDeadlineWouldBeExceeded(ctx, delay) {
			outcome, err := execution.exhausted(ctx, StopReasonDeadline)
			return zero, outcome, err
		}

		if execution.maxElapsedWouldBeExceeded(delay) {
			outcome, err := execution.exhausted(ctx, StopReasonMaxElapsed)
			return zero, outcome, err
		}

		execution.retryDelay(ctx, delay)

		if err := execution.waitDelay(ctx, delay); err != nil {
			outcome, err := execution.interrupted(ctx, err)
			return zero, outcome, err
		}
	}
}
