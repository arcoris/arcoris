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
	"testing"
	"time"
)

func TestBudgetHasBudget(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		b    Budget
		want bool
	}{
		{
			name: "no deadline",
			b:    Budget{},
			want: false,
		},
		{
			name: "expired deadline",
			b: Budget{
				HasDeadline: true,
				Expired:     true,
			},
			want: false,
		},
		{
			name: "zero remaining active deadline",
			b: Budget{
				HasDeadline: true,
			},
			want: false,
		},
		{
			name: "positive remaining active deadline",
			b: Budget{
				HasDeadline: true,
				Remaining:   time.Second,
			},
			want: true,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			if got := tt.b.HasBudget(); got != tt.want {
				t.Fatalf("HasBudget() = %v, want %v", got, tt.want)
			}
		})
	}
}
