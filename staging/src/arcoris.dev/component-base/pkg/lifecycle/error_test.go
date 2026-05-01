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

package lifecycle

import (
	"strings"
	"testing"
)

func TestSentinelErrorsNonNil(t *testing.T) {
	t.Parallel()

	for _, err := range lifecycleSentinelErrors() {
		if err == nil {
			t.Fatal("sentinel error is nil")
		}
	}
}

func TestSentinelErrorStringsStable(t *testing.T) {
	t.Parallel()

	// Stable sentinel text keeps wrapped lifecycle errors recognizable in logs
	// without depending on concrete error wrapper types.
	want := []string{
		"lifecycle: invalid transition",
		"lifecycle: terminal state",
		"lifecycle: failure cause required",
		"lifecycle: guard rejected",
		"lifecycle: invalid wait predicate",
		"lifecycle: invalid wait target",
		"lifecycle: wait target unreachable",
	}
	for i, err := range lifecycleSentinelErrors() {
		if got := err.Error(); got != want[i] {
			t.Fatalf("sentinel[%d] = %q, want %q", i, got, want[i])
		}
		if !strings.HasPrefix(err.Error(), "lifecycle:") {
			t.Fatalf("sentinel[%d] = %q, want lifecycle prefix", i, err.Error())
		}
	}
}

func lifecycleSentinelErrors() []error {
	return []error{
		ErrInvalidTransition,
		ErrTerminalState,
		ErrFailureCauseRequired,
		ErrGuardRejected,
		ErrInvalidWaitPredicate,
		ErrInvalidWaitTarget,
		ErrWaitTargetUnreachable,
	}
}
