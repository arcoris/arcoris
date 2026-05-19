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

// Effect classifies the side-effect semantics of one admission attempt.
//
// Effect is closed because Result validation must know whether a grant is
// required, forbidden, or optional.
type Effect uint8

const (
	// EffectUnknown is the zero value and is invalid.
	EffectUnknown Effect = iota

	// EffectNone means admission did not mutate resource or ownership state.
	EffectNone

	// EffectCommitted means admission committed a spend-only side effect, such as
	// consuming a retry budget or rate-limit token.
	EffectCommitted

	// EffectOwned means admission returned a caller-owned grant that must be
	// released, closed, committed, or rolled back according to the domain API.
	EffectOwned

	// EffectQueued means admission accepted ownership of waiting work.
	EffectQueued
)

// IsValid reports whether e is a defined non-zero effect.
func (e Effect) IsValid() bool {
	switch e {
	case EffectNone, EffectCommitted, EffectOwned, EffectQueued:
		return true
	default:
		return false
	}
}

// HasSideEffect reports whether e records a committed, owned, or queued effect.
func (e Effect) HasSideEffect() bool {
	return e == EffectCommitted || e == EffectOwned || e == EffectQueued
}

// RequiresGrant reports whether e requires a Result grant value.
//
// Only owned effects require a typed grant because the caller receives lifecycle
// responsibility that must later be released, committed, or rolled back by the
// domain package.
func (e Effect) RequiresGrant() bool {
	return e == EffectOwned
}

// AllowsGrant reports whether e may carry a Result grant value.
//
// Owned effects require a grant. Queued effects may optionally carry a queue
// handle, ticket, or cancellation token. Other effects must not carry caller
// ownership.
func (e Effect) AllowsGrant() bool {
	return e == EffectOwned || e == EffectQueued
}

// String returns the stable machine-readable effect name.
//
// Undefined values format as "unknown" so diagnostics remain safe even when a
// caller constructs an invalid Effect manually.
func (e Effect) String() string {
	switch e {
	case EffectNone:
		return "none"
	case EffectCommitted:
		return "committed"
	case EffectOwned:
		return "owned"
	case EffectQueued:
		return "queued"
	default:
		return "unknown"
	}
}
