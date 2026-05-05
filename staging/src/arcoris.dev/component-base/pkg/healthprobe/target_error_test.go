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

package healthprobe

import (
	"errors"
	"testing"

	"arcoris.dev/component-base/pkg/health"
)

func TestDuplicateTargetErrorClassifiesWithErrorsIs(t *testing.T) {
	t.Parallel()

	err := DuplicateTargetError{
		Target:        health.TargetReady,
		Index:         2,
		PreviousIndex: 0,
	}

	if !errors.Is(err, ErrDuplicateTarget) {
		t.Fatalf("errors.Is(%v, ErrDuplicateTarget) = false, want true", err)
	}
}

func TestDuplicateTargetErrorSupportsErrorsAs(t *testing.T) {
	t.Parallel()

	err := error(DuplicateTargetError{
		Target:        health.TargetReady,
		Index:         2,
		PreviousIndex: 0,
	})

	var duplicateErr DuplicateTargetError
	if !errors.As(err, &duplicateErr) {
		t.Fatalf("errors.As(%T, DuplicateTargetError) = false, want true", err)
	}
	if duplicateErr.Target != health.TargetReady {
		t.Fatalf("Target = %s, want %s", duplicateErr.Target, health.TargetReady)
	}
	if duplicateErr.Index != 2 {
		t.Fatalf("Index = %d, want 2", duplicateErr.Index)
	}
	if duplicateErr.PreviousIndex != 0 {
		t.Fatalf("PreviousIndex = %d, want 0", duplicateErr.PreviousIndex)
	}
}

func TestDuplicateTargetErrorMessage(t *testing.T) {
	t.Parallel()

	err := DuplicateTargetError{
		Target:        health.TargetReady,
		Index:         2,
		PreviousIndex: 0,
	}

	if err.Error() == "" {
		t.Fatal("Error() returned empty message")
	}
}
