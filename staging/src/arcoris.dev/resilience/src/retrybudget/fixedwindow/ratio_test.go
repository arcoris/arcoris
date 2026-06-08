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

package fixedwindow

import (
	"errors"
	"testing"

	"arcoris.dev/resilience/retrybudget"
)

func TestNewRatioDelegatesToRetryBudgetRatio(t *testing.T) {
	t.Parallel()

	got, err := NewRatio(2, 10)
	if err != nil {
		t.Fatalf("NewRatio() error = %v", err)
	}
	if want := retrybudget.MustRatio(1, 5); got != want {
		t.Fatalf("NewRatio() = %v, want %v", got, want)
	}
}

func TestNewRatioRejectsInvalidRatio(t *testing.T) {
	t.Parallel()

	_, err := NewRatio(2, 1)
	if !errors.Is(err, retrybudget.ErrInvalidRatio) {
		t.Fatalf("NewRatio() error = %v, want %v", err, retrybudget.ErrInvalidRatio)
	}
}

func TestMustRatioPanicsOnInvalidRatio(t *testing.T) {
	t.Parallel()

	requirePanicError(t, retrybudget.ErrInvalidRatio, func() {
		_ = MustRatio(2, 1)
	})
}
