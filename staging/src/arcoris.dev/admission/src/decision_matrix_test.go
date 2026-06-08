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

func TestDecisionOutcomeEffectMatrix(t *testing.T) {
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
			t.Run(outcome.String()+"_"+effect.String(), func(t *testing.T) {
				t.Parallel()

				decision := Decision{
					Outcome: outcome,
					Reason:  ReasonAdmitted,
					Effect:  effect,
				}
				if got, want := decision.IsValid(), validDecisionShape(outcome, effect); got != want {
					t.Fatalf("Decision{%s,%s}.IsValid() = %t, want %t", outcome, effect, got, want)
				}
			})
		}
	}
}

func TestDecisionInvalidReasonIsInvalidForEveryOutcomeEffect(t *testing.T) {
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
			t.Run(outcome.String()+"_"+effect.String(), func(t *testing.T) {
				t.Parallel()

				decision := Decision{
					Outcome: outcome,
					Reason:  "bad-reason",
					Effect:  effect,
				}
				if decision.IsValid() {
					t.Fatalf("decision with invalid reason is valid: %+v", decision)
				}
			})
		}
	}
}

func TestDecisionStateHelpersDelegateToOutcomeAndEffect(t *testing.T) {
	t.Parallel()

	decision := Decision{
		Outcome: OutcomeQueued,
		Reason:  "bad-reason",
		Effect:  EffectQueued,
	}

	if !decision.IsQueued() {
		t.Fatal("IsQueued() = false, want true")
	}
	if decision.IsTerminal() {
		t.Fatal("IsTerminal() = true, want false")
	}
	if !decision.HasSideEffect() {
		t.Fatal("HasSideEffect() = false, want true")
	}
	if !decision.AllowsGrant() {
		t.Fatal("AllowsGrant() = false, want true")
	}
	if decision.RequiresGrant() {
		t.Fatal("RequiresGrant() = true, want false")
	}
	if decision.IsValid() {
		t.Fatal("helper state should not make invalid reason valid")
	}
}

func validDecisionShape(outcome Outcome, effect Effect) bool {
	if outcome == OutcomeUnknown || effect == EffectUnknown {
		return false
	}
	switch outcome {
	case OutcomeAdmitted:
		return effect == EffectNone || effect == EffectCommitted || effect == EffectOwned
	case OutcomeDenied, OutcomeDeferred:
		return effect == EffectNone
	case OutcomeQueued:
		return effect == EffectQueued
	default:
		return false
	}
}
