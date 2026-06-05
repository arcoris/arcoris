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

// QueueDecision returns a queued decision with system-owned waiting semantics.
//
// Queued means the component accepted ownership of waiting work. A queued result
// may optionally include a queue handle, ticket, cancellation token, or other
// domain-specific grant.
func QueueDecision(reason Reason) Decision {
	return Decision{
		Outcome: OutcomeQueued,
		Reason:  reason,
		Effect:  EffectQueued,
	}
}

// QueuedResult returns a queued result with a queue grant and metadata.
func QueuedResult[G any, M any](reason Reason, grant G, metadata M) Result[G, M] {
	return resultWith(QueueDecision(reason), grant, true, metadata, true)
}

// QueuedNoGrantResult returns a queued result with metadata and no queue grant.
func QueuedNoGrantResult[M any](reason Reason, metadata M) Result[NoGrant, M] {
	var grant NoGrant
	return resultWith(QueueDecision(reason), grant, false, metadata, true)
}
