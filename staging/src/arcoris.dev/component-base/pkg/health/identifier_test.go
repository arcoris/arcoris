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

package health

import "testing"

func TestValidLowerSnakeIdentifier(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name  string
		value string
		max   int
		want  bool
	}{
		{name: "simple", value: "storage", max: 32, want: true},
		{name: "digits_after_first", value: "pool_1", max: 32, want: true},
		{name: "empty", value: "", max: 32, want: false},
		{name: "too_long", value: "storage", max: 3, want: false},
		{name: "digit_first", value: "1pool", max: 32, want: false},
		{name: "leading_underscore", value: "_pool", max: 32, want: false},
		{name: "trailing_underscore", value: "pool_", max: 32, want: false},
		{name: "repeated_underscore", value: "pool__main", max: 32, want: false},
		{name: "uppercase", value: "Pool", max: 32, want: false},
		{name: "dash", value: "pool-main", max: 32, want: false},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			if got := validLowerSnakeIdentifier(test.value, test.max); got != test.want {
				t.Fatalf("validLowerSnakeIdentifier(%q, %d) = %v, want %v", test.value, test.max, got, test.want)
			}
		})
	}
}
