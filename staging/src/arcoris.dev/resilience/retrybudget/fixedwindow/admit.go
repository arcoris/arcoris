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

package fixedwindow

import (
	"arcoris.dev/admission"
	"arcoris.dev/resilience/retrybudget"
	"arcoris.dev/snapshot"
)

// TryAdmit exposes Limiter through admission's generic result contract.
//
// TryAdmit delegates to the same atomic check-and-spend path as TryAdmitRetry.
// An admitted result means the retry attempt has already been recorded. A denied
// result means no retry was spent. No grant is returned because retry-budget
// admission is committed spend-only and has no release path.
//
// The request is intentionally empty. Limiter does not perform catalog lookup,
// queueing, waiting, context observation, policy orchestration, or retry
// execution in this adapter.
func (l *Limiter) TryAdmit(
	retrybudget.Request,
) admission.Result[
	admission.NoGrant,
	snapshot.Snapshot[retrybudget.Snapshot],
] {
	return l.TryAdmitRetry().AdmissionResult()
}
