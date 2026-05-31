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

// NewTimeOfDay constructs a wall-clock time without leap-second support.
//
// The accepted range is exactly 00:00:00.000000000 through
// 23:59:59.999999999. Calendar, timezone, and daylight-saving rules are outside
// this value type.
func NewTimeOfDay(hour, minute, second, nanosecond int) (TimeOfDay, error) {
	switch {
	case hour < 0 || hour > 23:
		return TimeOfDay{}, invalidTimeOfDay("hour %d is outside 0..23", hour)
	case minute < 0 || minute > 59:
		return TimeOfDay{}, invalidTimeOfDay("minute %d is outside 0..59", minute)
	case second < 0 || second > 59:
		return TimeOfDay{}, invalidTimeOfDay("second %d is outside 0..59", second)
	case nanosecond < 0 || nanosecond > 999999999:
		return TimeOfDay{}, invalidTimeOfDay("nanosecond %d is outside 0..999999999", nanosecond)
	default:
		return TimeOfDay{hour: hour, minute: minute, second: second, nanosecond: nanosecond}, nil
	}
}

// invalidTimeOfDay returns a structured time-of-day construction error.
//
// All time-of-day range failures share the same path and reason; the detail
// explains the specific component that failed.
func invalidTimeOfDay(format string, args ...any) error {
	return errorf(pathTimeOfDay, ErrInvalidTime, ErrorReasonInvalidTime, format, args...)
}
