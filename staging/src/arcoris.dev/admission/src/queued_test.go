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

func TestQueueDecision(t *testing.T) {
	t.Parallel()

	requireDecision(t, QueueDecision(ReasonQueued), Decision{
		Outcome: OutcomeQueued,
		Reason:  ReasonQueued,
		Effect:  EffectQueued,
	})
}

func TestQueuedResult(t *testing.T) {
	t.Parallel()

	result := QueuedResult(ReasonQueued, "ticket", "metadata")
	requireResultShape(t, result, QueueDecision(ReasonQueued), true, true)
}

func TestQueuedNoGrantResult(t *testing.T) {
	t.Parallel()

	result := QueuedNoGrantResult(ReasonQueued, "metadata")
	requireResultShape(t, result, QueueDecision(ReasonQueued), false, true)
}

func TestQueuedConstructorsWithInvalidReasonReturnInvalidValues(t *testing.T) {
	t.Parallel()

	invalid := Reason("bad-reason")
	if QueueDecision(invalid).IsValid() {
		t.Fatal("QueueDecision with invalid reason is valid")
	}
	if QueuedResult(invalid, "ticket", "metadata").IsValid() {
		t.Fatal("QueuedResult with invalid reason is valid")
	}
	if QueuedNoGrantResult(invalid, "metadata").IsValid() {
		t.Fatal("QueuedNoGrantResult with invalid reason is valid")
	}
}
