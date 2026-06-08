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
	"arcoris.dev/resilience/retrybudget"
)

// TryAdmit attempts to admit one retry through the wrapped retry budget.
//
// A successful result is admitted with EffectCommitted and carries no grant,
// because the retry attempt has already been spent by TryAdmitRetry. A denied
// result carries no grant and preserves the full retrybudget.Decision as
// metadata. Nil budgets panic with ErrNilRetryAdmitter.
func (a Admitter) TryAdmit(Request) admission.Result[
	admission.NoGrant,
	retrybudget.Decision,
] {
	a.requireBudget()
	return AdmissionResult(a.budget.TryAdmitRetry())
}
