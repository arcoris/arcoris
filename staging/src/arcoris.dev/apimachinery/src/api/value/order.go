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

// CompareOrdered compares two same-kind ordered scalar values.
//
// The bool result is false when the values have different kinds or when the
// kind has no descriptor-neutral ordering. Strings are deliberately excluded:
// objectquery string operators use prefix/suffix/contains instead of generic
// less/greater string ordering.
func CompareOrdered(left Value, right Value) (int, bool) {
	if left.Kind() != right.Kind() {
		return 0, false
	}

	switch left.Kind() {
	case KindInteger:
		l, _ := left.AsInteger()
		r, _ := right.AsInteger()
		return l.Compare(r), true
	case KindFloat:
		return compareFloats(left, right), true
	case KindDecimal:
		l, _ := left.AsDecimal()
		r, _ := right.AsDecimal()
		return l.Compare(r), true
	case KindTimestamp:
		l, _ := left.AsTimestamp()
		r, _ := right.AsTimestamp()
		return l.Compare(r), true
	case KindDate:
		l, _ := left.AsDate()
		r, _ := right.AsDate()
		return compareDates(l, r), true
	case KindTimeOfDay:
		l, _ := left.AsTimeOfDay()
		r, _ := right.AsTimeOfDay()
		return compareTimesOfDay(l, r), true
	case KindDuration:
		l, _ := left.AsDuration()
		r, _ := right.AsDuration()
		return compareOrderedValues(l, r), true
	default:
		return 0, false
	}
}

// compareFloats returns a three-way numeric comparison.
func compareFloats(left Value, right Value) int {
	l, _ := left.AsFloat()
	r, _ := right.AsFloat()
	return compareOrderedValues(l, r)
}

// compareDates compares calendar dates by stored calendar fields.
func compareDates(left Date, right Date) int {
	if cmp := compareOrderedValues(left.Year(), right.Year()); cmp != 0 {
		return cmp
	}
	if cmp := compareOrderedValues(left.Month(), right.Month()); cmp != 0 {
		return cmp
	}
	return compareOrderedValues(left.Day(), right.Day())
}

// compareTimesOfDay compares wall-clock times by stored time components.
func compareTimesOfDay(left TimeOfDay, right TimeOfDay) int {
	if cmp := compareOrderedValues(left.Hour(), right.Hour()); cmp != 0 {
		return cmp
	}
	if cmp := compareOrderedValues(left.Minute(), right.Minute()); cmp != 0 {
		return cmp
	}
	if cmp := compareOrderedValues(left.Second(), right.Second()); cmp != 0 {
		return cmp
	}
	return compareOrderedValues(left.Nanosecond(), right.Nanosecond())
}

// compareOrderedValues is the common three-way comparison for ordered Go
// primitives used by concrete Value payloads.
func compareOrderedValues[T ~int | ~int64 | ~float64](left T, right T) int {
	switch {
	case left < right:
		return -1
	case left > right:
		return 1
	default:
		return 0
	}
}
