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

// ComponentDescriptor describes an open-world admission component.
//
// Descriptors are stable component metadata. They are not global registrations
// and they do not identify runtime instances.
type ComponentDescriptor struct {
	// ID is the stable component identifier or ownership path.
	ID ComponentID

	// Kind is the coarse stable role of the component.
	Kind ComponentKind

	// Capabilities records the outcome and effect classes the component may
	// produce. A zero value is valid and means unspecified.
	Capabilities CapabilitySet
}

// IsValid reports whether d has valid component identity, kind, and capability
// metadata. A zero CapabilitySet is valid and means unspecified.
func (d ComponentDescriptor) IsValid() bool {
	return d.ID.IsValid() && d.Kind.IsValid() && d.Capabilities.IsValid()
}
