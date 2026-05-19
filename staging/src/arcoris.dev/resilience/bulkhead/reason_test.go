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

package bulkhead

import (
	"errors"
	"testing"
)

func TestReasonString(t *testing.T) {
	t.Parallel()

	if got := ReasonFull.String(); got != "full" {
		t.Fatalf("String() = %q, want full", got)
	}
}

func TestReasonIsValid(t *testing.T) {
	t.Parallel()

	tests := []struct {
		reason Reason
		want   bool
	}{
		{reason: ReasonAllowed, want: true},
		{reason: ReasonFull, want: true},
		{reason: ReasonUnknown, want: false},
		{reason: Reason("other"), want: false},
	}

	for _, tt := range tests {
		if got := tt.reason.IsValid(); got != tt.want {
			t.Fatalf("%q IsValid() = %v, want %v", tt.reason, got, tt.want)
		}
	}
}

func TestReasonIsDenied(t *testing.T) {
	t.Parallel()

	if !ReasonFull.IsDenied() {
		t.Fatal("ReasonFull.IsDenied() = false, want true")
	}
	if ReasonAllowed.IsDenied() {
		t.Fatal("ReasonAllowed.IsDenied() = true, want false")
	}
}

func TestReasonErr(t *testing.T) {
	t.Parallel()

	if err := ReasonAllowed.Err(); err != nil {
		t.Fatalf("ReasonAllowed.Err() = %v, want nil", err)
	}
	if err := ReasonFull.Err(); !errors.Is(err, ErrFull) {
		t.Fatalf("ReasonFull.Err() = %v, want %v", err, ErrFull)
	}
}
