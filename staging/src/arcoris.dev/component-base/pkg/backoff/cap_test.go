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

package backoff

import (
	"testing"
	"time"
)

func TestCapRejectsInvalidInput(t *testing.T) {
	mustPanicWith(t, errNilCapSchedule, func() {
		Cap(nil, time.Second)
	})
	mustPanicWith(t, errNegativeCapMaxDelay, func() {
		Cap(Fixed(time.Second), -time.Nanosecond)
	})
}

func TestCapCapsAvailableChildDelays(t *testing.T) {
	sequence := Cap(Delays(time.Second, 3*time.Second), 2*time.Second).NewSequence()

	mustNext(t, sequence, time.Second)
	mustNext(t, sequence, 2*time.Second)
}

func TestCapAllowsZeroMaximumDelay(t *testing.T) {
	sequence := Cap(Delays(time.Second), 0).NewSequence()

	mustNext(t, sequence, 0)
	mustExhausted(t, sequence)
}

func TestCapPreservesChildExhaustion(t *testing.T) {
	sequence := Cap(Delays(2*time.Second), time.Second).NewSequence()

	mustNext(t, sequence, time.Second)
	mustExhausted(t, sequence)
}

func TestCapRejectsNilChildSequence(t *testing.T) {
	mustPanicWith(t, errCapScheduleReturnedNilSequence, func() {
		Cap(nilSequenceSchedule{}, time.Second).NewSequence()
	})
}

func TestCapRejectsNegativeChildDelay(t *testing.T) {
	sequence := Cap(ScheduleFunc(func() Sequence { return negativeDelaySequence{} }), time.Second).NewSequence()

	mustPanicWith(t, errCapScheduleReturnedNegativeDelay, func() {
		sequence.Next()
	})
}
