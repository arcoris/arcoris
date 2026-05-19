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

// Commit returns an admitted decision with a committed spend-only side effect.
//
// Commit fits budgets, tokens, and counters that are consumed by the successful
// admission attempt and are not released by the caller later. The returned
// Decision is the semantic base for Committed and CommittedNoMetadata results.
func Commit(reason Reason) Decision {
	return Decision{
		Outcome: OutcomeAdmitted,
		Reason:  reason,
		Effect:  EffectCommitted,
	}
}

// Committed returns an admitted spend-only result with metadata.
//
// The result carries no grant because committed effects are not released by the
// caller. Retry-budget spends and token-bucket spends are examples of committed
// admission effects.
func Committed[M any](
	reason Reason,
	metadata M,
) Result[NoGrant, M] {
	return resultWith(
		Commit(reason),
		none[NoGrant](),
		some(metadata),
	)
}

// CommittedNoMetadata returns an admitted spend-only result without metadata.
func CommittedNoMetadata(reason Reason) Result[NoGrant, NoMetadata] {
	return resultWith(
		Commit(reason),
		none[NoGrant](),
		none[NoMetadata](),
	)
}
