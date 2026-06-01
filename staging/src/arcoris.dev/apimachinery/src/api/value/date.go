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

// Date stores a calendar date without time-of-day or timezone.
//
// Date is not backed by time.Time so it cannot accidentally carry a location,
// clock, or monotonic component. Calendar validity is enforced by NewDate.
type Date struct {
	// year stores the proleptic Gregorian calendar year.
	year int
	// month stores January through December.
	month time.Month
	// day stores the day of month after calendar validation.
	day int
}

// Year returns the calendar year.
//
// The value is stored exactly as constructed; this package does not impose a
// minimum or maximum year policy.
func (d Date) Year() int {
	return d.year
}

// Month returns the calendar month.
//
// Values produced by NewDate are always January through December.
func (d Date) Month() time.Month {
	return d.month
}

// Day returns the day of month.
//
// Values produced by NewDate are guaranteed to exist in the stored year/month,
// including leap-year handling.
func (d Date) Day() int {
	return d.day
}

// IsValid reports whether d could have been produced by NewDate.
//
// The zero Date is invalid because month zero is not a calendar month. This
// lets DateValue reject uninitialized Date values instead of turning them into
// concrete payloads.
func (d Date) IsValid() bool {
	return d.month >= time.January &&
		d.month <= time.December &&
		isExistingDate(d.year, d.month, d.day)
}

// String returns the canonical diagnostic date text.
//
// The YYYY-MM-DD text is for diagnostics and tests. JSON or other wire encoding
// belongs to a codec layer.
func (d Date) String() string {
	return formatDate(d.year, d.month, d.day)
}

// Equal reports whether d and other represent the same calendar date.
//
// Date equality ignores timezone by design because Date has no timezone member.
func (d Date) Equal(other Date) bool {
	return d.year == other.year && d.month == other.month && d.day == other.day
}
