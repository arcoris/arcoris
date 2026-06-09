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

// validateListMapKeyIdentityDescriptor checks the descriptor-level identity contract
// for one ListMap key field.
//
// ListMap keys are stronger than ordinary required object fields. Future value
// validation, diff, and apply layers must be able to extract them into stable
// selector identity values. This first pass limits that identity surface to
// non-nullable bool, string, and integer descriptors, including references that
// resolve to those descriptors.
func validateListMapKeyIdentityDescriptor(
	field FieldDescriptor,
	resolver Resolver,
	path string,
	resolving map[TypeName]bool,
) error {
	return validateListMapKeyDescriptor(field.Descriptor(), resolver, path, resolving)
}

// validateListMapKeyDescriptor validates a concrete or referenced key descriptor.
func validateListMapKeyDescriptor(
	descriptor Descriptor,
	resolver Resolver,
	path string,
	resolving map[TypeName]bool,
) error {
	if descriptor.Nullable() {
		return descriptorErrorf(
			path,
			ErrInvalidField,
			DescriptorErrorReasonInvalidListMapKeyDescriptor,
			"ListMap key field must be non-nullable",
		)
	}

	switch descriptor.code {
	case DescriptorBool,
		DescriptorString,
		DescriptorInt8,
		DescriptorInt16,
		DescriptorInt32,
		DescriptorInt64,
		DescriptorUint8,
		DescriptorUint16,
		DescriptorUint32,
		DescriptorUint64:
		return nil
	case DescriptorRef:
		return validateListMapKeyRef(descriptor.ref.name, resolver, path, resolving)
	default:
		return descriptorErrorf(
			path,
			ErrInvalidField,
			DescriptorErrorReasonInvalidListMapKeyDescriptor,
			"ListMap key field descriptor %s cannot be represented as a stable identity value",
			descriptor.code,
		)
	}
}

// validateListMapKeyRef resolves a referenced key descriptor without introducing a
// catalog dependency or a public identity-descriptor API.
func validateListMapKeyRef(
	name TypeName,
	resolver Resolver,
	path string,
	resolving map[TypeName]bool,
) error {
	if !name.IsValid() {
		return descriptorErrorf(
			path,
			ErrInvalidDescriptorReference,
			DescriptorErrorReasonInvalidReferenceName,
			"reference name %q is not a valid TypeName",
			name,
		)
	}

	if resolver == nil {
		return nil
	}

	if resolving[name] {
		return descriptorErrorf(
			path,
			ErrInvalidDescriptorReference,
			DescriptorErrorReasonReferenceCycle,
			"reference %q creates a recursive Definition graph",
			name,
		)
	}

	def, ok := resolver.Resolve(name)

	if !ok {
		return descriptorErrorf(
			path,
			ErrUnresolvedDescriptorReference,
			DescriptorErrorReasonUnknownReference,
			"reference %q was not found in resolver",
			name,
		)
	}

	next := copyResolving(resolving)
	next[name] = true

	return validateListMapKeyDescriptor(def.Descriptor(), resolver, path, next)
}
