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


package core

import "testing"

func TestRangeLenAndEmpty(t *testing.T) {
	tests := []struct {
		name      string
		r         Range
		wantLen   int
		wantEmpty bool
	}{
		{name: "normal", r: Range{Start: 2, End: 7}, wantLen: 5},
		{name: "point", r: Range{Start: 4, End: 4}, wantEmpty: true},
		{name: "inverted", r: Range{Start: 8, End: 3}, wantEmpty: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.r.Len(); got != tt.wantLen {
				t.Fatalf("Len() = %d, want %d", got, tt.wantLen)
			}
			if got := tt.r.Empty(); got != tt.wantEmpty {
				t.Fatalf("Empty() = %v, want %v", got, tt.wantEmpty)
			}
		})
	}
}
