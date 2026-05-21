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

package eval

import (
	"errors"
	"testing"

	"arcoris.dev/health"
)

func TestCheckResultErrorsClassifyWithErrorsIs(t *testing.T) {
	t.Parallel()

	invalid := InvalidCheckResultError{
		CheckName: "storage",
		Result: health.Result{
			Status: health.StatusHealthy,
			Reason: health.Reason("bad-reason"),
		},
	}
	mismatch := MismatchedCheckResultError{
		CheckName:  "storage",
		ResultName: "database",
	}

	if !errors.Is(invalid, ErrInvalidCheckResult) {
		t.Fatalf("errors.Is(%v, ErrInvalidCheckResult) = false, want true", invalid)
	}
	if !errors.Is(mismatch, ErrMismatchedCheckResult) {
		t.Fatalf("errors.Is(%v, ErrMismatchedCheckResult) = false, want true", mismatch)
	}
	if invalid.Error() == "" || mismatch.Error() == "" {
		t.Fatal("check result error messages must be non-empty")
	}
}
