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

import "fmt"

// validateList checks element descriptor, length limits, semantics, and map keys.
func validateList(desc Descriptor, resolver Resolver, path string, resolving map[TypeName]bool) error {
	if desc.list.elem == nil {
		return descriptorErrorf(
			path+".elem",
			ErrInvalidDescriptor,
			DescriptorErrorReasonMissingElement,
			"list descriptor must have an element descriptor",
		)
	}

	if err := validateDescriptor(*desc.list.elem, resolver, path+".elem", resolving); err != nil {
		return err
	}

	if err := validateLengthLimits(desc.list.minLen, desc.list.maxLen, path+".len"); err != nil {
		return err
	}

	if !desc.list.semantics.IsValid() {
		return descriptorErrorf(
			path+".semantics",
			ErrInvalidDescriptor,
			DescriptorErrorReasonInvalidSemantics,
			"list semantics %d is not supported",
			desc.list.semantics,
		)
	}

	if desc.list.semantics == ListSet {
		if err := validateListSetElementDescriptor(*desc.list.elem, resolver, path+".elem", resolving); err != nil {
			return err
		}
	}

	if desc.list.semantics != ListMap {
		if len(desc.list.mapKeys) > 0 {
			return descriptorErrorf(
				path+".mapKeys",
				ErrInvalidField,
				DescriptorErrorReasonInvalidListMapKey,
				"list map keys are only valid with ListMap semantics",
			)
		}

		return nil
	}

	if len(desc.list.mapKeys) == 0 {
		return descriptorErrorf(
			path+".mapKeys",
			ErrInvalidField,
			DescriptorErrorReasonMissingListMapKey,
			"ListMap semantics requires at least one key field",
		)
	}

	for i, key := range desc.list.mapKeys {
		keyPath := fmt.Sprintf("%s.mapKeys[%d]", path, i)

		if !key.IsValid() {
			return descriptorErrorf(
				keyPath,
				ErrInvalidField,
				DescriptorErrorReasonInvalidListMapKey,
				"ListMap key %q is not a valid field name",
				key,
			)
		}
	}

	if err := validateDuplicateListMapKeys(desc.list.mapKeys, path); err != nil {
		return err
	}

	if desc.list.elem.code == DescriptorRef && resolver == nil {
		return nil
	}

	object, ok := listMapObject(*desc.list.elem, resolver)

	if !ok {
		return descriptorErrorf(
			path+".elem",
			ErrInvalidDescriptor,
			DescriptorErrorReasonListMapElementNotObject,
			"ListMap element must be an object descriptor or a DescriptorRef resolving to an object",
		)
	}

	fields := make(map[FieldName]FieldDescriptor, len(object.fields))

	for _, field := range object.fields {
		fields[field.name] = field
	}

	for i, key := range desc.list.mapKeys {
		keyPath := fmt.Sprintf("%s.mapKeys[%d]", path, i)
		field, ok := fields[key]

		if !ok {
			return descriptorErrorf(
				keyPath,
				ErrInvalidField,
				DescriptorErrorReasonListMapKeyNotFound,
				"ListMap key %q is not present in the object element",
				key,
			)
		}

		if !field.IsRequired() {
			return descriptorErrorf(
				keyPath,
				ErrInvalidField,
				DescriptorErrorReasonListMapKeyOptional,
				"ListMap key field %q must be required",
				key,
			)
		}

		if err := validateListMapKeyIdentityDescriptor(field, resolver, keyPath, resolving); err != nil {
			return err
		}
	}

	return nil
}

// validateDuplicateListMapKeys rejects repeated declared ListMap keys.
func validateDuplicateListMapKeys(keys []FieldName, path string) error {
	firstIndexes := make(map[FieldName]int, len(keys))

	for i, key := range keys {
		if first, exists := firstIndexes[key]; exists {
			return descriptorErrorf(
				fmt.Sprintf("%s.mapKeys[%d]", path, i),
				ErrInvalidField,
				DescriptorErrorReasonDuplicateListMapKey,
				"ListMap key %q is duplicated at indexes %d and %d",
				key,
				first,
				i,
			)
		}

		firstIndexes[key] = i
	}

	return nil
}

// listMapObject resolves the object payload used by ListMap semantics.
func listMapObject(elem Descriptor, resolver Resolver) (objectPayload, bool) {
	if elem.code == DescriptorObject {
		return elem.object, true
	}

	if elem.code != DescriptorRef || resolver == nil {
		return objectPayload{}, false
	}

	def, ok := resolver.Resolve(elem.ref.name)

	if !ok {
		return objectPayload{}, false
	}

	descriptor := def.Descriptor()

	if descriptor.code != DescriptorObject {
		return objectPayload{}, false
	}

	return descriptor.object, true
}

// validateListSetElementDescriptor checks the stable identity contract for
// ListSet elements.
//
// Set elements need deterministic future identity for validation, comparison,
// field ownership, and apply. This first pass deliberately limits set elements
// to non-nullable bool, string, and exact-width integer descriptors, including
// references that resolve to those descriptors.
func validateListSetElementDescriptor(
	desc Descriptor,
	resolver Resolver,
	path string,
	resolving map[TypeName]bool,
) error {
	if desc.Nullable() {
		return descriptorErrorf(
			path,
			ErrInvalidDescriptor,
			DescriptorErrorReasonInvalidListSetElement,
			"ListSet element descriptor must be non-nullable",
		)
	}

	switch desc.code {
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
		return validateListSetElementRef(desc.ref.name, resolver, path, resolving)
	default:
		return descriptorErrorf(
			path,
			ErrInvalidDescriptor,
			DescriptorErrorReasonInvalidListSetElement,
			"ListSet element descriptor %s cannot be represented as a stable identity value",
			desc.code,
		)
	}
}

// validateListSetElementRef resolves a referenced ListSet element descriptor
// while preserving descriptor-reference cycle checks.
func validateListSetElementRef(name TypeName, resolver Resolver, path string, resolving map[TypeName]bool) error {
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

	return validateListSetElementDescriptor(def.Descriptor(), resolver, path, next)
}
