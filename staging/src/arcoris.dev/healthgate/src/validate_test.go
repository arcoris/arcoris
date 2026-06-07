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

package healthgate

import (
	"errors"
	"testing"

	"arcoris.dev/health"
)

func TestNormalizeGateResultValidatesNameAndStructure(t *testing.T) {
	t.Parallel()

	result, err := normalizeGateResult("ready_gate", health.Healthy(""))
	if err != nil {
		t.Fatalf("normalizeGateResult(empty name) = %v, want nil", err)
	}
	if result.Name != "ready_gate" {
		t.Fatalf("normalized name = %q, want ready_gate", result.Name)
	}

	if _, err := normalizeGateResult("ready_gate", health.Healthy("other_gate")); !errors.Is(err, ErrMismatchedGateResult) {
		t.Fatalf("normalizeGateResult(mismatch) = %v, want ErrMismatchedGateResult", err)
	}
	if _, err := normalizeGateResult("ready_gate", health.Result{Status: health.Status(99)}); !errors.Is(err, ErrInvalidGateResult) {
		t.Fatalf("normalizeGateResult(invalid) = %v, want ErrInvalidGateResult", err)
	}
}
