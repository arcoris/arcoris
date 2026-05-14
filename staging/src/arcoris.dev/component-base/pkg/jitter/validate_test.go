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

package jitter

import (
	"math"
	"testing"
	"time"

	"arcoris.dev/component-base/pkg/delay"
)

const testValidationMessage = "jitter: test validation failure"

func TestRequireValidationMessageRejectsEmptyMessage(t *testing.T) {
	mustPanicWith(t, errNilValidationMessage, func() {
		requireValidationMessage("")
	})
}

func TestRequireValidationHelpersRejectInvalidInput(t *testing.T) {
	tests := []struct {
		name string
		fn   func()
	}{
		{name: "schedule", fn: func() { requireSchedule(nil, testValidationMessage) }},
		{name: "sequence", fn: func() { requireSequence(nil, testValidationMessage) }},
		{name: "source", fn: func() { requireRandomSource(nil, testValidationMessage) }},
		{name: "random", fn: func() { requireRandom(nil, testValidationMessage) }},
		{name: "option", fn: func() { requireRandomOption(nil, testValidationMessage) }},
		{name: "config", fn: func() { requireRandomConfig(nil, testValidationMessage) }},
		{name: "negative duration", fn: func() { requireNonNegativeDuration(-time.Nanosecond, testValidationMessage) }},
		{name: "zero positive duration", fn: func() { requirePositiveDuration(0, testValidationMessage) }},
		{name: "reversed range", fn: func() { requireDurationNotBefore(time.Millisecond, time.Second, testValidationMessage) }},
		{name: "invalid multiplier", fn: func() { requireFloatGreaterThanOne(1, testValidationMessage) }},
		{name: "negative sequence delay", fn: func() { requireNonNegativeSequenceDelay(-time.Nanosecond, true, testValidationMessage) }},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mustPanicWith(t, testValidationMessage, tt.fn)
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

func TestValidationHelpersAllowValidValues(t *testing.T) {
	source := RandomSourceFunc(func() RandomGenerator { return fixedRandom(0) })
	requireValidationMessage(testValidationMessage)
	requireSchedule(delay.Immediate(), testValidationMessage)
	requireSequence(delay.Immediate().NewSequence(), testValidationMessage)
	requireRandomSource(source, testValidationMessage)
	requireRandom(fixedRandom(0), testValidationMessage)
	requireRandomOption(func(*randomConfig) {}, testValidationMessage)
	requireRandomConfig(&randomConfig{source: source}, testValidationMessage)
	requireNonNegativeDuration(0, testValidationMessage)
	requirePositiveDuration(time.Nanosecond, testValidationMessage)
	requireDurationNotBefore(time.Second, time.Second, testValidationMessage)
	requireFloatGreaterThanOne(1.1, testValidationMessage)
	requireJitterFactor(0)
	requireJitterRatio(1)
	requireNonNegativeSequenceDelay(0, true, testValidationMessage)
	requireNonNegativeSequenceDelay(-time.Nanosecond, false, testValidationMessage)
}
