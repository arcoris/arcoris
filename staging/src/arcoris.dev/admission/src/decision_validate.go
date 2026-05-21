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

// IsValid reports whether d satisfies admission decision invariants.
//
// The invariant ties Outcome and Effect together. Admitted decisions may have no
// side effect, a committed spend-only side effect, or an owned grant side
// effect. Denied and deferred decisions must be side-effect-free because the
// caller keeps responsibility for the rejected or postponed work. Queued
// decisions must use EffectQueued because the system has accepted ownership of
// waiting work.
func (d Decision) IsValid() bool {
	if !d.hasValidFields() {
		return false
	}
	if d.Outcome == OutcomeAdmitted {
		return d.hasAdmittedEffect()
	}
	if d.Outcome == OutcomeDenied || d.Outcome == OutcomeDeferred {
		return d.Effect == EffectNone
	}

	return d.Effect == EffectQueued
}

// hasValidFields reports whether all closed/open vocabularies are syntactically
// valid before cross-field semantic checks run.
func (d Decision) hasValidFields() bool {
	return d.Outcome.IsValid() && d.Reason.IsValid() && d.Effect.IsValid()
}

// hasAdmittedEffect reports whether an admitted decision uses one of the effect
// classes allowed for immediately admitted work.
func (d Decision) hasAdmittedEffect() bool {
	return d.Effect == EffectNone ||
		d.Effect == EffectCommitted ||
		d.Effect == EffectOwned
}
