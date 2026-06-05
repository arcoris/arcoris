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

// DenyDecision returns a denied decision with no side effect.
//
// Denied work is not admitted by the current component, and the component does
// not own waiting state or caller cleanup.
func DenyDecision(reason Reason) Decision {
	return Decision{
		Outcome: OutcomeDenied,
		Reason:  reason,
		Effect:  EffectNone,
	}
}

// DeniedResult returns a denied result with metadata.
func DeniedResult[M any](reason Reason, metadata M) Result[NoGrant, M] {
	var grant NoGrant
	return resultWith(DenyDecision(reason), grant, false, metadata, true)
}

// DeniedForResult returns a denied result typed for callers whose success path
// carries a grant.
//
// The returned result does not retain or manufacture a grant. Its grant type is
// present only to keep generic success and denial paths type-compatible.
func DeniedForResult[G any, M any](reason Reason, metadata M) Result[G, M] {
	var grant G
	return resultWith(DenyDecision(reason), grant, false, metadata, true)
}

// DeniedNoMetadataResult returns a denied result without metadata.
func DeniedNoMetadataResult(reason Reason) Result[NoGrant, NoMetadata] {
	var grant NoGrant
	var metadata NoMetadata
	return resultWith(DenyDecision(reason), grant, false, metadata, false)
}
