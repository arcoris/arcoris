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

package health

import (
	"errors"
	"testing"
)

func TestResultPredicateHelpersReflectResultFields(t *testing.T) {
	t.Parallel()

	cause := errors.New("cause")
	result := Degraded("storage", ReasonOverloaded, "overloaded").
		WithObserved(testObserved).
		WithCause(cause)

	if !result.IsNamed() || !result.IsObserved() || !result.HasCause() {
		t.Fatalf("field predicates mismatch for %+v", result)
	}
	if !result.HasReason(ReasonOverloaded) || result.IsNegative() || !result.IsOperational() {
		t.Fatalf("semantic predicates mismatch for %+v", result)
	}
}
