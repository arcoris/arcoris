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
)

func TestInvalidExecutionPolicyErrorClassifiesWithErrorsIs(t *testing.T) {
	t.Parallel()

	err := InvalidExecutionPolicyError{
		Field:          "max_concurrency",
		Mode:           ExecutionParallel,
		MaxConcurrency: 0,
	}

	if !errors.Is(err, ErrInvalidExecutionPolicy) {
		t.Fatalf("errors.Is(%v, ErrInvalidExecutionPolicy) = false, want true", err)
	}
}

func TestInvalidExecutionPolicyErrorSupportsErrorsAs(t *testing.T) {
	t.Parallel()

	err := error(InvalidExecutionPolicyError{
		Field:          "mode",
		Mode:           ExecutionMode(99),
		MaxConcurrency: 1,
	})

	var policyErr InvalidExecutionPolicyError
	if !errors.As(err, &policyErr) {
		t.Fatalf("errors.As(%T, InvalidExecutionPolicyError) = false, want true", err)
	}
	if policyErr.Field != "mode" {
		t.Fatalf("Field = %q, want mode", policyErr.Field)
	}
	if policyErr.Mode != ExecutionMode(99) {
		t.Fatalf("Mode = %s, want invalid mode", policyErr.Mode)
	}
}

func TestInvalidExecutionPolicyErrorMessage(t *testing.T) {
	t.Parallel()

	err := InvalidExecutionPolicyError{
		Field:          "max_concurrency",
		Mode:           ExecutionParallel,
		MaxConcurrency: 0,
	}

	if err.Error() == "" {
		t.Fatal("Error() returned empty message")
	}
}
