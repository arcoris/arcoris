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

package admissioncatalog

import "arcoris.dev/admission"

// newReasonStore returns a descriptor store keyed and sorted by reason.
func newReasonStore() descriptorStore[admission.Reason, ReasonDescriptor] {
	return newDescriptorStore(reasonDescriptorKey, reasonDescriptorLess)
}

// newKindStore returns a descriptor store keyed and sorted by component kind.
func newKindStore() descriptorStore[admission.ComponentKind, ComponentKindDescriptor] {
	return newDescriptorStore(kindDescriptorKey, kindDescriptorLess)
}

// newComponentStore returns a descriptor store keyed and sorted by component ID.
func newComponentStore() descriptorStore[admission.ComponentID, ComponentDescriptor] {
	return newDescriptorStore(componentDescriptorKey, componentDescriptorLess)
}

// initReasonStore makes a zero descriptor store usable for reason descriptors.
func initReasonStore(s *descriptorStore[admission.Reason, ReasonDescriptor]) {
	s.init(reasonDescriptorKey, reasonDescriptorLess)
}

// initKindStore makes a zero descriptor store usable for kind descriptors.
func initKindStore(s *descriptorStore[admission.ComponentKind, ComponentKindDescriptor]) {
	s.init(kindDescriptorKey, kindDescriptorLess)
}

// initComponentStore makes a zero descriptor store usable for component descriptors.
func initComponentStore(s *descriptorStore[admission.ComponentID, ComponentDescriptor]) {
	s.init(componentDescriptorKey, componentDescriptorLess)
}
