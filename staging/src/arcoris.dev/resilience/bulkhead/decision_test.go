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

func TestDecisionIsValid(t *testing.T) {
	t.Parallel()

	l, err := New(1)
	if err != nil {
		t.Fatalf("New returned error: %v", err)
	}

	permit, allowed := l.TryAcquire()
	defer permit.Release()
	_, denied := l.TryAcquire()

	tests := []struct {
		name string
		dec  Decision
		want bool
	}{
		{name: "allowed", dec: allowed, want: true},
		{name: "denied", dec: denied, want: true},
		{name: "zero", dec: Decision{}, want: false},
		{name: "allowed with denied reason", dec: Decision{Allowed: true, Reason: ReasonFull, Snapshot: allowed.Snapshot}, want: false},
		{name: "denied with allowed reason", dec: Decision{Allowed: false, Reason: ReasonAllowed, Snapshot: denied.Snapshot}, want: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.dec.IsValid(); got != tt.want {
				t.Fatalf("IsValid() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDecisionErr(t *testing.T) {
	t.Parallel()

	l, err := New(1)
	if err != nil {
		t.Fatalf("New returned error: %v", err)
	}

	permit, allowed := l.TryAcquire()
	defer permit.Release()
	_, denied := l.TryAcquire()

	if err := allowed.Err(); err != nil {
		t.Fatalf("allowed Err() = %v, want nil", err)
	}
	if err := denied.Err(); !errors.Is(err, ErrFull) {
		t.Fatalf("denied Err() = %v, want %v", err, ErrFull)
	}
	if err := (Decision{}).Err(); !errors.Is(err, ErrInvalidDecision) {
		t.Fatalf("zero Err() = %v, want %v", err, ErrInvalidDecision)
	}
}
