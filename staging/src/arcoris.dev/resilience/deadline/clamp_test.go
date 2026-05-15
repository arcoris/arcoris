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

func TestClampPanicsOnInvalidInput(t *testing.T) {
	t.Parallel()

	requirePanic(t, panicNilContext, func() {
		_, _ = Clamp(nil, time.Now(), 0)
	})
	requirePanic(t, panicNegativeDuration("requested"), func() {
		_, _ = Clamp(context.Background(), time.Now(), -time.Nanosecond)
	})
}

func TestClampBoundsRequestedDuration(t *testing.T) {
	t.Parallel()

	now := testNow()

	canceled, cancel := context.WithCancel(context.Background())
	cancel()

	tests := []struct {
		name      string
		ctx       context.Context
		requested time.Duration
		want      time.Duration
		wantOK    bool
	}{
		{
			name:      "no deadline",
			ctx:       context.Background(),
			requested: time.Second,
			want:      time.Second,
			wantOK:    true,
		},
		{
			name:      "canceled context",
			ctx:       canceled,
			requested: time.Second,
			wantOK:    false,
		},
		{
			name: "canceled context with future deadline",
			ctx: func() context.Context {
				ctx, cancel := context.WithDeadline(context.Background(), now.Add(time.Second))
				t.Cleanup(cancel)
				cancel()
				return ctx
			}(),
			requested: time.Second,
			wantOK:    false,
		},
		{
			name:      "expired deadline",
			ctx:       contextWithDeadline(t, now.Add(-time.Second)),
			requested: time.Second,
			wantOK:    false,
		},
		{
			name:      "requested below remaining",
			ctx:       contextWithDeadline(t, now.Add(10*time.Second)),
			requested: 5 * time.Second,
			want:      5 * time.Second,
			wantOK:    true,
		},
		{
			name:      "requested equal remaining",
			ctx:       contextWithDeadline(t, now.Add(5*time.Second)),
			requested: 5 * time.Second,
			want:      5 * time.Second,
			wantOK:    true,
		},
		{
			name:      "requested above remaining",
			ctx:       contextWithDeadline(t, now.Add(2*time.Second)),
			requested: 5 * time.Second,
			want:      2 * time.Second,
			wantOK:    true,
		},
		{
			name:      "zero requested",
			ctx:       contextWithDeadline(t, now.Add(time.Second)),
			requested: 0,
			want:      0,
			wantOK:    true,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got, ok := Clamp(tt.ctx, now, tt.requested)
			if got != tt.want || ok != tt.wantOK {
				t.Fatalf("Clamp() = (%v, %v), want (%v, %v)", got, ok, tt.want, tt.wantOK)
			}
		})
	}
}
