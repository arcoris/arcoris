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

func TestDecisionIsValidAcceptsConstructedDecisions(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		decision Decision
	}{
		{name: "admit", decision: Admit(ReasonAdmitted)},
		{name: "commit", decision: Commit(ReasonAdmitted)},
		{name: "grant", decision: Grant(ReasonAdmitted)},
		{name: "deny", decision: Deny(Reason("capacity_exhausted"))},
		{name: "queue", decision: Queue(ReasonQueued)},
		{name: "defer", decision: Defer(ReasonDeferred)},
		{name: "admitted default", decision: Admitted()},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			if !tt.decision.IsValid() {
				t.Fatalf("decision should be valid: %+v", tt.decision)
			}
		})
	}
}

func TestDecisionIsValidRejectsInvalidCombinations(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		decision Decision
	}{
		{name: "zero decision", decision: Decision{}},
		{
			name: "unknown outcome",
			decision: Decision{
				Outcome: OutcomeUnknown,
				Reason:  ReasonAdmitted,
				Effect:  EffectNone,
			},
		},
		{
			name: "invalid reason",
			decision: Decision{
				Outcome: OutcomeAdmitted,
				Reason:  "",
				Effect:  EffectNone,
			},
		},
		{
			name: "unknown effect",
			decision: Decision{
				Outcome: OutcomeAdmitted,
				Reason:  ReasonAdmitted,
				Effect:  EffectUnknown,
			},
		},
		{
			name: "denied committed effect",
			decision: Decision{
				Outcome: OutcomeDenied,
				Reason:  Reason("capacity_exhausted"),
				Effect:  EffectCommitted,
			},
		},
		{
			name: "deferred owned effect",
			decision: Decision{
				Outcome: OutcomeDeferred,
				Reason:  ReasonDeferred,
				Effect:  EffectOwned,
			},
		},
		{
			name: "queued without queued effect",
			decision: Decision{
				Outcome: OutcomeQueued,
				Reason:  ReasonQueued,
				Effect:  EffectNone,
			},
		},
		{
			name: "admitted with queued effect",
			decision: Decision{
				Outcome: OutcomeAdmitted,
				Reason:  ReasonAdmitted,
				Effect:  EffectQueued,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			if tt.decision.IsValid() {
				t.Fatalf("decision should be invalid: %+v", tt.decision)
			}
		})
	}
}
