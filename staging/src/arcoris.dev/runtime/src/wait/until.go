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

package wait

import (
	"context"
	"time"
)

// Until evaluates condition until it is satisfied, fails, or ctx stops.
//
// Until is a small fixed-interval wait loop. It owns only loop mechanics:
// condition evaluation, sleeping between unsuccessful evaluations, and mapping
// wait-owned context stops into wait-owned errors. It does not implement retry
// policy, backoff growth, randomized cadence, rate limiting, metrics, panic
// recovery, or scheduler policy.
//
// Evaluation is immediate: Until checks ctx and evaluates condition once before
// sleeping. If condition is already satisfied, Until returns nil without waiting
// for interval to elapse.
//
// The condition result is interpreted as follows:
//
//   - done=true, err=nil completes the wait successfully and returns nil;
//   - done=false, err=nil sleeps for interval and evaluates again;
//   - err!=nil stops the wait and returns err unchanged.
//
// Condition errors are condition-owned errors. Until does not reinterpret, wrap,
// classify, retry, or suppress them. In particular, if condition returns raw
// context.Canceled or context.DeadlineExceeded for condition-owned work, Until
// returns that raw error unchanged.
//
// Context stops are wait-owned errors. If ctx is cancelled before the condition
// is satisfied, Until returns an error classified as ErrInterrupted. If ctx ends
// because its deadline expired, Until returns an error classified as ErrTimeout
// and ErrInterrupted. Cancellation causes created with context.WithCancelCause,
// context.WithDeadlineCause, or context.WithTimeoutCause are preserved as the
// wrapped cause.
//
// If condition returns done=true during an evaluation, success wins for that
// evaluation even if ctx was cancelled while the condition was running. If
// condition returns done=false and ctx is already stopped after evaluation,
// Until returns the wait-owned context stop error without sleeping again.
//
// Until does not recover panics raised by condition. Panic recovery, if required,
// belongs to a higher-level runtime owner or to an explicit wrapper.
//
// Until panics when ctx is nil, interval is not positive, or condition is nil.
func Until(ctx context.Context, interval time.Duration, condition ConditionFunc) error {
	requireContext(ctx)
	requirePositiveInterval(interval)
	requireCondition(condition)

	done, err := evaluateUntilCondition(ctx, condition)
	if err != nil {
		return err
	}
	if done {
		return nil
	}

	timer := NewTimer(interval)
	defer timer.StopAndDrain()

	for {
		if err = timer.Wait(ctx); err != nil {
			return err
		}

		done, err = evaluateUntilCondition(ctx, condition)
		if err != nil {
			return err
		}
		if done {
			return nil
		}

		timer.Reset(interval)
	}
}

// evaluateUntilCondition evaluates condition once under ctx.
//
// The helper centralizes the context-before-evaluation rule used by Until. A
// context that is already stopped before evaluation prevents the condition from
// running and is returned as a wait-owned stop error. A context that stops during
// condition execution does not override a successful condition result from that
// evaluation.
func evaluateUntilCondition(ctx context.Context, condition ConditionFunc) (done bool, err error) {
	if err = contextStopError(ctx); err != nil {
		return false, err
	}

	done, err = condition(ctx)
	if err != nil {
		return false, err
	}
	if done {
		return true, nil
	}

	if err = contextStopError(ctx); err != nil {
		return false, err
	}

	return false, nil
}
