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

import (
	"testing"
	"time"
)

func TestFormatDate(t *testing.T) {
	tests := []struct {
		name  string
		year  int
		month time.Month
		day   int
		want  string
	}{
		{name: "positive", year: 2024, month: time.February, day: 9, want: "2024-02-09"},
		{name: "negative padded", year: -1, month: time.January, day: 2, want: "-001-01-02"},
		{name: "wide year", year: 12024, month: time.December, day: 31, want: "12024-12-31"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := formatDate(tt.year, tt.month, tt.day)

			requireEqual(t, got, tt.want)
		})
	}
}

func TestFormatDateMonth(t *testing.T) {
	requireEqual(t, formatDateMonth(2024, time.February), "2024-02")
	requireEqual(t, formatDateMonth(-1, time.January), "-001-01")
}
