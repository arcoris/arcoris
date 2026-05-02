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
	"errors"
	"strings"
	"testing"
)

func TestTaskErrorMatchesErrTaskFailedAndUnwrapsOriginal(t *testing.T) {
	t.Parallel()

	original := context.DeadlineExceeded
	err := TaskError{Name: "worker", Err: original}

	if !errors.Is(err, ErrTaskFailed) {
		t.Fatal("TaskError does not match ErrTaskFailed")
	}
	if !errors.Is(err, original) {
		t.Fatal("TaskError does not unwrap original error")
	}
	if !strings.Contains(err.Error(), "worker") {
		t.Fatalf("TaskError string = %q, want task name", err.Error())
	}
}

func TestTaskErrorHandlesEmptyName(t *testing.T) {
	t.Parallel()

	err := TaskError{Err: context.Canceled}

	if !strings.Contains(err.Error(), "run task failed") {
		t.Fatalf("TaskError string = %q, want generic task text", err.Error())
	}
}

func TestTaskErrorsExtractsJoinedTaskErrors(t *testing.T) {
	t.Parallel()

	first := TaskError{Name: "first", Err: context.Canceled}
	second := TaskError{Name: "second", Err: context.DeadlineExceeded}

	got := TaskErrors(errors.Join(first, second))
	if len(got) != 2 {
		t.Fatalf("TaskErrors len = %d, want 2", len(got))
	}
	if got[0].Name != "first" || got[1].Name != "second" {
		t.Fatalf("TaskErrors = %+v, want first then second", got)
	}
}

func TestTaskErrorsReturnsNilForNilAndNonTaskErrors(t *testing.T) {
	t.Parallel()

	if got := TaskErrors(nil); got != nil {
		t.Fatalf("TaskErrors(nil) = %+v, want nil", got)
	}

	if got := TaskErrors(context.Canceled); len(got) != 0 {
		t.Fatalf("TaskErrors(non-task) len = %d, want 0", len(got))
	}
}
