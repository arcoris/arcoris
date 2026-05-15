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

package deadline

import (
	"context"
	"testing"
	"time"
)

func TestCanStartPanicsOnInvalidInput(t *testing.T) {
	t.Parallel()

	requirePanic(t, panicNilContext, func() {
		_ = CanStart(nil, time.Now(), 0)
	})
	requirePanic(t, panicNegativeDuration("min"), func() {
		_ = CanStart(context.Background(), time.Now(), -time.Nanosecond)
	})
}

func TestCanStartBuildsDecision(t *testing.T) {
	t.Parallel()

	now := testNow()

	canceled, cancel := context.WithCancel(context.Background())
	cancel()

	tests := []struct {
		name string
		ctx  context.Context
		min  time.Duration
		want Decision
	}{
		{
			name: "canceled context",
			ctx:  canceled,
			want: Decision{Reason: ReasonContextDone},
		},
		{
			name: "canceled context with future deadline",
			ctx: func() context.Context {
				ctx, cancel := context.WithDeadline(context.Background(), now.Add(time.Second))
				t.Cleanup(cancel)
				cancel()
				return ctx
			}(),
			want: Decision{
				Remaining: time.Second,
				Reason:    ReasonContextDone,
			},
		},
		{
			name: "no deadline",
			ctx:  context.Background(),
			min:  time.Hour,
			want: Decision{Allowed: true, Reason: ReasonNoDeadline},
		},
		{
			name: "expired deadline",
			ctx:  contextWithDeadline(t, now.Add(-time.Second)),
			want: Decision{Reason: ReasonExpired},
		},
		{
			name: "expired deadline takes precedence over context done",
			ctx: func() context.Context {
				ctx, cancel := context.WithDeadline(context.Background(), now.Add(-time.Second))
				t.Cleanup(cancel)
				cancel()
				return ctx
			}(),
			want: Decision{Reason: ReasonExpired},
		},
		{
			name: "remaining below minimum",
			ctx:  contextWithDeadline(t, now.Add(10*time.Millisecond)),
			min:  20 * time.Millisecond,
			want: Decision{Remaining: 10 * time.Millisecond, Reason: ReasonInsufficientBudget},
		},
		{
			name: "remaining equal minimum",
			ctx:  contextWithDeadline(t, now.Add(20*time.Millisecond)),
			min:  20 * time.Millisecond,
			want: Decision{Allowed: true, Remaining: 20 * time.Millisecond, Reason: ReasonAllowed},
		},
		{
			name: "remaining above minimum",
			ctx:  contextWithDeadline(t, now.Add(30*time.Millisecond)),
			min:  20 * time.Millisecond,
			want: Decision{Allowed: true, Remaining: 30 * time.Millisecond, Reason: ReasonAllowed},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got := CanStart(tt.ctx, now, tt.min)
			if got != tt.want {
				t.Fatalf("CanStart() = %#v, want %#v", got, tt.want)
			}
		})
	}
}
