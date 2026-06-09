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

// validateMap checks dynamic-map key, value, and length rules.
func validateMap(desc Descriptor, resolver Resolver, path string, resolving map[TypeName]bool) error {
	if desc.mapType.key == nil {
		return descriptorErrorf(
			path+".key",
			ErrInvalidDescriptor,
			DescriptorErrorReasonInvalidMapKey,
			"map descriptor must have a key descriptor",
		)
	}

	if err := validateDescriptor(*desc.mapType.key, resolver, path+".key", resolving); err != nil {
		return err
	}

	if err := validateMapKeyDescriptor(*desc.mapType.key, resolver, path+".key", resolving); err != nil {
		return err
	}

	if desc.mapType.value == nil {
		return descriptorErrorf(
			path+".value",
			ErrInvalidDescriptor,
			DescriptorErrorReasonMissingValue,
			"map descriptor must have a value descriptor",
		)
	}

	if err := validateDescriptor(*desc.mapType.value, resolver, path+".value", resolving); err != nil {
		return err
	}

	return validateLengthLimits(desc.mapType.minLen, desc.mapType.maxLen, path+".len")
}

// validateMapKeyDescriptor verifies that key can describe concrete string map keys.
func validateMapKeyDescriptor(
	key Descriptor,
	resolver Resolver,
	path string,
	resolving map[TypeName]bool,
) error {
	if key.Nullable() {
		return descriptorErrorf(
			path,
			ErrInvalidField,
			DescriptorErrorReasonInvalidMapKey,
			"map key descriptor must be non-nullable",
		)
	}

	switch key.Code() {
	case DescriptorString:
		return nil
	case DescriptorRef:
		return validateMapKeyRef(key, resolver, path, resolving)
	default:
		return descriptorErrorf(
			path,
			ErrInvalidField,
			DescriptorErrorReasonInvalidMapKey,
			"map key descriptor must be string-like, got %s",
			key.Code(),
		)
	}
}

// validateMapKeyRef resolves a map key reference when resolved validation is active.
func validateMapKeyRef(
	key Descriptor,
	resolver Resolver,
	path string,
	resolving map[TypeName]bool,
) error {
	name := key.ref.name

	if resolver == nil {
		return nil
	}

	if resolving[name] {
		return descriptorErrorf(
			path,
			ErrInvalidDescriptorReference,
			DescriptorErrorReasonReferenceCycle,
			"map key reference %q creates a recursive Definition graph",
			name,
		)
	}

	def, ok := resolver.Resolve(name)
	if !ok {
		return descriptorErrorf(
			path,
			ErrUnresolvedDescriptorReference,
			DescriptorErrorReasonUnknownReference,
			"map key reference %q was not found in resolver",
			name,
		)
	}

	next := copyResolving(resolving)
	next[name] = true

	return validateMapKeyDescriptor(def.Descriptor(), resolver, path, next)
}
