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

import "testing"

func TestKindString(t *testing.T) {
	tests := []struct {
		kind Kind
		want string
	}{
		{KindUnknown, "unknown"},
		{KindNoop, "noop"},
		{KindFixedWindow, "fixed_window"},
		{Kind(99), "unknown"},
	}
	for _, tt := range tests {
		if got := tt.kind.String(); got != tt.want {
			t.Fatalf("%v.String() = %q, want %q", uint8(tt.kind), got, tt.want)
		}
	}
}

func TestKindIsValid(t *testing.T) {
	tests := []struct {
		kind Kind
		want bool
	}{
		{KindUnknown, false},
		{KindNoop, true},
		{KindFixedWindow, true},
		{Kind(99), false},
	}
	for _, tt := range tests {
		if got := tt.kind.IsValid(); got != tt.want {
			t.Fatalf("%v.IsValid() = %v, want %v", uint8(tt.kind), got, tt.want)
		}
	}
}
