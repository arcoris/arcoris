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

// descriptorFlags stores descriptor-wide boolean modifiers.
//
// Flags are deliberately private so new descriptor-wide semantics can be added
// without allowing callers to manufacture unsupported combinations. Validation
// remains the authority on whether a flag is legal for a DescriptorKind.
type descriptorFlags uint8

const (
	// descriptorFlagNullable records value nullability for every non-null descriptor.
	//
	// This flag does not describe field presence and does not turn DescriptorNull into
	// a nullable type. DescriptorNull is already the null literal.
	descriptorFlagNullable descriptorFlags = 1 << iota
)

// withNullable returns a copy of desc that admits null values.
//
// Builders use this helper instead of mutating shared state. The helper does
// not reject DescriptorNull so it can remain tiny and allocation-free; ValidateResolved
// enforces the exact DescriptorNull invariant that DescriptorNull cannot be nullable.
func (desc Descriptor) withNullable() Descriptor {
	desc.flags |= descriptorFlagNullable

	return desc
}
