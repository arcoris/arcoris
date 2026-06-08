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

package retrybudgetadmission

import (
	"arcoris.dev/admission"
	admissionbuiltin "arcoris.dev/admissioncatalog/builtin"
	"arcoris.dev/resilience/retrybudget"
)

// invalidAdmissionReason preserves invalid retry-budget decisions during
// conversion to admission.Result.
//
// Produced retrybudget.Decision values should be valid. If callers manually
// construct an invalid Decision, AdmissionResult must not hide that invalidity
// by returning a normal generic admission result.
const invalidAdmissionReason admission.Reason = ""

// AdmissionResult converts d into the generic admission result contract.
//
// Retry-budget admission is committed and grant-free. Allowed decisions have
// already spent one retry attempt. Denied decisions have not spent a retry. The
// complete retrybudget.Decision is carried as metadata so callers can inspect
// precise budget state without reverse-engineering generic admission reasons.
func AdmissionResult(d retrybudget.Decision) admission.Result[
	admission.NoGrant,
	retrybudget.Decision,
] {
	if !d.IsValid() {
		if d.IsAllowed() {
			return admission.CommittedResult(
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
		return admission.CommittedResult(
			admission.ReasonAdmitted,
			d,
		)
	}

	return admission.DeniedResult(
		admissionbuiltin.ReasonBudgetExhausted,
		d,
	)
}
