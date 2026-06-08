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

package bulkheadadmission

import (
	"arcoris.dev/admission"
	admissionbuiltin "arcoris.dev/admissioncatalog/builtin"
	"arcoris.dev/resilience/bulkhead"
)

// TryAdmit attempts to admit a bulkhead request without waiting.
//
// A successful result is admitted with EffectOwned and carries a live
// *bulkhead.Lease grant that the caller must release. A denied result carries no
// grant and preserves the exact bulkhead.Observation in metadata.
//
// Invalid request amounts and invalid bulkhead receivers are programmer or
// configuration errors. They panic through the core bulkhead path instead of
// being converted into denied admission results.
func (a Admitter) TryAdmit(req Request) admission.Result[*bulkhead.Lease, bulkhead.Observation] {
	lease, observation, ok := a.bulkhead.TryAcquireAmount(req.Amount)
	if !ok {
		return admission.DeniedForResult[*bulkhead.Lease](
			admissionbuiltin.ReasonCapacityExhausted,
			observation,
		)
	}

	return admission.GrantedResult(
		admission.ReasonAdmitted,
		lease,
		observation,
	)
}
