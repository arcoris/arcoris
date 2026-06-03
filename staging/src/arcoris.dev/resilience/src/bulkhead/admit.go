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

package bulkhead

import (
	"arcoris.dev/admission"
	admissionbuiltin "arcoris.dev/admissioncatalog/builtin"
	"arcoris.dev/snapshot"
)

// TryAdmit attempts to admit a bulkhead request without waiting.
//
// TryAdmit is the admission-compatible surface for the same capacity operation
// used by TryAcquireAmount. A successful result is admitted with EffectOwned and
// carries a Lease grant that the caller must release. Capacity exhaustion is a
// valid denied result with the observed capacity snapshot as metadata.
//
// The method deliberately performs no catalog lookup, policy evaluation, retry,
// queueing, context observation, or error classification. Bulkhead remains a
// local bounded in-flight primitive; admission only supplies the generic Result
// shape.
func (b *Bulkhead) TryAdmit(req Request) admission.Result[*Lease, snapshot.Snapshot[Snapshot]] {
	lease, snap, ok := b.TryAcquireAmount(req.Amount)
	if !ok {
		return admission.DeniedFor[*Lease](
			admissionbuiltin.ReasonCapacityExhausted,
			snap,
		)
	}

	return admission.Granted(
		admission.ReasonAdmitted,
		lease,
		snap,
	)
}
