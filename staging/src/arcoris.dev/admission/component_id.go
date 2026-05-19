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

// maxComponentIDLength bounds stable component paths so they remain suitable
// for logs, documentation, and higher-level observability dimensions.
const maxComponentIDLength = 128

// ComponentID is a stable open-world identifier for an admission component.
//
// ComponentID identifies a component kind or role, not a runtime instance. It
// uses a dot-separated lower_snake_case path such as resilience.bulkhead or
// scheduler.tenant_fairness.
type ComponentID string

// IsValid reports whether id is a valid dot-separated component identifier.
//
// IDs are stable component paths. They are not request IDs, tenant IDs, runtime
// instance names, metric labels, or transport addresses.
func (id ComponentID) IsValid() bool {
	return validDotPathIdentifier(string(id), maxComponentIDLength)
}

// String returns id as a string.
//
// The method intentionally performs no validation. It is safe for diagnostics
// and deterministic ordering even when the caller is holding an invalid value.
func (id ComponentID) String() string {
	return string(id)
}

// MustComponentID returns value as a ComponentID or panics when value is invalid.
//
// The helper is intended for package-level constants and tests where an invalid
// literal is a programming error.
func MustComponentID(value string) ComponentID {
	id := ComponentID(value)
	if !id.IsValid() {
		panic("admission.ComponentID: invalid component id")
	}
	return id
}
