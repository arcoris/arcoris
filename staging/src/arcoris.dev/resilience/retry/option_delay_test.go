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
	"testing"
	"time"

	"arcoris.dev/chrono/delay"
)

func TestWithDelaySchedule(t *testing.T) {
	schedule := delay.Fixed(250 * time.Millisecond)

	config := configOf(WithDelaySchedule(schedule))

	sequence := config.delay.NewSequence()
	delay, ok := sequence.Next()
	if !ok {
		t.Fatalf("configured delay sequence exhausted")
	}
	if delay != 250*time.Millisecond {
		t.Fatalf("configured delay = %s, want %s", delay, 250*time.Millisecond)
	}
}

func TestWithDelayScheduleLastWins(t *testing.T) {
	config := configOf(
		WithDelaySchedule(delay.Fixed(100*time.Millisecond)),
		WithDelaySchedule(delay.Fixed(200*time.Millisecond)),
	)

	delay, ok := config.delay.NewSequence().Next()
	if !ok {
		t.Fatalf("configured delay sequence exhausted")
	}
	if delay != 200*time.Millisecond {
		t.Fatalf("configured delay = %s, want %s", delay, 200*time.Millisecond)
	}
}

func TestWithDelaySchedulePanicsOnNilSchedule(t *testing.T) {
	expectPanic(t, panicNilDelaySchedule, func() {
		_ = WithDelaySchedule(nil)
	})
}
