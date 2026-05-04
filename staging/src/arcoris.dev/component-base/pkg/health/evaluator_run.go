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

package health

import (
	"context"
	"errors"
	"runtime/debug"
	"time"
)

// evaluateCheck executes one checker and applies evaluator-owned normalization.
func (e *Evaluator) evaluateCheck(ctx context.Context, check Checker, timeout time.Duration) Result {
	started := e.clock.Now()

	name := ""
	if !nilChecker(check) {
		name = check.Name()
	}

	result := e.runCheck(ctx, check, timeout)

	finished := e.clock.Now()
	duration := nonNegativeDuration(e.clock.Since(started))

	return normalizeEvaluatedResult(result, name, finished, duration)
}

// runCheck runs check and returns a raw Result.
//
// Timeout enforcement uses a goroutine so the evaluator can return when the
// timeout or parent context expires. Checkers still MUST observe ctx. A checker
// that ignores ctx may continue running after a timeout result has already been
// returned.
func (e *Evaluator) runCheck(ctx context.Context, check Checker, timeout time.Duration) Result {
	if nilChecker(check) {
		return Unknown(
			"",
			ReasonNotObserved,
			"health checker is nil",
		).WithCause(ErrNilChecker)
	}

	if timeout == 0 {
		return callCheck(ctx, check)
	}

	checkCtx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	resultCh := make(chan Result, 1)

	go func() {
		resultCh <- callCheck(checkCtx, check)
	}()

	select {
	case result := <-resultCh:
		return result

	case <-checkCtx.Done():
		return interruptedResult(check.Name(), checkCtx)
	}
}

// callCheck invokes check and converts panics into health results.
func callCheck(ctx context.Context, check Checker) (result Result) {
	defer func() {
		if recovered := recover(); recovered != nil {
			result = Unhealthy(
				check.Name(),
				ReasonPanic,
				"health check panicked",
			).WithCause(PanicError{
				Value: recovered,
				Stack: debug.Stack(),
			})
		}
	}()

	return check.Check(ctx)
}

// interruptedResult converts context interruption into an unknown health result.
func interruptedResult(name string, ctx context.Context) Result {
	err := ctx.Err()
	cause := context.Cause(ctx)
	if cause == nil {
		cause = err
	}

	if errors.Is(err, context.DeadlineExceeded) {
		return Unknown(
			name,
			ReasonTimeout,
			"health check timed out",
		).WithCause(cause)
	}

	return Unknown(
		name,
		ReasonCanceled,
		"health check canceled",
	).WithCause(cause)
}

// normalizeEvaluatedResult applies evaluator-owned boundary normalization.
//
// Evaluator owns checker identity in reports. A checker may leave Result.Name
// empty, but it must not return another checker's name. Invalid reasons are also
// converted into unknown misconfiguration results so Evaluator never returns an
// invalid Report because of malformed checker output.
func normalizeEvaluatedResult(result Result, defaultName string, observed time.Time, duration time.Duration) Result {
	duration = nonNegativeDuration(duration)

	if result.Name != "" && result.Name != defaultName {
		return Unknown(
			defaultName,
			ReasonMisconfigured,
			"health check returned a mismatched result name",
		).WithCause(MismatchedCheckResultError{
			CheckName:  defaultName,
			ResultName: result.Name,
		}).WithObserved(observed).WithDuration(duration)
	}

	if !result.Reason.IsValid() {
		return Unknown(
			defaultName,
			ReasonMisconfigured,
			"health check returned an invalid reason",
		).WithCause(InvalidCheckResultError{
			CheckName: defaultName,
			Result:    result,
		}).WithObserved(observed).WithDuration(duration)
	}

	result = result.Normalize(defaultName, observed)

	if result.Duration == 0 {
		result.Duration = duration
	}

	if result.Duration < 0 {
		result.Duration = 0
	}

	return result
}

// nonNegativeDuration returns duration unless it is negative.
//
// Negative durations can occur with mutable fake clocks. Runtime reports should
// remain conservative and never expose negative elapsed time.
func nonNegativeDuration(duration time.Duration) time.Duration {
	if duration < 0 {
		return 0
	}

	return duration
}
