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

package retrybudget

import (
	"math"
	"testing"
)

func TestAttemptsSnapshotTotal(t *testing.T) {
	tests := []struct {
		name string
		val  AttemptsSnapshot
		want uint64
	}{
		{name: "zero", val: AttemptsSnapshot{}, want: 0},
		{name: "sum", val: AttemptsSnapshot{Original: 2, Retry: 3}, want: 5},
		{name: "saturates", val: AttemptsSnapshot{Original: math.MaxUint64, Retry: 1}, want: math.MaxUint64},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.val.Total(); got != tt.want {
				t.Fatalf("Total() = %d, want %d", got, tt.want)
			}
		})
	}
}

func TestAttemptsSnapshotHasTraffic(t *testing.T) {
	tests := []struct {
		val  AttemptsSnapshot
		want bool
	}{
		{AttemptsSnapshot{}, false},
		{AttemptsSnapshot{Original: 1}, true},
		{AttemptsSnapshot{Retry: 1}, true},
	}
	for _, tt := range tests {
		if got := tt.val.HasTraffic(); got != tt.want {
			t.Fatalf("HasTraffic() = %v, want %v", got, tt.want)
		}
	}
}
