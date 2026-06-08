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

package retrybudgetadmission

import (
	"errors"
	"testing"
)

func TestAdmitterRequireBudgetRejectsNilBudget(t *testing.T) {
	t.Parallel()

	requirePanicError(t, ErrNilRetryAdmitter, func() {
		Admitter{}.requireBudget()
	})
}

func TestAdmitterRequireBudgetAcceptsConfiguredBudget(t *testing.T) {
	t.Parallel()

	New(&scriptedBudget{}).requireBudget()
}

func requirePanicError(t *testing.T, want error, fn func()) {
	t.Helper()

	defer func() {
		recovered := recover()
		if recovered == nil {
			t.Fatalf("panic = nil, want %v", want)
		}
		err, ok := recovered.(error)
		if !ok {
			t.Fatalf("panic = %T(%v), want error %v", recovered, recovered, want)
		}
		if !errors.Is(err, want) {
			t.Fatalf("panic = %v, want %v", err, want)
		}
	}()

	fn()
}
