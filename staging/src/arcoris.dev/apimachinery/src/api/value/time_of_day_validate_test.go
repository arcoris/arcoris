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

func TestNewTimeOfDay(t *testing.T) {
	timeOfDay, err := NewTimeOfDay(1, 2, 3, 4)
	requireNoError(t, err)

	requireEqual(t, timeOfDay.Hour(), 1)
	requireEqual(t, timeOfDay.Minute(), 2)
	requireEqual(t, timeOfDay.Second(), 3)
	requireEqual(t, timeOfDay.Nanosecond(), 4)
}

func TestNewTimeOfDayRejectsInvalidValues(t *testing.T) {
	tests := []struct {
		name       string
		hour       int
		minute     int
		second     int
		nanosecond int
	}{
		{name: "invalid hour", hour: 24},
		{name: "invalid minute", minute: 60},
		{name: "invalid second", second: 60},
		{name: "invalid nanosecond", nanosecond: 1000000000},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := NewTimeOfDay(tt.hour, tt.minute, tt.second, tt.nanosecond)
			requireValueError(t, err, ErrInvalidTime, pathTimeOfDay, ErrorReasonInvalidTime)
			requireErrorIs(t, err, ErrInvalidValue)
		})
	}
}
