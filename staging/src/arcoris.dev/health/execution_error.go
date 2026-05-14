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
	"fmt"
)

var (
	// ErrInvalidExecutionPolicy identifies an unsupported evaluator execution
	// policy.
	//
	// Execution policy is validated at Evaluator construction time so Evaluate
	// never has to guess how to schedule checks for a target.
	ErrInvalidExecutionPolicy = errors.New("health: invalid execution policy")
)

// InvalidExecutionPolicyError describes an invalid evaluator execution policy.
//
// InvalidExecutionPolicyError is classified as ErrInvalidExecutionPolicy. Callers
// should use errors.Is for classification and inspect fields only for diagnostics.
type InvalidExecutionPolicyError struct {
	// Field identifies the invalid execution policy field.
	//
	// Expected values are "mode" and "max_concurrency".
	Field string

	// Mode is the configured execution mode.
	Mode ExecutionMode

	// MaxConcurrency is the configured parallel concurrency limit.
	MaxConcurrency int
}

// Error returns the invalid execution policy message.
func (e InvalidExecutionPolicyError) Error() string {
	return fmt.Sprintf(
		"%v: field=%s mode=%s max_concurrency=%d",
		ErrInvalidExecutionPolicy,
		e.Field,
		e.Mode,
		e.MaxConcurrency,
	)
}

// Is reports whether target matches the invalid execution policy classification.
func (e InvalidExecutionPolicyError) Is(target error) bool {
	return target == ErrInvalidExecutionPolicy
}
