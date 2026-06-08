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

package deadlineadmission

import (
	"arcoris.dev/admission"
	admissionbuiltin "arcoris.dev/admissioncatalog/builtin"
	"arcoris.dev/resilience/deadline"
)

// invalidAdmissionReason preserves invalid deadline decisions during conversion
// to admission.Result.
//
// Produced deadline.Decision values should be valid. If callers manually
// construct an invalid Decision, AdmissionResult must not hide that invalidity by
// returning a valid generic admission result.
const invalidAdmissionReason admission.Reason = ""

// AdmissionResult converts d into the generic admission result contract.
//
// Deadline admission is pure start-decision metadata. Allowed decisions become
// admitted no-side-effect results. Denied decisions remain no-side-effect
// denials. No grant is returned, no release or rollback path exists, and the
// deadline.Decision itself is carried as metadata.
func AdmissionResult(d deadline.Decision) admission.Result[
	admission.NoGrant,
	deadline.Decision,
] {
	if !d.IsValid() {
		if d.IsAllowed() {
			return admission.AcceptedResult(
				invalidAdmissionReason,
				d,
			)
		}
		return admission.DeniedResult(
			invalidAdmissionReason,
			d,
		)
	}

	if d.IsAllowed() {
		return admission.AcceptedResult(
			admission.ReasonAdmitted,
			d,
		)
	}

	switch d.Reason {
	case deadline.ReasonContextDone:
		return admission.DeniedResult(
			admissionbuiltin.ReasonCanceled,
			d,
		)
	case deadline.ReasonExpired, deadline.ReasonInsufficientBudget:
		return admission.DeniedResult(
			admissionbuiltin.ReasonDeadlineExceeded,
			d,
		)
	default:
		return admission.DeniedResult(
			invalidAdmissionReason,
			d,
		)
	}
}
