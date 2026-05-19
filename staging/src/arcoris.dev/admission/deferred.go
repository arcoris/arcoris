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

// Defer returns a deferred decision with no side effect.
//
// Deferred work is not accepted now, but the caller retains responsibility for
// deciding whether and when to try again. The returned Decision is the semantic
// base for Deferred, DeferredFor, and DeferredNoMetadata results.
func Defer(reason Reason) Decision {
	return Decision{
		Outcome: OutcomeDeferred,
		Reason:  reason,
		Effect:  EffectNone,
	}
}

// Deferred returns a deferred result with metadata.
//
// Deferred means the component did not accept the work now and did not create
// system-owned waiting. The caller keeps responsibility for any later retry.
func Deferred[M any](
	reason Reason,
	metadata M,
) Result[NoGrant, M] {
	return resultWith(
		Defer(reason),
		none[NoGrant](),
		some(metadata),
	)
}

// DeferredFor returns a deferred result typed for callers whose success path has
// a grant.
func DeferredFor[G any, M any](
	reason Reason,
	metadata M,
) Result[G, M] {
	return resultWith(
		Defer(reason),
		none[G](),
		some(metadata),
	)
}

// DeferredNoMetadata returns a deferred result without metadata.
func DeferredNoMetadata(reason Reason) Result[NoGrant, NoMetadata] {
	return resultWith(
		Defer(reason),
		none[NoGrant](),
		none[NoMetadata](),
	)
}
