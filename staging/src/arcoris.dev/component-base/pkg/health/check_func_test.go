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

package health

import (
	"context"
	"errors"
	"testing"
)

func TestNewCheckValidatesInputs(t *testing.T) {
	t.Parallel()

	if _, err := NewCheck("bad-name", func(context.Context) Result { return Healthy("") }); !errors.Is(err, ErrInvalidCheckName) {
		t.Fatalf("NewCheck(invalid name) = %v, want ErrInvalidCheckName", err)
	}

	if _, err := NewCheck("storage", nil); !errors.Is(err, ErrNilCheckFunc) {
		t.Fatalf("NewCheck(nil func) = %v, want ErrNilCheckFunc", err)
	}
}

func TestCheckFuncCheckerFillsEmptyResultName(t *testing.T) {
	t.Parallel()

	checker, err := NewCheck("storage", func(context.Context) Result {
		return Healthy("")
	})
	if err != nil {
		t.Fatalf("NewCheck() = %v, want nil", err)
	}

	if got := checker.Name(); got != "storage" {
		t.Fatalf("Name() = %q, want storage", got)
	}

	result := checker.Check(context.Background())
	if result.Name != "storage" {
		t.Fatalf("result name = %q, want storage", result.Name)
	}
}

func TestNewErrorCheckMapsErrors(t *testing.T) {
	t.Parallel()

	cause := errors.New("disk failed")
	checker, err := NewErrorCheck("storage", func(context.Context) error {
		return cause
	})
	if err != nil {
		t.Fatalf("NewErrorCheck() = %v, want nil", err)
	}

	result := checker.Check(context.Background())
	if result.Status != StatusUnhealthy {
		t.Fatalf("status = %s, want unhealthy", result.Status)
	}
	if result.Reason != ReasonFatal {
		t.Fatalf("reason = %s, want fatal", result.Reason)
	}
	if !errors.Is(result.Cause, cause) {
		t.Fatalf("cause = %v, want %v", result.Cause, cause)
	}
}

func TestNewErrorCheckValidatesNameBeforeFunc(t *testing.T) {
	t.Parallel()

	_, err := NewErrorCheck("bad-name", nil)
	if !errors.Is(err, ErrInvalidCheckName) {
		t.Fatalf("NewErrorCheck(invalid name, nil func) = %v, want ErrInvalidCheckName", err)
	}
}

func TestNewErrorCheckMapsNilErrorToHealthy(t *testing.T) {
	t.Parallel()

	checker, err := NewErrorCheck("storage", func(context.Context) error {
		return nil
	})
	if err != nil {
		t.Fatalf("NewErrorCheck() = %v, want nil", err)
	}

	if result := checker.Check(context.Background()); result.Status != StatusHealthy {
		t.Fatalf("status = %s, want healthy", result.Status)
	}
}

func TestMustCheckPanicsOnInvalidInput(t *testing.T) {
	t.Parallel()

	mustPanicWith(t, ErrNilCheckFunc, func() {
		MustCheck("storage", nil)
	})
}

func TestMustCheckReturnsCheckerOnValidInput(t *testing.T) {
	t.Parallel()

	checker := MustCheck("storage", func(context.Context) Result {
		return Healthy("storage")
	})
	if checker.Name() != "storage" {
		t.Fatalf("Name() = %q, want storage", checker.Name())
	}
}

func TestMustErrorCheckPanicsOnInvalidInput(t *testing.T) {
	t.Parallel()

	mustPanicWith(t, ErrNilCheckFunc, func() {
		MustErrorCheck("storage", nil)
	})
}

func TestMustErrorCheckReturnsCheckerOnValidInput(t *testing.T) {
	t.Parallel()

	checker := MustErrorCheck("storage", func(context.Context) error {
		return nil
	})
	if result := checker.Check(context.Background()); result.Status != StatusHealthy {
		t.Fatalf("status = %s, want healthy", result.Status)
	}
}
