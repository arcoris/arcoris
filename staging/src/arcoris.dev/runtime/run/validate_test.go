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

package run

import (
	"context"
	"testing"
)

func TestValidationHelpersAcceptValidInputs(t *testing.T) {
	t.Parallel()

	requireValidationMessage("valid")
	requireContext(context.Background(), "valid context")
	requireGroup(NewGroup(context.Background()))
	requireTask(func(ctx context.Context) error { return nil })
	requireTaskName("worker")
	requireErrorMode(ErrorModeJoin)
	requireErrorMode(ErrorModeFirst)
	requireGroupOption(WithCancelOnError(true))
}

func TestValidationHelpersRejectMissingDiagnosticText(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		fn   func()
	}{
		{name: "message", fn: func() { requireValidationMessage("") }},
		{name: "context", fn: func() { requireContext(context.Background(), "") }},
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
		{name: "nil group", message: errNilGroup, fn: func() { requireGroup(nil) }},
		{name: "uninitialized group", message: errUninitializedGroup, fn: func() { requireGroup(&Group{}) }},
		{name: "nil task", message: errNilTask, fn: func() { requireTask(nil) }},
		{name: "empty task name", message: errEmptyTaskName, fn: func() { requireTaskName("") }},
		{name: "untrimmed task name", message: errUntrimmedTaskName, fn: func() { requireTaskName(" worker") }},
		{name: "invalid error mode", message: errInvalidErrorMode, fn: func() { requireErrorMode(ErrorMode(99)) }},
		{name: "nil group option", message: errNilGroupOption, fn: func() { requireGroupOption(nil) }},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			mustPanicWith(t, tc.message, tc.fn)
		})
	}
}

func TestNoCopyMarkerMethodsAreCallable(t *testing.T) {
	t.Parallel()

	var marker noCopy
	marker.Lock()
	marker.Unlock()
}
