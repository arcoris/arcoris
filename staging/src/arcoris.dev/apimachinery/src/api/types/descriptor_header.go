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

package types

// descriptorHeader is the shared working state carried by public builders.
//
// Builders keep only the DescriptorKind, common flags, and their exact local payload.
// They do not keep a full Descriptor while building. Descriptor remains the normalized
// descriptor/IR assembled at Descriptor() boundaries.
type descriptorHeader struct {
	// code selects the descriptor kind being built.
	code DescriptorKind
	// flags records descriptor-wide builder flags such as nullability.
	flags descriptorFlags
}

// newHeader creates builder header state for code.
func newHeader(code DescriptorKind) descriptorHeader {
	return descriptorHeader{code: code}
}

// withNullable returns a copy of h that admits null values.
func (h descriptorHeader) withNullable() descriptorHeader {
	h.flags |= descriptorFlagNullable

	return h
}

// descriptorFromHeader creates the normalized descriptor shell for a builder.
func descriptorFromHeader(h descriptorHeader) Descriptor {
	return Descriptor{code: h.code, flags: h.flags}
}
