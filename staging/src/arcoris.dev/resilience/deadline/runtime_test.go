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

func TestActiveBudget(t *testing.T) {
	t.Parallel()

	now := testNow()

	canceled, cancel := context.WithCancel(context.Background())
	cancel()

	tests := []struct {
		name        string
		ctx         context.Context
		wantOK      bool
		wantBudget  Budget
		checkBudget func(t *testing.T, budget Budget)
	}{
		{
			name:   "active context without deadline",
			ctx:    context.Background(),
			wantOK: true,
		},
		{
			name:   "active context with future deadline",
			ctx:    contextWithDeadline(t, now.Add(5*time.Second)),
			wantOK: true,
			wantBudget: Budget{
				Deadline:    now.Add(5 * time.Second),
				Remaining:   5 * time.Second,
				HasDeadline: true,
			},
		},
		{
			name:   "expired deadline",
			ctx:    contextWithDeadline(t, now.Add(-time.Second)),
			wantOK: false,
			wantBudget: Budget{
				Deadline:    now.Add(-time.Second),
				HasDeadline: true,
				Expired:     true,
			},
		},
		{
			name:   "canceled context without deadline",
			ctx:    canceled,
			wantOK: false,
		},
		{
			name: "canceled context with future deadline",
			ctx: func() context.Context {
				ctx, cancel := context.WithDeadline(context.Background(), now.Add(5*time.Second))
				t.Cleanup(cancel)
				cancel()
				return ctx
			}(),
			wantOK: false,
			wantBudget: Budget{
				Deadline:    now.Add(5 * time.Second),
				Remaining:   5 * time.Second,
				HasDeadline: true,
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got, ok := activeBudget(tt.ctx, now)
			if got != tt.wantBudget || ok != tt.wantOK {
				t.Fatalf("activeBudget() = (%#v, %v), want (%#v, %v)", got, ok, tt.wantBudget, tt.wantOK)
			}
		})
	}
}
