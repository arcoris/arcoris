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
	grantPresence := []bool{false, true}
	metadataPresence := []bool{false, true}

	for _, outcome := range outcomes {
		for _, effect := range effects {
			for _, withGrant := range grantPresence {
				for _, withMetadata := range metadataPresence {
					t.Run(matrixName(outcome, effect, withGrant, withMetadata), func(t *testing.T) {
						t.Parallel()

						decision := Decision{
							Outcome: outcome,
							Reason:  ReasonAdmitted,
							Effect:  effect,
						}
						result := matrixResult(withGrant, withMetadata, decision)
						want := validResultShape(outcome, effect, withGrant)
						if got := result.IsValid(); got != want {
							t.Fatalf("Result{%s,%s,grant=%t,metadata=%t}.IsValid() = %t, want %t",
								outcome,
								effect,
								withGrant,
								withMetadata,
								got,
								want,
							)
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
		name  string
		grant Maybe[string]
		dec   Decision
	}{
		{name: "admitted none", grant: noneString(), dec: Admit(ReasonAdmitted)},
		{name: "admitted committed", grant: noneString(), dec: Commit(ReasonAdmitted)},
		{name: "admitted owned", grant: someString("grant"), dec: Grant(ReasonAdmitted)},
		{name: "denied none", grant: noneString(), dec: Deny(ReasonDenied)},
		{name: "deferred none", grant: noneString(), dec: Defer(ReasonDeferred)},
		{name: "queued with grant", grant: someString("ticket"), dec: Queue(ReasonQueued)},
		{name: "queued without grant", grant: noneString(), dec: Queue(ReasonQueued)},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			withMetadata := resultWith(tt.dec, tt.grant, someString("metadata"))
			withoutMetadata := resultWith(tt.dec, tt.grant, noneString())
			if !withMetadata.IsValid() {
				t.Fatalf("with metadata should be valid: %+v", withMetadata.Decision())
			}
			if !withoutMetadata.IsValid() {
				t.Fatalf("without metadata should be valid: %+v", withoutMetadata.Decision())
			}
		})
	}
}

func TestResultMetadataAbsenceIsValidForAllValidCoreShapes(t *testing.T) {
	t.Parallel()

	tests := []Result[string, string]{
		resultWith(Admit(ReasonAdmitted), noneString(), noneString()),
		resultWith(Commit(ReasonAdmitted), noneString(), noneString()),
		resultWith(Grant(ReasonAdmitted), someString("grant"), noneString()),
		resultWith(Deny(ReasonDenied), noneString(), noneString()),
		resultWith(Defer(ReasonDeferred), noneString(), noneString()),
		resultWith(Queue(ReasonQueued), someString("ticket"), noneString()),
		resultWith(Queue(ReasonQueued), noneString(), noneString()),
	}

	for i, result := range tests {
		if !result.IsValid() {
			t.Fatalf("result[%d] without metadata is invalid: %+v", i, result.Decision())
		}
	}
}

func TestResultZeroValueIsInvalid(t *testing.T) {
	t.Parallel()

	var result Result[string, string]
	if result.IsValid() {
		t.Fatal("zero-value Result is valid, want invalid")
	}
}

func TestResultInvalidReasonIsInvalidRegardlessOfGrantShape(t *testing.T) {
	t.Parallel()

	grantPresence := []bool{false, true}
	metadataPresence := []bool{false, true}
	decisions := []Decision{
		{Outcome: OutcomeAdmitted, Reason: "bad-reason", Effect: EffectNone},
		{Outcome: OutcomeAdmitted, Reason: "bad-reason", Effect: EffectCommitted},
		{Outcome: OutcomeAdmitted, Reason: "bad-reason", Effect: EffectOwned},
		{Outcome: OutcomeDenied, Reason: "bad-reason", Effect: EffectNone},
		{Outcome: OutcomeDeferred, Reason: "bad-reason", Effect: EffectNone},
		{Outcome: OutcomeQueued, Reason: "bad-reason", Effect: EffectQueued},
	}

	for _, decision := range decisions {
		for _, withGrant := range grantPresence {
			for _, withMetadata := range metadataPresence {
				t.Run(matrixName(decision.Outcome, decision.Effect, withGrant, withMetadata), func(t *testing.T) {
					t.Parallel()

					result := matrixResult(withGrant, withMetadata, decision)
					if result.IsValid() {
						t.Fatalf("result with invalid reason is valid: %+v", result.Decision())
					}
				})
			}
		}
	}
}

func matrixResult(grantPresent bool, metadataPresent bool, decision Decision) Result[string, string] {
	grant := noneString()
	if grantPresent {
		grant = someString("grant")
	}

	metadata := noneString()
	if metadataPresent {
		metadata = someString("metadata")
	}

	return resultWith(decision, grant, metadata)
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

func matrixName(outcome Outcome, effect Effect, grantPresent bool, metadataPresent bool) string {
	name := outcome.String() + "_" + effect.String()
	if grantPresent {
		name += "_grant"
	} else {
		name += "_no_grant"
	}
	if metadataPresent {
		name += "_metadata"
	} else {
		name += "_no_metadata"
	}
	return name
}
