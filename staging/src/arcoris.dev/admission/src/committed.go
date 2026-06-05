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

// CommitDecision returns an admitted decision with a committed side effect.
//
// Committed decisions fit budgets, tokens, and counters that are consumed by the
// successful admission attempt and are not released by the caller later.
func CommitDecision(reason Reason) Decision {
	return Decision{
		Outcome: OutcomeAdmitted,
		Reason:  reason,
		Effect:  EffectCommitted,
	}
}

// CommittedResult returns an admitted committed-side-effect result with
// metadata.
func CommittedResult[M any](reason Reason, metadata M) Result[NoGrant, M] {
	var grant NoGrant
	return resultWith(CommitDecision(reason), grant, false, metadata, true)
}

// CommittedNoMetadataResult returns an admitted committed-side-effect result
// without metadata.
func CommittedNoMetadataResult(reason Reason) Result[NoGrant, NoMetadata] {
	var grant NoGrant
	var metadata NoMetadata
	return resultWith(CommitDecision(reason), grant, false, metadata, false)
}
