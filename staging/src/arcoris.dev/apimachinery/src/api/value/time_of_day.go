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

// TimeOfDay stores a time of day without date or timezone.
//
// TimeOfDay models a wall-clock time within a single day. It intentionally has
// no date, location, UTC offset, or leap-second representation.
type TimeOfDay struct {
	// hour stores the 24-hour clock hour in 0..23.
	hour int
	// minute stores the minute within the hour in 0..59.
	minute int
	// second stores the second within the minute in 0..59.
	second int
	// nanosecond stores fractional seconds in 0..999999999.
	nanosecond int
}

// Hour returns the 24-hour clock hour.
func (t TimeOfDay) Hour() int {
	return t.hour
}

// Minute returns the minute within the hour.
func (t TimeOfDay) Minute() int {
	return t.minute
}

// Second returns the second within the minute.
//
// Leap seconds are not modeled in this first pass, so constructed values are
// always in 0..59.
func (t TimeOfDay) Second() int {
	return t.second
}

// Nanosecond returns the fractional second in nanoseconds.
func (t TimeOfDay) Nanosecond() int {
	return t.nanosecond
}

// IsValid reports whether t is inside the supported wall-clock range.
//
// Unlike Date, the zero TimeOfDay is valid: it represents midnight exactly.
func (t TimeOfDay) IsValid() bool {
	return t.hour >= 0 &&
		t.hour <= 23 &&
		t.minute >= 0 &&
		t.minute <= 59 &&
		t.second >= 0 &&
		t.second <= 59 &&
		t.nanosecond >= 0 &&
		t.nanosecond <= 999999999
}

// String returns the canonical diagnostic time-of-day text.
//
// Fractional seconds are emitted only when non-zero. The returned text is not a
// package-level wire codec.
func (t TimeOfDay) String() string {
	return formatTimeOfDay(t)
}

// Equal reports whether t and other represent the same time of day.
//
// Equality compares only the four stored members because TimeOfDay has no date
// or timezone state.
func (t TimeOfDay) Equal(other TimeOfDay) bool {
	return t.hour == other.hour &&
		t.minute == other.minute &&
		t.second == other.second &&
		t.nanosecond == other.nanosecond
}
