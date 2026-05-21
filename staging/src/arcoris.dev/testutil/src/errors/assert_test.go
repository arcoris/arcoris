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

package errorassert

import (
	"errors"
	"fmt"
	"testing"
)

type valueError struct {
	message string
}

func (e valueError) Error() string {
	return e.message
}

type pointerError struct {
	message string
}

func (e *pointerError) Error() string {
	return e.message
}

func TestRequire(t *testing.T) {
	t.Parallel()

	Require(t, errors.New("boom"))
}

func TestRequireNone(t *testing.T) {
	t.Parallel()

	RequireNone(t, nil)
}

func TestRequireIs(t *testing.T) {
	t.Parallel()

	sentinel := errors.New("sentinel")
	RequireIs(t, sentinel, sentinel)
	RequireIs(t, fmt.Errorf("wrapped: %w", sentinel), sentinel)
}

func TestRequireIsNot(t *testing.T) {
	t.Parallel()

	RequireIsNot(t, errors.New("left"), errors.New("right"))
}

func TestRequireMessage(t *testing.T) {
	t.Parallel()

	RequireMessage(t, errors.New("boom"), "boom")
}

func TestRequireUnwrapsTo(t *testing.T) {
	t.Parallel()

	cause := errors.New("cause")
	RequireUnwrapsTo(t, fmt.Errorf("wrapped: %w", cause), cause)
}

func TestRequireAs(t *testing.T) {
	t.Parallel()

	t.Run("value error", func(t *testing.T) {
		t.Parallel()

		want := valueError{message: "value"}
		got := RequireAs[valueError](t, fmt.Errorf("wrapped: %w", want))
		if got != want {
			t.Fatalf("RequireAs() = %#v, want %#v", got, want)
		}
	})

	t.Run("pointer error", func(t *testing.T) {
		t.Parallel()

		want := &pointerError{message: "pointer"}
		got := RequireAs[*pointerError](t, fmt.Errorf("wrapped: %w", want))
		if got != want {
			t.Fatalf("RequireAs() = %#v, want %#v", got, want)
		}
	})

	t.Run("interface error", func(t *testing.T) {
		t.Parallel()

		want := errors.New("boom")
		got := RequireAs[error](t, fmt.Errorf("wrapped: %w", want))
		if !errors.Is(got, want) {
			t.Fatalf("RequireAs() = %v, want error matching %v", got, want)
		}
	})
}
