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

package eval

import (
	"context"
	"runtime/debug"
	"time"

	"arcoris.dev/health"
)

// runCheck runs check and returns a raw health.Result.
//
// name is the already-validated checker identity read at the evaluation
// boundary. Timeout and panic paths use this value instead of calling Name again
// after Check has started.
//
// Timeout enforcement uses a goroutine so the evaluator can return when the
// timeout or parent context expires. Checkers still MUST observe ctx. A checker
// that ignores ctx may continue running after a timeout result has already been
// returned.
func (e *Evaluator) runCheck(ctx context.Context, check health.Checker, name string, timeout time.Duration) health.Result {
	if timeout == 0 {
		return callCheck(ctx, check, name)
	}

	checkCtx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	resultCh := make(chan health.Result, 1)

	go func() {
		resultCh <- callCheck(checkCtx, check, name)
	}()

	select {
	case res := <-resultCh:
		return res

	case <-checkCtx.Done():
		return interruptedResult(name, checkCtx)
	}
}

// callCheck invokes check and converts panics into health results.
func callCheck(ctx context.Context, check health.Checker, name string) (res health.Result) {
	defer func() {
		if recovered := recover(); recovered != nil {
			res = health.Unhealthy(
				name,
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
