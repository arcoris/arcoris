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

import "arcoris.dev/resilience/retrybudget"

// Admitter adapts a retry-budget RetryAdmitter to admission.Result.
//
// Admitter owns no accounting state. The wrapped RetryAdmitter remains the only
// component that decides and records retry attempts. Admitter is safe to copy by
// value when the wrapped budget is safe to share.
type Admitter struct {
	// budget is the direct retry-budget spend primitive being adapted.
	//
	// A nil budget is programmer misuse. TryAdmit panics with
	// ErrNilRetryAdmitter instead of manufacturing a denied result, because nil
	// dependencies are configuration errors rather than budget pressure.
	budget retrybudget.RetryAdmitter
}

// New returns an admission adapter for budget.
//
// New does not call the budget. Nil budgets are detected when TryAdmit is called
// so adapter construction stays cheap and side-effect free.
func New(budget retrybudget.RetryAdmitter) Admitter {
	return Admitter{budget: budget}
}
