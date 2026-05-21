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

func TestRemainingPanicsOnNilContext(t *testing.T) {
	t.Parallel()

	requirePanic(t, panicNilContext, func() {
		_, _ = Remaining(nil, time.Now())
	})
}

func TestRemainingReportsDeadlineBudget(t *testing.T) {
	t.Parallel()

	now := testNow()

	tests := []struct {
		name         string
		ctx          context.Context
		want         time.Duration
		wantDeadline bool
	}{
		{
			name:         "no deadline",
			ctx:          context.Background(),
			want:         0,
			wantDeadline: false,
		},
		{
			name:         "active deadline",
			ctx:          contextWithDeadline(t, now.Add(10*time.Second)),
			want:         10 * time.Second,
			wantDeadline: true,
		},
		{
			name:         "expired deadline",
			ctx:          contextWithDeadline(t, now.Add(-time.Second)),
			want:         0,
			wantDeadline: true,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got, ok := Remaining(tt.ctx, now)
			if got != tt.want || ok != tt.wantDeadline {
				t.Fatalf("Remaining() = (%v, %v), want (%v, %v)", got, ok, tt.want, tt.wantDeadline)
			}
		})
	}
}
