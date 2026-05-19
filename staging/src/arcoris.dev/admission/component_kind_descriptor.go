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

// ComponentKindDescriptor describes stable open-world metadata for a component
// kind.
//
// The descriptor describes a role, not a concrete runtime component instance.
// It is catalog metadata owned by an explicit KindRegistry; creating a
// descriptor never registers it globally and never changes process-wide state.
// Capabilities are descriptive metadata for documentation, config validation,
// and operator-facing catalogs. They do not enforce Result validity directly.
type ComponentKindDescriptor struct {
	// Kind is the stable open-world component role described by this catalog
	// entry.
	Kind ComponentKind

	// Capabilities records the outcome and effect surface commonly supported by
	// the kind. The zero set is valid and means unspecified capabilities.
	Capabilities CapabilitySet
}

// IsValid reports whether d is syntactically valid catalog metadata.
//
// Validation is local to the value: it checks the kind identifier and
// capability bits, but it does not check duplicate registration or membership in
// any registry. Those catalog-level checks belong to KindRegistry.
func (d ComponentKindDescriptor) IsValid() bool {
	return d.Kind.IsValid() && d.Capabilities.IsValid()
}
