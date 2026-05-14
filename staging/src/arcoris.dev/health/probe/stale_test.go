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

package probe

import (
	"errors"
	"testing"
	"time"
)

func TestIsStale(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name       string
		age        time.Duration
		staleAfter time.Duration
		want       bool
	}{
		{
			name:       "disabled with zero stale after",
			age:        time.Hour,
			staleAfter: 0,
			want:       false,
		},
		{
			name:       "defensive disabled with negative stale after",
			age:        time.Hour,
			staleAfter: -time.Second,
			want:       false,
		},
		{
			name:       "negative age is not stale",
			age:        -time.Second,
			staleAfter: time.Second,
			want:       false,
		},
		{
			name:       "zero age is fresh",
			age:        0,
			staleAfter: time.Second,
			want:       false,
		},
		{
			name:       "age below stale after is fresh",
			age:        999 * time.Millisecond,
			staleAfter: time.Second,
			want:       false,
		},
		{
			name:       "age equal stale after is still fresh",
			age:        time.Second,
			staleAfter: time.Second,
			want:       false,
		},
		{
			name:       "age above stale after is stale",
			age:        time.Second + time.Nanosecond,
			staleAfter: time.Second,
			want:       true,
		},
		{
			name:       "large age above stale after is stale",
			age:        time.Hour,
			staleAfter: time.Second,
			want:       true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			if got := isStale(tc.age, tc.staleAfter); got != tc.want {
				t.Fatalf(
					"isStale(%s, %s) = %v, want %v",
					tc.age,
					tc.staleAfter,
					got,
					tc.want,
				)
			}
		})
	}
}

func TestValidateStaleAfter(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name       string
		staleAfter time.Duration
		wantErr    bool
	}{
		{
			name:       "zero disables stale detection",
			staleAfter: 0,
			wantErr:    false,
		},
		{
			name:       "positive stale after",
			staleAfter: time.Second,
			wantErr:    false,
		},
		{
			name:       "negative stale after",
			staleAfter: -time.Nanosecond,
			wantErr:    true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			err := validateStaleAfter(tc.staleAfter)

			if tc.wantErr {
				if !errors.Is(err, ErrInvalidStaleAfter) {
					t.Fatalf(
						"validateStaleAfter(%s) = %v, want ErrInvalidStaleAfter",
						tc.staleAfter,
						err,
					)
				}
				return
			}

			if err != nil {
				t.Fatalf("validateStaleAfter(%s) = %v, want nil", tc.staleAfter, err)
			}
		})
	}
}

func TestValidateStaleAfterReturnsTypedError(t *testing.T) {
	t.Parallel()

	const staleAfter = -time.Second

	err := validateStaleAfter(staleAfter)

	var staleErr InvalidStaleAfterError
	if !errors.As(err, &staleErr) {
		t.Fatalf(
			"errors.As(%T, InvalidStaleAfterError) = false, want true; err=%v",
			err,
			err,
		)
	}
	if staleErr.StaleAfter != staleAfter {
		t.Fatalf("StaleAfter = %s, want %s", staleErr.StaleAfter, staleAfter)
	}
}
