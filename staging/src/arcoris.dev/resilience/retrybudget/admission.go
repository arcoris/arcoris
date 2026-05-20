/*
  Copyright 2026 The ARCORIS Authors

  Licensed under the Apache License, Version 2.0 (the "License");
  you may not use this file except in compliance with the License.
  You may obtain a copy of the License at

      http://www.apache.org/licenses/LICENSE-2.0

  Unless required by applicable law or agreed to in writing, software
  distributed under the License is distributed on an "AS IS" BASIS,
  WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
  See the License for the specific language governing permissions and
  limitations under the License.
*/

package retrybudget

import (
	"arcoris.dev/admission"
	"arcoris.dev/snapshot"
)

// AdmissionAdmitter is the admission-compatible retry-budget surface.
//
// AdmissionAdmitter does not replace RetryAdmitter and is intentionally not
// embedded into Budget. The existing Budget contract remains the direct
// retry-budget API, while implementations may additionally expose TryAdmit for
// generic admission pipelines that understand admission.Result.
//
// Successful retry-budget admission is a committed spend-only effect: no grant
// is returned and there is no release path. Denied retry-budget admission is
// ordinary budget back-pressure with a snapshot describing the observed state.
type AdmissionAdmitter interface {
	TryAdmit(Request) admission.Result[
		admission.NoGrant,
		snapshot.Snapshot[Snapshot],
	]
}
