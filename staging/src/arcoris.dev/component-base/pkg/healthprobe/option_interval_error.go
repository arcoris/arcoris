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

package healthprobe

import (
	"errors"
	"fmt"
	"time"
)

var (
	// ErrInvalidInterval identifies an invalid fixed probe interval.
	//
	// Runner intervals must be positive. Zero and negative values cannot drive a
	// meaningful ticker loop.
	ErrInvalidInterval = errors.New("healthprobe: invalid interval")
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
