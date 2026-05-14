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

package delay

import (
	"math"
	"testing"
	"time"
)

const testValidationMessage = "delay: test validation failure"

func TestRequireValidationMessageRejectsEmptyMessage(t *testing.T) {
	mustPanicWith(t, errNilValidationMessage, func() {
		requireValidationMessage("")
	})
}

func TestRequireScheduleRejectsNilSchedule(t *testing.T) {
	mustPanicWith(t, testValidationMessage, func() {
		requireSchedule(nil, testValidationMessage)
	})
}

func TestRequireSequenceRejectsNilSequence(t *testing.T) {
	mustPanicWith(t, testValidationMessage, func() {
		requireSequence(nil, testValidationMessage)
	})
}

func TestRequireNonNegativeDurationRejectsNegativeDuration(t *testing.T) {
	mustPanicWith(t, testValidationMessage, func() {
		requireNonNegativeDuration(-time.Nanosecond, testValidationMessage)
	})
}

func TestRequirePositiveDurationRejectsNonPositiveDuration(t *testing.T) {
	mustPanicWith(t, testValidationMessage, func() {
		requirePositiveDuration(0, testValidationMessage)
	})
	mustPanicWith(t, testValidationMessage, func() {
		requirePositiveDuration(-time.Nanosecond, testValidationMessage)
	})
}

func TestRequireNonNegativeCountRejectsNegativeCount(t *testing.T) {
	mustPanicWith(t, testValidationMessage, func() {
		requireNonNegativeCount(-1, testValidationMessage)
	})
}

func TestRequireFloatGreaterThanOneRejectsInvalidValue(t *testing.T) {
	for _, v := range []float64{1, 0.5, math.NaN(), math.Inf(1), math.Inf(-1)} {
		mustPanicWith(t, testValidationMessage, func() {
			requireFloatGreaterThanOne(v, testValidationMessage)
		})
	}
}

func TestRequireNonNegativeSequenceDelayRejectsNegativeAvailableDelay(t *testing.T) {
	mustPanicWith(t, testValidationMessage, func() {
		requireNonNegativeSequenceDelay(-time.Nanosecond, true, testValidationMessage)
	})
}

func TestValidationHelpersAllowValidValues(t *testing.T) {
	requireValidationMessage(testValidationMessage)
	requireSchedule(Immediate(), testValidationMessage)
	requireSequence(Immediate().NewSequence(), testValidationMessage)
	requireNonNegativeDuration(0, testValidationMessage)
	requirePositiveDuration(time.Nanosecond, testValidationMessage)
	requireNonNegativeCount(0, testValidationMessage)
	requireFloatGreaterThanOne(1.1, testValidationMessage)
	requireNonNegativeSequenceDelay(0, true, testValidationMessage)
	requireNonNegativeSequenceDelay(-time.Nanosecond, false, testValidationMessage)
}
