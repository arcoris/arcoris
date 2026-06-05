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

// AdmitDecision returns an admitted decision with no side effect.
//
// Use AdmitDecision when a component only reports that work may proceed and does
// not transfer a grant or commit a budget/token spend.
func AdmitDecision(reason Reason) Decision {
	return Decision{
		Outcome: OutcomeAdmitted,
		Reason:  reason,
		Effect:  EffectNone,
	}
}

// AdmittedDecision returns the generic admitted decision with no side effect.
func AdmittedDecision() Decision {
	return AdmitDecision(ReasonAdmitted)
}

// AcceptedResult returns an admitted no-side-effect result with metadata.
func AcceptedResult[M any](reason Reason, metadata M) Result[NoGrant, M] {
	var grant NoGrant
	return resultWith(AdmitDecision(reason), grant, false, metadata, true)
}

// AcceptedNoMetadataResult returns an admitted no-side-effect result without
// metadata.
func AcceptedNoMetadataResult(reason Reason) Result[NoGrant, NoMetadata] {
	var grant NoGrant
	var metadata NoMetadata
	return resultWith(AdmitDecision(reason), grant, false, metadata, false)
}
