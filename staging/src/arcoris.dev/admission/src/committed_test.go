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

func TestCommitDecision(t *testing.T) {
	t.Parallel()

	requireDecision(t, CommitDecision(ReasonAdmitted), Decision{
		Outcome: OutcomeAdmitted,
		Reason:  ReasonAdmitted,
		Effect:  EffectCommitted,
	})
}

func TestCommittedResult(t *testing.T) {
	t.Parallel()

	result := CommittedResult(ReasonAdmitted, "metadata")
	requireResultShape(t, result, CommitDecision(ReasonAdmitted), false, true)
}

func TestCommittedNoMetadataResult(t *testing.T) {
	t.Parallel()

	result := CommittedNoMetadataResult(ReasonAdmitted)
	requireResultShape(t, result, CommitDecision(ReasonAdmitted), false, false)
}

func TestCommittedConstructorsWithInvalidReasonReturnInvalidValues(t *testing.T) {
	t.Parallel()

	invalid := Reason("bad-reason")
	if CommitDecision(invalid).IsValid() {
		t.Fatal("CommitDecision with invalid reason is valid")
	}
	if CommittedResult(invalid, "metadata").IsValid() {
		t.Fatal("CommittedResult with invalid reason is valid")
	}
	if CommittedNoMetadataResult(invalid).IsValid() {
		t.Fatal("CommittedNoMetadataResult with invalid reason is valid")
	}
}
