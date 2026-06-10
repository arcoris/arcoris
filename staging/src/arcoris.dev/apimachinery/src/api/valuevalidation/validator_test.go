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

package valuevalidation

import (
	"errors"
	"testing"

	"arcoris.dev/apimachinery/api/fieldpath"
)

func TestValidatorResultReturnsNilWhenNoDiagnosticsCollected(t *testing.T) {
	run := newValidator(Options{})

	if err := run.result(); err != nil {
		t.Fatalf("result() = %v, want nil", err)
	}
}

func TestValidatorResultReturnsCollectedDiagnostics(t *testing.T) {
	run := newValidator(Options{})
	run.add(
		fieldpath.Root(),
		ErrInvalidValue,
		ErrorReasonInvalidZero,
		"value is invalid",
	)

	if !errors.Is(run.result(), ErrInvalidValue) {
		t.Fatalf("errors.Is(ErrInvalidValue) = false")
	}
}

func TestValidatorStopsAtConfiguredDiagnosticBudget(t *testing.T) {
	run := newValidator(Options{MaxErrors: 1})
	run.add(
		fieldpath.Root().Field(fieldpath.MustFieldName("first")),
		ErrMissingField,
		ErrorReasonMissingField,
		"first field is missing",
	)
	run.add(
		fieldpath.Root().Field(fieldpath.MustFieldName("second")),
		ErrMissingField,
		ErrorReasonMissingField,
		"second field is missing",
	)

	if got := run.errors.Len(); got != 1 {
		t.Fatalf("diagnostic count = %d, want 1", got)
	}
	if !run.shouldStop() {
		t.Fatalf("shouldStop() = false")
	}
}
