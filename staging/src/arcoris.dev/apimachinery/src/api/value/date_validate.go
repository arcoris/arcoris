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

import "time"

// NewDate constructs a calendar date after verifying the day exists.
//
// The check uses Go's calendar normalization to reject impossible dates such as
// April 31 or February 29 in a non-leap year. It does not enforce business
// bounds such as "not before 1970"; those are descriptor or policy concerns.
func NewDate(year int, month time.Month, day int) (Date, error) {
	if month < time.January || month > time.December {
		return Date{}, errorf(
			pathDate,
			ErrInvalidDate,
			ErrorReasonInvalidDate,
			"month %d is outside January..December",
			month,
		)
	}

	if !isExistingDate(year, month, day) {
		return Date{}, errorf(
			pathDate,
			ErrInvalidDate,
			ErrorReasonInvalidDate,
			"day %d does not exist in %04d-%02d",
			day,
			year,
			month,
		)
	}

	return Date{year: year, month: month, day: day}, nil
}

// isExistingDate checks calendar existence through time.Date normalization.
//
// time.Date normalizes invalid dates instead of rejecting them. Round-tripping
// the components lets us detect that normalization deterministically.
func isExistingDate(year int, month time.Month, day int) bool {
	candidate := time.Date(year, month, day, 0, 0, 0, 0, time.UTC)

	return candidate.Year() == year &&
		candidate.Month() == month &&
		candidate.Day() == day
}
