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
	requireEqual(t, date.Equal(Date{year: 2024, month: time.February, day: 29}), true)
	requireEqual(t, date.Equal(Date{year: 2023, month: time.February, day: 28}), false)
}
