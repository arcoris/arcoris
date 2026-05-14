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
)

func TestDurationHelpersSaturateInsteadOfWrapping(t *testing.T) {
	tests := []struct {
		name string
		got  time.Duration
		want time.Duration
	}{
		{name: "add overflow", got: saturatingDurationAdd(maxDuration, time.Nanosecond), want: maxDuration},
		{name: "subtract floors at zero", got: saturatingDurationSub(time.Nanosecond, time.Second), want: 0},
		{name: "float multiply overflow", got: saturatingDurationMulFloat(time.Second, math.Inf(1)), want: maxDuration},
		{name: "float multiply NaN", got: saturatingDurationMulFloat(time.Second, math.NaN()), want: 0},
		{name: "min", got: minDuration(time.Second, time.Millisecond), want: time.Millisecond},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.got != tt.want {
				t.Fatalf("got %s, want %s", tt.got, tt.want)
			}
		})
	}
}
