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

func TestInspectPanicsOnNilContext(t *testing.T) {
	t.Parallel()

	requirePanic(t, panicNilContext, func() {
		_ = Inspect(nil, time.Now())
	})
}

func TestInspectBuildsBudget(t *testing.T) {
	t.Parallel()

	now := testNow()
	future := now.Add(5 * time.Second)
	past := now.Add(-time.Second)

	tests := []struct {
		name string
		ctx  context.Context
		want Budget
	}{
		{
			name: "no deadline",
			ctx:  context.Background(),
			want: Budget{},
		},
		{
			name: "future deadline",
			ctx:  contextWithDeadline(t, future),
			want: Budget{
				Deadline:    future,
				Remaining:   5 * time.Second,
				HasDeadline: true,
			},
		},
		{
			name: "deadline equal now",
			ctx:  contextWithDeadline(t, now),
			want: Budget{
				Deadline:    now,
				HasDeadline: true,
				Expired:     true,
			},
		},
		{
			name: "expired deadline",
			ctx:  contextWithDeadline(t, past),
			want: Budget{
				Deadline:    past,
				HasDeadline: true,
				Expired:     true,
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got := Inspect(tt.ctx, now)
			if got != tt.want {
				t.Fatalf("Inspect() = %#v, want %#v", got, tt.want)
			}
		})
	}
}
