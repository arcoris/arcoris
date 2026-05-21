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

package fixedwindow

import (
	"math"
	"testing"
)

func TestAllowedRetries(t *testing.T) {
	tests := []struct {
		name     string
		original uint64
		ratio    float64
		min      uint64
		want     uint64
	}{
		{name: "no original", original: 0, ratio: 0.2, min: 10, want: 10},
		{name: "zero ratio", original: 100, ratio: 0, min: 7, want: 7},
		{name: "floor", original: 9, ratio: 0.2, min: 1, want: 2},
		{name: "exact", original: 10, ratio: 0.2, min: 1, want: 3},
		{name: "one ratio", original: 10, ratio: 1, min: 1, want: 11},
		{name: "saturates", original: math.MaxUint64, ratio: 1, min: 1, want: math.MaxUint64},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := allowedRetries(tt.original, tt.ratio, tt.min)
			if got != tt.want {
				t.Fatalf("allowedRetries() = %d, want %d", got, tt.want)
			}
		})
	}
}

func TestAvailableRetries(t *testing.T) {
	tests := []struct {
		allowed uint64
		used    uint64
		want    uint64
	}{
		{allowed: 10, used: 0, want: 10},
		{allowed: 10, used: 3, want: 7},
		{allowed: 10, used: 10, want: 0},
		{allowed: 10, used: 11, want: 0},
	}

	for _, tt := range tests {
		got := availableRetries(tt.allowed, tt.used)
		if got != tt.want {
			t.Fatalf("availableRetries(%d, %d) = %d, want %d", tt.allowed, tt.used, got, tt.want)
		}
	}
}

func TestSaturatingInc(t *testing.T) {
	if got := saturatingInc(0); got != 1 {
		t.Fatalf("saturatingInc(0) = %d, want 1", got)
	}
	if got := saturatingInc(math.MaxUint64 - 1); got != math.MaxUint64 {
		t.Fatalf("saturatingInc(max-1) = %d, want max", got)
	}
	if got := saturatingInc(math.MaxUint64); got != math.MaxUint64 {
		t.Fatalf("saturatingInc(max) = %d, want max", got)
	}
}
