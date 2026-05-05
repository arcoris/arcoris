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
	"errors"
	"testing"
)

func TestRegistryErrorsClassifyWithErrorsIs(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name   string
		err    error
		target error
	}{
		{
			name:   "nil_checker",
			err:    NilCheckerError{Target: TargetReady, Index: 1},
			target: ErrNilChecker,
		},
		{
			name: "invalid_check_name",
			err: InvalidCheckNameError{
				Target: TargetReady,
				Index:  2,
				Name:   "bad-name",
				Err:    ErrInvalidCheckName,
			},
			target: ErrInvalidCheckName,
		},
		{
			name: "duplicate_check",
			err: DuplicateCheckError{
				Target:        TargetReady,
				Name:          "storage",
				Index:         3,
				PreviousIndex: -1,
			},
			target: ErrDuplicateCheck,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			if !errors.Is(tc.err, tc.target) {
				t.Fatalf("errors.Is(%v, %v) = false, want true", tc.err, tc.target)
			}
			if tc.err.Error() == "" {
				t.Fatal("Error() returned empty message")
			}
		})
	}
}

func TestInvalidCheckNameErrorUnwrapsNameCause(t *testing.T) {
	t.Parallel()

	err := InvalidCheckNameError{
		Target: TargetReady,
		Index:  0,
		Name:   "",
		Err:    ErrEmptyCheckName,
	}

	if !errors.Is(err, ErrEmptyCheckName) {
		t.Fatalf("errors.Is(%v, ErrEmptyCheckName) = false, want true", err)
	}
}

func TestDuplicateCheckErrorFormatsBatchAndExistingConflicts(t *testing.T) {
	t.Parallel()

	batch := DuplicateCheckError{
		Target:        TargetReady,
		Name:          "storage",
		Index:         2,
		PreviousIndex: 0,
	}
	existing := DuplicateCheckError{
		Target:        TargetReady,
		Name:          "storage",
		Index:         2,
		PreviousIndex: -1,
	}

	if batch.Error() == "" || existing.Error() == "" {
		t.Fatal("DuplicateCheckError messages must be non-empty")
	}
}
