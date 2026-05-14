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

package delay

import (
	"testing"
	"time"
)

func TestDelaysRejectsNegativeDelay(t *testing.T) {
	mustPanicWith(t, errNegativeExplicitDelay, func() {
		Delays(time.Second, -time.Nanosecond)
	})
}

func TestDelaysReturnsConfiguredValuesThenExhausts(t *testing.T) {
	sequence := Delays(0, time.Second, 2*time.Second).NewSequence()

	mustNext(t, sequence, 0)
	mustNext(t, sequence, time.Second)
	mustNext(t, sequence, 2*time.Second)
	mustExhausted(t, sequence)
	mustExhausted(t, sequence)
}

func TestDelaysAllowsEmptyInput(t *testing.T) {
	mustExhausted(t, Delays().NewSequence())
}

func TestDelaysCopiesCallerSlice(t *testing.T) {
	delays := []time.Duration{time.Second, 2 * time.Second}
	schedule := Delays(delays...)
	delays[0] = time.Hour

	sequence := schedule.NewSequence()
	mustNext(t, sequence, time.Second)
	mustNext(t, sequence, 2*time.Second)
	mustExhausted(t, sequence)
}

func TestDelaysSequencesHaveIndependentCursors(t *testing.T) {
	schedule := Delays(time.Second, 2*time.Second)

	left := schedule.NewSequence()
	right := schedule.NewSequence()

	mustNext(t, left, time.Second)
	mustNext(t, left, 2*time.Second)
	mustNext(t, right, time.Second)
	mustNext(t, right, 2*time.Second)
}
