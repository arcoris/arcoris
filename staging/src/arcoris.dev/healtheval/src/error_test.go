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

package eval

import (
	"errors"
	"testing"

	"arcoris.dev/health"
)

func TestEvaluatorErrorMessages(t *testing.T) {
	t.Parallel()

	tests := []struct {
		err  error
		want string
	}{
		{ErrNilResolver, "healtheval: nil check resolver"},
		{ErrNilEvaluatorOption, "healtheval: nil evaluator option"},
		{ErrNilClock, "healtheval: nil clock"},
		{ErrInvalidTimeout, "healtheval: invalid timeout"},
		{ErrInvalidCheckResult, "healtheval: invalid check result"},
		{ErrMismatchedCheckResult, "healtheval: mismatched check result"},
		{ErrMismatchedResolvedTarget, "healtheval: mismatched resolved target"},
		{ErrInvalidExecutionPolicy, "healtheval: invalid execution policy"},
	}

	for _, tc := range tests {
		t.Run(tc.want, func(t *testing.T) {
			t.Parallel()

			if got := tc.err.Error(); got != tc.want {
				t.Fatalf("Error() = %q, want %q", got, tc.want)
			}
		})
	}
}

func TestMismatchedResolvedTargetErrorClassification(t *testing.T) {
	t.Parallel()

	err := MismatchedResolvedTargetError{
		Requested: health.TargetReady,
		Resolved:  health.TargetLive,
	}

	if !errors.Is(err, ErrMismatchedResolvedTarget) {
		t.Fatal("MismatchedResolvedTargetError should match ErrMismatchedResolvedTarget")
	}
}
