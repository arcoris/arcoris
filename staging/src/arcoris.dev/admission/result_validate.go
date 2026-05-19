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

// IsValid reports whether r satisfies admission result invariants.
//
// Result validity starts with Decision validity, then checks the typed grant
// shape implied by the decision effect. Owned results must contain a grant.
// Denied and deferred results must not contain a grant. Queued results may carry
// a domain-specific queue handle, and no-side-effect or committed results must
// not transfer caller-owned grants.
func (r Result[G, M]) IsValid() bool {
	if !r.decision.IsValid() {
		return false
	}
	if r.decision.RequiresGrant() {
		return r.HasGrant()
	}
	if !r.decision.AllowsGrant() && r.HasGrant() {
		return false
	}
	if r.decision.IsDenied() || r.decision.IsDeferred() {
		return !r.HasGrant()
	}

	return true
}
