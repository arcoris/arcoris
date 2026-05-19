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

// IsAdmitted reports whether d allows work to proceed immediately.
func (d Decision) IsAdmitted() bool {
	return d.Outcome.IsAdmitted()
}

// IsDenied reports whether d rejects the current admission attempt.
func (d Decision) IsDenied() bool {
	return d.Outcome.IsDenied()
}

// IsQueued reports whether d accepted system-owned waiting work.
func (d Decision) IsQueued() bool {
	return d.Outcome.IsQueued()
}

// IsDeferred reports whether d leaves retry ownership with the caller.
func (d Decision) IsDeferred() bool {
	return d.Outcome.IsDeferred()
}

// IsTerminal reports whether d completes the current admission attempt without
// leaving system-owned waiting behind.
func (d Decision) IsTerminal() bool {
	return d.Outcome.IsTerminal()
}

// HasSideEffect reports whether d records committed, owned, or queued state.
func (d Decision) HasSideEffect() bool {
	return d.Effect.HasSideEffect()
}

// RequiresGrant reports whether a typed Result carrying d must include a grant.
func (d Decision) RequiresGrant() bool {
	return d.Effect.RequiresGrant()
}

// AllowsGrant reports whether a typed Result carrying d may include a grant.
func (d Decision) AllowsGrant() bool {
	return d.Effect.AllowsGrant()
}
