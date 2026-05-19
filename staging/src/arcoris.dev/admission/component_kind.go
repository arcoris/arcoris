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

// maxComponentKindLength keeps kind values short enough for stable diagnostic
// surfaces while leaving room for open-world domain-specific names.
const maxComponentKindLength = 64

// ComponentKind is an open-world coarse class for admission components.
//
// Kind groups a component by role, while ComponentID identifies the component
// more precisely. The value is stable metadata and must not contain runtime
// instance data.
type ComponentKind string

const (
	// KindBulkhead describes bounded in-flight isolation components.
	KindBulkhead ComponentKind = "bulkhead"

	// KindRetryBudget describes components that limit retry amplification.
	KindRetryBudget ComponentKind = "retry_budget"

	// KindDeadline describes components that derive or enforce execution
	// budgets from deadlines.
	KindDeadline ComponentKind = "deadline"

	// KindRateLimiter describes components that gate work by rate or token
	// availability.
	KindRateLimiter ComponentKind = "rate_limiter"

	// KindQueue describes components that can take ownership of waiting work.
	KindQueue ComponentKind = "queue"

	// KindScheduler describes components that choose when or where work runs.
	KindScheduler ComponentKind = "scheduler"

	// KindWorkerPool describes components that own worker execution capacity.
	KindWorkerPool ComponentKind = "worker_pool"

	// KindOverloadGate describes components that deny or defer work under
	// overload.
	KindOverloadGate ComponentKind = "overload_gate"

	// KindTenantIsolation describes components that isolate tenants or other
	// workload classes from each other.
	KindTenantIsolation ComponentKind = "tenant_isolation"
)

// IsValid reports whether k is a valid lower_snake_case component kind.
//
// Custom kinds are allowed, but they must remain short, stable, ASCII
// identifiers suitable for logs, docs, and metrics dimensions owned elsewhere.
func (k ComponentKind) IsValid() bool {
	return validLowerSnakeIdentifier(string(k), maxComponentKindLength)
}

// String returns k as a string.
func (k ComponentKind) String() string {
	return string(k)
}
