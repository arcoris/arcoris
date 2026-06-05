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

func TestResultOutcomeEffectGrantMatrix(t *testing.T) {
	t.Parallel()

	outcomes := []Outcome{
		OutcomeUnknown,
		OutcomeAdmitted,
		OutcomeDenied,
		OutcomeQueued,
		OutcomeDeferred,
	}
	effects := []Effect{
		EffectUnknown,
		EffectNone,
		EffectCommitted,
		EffectOwned,
		EffectQueued,
	}

	for _, outcome := range outcomes {
		for _, effect := range effects {
			for _, withGrant := range []bool{false, true} {
				for _, withMetadata := range []bool{false, true} {
					t.Run(matrixName(outcome, effect, withGrant, withMetadata), func(t *testing.T) {
						t.Parallel()

						decision := Decision{
							Outcome: outcome,
							Reason:  ReasonAdmitted,
							Effect:  effect,
						}
						result := resultWith(decision, "grant", withGrant, "metadata", withMetadata)
						want := validResultShape(outcome, effect, withGrant)
						if got := result.IsValid(); got != want {
							t.Fatalf("IsValid() = %t, want %t for %+v grant=%t metadata=%t", got, want, decision, withGrant, withMetadata)
						}
					})
				}
			}
		}
	}
}

func TestResultMetadataPresenceDoesNotAffectCoreValidity(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		decision Decision
		grant    bool
	}{
		{name: "admitted none", decision: AdmitDecision(ReasonAdmitted)},
		{name: "admitted committed", decision: CommitDecision(ReasonAdmitted)},
		{name: "admitted owned", decision: GrantDecision(ReasonAdmitted), grant: true},
		{name: "denied none", decision: DenyDecision(ReasonDenied)},
		{name: "deferred none", decision: DeferDecision(ReasonDeferred)},
		{name: "queued with grant", decision: QueueDecision(ReasonQueued), grant: true},
		{name: "queued without grant", decision: QueueDecision(ReasonQueued)},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			withMetadata := resultWith(tt.decision, "grant", tt.grant, "metadata", true)
			withoutMetadata := resultWith(tt.decision, "grant", tt.grant, "", false)
			if !withMetadata.IsValid() {
				t.Fatalf("with metadata is invalid: %+v", withMetadata.Decision())
			}
			if !withoutMetadata.IsValid() {
				t.Fatalf("without metadata is invalid: %+v", withoutMetadata.Decision())
			}
		})
	}
}

func TestResultInvalidReasonIsInvalidRegardlessOfPresence(t *testing.T) {
	t.Parallel()

	decisions := []Decision{
		{Outcome: OutcomeAdmitted, Reason: "bad-reason", Effect: EffectNone},
		{Outcome: OutcomeAdmitted, Reason: "bad-reason", Effect: EffectCommitted},
		{Outcome: OutcomeAdmitted, Reason: "bad-reason", Effect: EffectOwned},
		{Outcome: OutcomeDenied, Reason: "bad-reason", Effect: EffectNone},
		{Outcome: OutcomeDeferred, Reason: "bad-reason", Effect: EffectNone},
		{Outcome: OutcomeQueued, Reason: "bad-reason", Effect: EffectQueued},
	}

	for _, decision := range decisions {
		for _, withGrant := range []bool{false, true} {
			for _, withMetadata := range []bool{false, true} {
				t.Run(matrixName(decision.Outcome, decision.Effect, withGrant, withMetadata), func(t *testing.T) {
					t.Parallel()

					result := resultWith(decision, "grant", withGrant, "metadata", withMetadata)
					if result.IsValid() {
						t.Fatalf("result with invalid reason is valid: %+v", result.Decision())
					}
				})
			}
		}
	}
}

func validResultShape(outcome Outcome, effect Effect, grantPresent bool) bool {
	if !validDecisionShape(outcome, effect) {
		return false
	}
	switch effect {
	case EffectOwned:
		return grantPresent
	case EffectQueued:
		return true
	default:
		return !grantPresent
	}
}
