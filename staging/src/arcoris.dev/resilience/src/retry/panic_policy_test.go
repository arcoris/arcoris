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
	"context"
	"errors"
	"testing"
	"time"

	"arcoris.dev/chrono/clock"
	"arcoris.dev/chrono/delay"
	panicassert "arcoris.dev/testutil/panic"
)

func TestRetryDoesNotRecoverOperationPanic(t *testing.T) {
	t.Parallel()

	errPanic := errors.New("operation panic")

	panicassert.RequireErrorIs(t, errPanic, func() {
		_, _ = DoObserved(context.Background(), func(context.Context) error {
			panic(errPanic)
		})
	})
}

func TestRetryDoesNotRecoverValueOperationPanic(t *testing.T) {
	t.Parallel()

	errPanic := errors.New("value operation panic")

	panicassert.RequireErrorIs(t, errPanic, func() {
		_, _, _ = DoValueObserved(context.Background(), func(context.Context) (int, error) {
			panic(errPanic)
		})
	})
}

func TestRetryDoesNotRecoverClassifierPanic(t *testing.T) {
	t.Parallel()

	errPanic := errors.New("classifier panic")

	panicassert.RequireErrorIs(t, errPanic, func() {
		_, _ = DoObserved(
			context.Background(),
			func(context.Context) error { return benchmarkErrBoom },
			WithClassifier(ClassifierFunc(func(error) bool {
				panic(errPanic)
			})),
		)
	})
}

func TestRetryDoesNotRecoverObserverPanic(t *testing.T) {
	t.Parallel()

	errPanic := errors.New("observer panic")

	panicassert.RequireErrorIs(t, errPanic, func() {
		_, _ = DoObserved(
			context.Background(),
			func(context.Context) error { return nil },
			WithObserverFunc(func(context.Context, Event) {
				panic(errPanic)
			}),
		)
	})
}

func TestRetryDoesNotRecoverClockPanic(t *testing.T) {
	t.Parallel()

	errPanic := errors.New("clock panic")

	panicassert.RequireErrorIs(t, errPanic, func() {
		_, _ = DoObserved(
			context.Background(),
			func(context.Context) error { return nil },
			WithClock(panicClock{err: errPanic}),
		)
	})
}

func TestRetryDoesNotRecoverDelaySchedulePanic(t *testing.T) {
	t.Parallel()

	errPanic := errors.New("delay schedule panic")

	panicassert.RequireErrorIs(t, errPanic, func() {
		_, _ = DoObserved(
			context.Background(),
			func(context.Context) error { return nil },
			WithDelaySchedule(delay.ScheduleFunc(func() delay.Sequence {
				panic(errPanic)
			})),
		)
	})
}

type panicClock struct {
	err error
}

func (c panicClock) Now() time.Time {
	panic(c.err)
}

func (c panicClock) Since(time.Time) time.Duration {
	panic(c.err)
}

func (c panicClock) Until(time.Time) time.Duration {
	panic(c.err)
}

func (c panicClock) After(time.Duration) <-chan time.Time {
	panic(c.err)
}

func (c panicClock) NewTimer(time.Duration) clock.Timer {
	panic(c.err)
}

func (c panicClock) NewTicker(time.Duration) clock.Ticker {
	panic(c.err)
}

func (c panicClock) Sleep(time.Duration) {
	panic(c.err)
}
