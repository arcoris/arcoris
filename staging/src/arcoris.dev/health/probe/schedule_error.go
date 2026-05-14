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

package probe

import (
	"errors"
	"fmt"
	"reflect"
	"time"

	"arcoris.dev/chrono/delay"
)

var (
	// ErrInvalidInterval identifies an invalid fixed probe interval.
	//
	// Runner intervals must be positive. Zero and negative values are reserved
	// for direct schedule configuration.
	ErrInvalidInterval = errors.New("healthprobe: invalid interval")

	// ErrNilSchedule identifies a nil probe schedule.
	ErrNilSchedule = errors.New("healthprobe: nil schedule")

	// ErrNilSequence identifies a schedule that returned a nil sequence.
	ErrNilSequence = errors.New("healthprobe: nil schedule sequence")

	// ErrExhaustedSchedule identifies a schedule sequence that ended before
	// context cancellation.
	ErrExhaustedSchedule = errors.New("healthprobe: exhausted schedule")

	// ErrInvalidScheduleDelay identifies a schedule sequence that returned a
	// negative delay.
	ErrInvalidScheduleDelay = errors.New("healthprobe: invalid schedule delay")
)

// InvalidIntervalError describes an invalid fixed probe interval.
//
// InvalidIntervalError is classified as ErrInvalidInterval. Callers should use
// errors.Is for classification and inspect Interval only for diagnostics.
type InvalidIntervalError struct {
	// Interval is the invalid interval value.
	Interval time.Duration
}

// Error returns the invalid interval message.
func (e InvalidIntervalError) Error() string {
	return fmt.Sprintf("%v: interval=%s", ErrInvalidInterval, e.Interval)
}

// Is reports whether target matches the invalid interval classification.
func (e InvalidIntervalError) Is(target error) bool {
	return target == ErrInvalidInterval
}

// validateInterval validates the fixed probe interval.
func validateInterval(interval time.Duration) error {
	if interval <= 0 {
		return InvalidIntervalError{Interval: interval}
	}

	return nil
}

// InvalidScheduleDelayError describes a negative delay returned by a schedule
// sequence.
//
// InvalidScheduleDelayError is classified as ErrInvalidScheduleDelay. Negative
// delays violate the delay.Sequence contract and cannot be safely passed to
// clock timers.
type InvalidScheduleDelayError struct {
	// Delay is the invalid delay value.
	Delay time.Duration
}

// Error returns the invalid schedule delay message.
func (e InvalidScheduleDelayError) Error() string {
	return fmt.Sprintf("%v: delay=%s", ErrInvalidScheduleDelay, e.Delay)
}

// Is reports whether target matches the invalid schedule delay classification.
func (e InvalidScheduleDelayError) Is(target error) bool {
	return target == ErrInvalidScheduleDelay
}

// validateSchedule validates the reusable probe schedule.
func validateSchedule(schedule delay.Schedule) error {
	if nilSchedule(schedule) {
		return ErrNilSchedule
	}

	return nil
}

// newSequence creates and validates one per-Run schedule sequence.
func newSequence(schedule delay.Schedule) (delay.Sequence, error) {
	if err := validateSchedule(schedule); err != nil {
		return nil, err
	}

	sequence := schedule.NewSequence()
	if sequence == nil {
		return nil, ErrNilSequence
	}

	return sequence, nil
}

func nilSchedule(schedule delay.Schedule) bool {
	if schedule == nil {
		return true
	}

	value := reflect.ValueOf(schedule)
	switch value.Kind() {
	case reflect.Chan, reflect.Func, reflect.Interface, reflect.Map, reflect.Pointer, reflect.Slice:
		return value.IsNil()
	default:
		return false
	}
}
