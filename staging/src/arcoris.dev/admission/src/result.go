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

// Result is a typed admission result.
//
// G is the grant type returned when an admission stage transfers ownership to
// the caller. M is the metadata type used for snapshots, diagnostics, or other
// read models associated with the decision.
type Result[G any, M any] struct {
	// decision is the closed semantic core of the admission result.
	decision Decision

	// grant is the caller-owned value returned by grant-bearing result shapes.
	//
	// The value is meaningful only when hasGrant is true. Constructors must leave
	// grant as the zero value when no grant is present, so omitted pointer-like
	// grants do not retain references supplied to unrelated constructor paths.
	grant G

	// hasGrant records grant presence separately from grant's zero value.
	hasGrant bool

	// metadata is domain-owned read-model or diagnostic data associated with the
	// result.
	//
	// The value is meaningful only when hasMetadata is true. Constructors must
	// leave metadata as the zero value when no metadata is present.
	metadata M

	// hasMetadata records metadata presence separately from metadata's zero
	// value.
	hasMetadata bool
}

// HasGrant reports whether r contains a typed grant value.
//
// Presence alone does not prove the Result is valid. For example, a denied
// result with a grant is invalid even though HasGrant reports true.
func (r Result[G, M]) HasGrant() bool {
	return r.hasGrant
}

// HasMetadata reports whether r contains typed metadata.
//
// Metadata is always optional from the core admission perspective. Domain
// packages may impose stronger expectations on their own Result aliases.
func (r Result[G, M]) HasMetadata() bool {
	return r.hasMetadata
}
