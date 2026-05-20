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

package deadline

import "arcoris.dev/admission"

// TryAdmit returns an admission-compatible start decision.
//
// TryAdmit is a thin, stateless adapter over CanStart. It preserves CanStart's
// validation behavior: nil contexts and negative minimum budgets panic at the
// API boundary instead of being converted into denied admission results.
//
// The result has no side effect and carries no grant. Admission metadata is the
// original deadline Decision so callers can inspect the deadline-specific reason
// and remaining budget without a catalog lookup, timer, goroutine, or stateful
// checker object.
func TryAdmit(req Request) admission.Result[
	admission.NoGrant,
	Decision,
] {
	return CanStart(
		req.Context,
		req.Now,
		req.Min,
	).AdmissionResult()
}
