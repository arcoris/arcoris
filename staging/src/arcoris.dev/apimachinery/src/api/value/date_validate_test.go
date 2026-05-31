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

import (
	"testing"
	"time"
)

func TestNewDate(t *testing.T) {
	date, err := NewDate(2024, time.February, 29)
	requireNoError(t, err)

	requireEqual(t, date.Year(), 2024)
	requireEqual(t, date.Month(), time.February)
	requireEqual(t, date.Day(), 29)
}

func TestNewDateRejectsInvalidCalendarValues(t *testing.T) {
	tests := []struct {
		name  string
		year  int
		month time.Month
		day   int
	}{
		{name: "invalid month", year: 2024, month: 13, day: 1},
		{name: "invalid day", year: 2024, month: time.January, day: 32},
		{name: "non leap feb 29", year: 2023, month: time.February, day: 29},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := NewDate(tt.year, tt.month, tt.day)
			requireValueError(t, err, ErrInvalidDate, pathDate, ErrorReasonInvalidDate)
			requireErrorIs(t, err, ErrInvalidValue)
		})
	}
}
