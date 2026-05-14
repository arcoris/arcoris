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
	"context"
	"errors"
	"testing"
	"time"

	"arcoris.dev/component-base/pkg/clock"
	"arcoris.dev/component-base/pkg/delay"
)

func TestNewRetryExecutionCreatesOwnedDelaySequence(t *testing.T) {
	now := time.Unix(10, 0)
	fake := clock.NewFakeClock(now)
	config := configOf(WithClock(fake))
	config.delay = retryTestSchedule{
		sequence: &retryTestSequence{delays: []time.Duration{time.Second}},
	}

	execution := newRetryExecution(config)

	if execution.startedAt != now {
		t.Fatalf("startedAt = %v, want %v", execution.startedAt, now)
	}

	got, ok := execution.nextDelay()
	if !ok {
		t.Fatalf("nextDelay ok = false, want true")
	}
	if got != time.Second {
		t.Fatalf("nextDelay = %v, want %v", got, time.Second)
	}
}

func TestNewRetryExecutionPanicsWhenDelayScheduleReturnsNilSequence(t *testing.T) {
	config := configOf()
	config.delay = retryTestSchedule{}

	expectPanic(t, panicNilDelaySequence, func() {
		_ = newRetryExecution(config)
	})
}

func TestRetryExecutionRecordsAttemptsFailuresAndDelayEvents(t *testing.T) {
	var events []Event

	fake := clock.NewFakeClock(time.Unix(20, 0))
	config := configOf(
		WithClock(fake),
		WithClassifier(RetryAll()),
		WithMaxAttempts(2),
		WithDelaySchedule(delay.Immediate()),
		WithObserverFunc(func(_ context.Context, event Event) {
			if !event.IsValid() {
				t.Fatalf("observer received invalid event: %+v", event)
			}
			events = append(events, event)
		}),
	)
	execution := newRetryExecution(config)
	errBoom := errors.New("boom")

	attempt := execution.nextAttempt(context.Background())
	if attempt.Number != 1 {
		t.Fatalf("attempt.Number = %d, want 1", attempt.Number)
	}
	if attempt.StartedAt != fake.Now() {
		t.Fatalf("attempt.StartedAt = %v, want %v", attempt.StartedAt, fake.Now())
	}

	execution.recordFailure(context.Background(), attempt, errBoom)
	if execution.lastErr != errBoom {
		t.Fatalf("lastErr = %v, want %v", execution.lastErr, errBoom)
	}
	if !execution.retryable(errBoom) {
		t.Fatalf("retryable returned false, want true")
	}
	if execution.maxAttemptsReached() {
		t.Fatalf("maxAttemptsReached returned true, want false")
	}

	got, ok := execution.nextDelay()
	if !ok {
		t.Fatalf("nextDelay ok = false, want true")
	}
	if got != 0 {
		t.Fatalf("nextDelay = %v, want 0", got)
	}

	execution.retryDelay(context.Background(), got)

	wantKinds := []EventKind{
		EventAttemptStart,
		EventAttemptFailure,
		EventRetryDelay,
	}
	if len(events) != len(wantKinds) {
		t.Fatalf("events len = %d, want %d: %+v", len(events), len(wantKinds), events)
	}
	for i, want := range wantKinds {
		if events[i].Kind != want {
			t.Fatalf("events[%d].Kind = %s, want %s", i, events[i].Kind, want)
		}
	}
}

func TestRetryExecutionNextDelayPanicsOnNegativeDelay(t *testing.T) {
	execution := &retryExecution{
		sequence: &retryTestSequence{delays: []time.Duration{-time.Nanosecond}},
	}

	expectPanic(t, panicNegativeDelay, func() {
		_, _ = execution.nextDelay()
	})
}

func TestRetryExecutionContextStopReturnsInterruptedError(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	execution := &retryExecution{}
	err := execution.contextStop(ctx)

	if !errors.Is(err, ErrInterrupted) {
		t.Fatalf("contextStop error = %v, want ErrInterrupted", err)
	}
	if !errors.Is(err, context.Canceled) {
		t.Fatalf("contextStop error = %v, want context.Canceled", err)
	}
}
