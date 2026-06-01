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

package value

import "testing"

func TestFormatTimeOfDay(t *testing.T) {
	tests := []struct {
		name string
		time TimeOfDay
		want string
	}{
		{name: "midnight", time: TimeOfDay{}, want: "00:00:00"},
		{name: "whole second", time: TimeOfDay{hour: 1, minute: 2, second: 3}, want: "01:02:03"},
		{
			name: "with nanoseconds",
			time: TimeOfDay{hour: 1, minute: 2, second: 3, nanosecond: 4},
			want: "01:02:03.000000004",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := formatTimeOfDay(tt.time)

			requireEqual(t, got, tt.want)
		})
	}
}
