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

// Queue returns a queued decision with system-owned waiting semantics.
//
// Queue means the component accepted ownership of waiting work. The returned
// Result may include a queue handle, ticket, or other domain-specific reference
// if the component exposes one. The returned Decision is the semantic base for
// Queued and QueuedNoGrant results.
func Queue(reason Reason) Decision {
	return Decision{
		Outcome: OutcomeQueued,
		Reason:  reason,
		Effect:  EffectQueued,
	}
}

// Queued returns a queued result with a queue handle and metadata.
//
// Queued means the system accepted waiting ownership. The grant value should be
// a domain-owned handle, ticket, or cancellation/reconciliation value when the
// queue exposes one.
func Queued[G any, M any](
	reason Reason,
	handle G,
	metadata M,
) Result[G, M] {
	return resultWith(
		Queue(reason),
		some(handle),
		some(metadata),
	)
}

// QueuedNoGrant returns a queued result without a queue handle.
//
// Some queues accept waiting work without exposing a caller-owned handle. The
// effect still records queued ownership, but the Result carries no grant.
func QueuedNoGrant[M any](
	reason Reason,
	metadata M,
) Result[NoGrant, M] {
	return resultWith(
		Queue(reason),
		none[NoGrant](),
		some(metadata),
	)
}
