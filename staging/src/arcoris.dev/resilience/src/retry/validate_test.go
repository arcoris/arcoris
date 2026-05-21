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

package retry

import (
	panicassert "arcoris.dev/testutil/panic"
	"context"
	"testing"
	"time"

	"arcoris.dev/chrono/clock"
	"arcoris.dev/chrono/delay"
)

func TestRequireContext(t *testing.T) {
	requireContext(context.Background())

	panicassert.RequireValue(t, panicNilContext, func() {
		requireContext(nil)
	})
}

func TestRequireOperation(t *testing.T) {
	requireOperation(func(context.Context) error {
		return nil
	})

	panicassert.RequireValue(t, panicNilOperation, func() {
		requireOperation(nil)
	})
}

func TestRequireValueOperation(t *testing.T) {
	requireValueOperation(func(context.Context) (int, error) {
		return 0, nil
	})

	panicassert.RequireValue(t, panicNilValueOperation, func() {
		requireValueOperation[int](nil)
	})
}

func TestRequireClock(t *testing.T) {
	requireClock(clock.RealClock{})

	panicassert.RequireValue(t, panicNilClock, func() {
		requireClock(nil)
	})
}

func TestRequireDelaySchedule(t *testing.T) {
	requireDelaySchedule(delay.Immediate())

	panicassert.RequireValue(t, panicNilDelaySchedule, func() {
		requireDelaySchedule(nil)
	})
}

func TestRequireDelayScheduleSequence(t *testing.T) {
	requireDelaySequence(delay.Immediate().NewSequence())

	panicassert.RequireValue(t, panicNilDelaySequence, func() {
		requireDelaySequence(nil)
	})
}

func TestRequireDelayScheduleDelay(t *testing.T) {
	requireDelay(0, true)
	requireDelay(time.Nanosecond, true)
	requireDelay(-time.Nanosecond, false)

	panicassert.RequireValue(t, panicNegativeDelay, func() {
		requireDelay(-time.Nanosecond, true)
	})
}

func TestRequireClassifier(t *testing.T) {
	requireClassifier(NeverRetry())

	panicassert.RequireValue(t, panicNilClassifier, func() {
		requireClassifier(nil)
	})
}

func TestRequireRetryableFunc(t *testing.T) {
	requireRetryableFunc(func(error) bool {
		return false
	})

	panicassert.RequireValue(t, panicNilClassifierFunc, func() {
		requireRetryableFunc(nil)
	})
}

func TestRequireMaxAttempts(t *testing.T) {
	requireMaxAttempts(1)
	requireMaxAttempts(10)

	panicassert.RequireValue(t, panicZeroMaxAttempts, func() {
		requireMaxAttempts(0)
	})
}

func TestRequireMaxElapsed(t *testing.T) {
	requireMaxElapsed(0)
	requireMaxElapsed(time.Nanosecond)

	panicassert.RequireValue(t, panicNegativeMaxElapsed, func() {
		requireMaxElapsed(-time.Nanosecond)
	})
}

func TestRequireObserver(t *testing.T) {
	requireObserver(ObserverFunc(func(context.Context, Event) {}))

	panicassert.RequireValue(t, panicNilObserver, func() {
		requireObserver(nil)
	})
}

func TestRequireObserverFunc(t *testing.T) {
	requireObserverFunc(func(context.Context, Event) {})

	panicassert.RequireValue(t, panicNilObserverFunc, func() {
		requireObserverFunc(nil)
	})
}

func TestRequireOption(t *testing.T) {
	requireOption(func(*config) {})

	panicassert.RequireValue(t, panicNilOption, func() {
		requireOption(nil)
	})
}
