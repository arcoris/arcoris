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

	panicassert.RequireErrorIs(t, ErrNilContext, func() {
		requireContext(nil)
	})
}

func TestRequireOperation(t *testing.T) {
	requireOperation(func(context.Context) error {
		return nil
	})

	panicassert.RequireErrorIs(t, ErrNilOperation, func() {
		requireOperation(nil)
	})
}

func TestRequireValueOperation(t *testing.T) {
	requireValueOperation(func(context.Context) (int, error) {
		return 0, nil
	})

	panicassert.RequireErrorIs(t, ErrNilValueOperation, func() {
		requireValueOperation[int](nil)
	})
}

func TestRequireClock(t *testing.T) {
	requireClock(clock.RealClock{})

	panicassert.RequireErrorIs(t, ErrNilClock, func() {
		requireClock(nil)
	})
}

func TestRequireDelaySchedule(t *testing.T) {
	requireDelaySchedule(delay.Immediate())

	panicassert.RequireErrorIs(t, ErrNilDelaySchedule, func() {
		requireDelaySchedule(nil)
	})
}

func TestRequireDelayScheduleSequence(t *testing.T) {
	requireDelaySequence(delay.Immediate().NewSequence())

	panicassert.RequireErrorIs(t, ErrNilDelaySequence, func() {
		requireDelaySequence(nil)
	})
}

func TestRequireDelayScheduleDelay(t *testing.T) {
	requireDelay(0, true)
	requireDelay(time.Nanosecond, true)
	requireDelay(-time.Nanosecond, false)

	panicassert.RequireErrorIs(t, ErrNegativeDelay, func() {
		requireDelay(-time.Nanosecond, true)
	})
}

func TestRequireClassifier(t *testing.T) {
	requireClassifier(NeverRetry())

	panicassert.RequireErrorIs(t, ErrNilClassifier, func() {
		requireClassifier(nil)
	})
}

func TestRequireRetryableFunc(t *testing.T) {
	requireRetryableFunc(func(error) bool {
		return false
	})

	panicassert.RequireErrorIs(t, ErrNilClassifierFunc, func() {
		requireRetryableFunc(nil)
	})
}

func TestRequireMaxAttempts(t *testing.T) {
	requireMaxAttempts(1)
	requireMaxAttempts(10)

	panicassert.RequireErrorIs(t, ErrZeroMaxAttempts, func() {
		requireMaxAttempts(0)
	})
}

func TestRequireMaxElapsed(t *testing.T) {
	requireMaxElapsed(0)
	requireMaxElapsed(time.Nanosecond)

	panicassert.RequireErrorIs(t, ErrNegativeMaxElapsed, func() {
		requireMaxElapsed(-time.Nanosecond)
	})
}

func TestRequireObserver(t *testing.T) {
	requireObserver(ObserverFunc(func(context.Context, Event) {}))

	panicassert.RequireErrorIs(t, ErrNilObserver, func() {
		requireObserver(nil)
	})
}

func TestRequireObserverFunc(t *testing.T) {
	requireObserverFunc(func(context.Context, Event) {})

	panicassert.RequireErrorIs(t, ErrNilObserverFunc, func() {
		requireObserverFunc(nil)
	})
}

func TestRequireOption(t *testing.T) {
	requireOption(func(*config) {})

	panicassert.RequireErrorIs(t, ErrNilOption, func() {
		requireOption(nil)
	})
}
