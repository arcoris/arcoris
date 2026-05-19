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

// Deny returns a denied decision with no side effect.
//
// Denied work is not accepted by the current component, and the component does
// not own waiting state or caller cleanup. The returned Decision is the semantic
// base for Denied, DeniedFor, and DeniedNoMetadata results.
func Deny(reason Reason) Decision {
	return Decision{
		Outcome: OutcomeDenied,
		Reason:  reason,
		Effect:  EffectNone,
	}
}

// Denied returns a denied result with metadata.
//
// Denied results never carry grants because the component did not accept the
// work and transferred no caller-owned lifecycle state.
func Denied[M any](
	reason Reason,
	metadata M,
) Result[NoGrant, M] {
	return resultWith(
		Deny(reason),
		none[NoGrant](),
		some(metadata),
	)
}

// DeniedFor returns a denied result typed for callers whose success path has a
// grant.
//
// This keeps generic call sites ergonomic: a function that usually returns
// Result[Lease, Snapshot] can deny without manufacturing a zero Lease value.
func DeniedFor[G any, M any](
	reason Reason,
	metadata M,
) Result[G, M] {
	return resultWith(
		Deny(reason),
		none[G](),
		some(metadata),
	)
}

// DeniedNoMetadata returns a denied result without metadata.
func DeniedNoMetadata(reason Reason) Result[NoGrant, NoMetadata] {
	return resultWith(
		Deny(reason),
		none[NoGrant](),
		none[NoMetadata](),
	)
}
