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
	seq := Delays(0, time.Second, 2*time.Second).NewSequence()

	mustNext(t, seq, 0)
	mustNext(t, seq, time.Second)
	mustNext(t, seq, 2*time.Second)
	mustExhausted(t, seq)
	mustExhausted(t, seq)
}

func TestDelaysAllowsEmptyInput(t *testing.T) {
	mustExhausted(t, Delays().NewSequence())
}

func TestDelaysCopiesCallerSlice(t *testing.T) {
	delays := []time.Duration{time.Second, 2 * time.Second}
	sched := Delays(delays...)
	delays[0] = time.Hour

	seq := sched.NewSequence()
	mustNext(t, seq, time.Second)
	mustNext(t, seq, 2*time.Second)
	mustExhausted(t, seq)
}

func TestDelaysSequencesHaveIndependentCursors(t *testing.T) {
	sched := Delays(time.Second, 2*time.Second)

	l := sched.NewSequence()
	r := sched.NewSequence()

	mustNext(t, l, time.Second)
	mustNext(t, l, 2*time.Second)
	mustNext(t, r, time.Second)
	mustNext(t, r, 2*time.Second)
}
