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

package retry

import (
	"testing"
	"time"
)

func TestAttemptZeroValueIsNotValid(t *testing.T) {
	var attempt Attempt

	if !attempt.IsZero() {
		t.Fatalf("zero Attempt IsZero() = false, want true")
	}
	if attempt.IsValid() {
		t.Fatalf("zero Attempt IsValid() = true, want false")
	}
	if attempt.IsFirst() {
		t.Fatalf("zero Attempt IsFirst() = true, want false")
	}
	if attempt.IsRetry() {
		t.Fatalf("zero Attempt IsRetry() = true, want false")
	}
}

func TestAttemptRequiresNonZeroNumber(t *testing.T) {
	attempt := Attempt{
		Number:    0,
		StartedAt: time.Unix(1, 0),
	}

	if attempt.IsZero() {
		t.Fatalf("Attempt with StartedAt set IsZero() = true, want false")
	}
	if attempt.IsValid() {
		t.Fatalf("Attempt with zero Number IsValid() = true, want false")
	}
	if attempt.IsFirst() {
		t.Fatalf("Attempt with zero Number IsFirst() = true, want false")
	}
	if attempt.IsRetry() {
		t.Fatalf("Attempt with zero Number IsRetry() = true, want false")
	}
}

func TestAttemptRequiresNonZeroStartedAt(t *testing.T) {
	attempt := Attempt{
		Number: 1,
	}

	if attempt.IsZero() {
		t.Fatalf("Attempt with Number set IsZero() = true, want false")
	}
	if attempt.IsValid() {
		t.Fatalf("Attempt with zero StartedAt IsValid() = true, want false")
	}
	if attempt.IsFirst() {
		t.Fatalf("Attempt with zero StartedAt IsFirst() = true, want false")
	}
	if attempt.IsRetry() {
		t.Fatalf("Attempt with zero StartedAt IsRetry() = true, want false")
	}
}

func TestAttemptFirstAttempt(t *testing.T) {
	attempt := Attempt{
		Number:    1,
		StartedAt: time.Unix(10, 0),
	}

	if attempt.IsZero() {
		t.Fatalf("valid first Attempt IsZero() = true, want false")
	}
	if !attempt.IsValid() {
		t.Fatalf("valid first Attempt IsValid() = false, want true")
	}
	if !attempt.IsFirst() {
		t.Fatalf("valid first Attempt IsFirst() = false, want true")
	}
	if attempt.IsRetry() {
		t.Fatalf("valid first Attempt IsRetry() = true, want false")
	}
}

func TestAttemptRetryAttempt(t *testing.T) {
	attempt := Attempt{
		Number:    2,
		StartedAt: time.Unix(20, 0),
	}

	if attempt.IsZero() {
		t.Fatalf("valid retry Attempt IsZero() = true, want false")
	}
	if !attempt.IsValid() {
		t.Fatalf("valid retry Attempt IsValid() = false, want true")
	}
	if attempt.IsFirst() {
		t.Fatalf("valid retry Attempt IsFirst() = true, want false")
	}
	if !attempt.IsRetry() {
		t.Fatalf("valid retry Attempt IsRetry() = false, want true")
	}
}

func TestAttemptHighNumberIsRetry(t *testing.T) {
	attempt := Attempt{
		Number:    100,
		StartedAt: time.Unix(30, 0),
	}

	if !attempt.IsValid() {
		t.Fatalf("high-number Attempt IsValid() = false, want true")
	}
	if attempt.IsFirst() {
		t.Fatalf("high-number Attempt IsFirst() = true, want false")
	}
	if !attempt.IsRetry() {
		t.Fatalf("high-number Attempt IsRetry() = false, want true")
	}
}
