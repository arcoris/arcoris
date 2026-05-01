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

	"arcoris.dev/component-base/pkg/backoff"
)

func expectPanic(t *testing.T, want any, fn func()) {
	t.Helper()

	defer func() {
		recovered := recover()
		if recovered == nil {
			t.Fatalf("function did not panic")
		}
		if !errors.Is(asError(recovered), asError(want)) && recovered != want {
			t.Fatalf("panic = %v, want %v", recovered, want)
		}
	}()

	fn()
}

func asError(value any) error {
	err, ok := value.(error)
	if !ok {
		return nil
	}
	return err
}

func retryTestAttempt(number uint) Attempt {
	return Attempt{
		Number:    number,
		StartedAt: time.Unix(int64(number), 0),
	}
}

func retryTestStartedAt() time.Time {
	return time.Unix(100, 0)
}

func retryTestFinishedAt() time.Time {
	return time.Unix(101, 0)
}

func retryTestSuccessOutcome(attempts uint) Outcome {
	return Outcome{
		Attempts:   attempts,
		StartedAt:  retryTestStartedAt(),
		FinishedAt: retryTestFinishedAt(),
		Reason:     StopReasonSucceeded,
	}
}

func retryTestFailureOutcome(attempts uint, reason StopReason, err error) Outcome {
	return Outcome{
		Attempts:   attempts,
		StartedAt:  retryTestStartedAt(),
		FinishedAt: retryTestFinishedAt(),
		LastErr:    err,
		Reason:     reason,
	}
}

type retryTestSchedule struct {
	sequence backoff.Sequence
}

func (s retryTestSchedule) NewSequence() backoff.Sequence {
	return s.sequence
}

type retryTestSequence struct {
	delays []time.Duration
}

func (s *retryTestSequence) Next() (time.Duration, bool) {
	if len(s.delays) == 0 {
		return 0, false
	}

	delay := s.delays[0]
	s.delays = s.delays[1:]
	return delay, true
}

type retryObserverRecorder struct {
	events []Event
	order  []string
	name   string
}

func (r *retryObserverRecorder) observe(_ context.Context, event Event) {
	r.events = append(r.events, event)
	if r.name != "" {
		r.order = append(r.order, r.name)
	}
}
