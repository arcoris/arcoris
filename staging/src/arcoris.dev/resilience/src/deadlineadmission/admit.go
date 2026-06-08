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
	"arcoris.dev/resilience/deadline"
)

// TryAdmit returns an admission-compatible deadline start decision.
//
// TryAdmit delegates to deadline.CanStart. Nil contexts and negative minimum
// budgets preserve core deadline panic behavior instead of being converted into
// denied admission results. The returned result has no side effect and carries no
// grant.
func TryAdmit(req Request) admission.Result[
	admission.NoGrant,
	deadline.Decision,
] {
	return AdmissionResult(deadline.CanStart(
		req.Context,
		req.Now,
		req.Min,
	))
}
