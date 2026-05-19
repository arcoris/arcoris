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
	"time"
)

func TestNewPanicsOnNilOption(t *testing.T) {
	t.Parallel()

	defer func() {
		if got := recover(); !errors.Is(asError(got), ErrNilOption) {
			t.Fatalf("panic = %v, want %v", got, ErrNilOption)
		}
	}()

	_, _ = New(1, nil)
}

func TestWithClockPanicsOnNilClock(t *testing.T) {
	t.Parallel()

	defer func() {
		if got := recover(); !errors.Is(asError(got), ErrNilClock) {
			t.Fatalf("panic = %v, want %v", got, ErrNilClock)
		}
	}()

	_ = WithClock(nil)
}

func TestWithClockControlsStampedPublicationTime(t *testing.T) {
	t.Parallel()

	now := time.Unix(123, 456)
	l, err := New(1, WithClock(fakeClock{now: now}))
	if err != nil {
		t.Fatalf("New returned error: %v", err)
	}

	stamped := l.Stamped()
	if !stamped.Updated.Equal(now) {
		t.Fatalf("Updated = %v, want %v", stamped.Updated, now)
	}
}

func asError(v any) error {
	if err, ok := v.(error); ok {
		return err
	}

	return nil
}
