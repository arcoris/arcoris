// Copyright 2026 The ARCORIS Authors.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package deadline

import (
	"context"
	"testing"
	"time"

	panicassert "arcoris.dev/testutil/panic"
)

func TestReservePanicsOnInvalidInput(t *testing.T) {
	t.Parallel()

	panicassert.RequireErrorIs(t, ErrNilContext, func() {
		_, _, _ = Reserve(nil, time.Now(), 0)
	})
	panicassert.RequireErrorIs(t, ErrNegativeDuration, func() {
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

func TestReserveBudgetBuildsNamedResult(t *testing.T) {
	t.Parallel()

	now := testNow()

	canceled, cancel := context.WithCancel(context.Background())
	cancel()

	tests := []struct {
		name string
		ctx  context.Context
		res  time.Duration
		want ReserveResult
	}{
		{
			name: "no deadline",
			ctx:  context.Background(),
			want: ReserveResult{
				OK:     true,
				Reason: ReasonNoDeadline,
			},
		},
		{
			name: "context done",
			ctx:  canceled,
			want: ReserveResult{
				Reason: ReasonContextDone,
			},
		},
		{
			name: "expired deadline",
			ctx:  contextWithDeadline(t, now.Add(-time.Second)),
			res:  time.Second,
			want: ReserveResult{
				Bounded: true,
				Reason:  ReasonExpired,
				Budget: Budget{
					Deadline:    now.Add(-time.Second),
					HasDeadline: true,
					Expired:     true,
				},
			},
		},
		{
			name: "insufficient budget",
			ctx:  contextWithDeadline(t, now.Add(time.Second)),
			res:  time.Second,
			want: ReserveResult{
				Bounded: true,
				Reason:  ReasonInsufficientBudget,
				Budget: Budget{
					Deadline:    now.Add(time.Second),
					Remaining:   time.Second,
					HasDeadline: true,
				},
			},
		},
		{
			name: "allowed bounded",
			ctx:  contextWithDeadline(t, now.Add(5*time.Second)),
			res:  time.Second,
			want: ReserveResult{
				Duration: 4 * time.Second,
				Bounded:  true,
				OK:       true,
				Reason:   ReasonAllowed,
				Budget: Budget{
					Deadline:    now.Add(5 * time.Second),
					Remaining:   5 * time.Second,
					HasDeadline: true,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got := ReserveBudget(tt.ctx, now, tt.res)
			if got != tt.want {
				t.Fatalf("ReserveBudget() = %+v, want %+v", got, tt.want)
			}
			if !got.IsValid() {
				t.Fatalf("ReserveBudget result is invalid: %+v", got)
			}
		})
	}
}

func TestReserveResultIsValidRejectsInvalidShapes(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name   string
		result ReserveResult
	}{
		{name: "zero", result: ReserveResult{}},
		{
			name: "negative duration",
			result: ReserveResult{
				Duration: -time.Nanosecond,
				Bounded:  true,
				OK:       true,
				Reason:   ReasonAllowed,
			},
		},
		{
			name: "allowed unbounded with duration",
			result: ReserveResult{
				Duration: time.Second,
				OK:       true,
				Reason:   ReasonNoDeadline,
			},
		},
		{
			name: "denied with allowed reason",
			result: ReserveResult{
				Reason: ReasonAllowed,
			},
		},
		{
			name: "context done bounded without deadline",
			result: ReserveResult{
				Bounded: true,
				Reason:  ReasonContextDone,
			},
		},
		{
			name: "context done with expired budget",
			result: ReserveResult{
				Bounded: true,
				Reason:  ReasonContextDone,
				Budget: Budget{
					HasDeadline: true,
					Expired:     true,
				},
			},
		},
		{
			name: "unknown reason",
			result: ReserveResult{
				Reason: Reason(255),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			if tt.result.IsValid() {
				t.Fatalf("IsValid() = true, want false for %+v", tt.result)
			}
		})
	}
}
