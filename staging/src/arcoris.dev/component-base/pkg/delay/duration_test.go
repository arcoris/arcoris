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

func TestDurationHelpersSaturateInsteadOfWrapping(t *testing.T) {
	tests := []struct {
		name string
		got  time.Duration
		want time.Duration
	}{
		{name: "add overflow", got: saturatingDurationAdd(maxDuration, time.Nanosecond), want: maxDuration},
		{name: "subtract floors at zero", got: saturatingDurationSub(time.Nanosecond, time.Second), want: 0},
		{name: "multiply overflow", got: saturatingDurationMul(maxDuration, 2), want: maxDuration},
		{name: "float conversion overflow", got: durationFromFloat(maxDurationFloat * 2), want: maxDuration},
		{name: "float conversion infinity", got: durationFromFloat(math.Inf(1)), want: maxDuration},
		{name: "float conversion NaN", got: durationFromFloat(math.NaN()), want: 0},
		{name: "cap", got: capDuration(2*time.Second, time.Second), want: time.Second},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.got != tt.want {
				t.Fatalf("got %s, want %s", tt.got, tt.want)
			}
		})
	}
}

func TestDurationSignPredicates(t *testing.T) {
	if !isNegativeDuration(-time.Nanosecond) {
		t.Fatal("negative duration was not reported negative")
	}
	if !isNonNegativeDuration(0) {
		t.Fatal("zero duration was not reported non-negative")
	}
}
