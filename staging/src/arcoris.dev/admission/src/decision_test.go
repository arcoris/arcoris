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

package admission

import "testing"

func TestDecisionIsCopyableValue(t *testing.T) {
	t.Parallel()

	original := Decision{
		Outcome: OutcomeAdmitted,
		Reason:  ReasonAdmitted,
		Effect:  EffectOwned,
	}
	copied := original

	if copied != original {
		t.Fatalf("copied decision = %+v, want %+v", copied, original)
	}
}

func TestDecisionZeroValueIsInvalid(t *testing.T) {
	t.Parallel()

	var decision Decision
	if decision.IsValid() {
		t.Fatal("zero-value Decision is valid, want invalid")
	}
}
