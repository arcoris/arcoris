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

package value

import "testing"

func TestTimeOfDayValue(t *testing.T) {
	timeOfDay := TimeOfDay{hour: 1, minute: 2, second: 3, nanosecond: 4}
	value, err := TimeOfDayValue(timeOfDay)
	requireNoError(t, err)

	requireEqual(t, value.Kind(), KindTimeOfDay)
	requireEqual(t, value.timeOfDayValue.Equal(timeOfDay), true)
}

func TestTimeOfDayValueAcceptsMidnight(t *testing.T) {
	value, err := TimeOfDayValue(TimeOfDay{})
	requireNoError(t, err)

	requireEqual(t, value.Kind(), KindTimeOfDay)
	requireEqual(t, value.timeOfDayValue.Equal(TimeOfDay{}), true)
}

func TestTimeOfDayValueRejectsInvalidValue(t *testing.T) {
	_, err := TimeOfDayValue(TimeOfDay{hour: 24})

	requireValueError(
		t,
		err,
		ErrInvalidTime,
		pathTimeOfDay,
		ErrorReasonInvalidTime,
	)
}

func TestMustTimeOfDayValuePanicsOnInvalidTime(t *testing.T) {
	requirePanic(t, func() {
		MustTimeOfDayValue(TimeOfDay{hour: 24})
	})
}
