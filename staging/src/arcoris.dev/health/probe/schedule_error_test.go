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

package probe

import (
	"errors"
	"testing"
	"time"

	"arcoris.dev/chrono/delay"
)

func TestInvalidIntervalError(t *testing.T) {
	t.Parallel()

	err := InvalidIntervalError{Interval: -time.Second}

	if !errors.Is(err, ErrInvalidInterval) {
		t.Fatalf("errors.Is(%v, ErrInvalidInterval) = false, want true", err)
	}
	if err.Error() == "" {
		t.Fatal("Error() returned empty message")
	}
}

func TestInvalidIntervalErrorSupportsErrorsAs(t *testing.T) {
	t.Parallel()

	err := error(InvalidIntervalError{Interval: -time.Second})

	var intervalErr InvalidIntervalError
	if !errors.As(err, &intervalErr) {
		t.Fatalf("errors.As(%T, InvalidIntervalError) = false, want true", err)
	}
	if intervalErr.Interval != -time.Second {
		t.Fatalf("Interval = %s, want %s", intervalErr.Interval, -time.Second)
	}
}

func TestInvalidScheduleDelayError(t *testing.T) {
	t.Parallel()

	err := InvalidScheduleDelayError{Delay: -time.Nanosecond}

	if !errors.Is(err, ErrInvalidScheduleDelay) {
		t.Fatalf("errors.Is(%v, ErrInvalidScheduleDelay) = false, want true", err)
	}
	if err.Error() == "" {
		t.Fatal("Error() returned empty message")
	}
}

func TestInvalidScheduleDelayErrorSupportsErrorsAs(t *testing.T) {
	t.Parallel()

	err := error(InvalidScheduleDelayError{Delay: -time.Nanosecond})

	var delayErr InvalidScheduleDelayError
	if !errors.As(err, &delayErr) {
		t.Fatalf("errors.As(%T, InvalidScheduleDelayError) = false, want true", err)
	}
	if delayErr.Delay != -time.Nanosecond {
		t.Fatalf("Delay = %s, want %s", delayErr.Delay, -time.Nanosecond)
	}
}

func TestValidateScheduleRejectsNil(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name  string
		sched delay.Schedule
	}{
		{
			name:  "nil interface",
			sched: nil,
		},
		{
			name:  "typed nil",
			sched: (*nilSequenceSchedule)(nil),
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			if err := validateSchedule(tc.sched); !errors.Is(err, ErrNilSchedule) {
				t.Fatalf("validateSchedule() = %v, want ErrNilSchedule", err)
			}
		})
	}
}

func TestNewSequenceRejectsNilSequence(t *testing.T) {
	t.Parallel()

	_, err := newSequence(nilSequenceSchedule{})

	if !errors.Is(err, ErrNilSequence) {
		t.Fatalf("newSequence(nil sequence) = %v, want ErrNilSequence", err)
	}
}
