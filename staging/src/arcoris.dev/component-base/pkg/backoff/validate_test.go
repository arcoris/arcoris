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

package backoff

import (
	"math"
	"testing"
	"time"
)

const testValidationMessage = "backoff: test validation failure"

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

func TestRequireRandomSourceRejectsNilSource(t *testing.T) {
	mustPanicWith(t, testValidationMessage, func() {
		requireRandomSource(nil, testValidationMessage)
	})
}

func TestRequireRandomRejectsNilRandom(t *testing.T) {
	mustPanicWith(t, testValidationMessage, func() {
		requireRandom(nil, testValidationMessage)
	})
}

func TestRequireRandomOptionRejectsNilOption(t *testing.T) {
	mustPanicWith(t, testValidationMessage, func() {
		requireRandomOption(nil, testValidationMessage)
	})
}

func TestRequireRandomConfigRejectsNilConfig(t *testing.T) {
	mustPanicWith(t, testValidationMessage, func() {
		requireRandomConfig(nil, testValidationMessage)
	})
}

func TestRequireNonNegativeDurationRejectsNegativeDuration(t *testing.T) {
	mustPanicWith(t, testValidationMessage, func() {
		requireNonNegativeDuration(-time.Nanosecond, testValidationMessage)
	})
}

func TestRequirePositiveDurationRejectsZeroAndNegativeDuration(t *testing.T) {
	mustPanicWith(t, testValidationMessage, func() {
		requirePositiveDuration(0, testValidationMessage)
	})
	mustPanicWith(t, testValidationMessage, func() {
		requirePositiveDuration(-time.Nanosecond, testValidationMessage)
	})
}

func TestRequireDurationNotBeforeRejectsUpperBeforeLower(t *testing.T) {
	mustPanicWith(t, testValidationMessage, func() {
		requireDurationNotBefore(time.Millisecond, time.Second, testValidationMessage)
	})
}

func TestRequireNonNegativeCountRejectsNegativeCount(t *testing.T) {
	mustPanicWith(t, testValidationMessage, func() {
		requireNonNegativeCount(-1, testValidationMessage)
	})
}

func TestRequireFiniteFloatRejectsNaNAndInfinity(t *testing.T) {
	mustPanicWith(t, testValidationMessage, func() {
		requireFiniteFloat(math.NaN(), testValidationMessage)
	})
	mustPanicWith(t, testValidationMessage, func() {
		requireFiniteFloat(math.Inf(1), testValidationMessage)
	})
}

func TestRequireFloatGreaterThanOneRejectsInvalidValues(t *testing.T) {
	for _, value := range []float64{1, 0, -1, math.NaN(), math.Inf(1)} {
		value := value
		t.Run("invalid", func(t *testing.T) {
			mustPanicWith(t, testValidationMessage, func() {
				requireFloatGreaterThanOne(value, testValidationMessage)
			})
		})
	}
}

func TestRequireJitterFactorRejectsInvalidValues(t *testing.T) {
	for _, value := range []float64{-1, math.NaN(), math.Inf(1)} {
		value := value
		t.Run("invalid", func(t *testing.T) {
			mustPanicWith(t, errInvalidJitterFactor, func() {
				requireJitterFactor(value)
			})
		})
	}
}

func TestRequireJitterRatioRejectsInvalidValues(t *testing.T) {
	for _, value := range []float64{-1, 1.1, math.NaN(), math.Inf(1)} {
		value := value
		t.Run("invalid", func(t *testing.T) {
			mustPanicWith(t, errInvalidJitterRatio, func() {
				requireJitterRatio(value)
			})
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
	requireRandomSource(defaultRandomSource(), testValidationMessage)
	requireRandom(fixedRandom(0), testValidationMessage)
	requireRandomOption(func(*randomConfig) {}, testValidationMessage)
	requireRandomConfig(&randomConfig{source: defaultRandomSource()}, testValidationMessage)
	requireNonNegativeDuration(0, testValidationMessage)
	requirePositiveDuration(time.Nanosecond, testValidationMessage)
	requireDurationNotBefore(time.Second, time.Second, testValidationMessage)
	requireNonNegativeCount(0, testValidationMessage)
	requireFiniteFloat(-1, testValidationMessage)
	requireFloatGreaterThanOne(1.1, testValidationMessage)
	requireJitterFactor(0)
	requireJitterRatio(1)
	requireNonNegativeSequenceDelay(0, true, testValidationMessage)
	requireNonNegativeSequenceDelay(-time.Nanosecond, false, testValidationMessage)
}
