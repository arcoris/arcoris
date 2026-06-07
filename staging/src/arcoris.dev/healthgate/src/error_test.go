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

func TestGateErrorsClassifySentinels(t *testing.T) {
	t.Parallel()

	invalid := InvalidGateResultError{
		GateName: "ready_gate",
		Result:   health.Result{Status: health.Status(99)},
	}
	if !errors.Is(invalid, ErrInvalidGateResult) {
		t.Fatal("InvalidGateResultError should match ErrInvalidGateResult")
	}

	mismatch := MismatchedGateResultError{
		GateName:   "ready_gate",
		ResultName: "other_gate",
	}
	if !errors.Is(mismatch, ErrMismatchedGateResult) {
		t.Fatal("MismatchedGateResultError should match ErrMismatchedGateResult")
	}
}
