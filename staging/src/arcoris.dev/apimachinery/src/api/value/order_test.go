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

func TestCompareOrderedDurationUsesNumericDuration(t *testing.T) {
	cmp, ok := CompareOrdered(DurationValue(10*time.Second), DurationValue(2*time.Second))
	if !ok || cmp <= 0 {
		t.Fatalf("CompareOrdered(10s, 2s) = (%d, %v); want greater/true", cmp, ok)
	}
}

func TestCompareOrderedDateTimeAndDecimal(t *testing.T) {
	earlyDate := MustDateValue(mustDate(2024, time.January, 1))
	lateDate := MustDateValue(mustDate(2024, time.January, 2))
	if cmp, ok := CompareOrdered(earlyDate, lateDate); !ok || cmp >= 0 {
		t.Fatalf("date compare = (%d, %v); want less/true", cmp, ok)
	}

	earlyTime := MustTimeOfDayValue(mustTimeOfDay(1, 2, 3, 0))
	lateTime := MustTimeOfDayValue(mustTimeOfDay(1, 2, 4, 0))
	if cmp, ok := CompareOrdered(earlyTime, lateTime); !ok || cmp >= 0 {
		t.Fatalf("time compare = (%d, %v); want less/true", cmp, ok)
	}

	if cmp, ok := CompareOrdered(DecimalValue(MustParseDecimal("1.2")), DecimalValue(MustParseDecimal("1.10"))); !ok || cmp <= 0 {
		t.Fatalf("decimal compare = (%d, %v); want greater/true", cmp, ok)
	}
}

func TestCompareOrderedRejectsUnorderedOrMismatchedValues(t *testing.T) {
	if _, ok := CompareOrdered(StringValue("a"), StringValue("b")); ok {
		t.Fatal("string CompareOrdered ok = true; want false")
	}
	if _, ok := CompareOrdered(Int64Value(1), StringValue("1")); ok {
		t.Fatal("mismatched CompareOrdered ok = true; want false")
	}
}

func mustDate(year int, month time.Month, day int) Date {
	date, err := NewDate(year, month, day)
	if err != nil {
		panic(err)
	}
	return date
}

func mustTimeOfDay(hour, minute, second, nanosecond int) TimeOfDay {
	timeOfDay, err := NewTimeOfDay(hour, minute, second, nanosecond)
	if err != nil {
		panic(err)
	}
	return timeOfDay
}
