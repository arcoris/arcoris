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

package admission

// maxReasonLength bounds reason codes so they remain stable identifiers instead
// of becoming places where callers embed dynamic request data.
const maxReasonLength = 128

// Reason is a stable machine-readable explanation for an admission outcome.
//
// Reason is open-world. Domain packages may define custom reasons without
// changing this package. Custom reasons must be stable lower_snake_case
// identifiers and must not contain dynamic data, secrets, raw errors, object
// names, tenant IDs, timestamps, addresses, stack traces, or request IDs.
type Reason string

const (
	// ReasonAdmitted is the generic success reason for admitted work.
	ReasonAdmitted Reason = "admitted"

	// ReasonDenied is the generic rejection reason when a more specific reason
	// is not available.
	ReasonDenied Reason = "denied"

	// ReasonQueued is the generic reason for work accepted into system-owned
	// waiting state.
	ReasonQueued Reason = "queued"

	// ReasonDeferred is the generic reason for work left with the caller for a
	// later retry or reconsideration.
	ReasonDeferred Reason = "deferred"
)

// IsValid reports whether r is a stable lower_snake_case admission reason.
//
// Reason intentionally validates syntax only. The vocabulary is open so domain
// packages can add stable reason codes without changing this package.
func (r Reason) IsValid() bool {
	return validLowerSnakeIdentifier(string(r), maxReasonLength)
}

// String returns r as a string.
//
// The method intentionally performs no validation. It is a stable formatting
// helper for diagnostics and tests, not a gate for using a reason in a Result.
func (r Reason) String() string {
	return string(r)
}
