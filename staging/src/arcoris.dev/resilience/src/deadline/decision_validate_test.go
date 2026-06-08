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
	"testing"
	"time"
)

func TestDecisionIsValid(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		decision Decision
		want     bool
	}{
		{
			name: "allowed with budget",
			decision: Decision{
				Allowed:   true,
				Remaining: time.Second,
				Reason:    ReasonAllowed,
			},
			want: true,
		},
		{
			name: "allowed without deadline",
			decision: Decision{
				Allowed: true,
				Reason:  ReasonNoDeadline,
			},
			want: true,
		},
		{
			name: "context done denial without deadline budget",
			decision: Decision{
				Reason: ReasonContextDone,
			},
			want: true,
		},
		{
			name: "context done denial with future deadline budget",
			decision: Decision{
				Remaining: time.Second,
				Reason:    ReasonContextDone,
			},
			want: true,
		},
		{
			name: "expired denial without remaining budget",
			decision: Decision{
				Reason: ReasonExpired,
			},
			want: true,
		},
		{
			name: "insufficient budget denial",
			decision: Decision{
				Remaining: time.Second,
				Reason:    ReasonInsufficientBudget,
			},
			want: true,
		},
		{
			name:     "zero decision",
			decision: Decision{},
		},
		{
			name: "allowed with denied reason",
			decision: Decision{
				Allowed:   true,
				Remaining: time.Second,
				Reason:    ReasonExpired,
			},
		},
		{
			name: "denied with allowed reason",
			decision: Decision{
				Reason: ReasonAllowed,
			},
		},
		{
			name: "denied with no-deadline reason",
			decision: Decision{
				Reason: ReasonNoDeadline,
			},
		},
		{
			name: "allowed without positive budget",
			decision: Decision{
				Allowed: true,
				Reason:  ReasonAllowed,
			},
		},
		{
			name: "no-deadline decision with remaining budget",
			decision: Decision{
				Allowed:   true,
				Remaining: time.Second,
				Reason:    ReasonNoDeadline,
			},
		},
		{
			name: "expired denial with remaining budget",
			decision: Decision{
				Remaining: time.Second,
				Reason:    ReasonExpired,
			},
		},
		{
			name: "insufficient budget denial without remaining budget",
			decision: Decision{
				Reason: ReasonInsufficientBudget,
			},
		},
		{
			name: "negative remaining",
			decision: Decision{
				Remaining: -time.Nanosecond,
				Reason:    ReasonExpired,
			},
		},
		{
			name: "insufficient budget denial with negative remaining",
			decision: Decision{
				Remaining: -time.Nanosecond,
				Reason:    ReasonInsufficientBudget,
			},
		},
		{
			name: "unknown reason",
			decision: Decision{
				Reason: Reason(255),
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			if got := test.decision.IsValid(); got != test.want {
				t.Fatalf("IsValid() = %v, want %v", got, test.want)
			}
		})
	}
}
