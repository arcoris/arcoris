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

package retry_test

import (
	"context"
	"errors"
	"fmt"
	"time"

	"arcoris.dev/chrono/clock"
	"arcoris.dev/chrono/delay"
	"arcoris.dev/resilience/retry"
)

func ExampleDo_defaultSingleAttempt() {
	errTransient := errors.New("transient")
	calls := 0

	err := retry.Do(context.Background(), func(context.Context) error {
		calls++
		return errTransient
	})

	fmt.Println("calls:", calls)
	fmt.Println("operation error:", errors.Is(err, errTransient))

	// Output:
	// calls: 1
	// operation error: true
}

func ExampleDo_withRetryAllAndMaxAttempts() {
	errTransient := errors.New("transient")
	calls := 0

	err := retry.Do(
		context.Background(),
		func(context.Context) error {
			calls++
			if calls < 3 {
				return errTransient
			}
			return nil
		},
		retry.WithClassifier(retry.RetryAll()),
		retry.WithMaxAttempts(3),
		retry.WithDelaySchedule(delay.Immediate()),
	)

	fmt.Println("calls:", calls)
	fmt.Println("success:", err == nil)

	// Output:
	// calls: 3
	// success: true
}

func ExampleDoValue() {
	value, err := retry.DoValue(context.Background(), func(context.Context) (string, error) {
		return "ready", nil
	})

	fmt.Println(value)
	fmt.Println(err == nil)

	// Output:
	// ready
	// true
}

func ExampleDoObserved() {
	errTransient := errors.New("transient")

	outcome, err := retry.DoObserved(
		context.Background(),
		func(context.Context) error {
			return errTransient
		},
		retry.WithClassifier(retry.RetryAll()),
		retry.WithMaxAttempts(2),
		retry.WithDelaySchedule(delay.Immediate()),
	)

	fmt.Println("exhausted:", retry.Exhausted(err))
	fmt.Println("attempts:", outcome.Attempts)
	fmt.Println("reason:", outcome.Reason)

	// Output:
	// exhausted: true
	// attempts: 2
	// reason: max_attempts
}

func ExampleDo_contextCancellation() {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	outcome, err := retry.DoObserved(ctx, func(context.Context) error {
		return nil
	})

	fmt.Println("interrupted:", retry.Interrupted(err))
	fmt.Println("attempts:", outcome.Attempts)

	// Output:
	// interrupted: true
	// attempts: 0
}

func ExampleClassifier() {
	errRetryable := errors.New("retryable")

	classifier := retry.ClassifierFunc(func(err error) bool {
		return errors.Is(err, errRetryable)
	})

	fmt.Println(classifier.Retryable(errRetryable))
	fmt.Println(classifier.Retryable(errors.New("permanent")))

	// Output:
	// true
	// false
}

func ExampleObserver() {
	errTransient := errors.New("transient")
	var events []retry.EventKind
	calls := 0

	_ = retry.Do(
		context.Background(),
		func(context.Context) error {
			calls++
			if calls == 1 {
				return errTransient
			}
			return nil
		},
		retry.WithClassifier(retry.RetryAll()),
		retry.WithMaxAttempts(2),
		retry.WithDelaySchedule(delay.Immediate()),
		retry.WithObserverFunc(func(_ context.Context, event retry.Event) {
			events = append(events, event.Kind)
		}),
	)

	for _, kind := range events {
		fmt.Println(kind)
	}

	// Output:
	// attempt_start
	// attempt_failure
	// retry_delay
	// attempt_start
	// retry_stop
}

func Example_maxElapsedIsSchedulingBoundary() {
	now := time.Date(2099, 6, 8, 12, 0, 0, 0, time.UTC)
	fake := clock.NewFakeClock(now)

	outcome, err := retry.DoObserved(
		context.Background(),
		func(context.Context) error {
			fake.Step(time.Hour)
			return nil
		},
		retry.WithClock(fake),
		retry.WithMaxElapsed(time.Nanosecond),
	)

	fmt.Println("success:", err == nil)
	fmt.Println("reason:", outcome.Reason)

	// Output:
	// success: true
	// reason: succeeded
}
