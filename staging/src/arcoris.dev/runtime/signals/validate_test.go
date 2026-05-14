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

package signals

import (
	"context"
	"os"
	"testing"
)

func TestValidationHelpersAcceptValidInputs(t *testing.T) {
	t.Parallel()

	requireValidationMessage("valid")
	requireContext(context.Background(), "valid context")
	requireSignal(testSIGINT, "valid signal")
	requireSignals([]os.Signal{testSIGINT, testSIGTERM}, "valid signals")
	requireNonEmptySignals([]os.Signal{testSIGINT}, "valid non-empty signals")
	requirePositiveBuffer(1, "valid positive buffer")
	requireNonNegativeBuffer(0, "valid non-negative buffer")
}

func TestValidationHelpersRejectMissingDiagnosticText(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		fn   func()
	}{
		{name: "message", fn: func() { requireValidationMessage("") }},
		{name: "context", fn: func() { requireContext(context.Background(), "") }},
		{name: "signal", fn: func() { requireSignal(testSIGINT, "") }},
		{name: "signals", fn: func() { requireSignals([]os.Signal{testSIGINT}, "") }},
		{name: "non-empty signals", fn: func() { requireNonEmptySignals([]os.Signal{testSIGINT}, "") }},
		{name: "positive buffer", fn: func() { requirePositiveBuffer(1, "") }},
		{name: "non-negative buffer", fn: func() { requireNonNegativeBuffer(0, "") }},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			mustPanicWith(t, errNilValidationMessage, tc.fn)
		})
	}
}

func TestValidationHelpersRejectInvalidInputs(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		message string
		fn      func()
	}{
		{name: "nil context", message: "nil context", fn: func() { requireContext(nil, "nil context") }},
		{name: "nil signal", message: "nil signal", fn: func() { requireSignal(nil, "nil signal") }},
		{name: "nil signals", message: "nil signals", fn: func() { requireSignals([]os.Signal{nil}, "nil signals") }},
		{name: "empty signals", message: "empty signals", fn: func() { requireNonEmptySignals(nil, "empty signals") }},
		{name: "nil non-empty signals", message: "nil non-empty signals", fn: func() {
			requireNonEmptySignals([]os.Signal{nil}, "nil non-empty signals")
		}},
		{name: "zero positive buffer", message: "positive buffer", fn: func() { requirePositiveBuffer(0, "positive buffer") }},
		{name: "negative positive buffer", message: "positive buffer", fn: func() { requirePositiveBuffer(-1, "positive buffer") }},
		{name: "negative buffer", message: "non-negative buffer", fn: func() { requireNonNegativeBuffer(-1, "non-negative buffer") }},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			mustPanicWith(t, tc.message, tc.fn)
		})
	}
}
