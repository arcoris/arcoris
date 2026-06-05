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

// DeferDecision returns a deferred decision with no side effect.
//
// Deferred work is not admitted now, but the caller retains responsibility for
// deciding whether and when to try again.
func DeferDecision(reason Reason) Decision {
	return Decision{
		Outcome: OutcomeDeferred,
		Reason:  reason,
		Effect:  EffectNone,
	}
}

// DeferredResult returns a deferred result with metadata.
func DeferredResult[M any](reason Reason, metadata M) Result[NoGrant, M] {
	var grant NoGrant
	return resultWith(DeferDecision(reason), grant, false, metadata, true)
}

// DeferredForResult returns a deferred result typed for callers whose success
// path carries a grant.
func DeferredForResult[G any, M any](reason Reason, metadata M) Result[G, M] {
	var grant G
	return resultWith(DeferDecision(reason), grant, false, metadata, true)
}

// DeferredNoMetadataResult returns a deferred result without metadata.
func DeferredNoMetadataResult(reason Reason) Result[NoGrant, NoMetadata] {
	var grant NoGrant
	var metadata NoMetadata
	return resultWith(DeferDecision(reason), grant, false, metadata, false)
}
