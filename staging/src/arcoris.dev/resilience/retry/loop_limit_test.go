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

package retry

import (
	"testing"
	"time"

	"arcoris.dev/chrono/clock"
)

func TestRetryExecutionMaxElapsedWouldBeExceeded(t *testing.T) {
	startedAt := time.Unix(100, 0)

	tests := []struct {
		name       string
		elapsed    time.Duration
		maxElapsed time.Duration
		delay      time.Duration
		want       bool
	}{
		{
			name:       "disabled",
			elapsed:    time.Hour,
			maxElapsed: 0,
			delay:      time.Hour,
			want:       false,
		},
		{
			name:       "elapsed already reached",
			elapsed:    time.Second,
			maxElapsed: time.Second,
			delay:      time.Nanosecond,
			want:       true,
		},
		{
			name:       "delay before remaining budget",
			elapsed:    time.Second,
			maxElapsed: 3 * time.Second,
			delay:      time.Second,
			want:       false,
		},
		{
			name:       "delay equals remaining budget",
			elapsed:    time.Second,
			maxElapsed: 2 * time.Second,
			delay:      time.Second,
			want:       true,
		},
		{
			name:       "delay after remaining budget",
			elapsed:    time.Second,
			maxElapsed: 2 * time.Second,
			delay:      time.Second + time.Nanosecond,
			want:       true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fake := clock.NewFakeClock(startedAt.Add(tt.elapsed))
			execution := &retryExecution{
				config: config{
					clock:      fake,
					maxElapsed: tt.maxElapsed,
				},
				startedAt: startedAt,
			}

			got := execution.maxElapsedWouldBeExceeded(tt.delay)
			if got != tt.want {
				t.Fatalf("maxElapsedWouldBeExceeded(%v) = %v, want %v", tt.delay, got, tt.want)
			}
		})
	}
}
