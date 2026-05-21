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

func TestReservePanicsOnInvalidInput(t *testing.T) {
	t.Parallel()

	requirePanic(t, panicNilContext, func() {
		_, _, _ = Reserve(nil, time.Now(), 0)
	})
	requirePanic(t, panicNegativeDuration("reserve"), func() {
		_, _, _ = Reserve(context.Background(), time.Now(), -time.Nanosecond)
	})
}

func TestReserveSubtractsTailBudget(t *testing.T) {
	t.Parallel()

	now := testNow()

	canceled, cancel := context.WithCancel(context.Background())
	cancel()

	tests := []struct {
		name        string
		ctx         context.Context
		reserve     time.Duration
		want        time.Duration
		wantBounded bool
		wantOK      bool
	}{
		{
			name:        "no deadline",
			ctx:         context.Background(),
			want:        0,
			wantBounded: false,
			wantOK:      true,
		},
		{
			name:        "canceled context",
			ctx:         canceled,
			reserve:     time.Second,
			wantBounded: false,
			wantOK:      false,
		},
		{
			name: "canceled context with future deadline",
			ctx: func() context.Context {
				ctx, cancel := context.WithDeadline(context.Background(), now.Add(5*time.Second))
				t.Cleanup(cancel)
				cancel()
				return ctx
			}(),
			reserve:     time.Second,
			wantBounded: true,
			wantOK:      false,
		},
		{
			name:        "expired deadline",
			ctx:         contextWithDeadline(t, now.Add(-time.Second)),
			reserve:     time.Second,
			wantBounded: true,
			wantOK:      false,
		},
		{
			name:        "zero reserve",
			ctx:         contextWithDeadline(t, now.Add(5*time.Second)),
			reserve:     0,
			want:        5 * time.Second,
			wantBounded: true,
			wantOK:      true,
		},
		{
			name:        "reserve below remaining",
			ctx:         contextWithDeadline(t, now.Add(5*time.Second)),
			reserve:     2 * time.Second,
			want:        3 * time.Second,
			wantBounded: true,
			wantOK:      true,
		},
		{
			name:        "reserve equal remaining",
			ctx:         contextWithDeadline(t, now.Add(5*time.Second)),
			reserve:     5 * time.Second,
			wantBounded: true,
			wantOK:      false,
		},
		{
			name:        "reserve above remaining",
			ctx:         contextWithDeadline(t, now.Add(5*time.Second)),
			reserve:     6 * time.Second,
			wantBounded: true,
			wantOK:      false,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got, bounded, ok := Reserve(tt.ctx, now, tt.reserve)
			if got != tt.want || bounded != tt.wantBounded || ok != tt.wantOK {
				t.Fatalf(
					"Reserve() = (%v, %v, %v), want (%v, %v, %v)",
					got,
					bounded,
					ok,
					tt.want,
					tt.wantBounded,
					tt.wantOK,
				)
			}
		})
	}
}
