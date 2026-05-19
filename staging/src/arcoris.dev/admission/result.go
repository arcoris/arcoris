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

	// grant is present only when the decision effect allows or requires a typed
	// domain grant. It stays private so Result invariants are preserved through
	// constructors and validation helpers.
	grant Maybe[G]

	// metadata is optional typed read-model data associated with the decision.
	// It is private for the same reason as grant: absence and presence are part
	// of the Result shape, not raw struct mutation.
	metadata Maybe[M]
}

// IsAdmitted reports whether r allows work to proceed immediately.
func (r Result[G, M]) IsAdmitted() bool {
	return r.decision.IsAdmitted()
}

// IsDenied reports whether r rejects the current admission attempt.
func (r Result[G, M]) IsDenied() bool {
	return r.decision.IsDenied()
}

// IsQueued reports whether r accepted system-owned waiting work.
func (r Result[G, M]) IsQueued() bool {
	return r.decision.IsQueued()
}

// IsDeferred reports whether r leaves retry ownership with the caller.
func (r Result[G, M]) IsDeferred() bool {
	return r.decision.IsDeferred()
}

// HasSideEffect reports whether r records committed, owned, or queued state.
func (r Result[G, M]) HasSideEffect() bool {
	return r.decision.HasSideEffect()
}

// HasGrant reports whether r contains a typed grant value.
func (r Result[G, M]) HasGrant() bool {
	return r.grant.IsSome()
}

// HasMetadata reports whether r contains typed metadata.
func (r Result[G, M]) HasMetadata() bool {
	return r.metadata.IsSome()
}
