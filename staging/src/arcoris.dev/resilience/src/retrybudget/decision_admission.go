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

package retrybudget

import (
	"arcoris.dev/admission"
	admissionbuiltin "arcoris.dev/admissioncatalog/builtin"
	"arcoris.dev/snapshot"
)

// invalidAdmissionReason preserves invalid retry-budget decisions during
// conversion to admission.Result.
//
// Produced Decision values should be valid. When callers manually construct an
// invalid Decision, AdmissionResult must not hide that invalidity by returning a
// valid generic admission result.
const invalidAdmissionReason admission.Reason = ""

// AdmissionResult converts d into the generic admission result contract.
//
// Allowed retry-budget decisions are committed spend-only effects: the retry has
// already been accounted for and no caller-owned grant or release path exists.
// Denied decisions are no-effect budget back-pressure. In both cases the
// retry-budget snapshot is carried as admission metadata.
//
// Invalid retry-budget decisions convert to invalid admission results instead of
// being normalized. This keeps hand-built invalid Decision values visible to
// tests and defensive callers.
func (d Decision) AdmissionResult() admission.Result[
	admission.NoGrant,
	snapshot.Snapshot[Snapshot],
] {
	if !d.IsValid() {
		if d.Allowed {
			return admission.CommittedResult(
				invalidAdmissionReason,
				d.Snapshot,
			)
		}
		return admission.DeniedResult(
			invalidAdmissionReason,
			d.Snapshot,
		)
	}

	if d.IsAllowed() {
		return admission.CommittedResult(
			admission.ReasonAdmitted,
			d.Snapshot,
		)
	}

	return admission.DeniedResult(
		admissionbuiltin.ReasonBudgetExhausted,
		d.Snapshot,
	)
}
