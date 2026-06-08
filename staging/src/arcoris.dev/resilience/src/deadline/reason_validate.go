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

package deadline

// IsValid reports whether r is one of the defined deadline reasons.
func (r Reason) IsValid() bool {
	switch r {
	case ReasonAllowed,
		ReasonContextDone,
		ReasonNoDeadline,
		ReasonExpired,
		ReasonInsufficientBudget:
		return true
	default:
		return false
	}
}

// IsAllowedReason reports whether r can explain an allowed decision.
func (r Reason) IsAllowedReason() bool {
	return r == ReasonAllowed || r == ReasonNoDeadline
}

// IsDeniedReason reports whether r can explain a denied decision.
func (r Reason) IsDeniedReason() bool {
	switch r {
	case ReasonContextDone, ReasonExpired, ReasonInsufficientBudget:
		return true
	default:
		return false
	}
}
