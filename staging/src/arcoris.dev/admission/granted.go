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

// Grant returns an admitted decision that requires a caller-owned grant.
//
// Grant fits lease-style admission such as bulkheads where successful admission
// transfers ownership that the caller must later release, close, commit, or roll
// back according to the domain API. The returned Decision is the semantic base
// for Granted and GrantedNoMetadata results.
func Grant(reason Reason) Decision {
	return Decision{
		Outcome: OutcomeAdmitted,
		Reason:  reason,
		Effect:  EffectOwned,
	}
}

// Granted returns an admitted result with a caller-owned grant and metadata.
//
// The grant is required because the decision uses EffectOwned. The caller is
// responsible for following the grant's domain lifecycle, such as releasing a
// bulkhead lease or rolling back a staged allocation.
func Granted[G any, M any](
	reason Reason,
	grant G,
	metadata M,
) Result[G, M] {
	return resultWith(
		Grant(reason),
		some(grant),
		some(metadata),
	)
}

// GrantedNoMetadata returns an admitted result with a caller-owned grant and no
// metadata.
func GrantedNoMetadata[G any](
	reason Reason,
	grant G,
) Result[G, NoMetadata] {
	return resultWith(
		Grant(reason),
		some(grant),
		none[NoMetadata](),
	)
}
