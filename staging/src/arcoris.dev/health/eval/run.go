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

package eval

import (
	"context"
	"errors"
	"reflect"
	"runtime/debug"
	"time"

	"arcoris.dev/health"
)

// evaluateCheck executes one checker and applies evaluator-owned normalization.
func (e *Evaluator) evaluateCheck(ctx context.Context, check health.Checker, timeout time.Duration) health.Result {
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

// runCheck runs check and returns a raw health.Result.
//
// Timeout enforcement uses a goroutine so the evaluator can return when the
// timeout or parent context expires. Checkers still MUST observe ctx. A checker
// that ignores ctx may continue running after a timeout result has already been
// returned.
func (e *Evaluator) runCheck(ctx context.Context, check health.Checker, timeout time.Duration) health.Result {
	if nilChecker(check) {
		return health.Unknown(
			"",
			health.ReasonNotObserved,
			"health checker is nil",
		).WithCause(health.ErrNilChecker)
	}

	if timeout == 0 {
		return callCheck(ctx, check)
	}

	checkCtx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	resultCh := make(chan health.Result, 1)

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
func callCheck(ctx context.Context, check health.Checker) (result health.Result) {
	defer func() {
		if recovered := recover(); recovered != nil {
			result = health.Unhealthy(
				check.Name(),
				health.ReasonPanic,
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
func interruptedResult(name string, ctx context.Context) health.Result {
	err := ctx.Err()
	cause := context.Cause(ctx)
	if cause == nil {
		cause = err
	}

	if errors.Is(err, context.DeadlineExceeded) {
		return health.Unknown(
			name,
			health.ReasonTimeout,
			"health check timed out",
		).WithCause(cause)
	}

	return health.Unknown(
		name,
		health.ReasonCanceled,
		"health check canceled",
	).WithCause(cause)
}

// normalizeEvaluatedResult applies evaluator-owned boundary normalization.
//
// Evaluator owns checker identity in reports. A checker may leave health.Result.Name
// empty, but it must not return another checker's name. Invalid reasons are also
// converted into unknown misconfiguration results so Evaluator never returns an
// invalid health.Report because of malformed checker output.
func normalizeEvaluatedResult(result health.Result, defaultName string, observed time.Time, duration time.Duration) health.Result {
	duration = nonNegativeDuration(duration)

	if result.Name != "" && result.Name != defaultName {
		return health.Unknown(
			defaultName,
			health.ReasonMisconfigured,
			"health check returned a mismatched result name",
		).WithCause(MismatchedCheckResultError{
			CheckName:  defaultName,
			ResultName: result.Name,
		}).WithObserved(observed).WithDuration(duration)
	}

	if !result.Reason.IsValid() {
		return health.Unknown(
			defaultName,
			health.ReasonMisconfigured,
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

// nilChecker reports whether chk is nil, including typed nil interface values.
func nilChecker(chk health.Checker) bool {
	if chk == nil {
		return true
	}

	value := reflect.ValueOf(chk)
	switch value.Kind() {
	case reflect.Chan, reflect.Func, reflect.Interface, reflect.Map, reflect.Pointer, reflect.Slice:
		return value.IsNil()
	default:
		return false
	}
}
