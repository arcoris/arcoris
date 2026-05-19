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

// Admit returns an admitted decision with no side effect.
//
// Use Admit when the component only reports that work may proceed and does not
// transfer a grant or commit a budget/token spend. The returned Decision is the
// semantic base for Accepted and AcceptedNoMetadata results.
func Admit(reason Reason) Decision {
	return Decision{
		Outcome: OutcomeAdmitted,
		Reason:  reason,
		Effect:  EffectNone,
	}
}

// Admitted returns the default admitted decision.
//
// It is a convenience for the common success case where ReasonAdmitted is
// precise enough and no typed Result is needed yet.
func Admitted() Decision {
	return Admit(ReasonAdmitted)
}

// Accepted returns an admitted no-side-effect result with metadata.
//
// The result carries no grant because the admission stage did not transfer
// ownership to the caller.
func Accepted[M any](
	reason Reason,
	metadata M,
) Result[NoGrant, M] {
	return resultWith(
		Admit(reason),
		none[NoGrant](),
		some(metadata),
	)
}

// AcceptedNoMetadata returns an admitted no-side-effect result without metadata.
func AcceptedNoMetadata(reason Reason) Result[NoGrant, NoMetadata] {
	return resultWith(
		Admit(reason),
		none[NoGrant](),
		none[NoMetadata](),
	)
}
