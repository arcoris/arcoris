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

func TestDateAccessorsAndString(t *testing.T) {
	date := Date{year: 2024, month: time.February, day: 29}

	requireEqual(t, date.Year(), 2024)
	requireEqual(t, date.Month(), time.February)
	requireEqual(t, date.Day(), 29)
	requireEqual(t, date.String(), "2024-02-29")
	requireEqual(
		t,
		date.Equal(Date{year: 2024, month: time.February, day: 29}),
		true,
	)
	requireEqual(
		t,
		date.Equal(Date{year: 2023, month: time.February, day: 28}),
		false,
	)
}

func TestDateIsValid(t *testing.T) {
	tests := []struct {
		name string
		date Date
		want bool
	}{
		{
			name: "zero",
			date: Date{},
			want: false,
		},
		{
			name: "valid leap day",
			date: Date{year: 2024, month: time.February, day: 29},
			want: true,
		},
		{
			name: "invalid month",
			date: Date{year: 2024, month: 13, day: 1},
			want: false,
		},
		{
			name: "invalid day",
			date: Date{year: 2023, month: time.February, day: 29},
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			requireEqual(t, tt.date.IsValid(), tt.want)
		})
	}
}
