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

// validateList checks element type, length limits, semantics, and map keys.
func validateList(t Type, resolver Resolver, path string, resolving map[TypeName]bool) error {
	if t.list.elem == nil {
		return typeErrorf(
			path+".elem",
			ErrInvalidType,
			TypeErrorReasonMissingElement,
			"list descriptor must have an element type",
		)
	}
	if err := validateType(*t.list.elem, resolver, path+".elem", resolving); err != nil {
		return err
	}
	if err := validateLengthLimits(t.list.minLen, t.list.maxLen, path+".len"); err != nil {
		return err
	}
	if !t.list.semantics.IsValid() {
		return typeErrorf(
			path+".semantics",
			ErrInvalidType,
			TypeErrorReasonInvalidSemantics,
			"list semantics %d is not supported",
			t.list.semantics,
		)
	}
	if t.list.semantics != ListMap {
		return nil
	}
	if len(t.list.mapKeys) == 0 {
		return typeErrorf(
			path+".mapKeys",
			ErrInvalidField,
			TypeErrorReasonMissingListMapKey,
			"ListMap semantics requires at least one key field",
		)
	}
	object, ok := listMapObject(*t.list.elem, resolver)
	if !ok {
		return typeErrorf(
			path+".elem",
			ErrInvalidType,
			TypeErrorReasonListMapElementNotObject,
			"ListMap element must be an object descriptor or a TypeRef resolving to an object",
		)
	}
	fields := make(map[FieldName]FieldDescriptor, len(object.fields))
	for _, field := range object.fields {
		fields[field.name] = field
	}
	for i, key := range t.list.mapKeys {
		keyPath := fmt.Sprintf("%s.mapKeys[%d]", path, i)
		if !key.IsValid() {
			return typeErrorf(
				keyPath,
				ErrInvalidField,
				TypeErrorReasonInvalidListMapKey,
				"ListMap key %q is not a valid field name",
				key,
			)
		}
		field, ok := fields[key]
		if !ok {
			return typeErrorf(
				keyPath,
				ErrInvalidField,
				TypeErrorReasonListMapKeyNotFound,
				"ListMap key %q is not present in the object element",
				key,
			)
		}
		if !field.IsRequired() {
			return typeErrorf(
				keyPath,
				ErrInvalidField,
				TypeErrorReasonListMapKeyOptional,
				"ListMap key field %q must be required",
				key,
			)
		}
	}
	return nil
}

// listMapObject resolves the object payload used by ListMap semantics.
func listMapObject(elem Type, resolver Resolver) (objectPayload, bool) {
	if elem.code == TypeObject {
		return elem.object, true
	}
	if elem.code != TypeRef || resolver == nil {
		return objectPayload{}, false
	}
	def, ok := resolver.ResolveType(elem.ref.name)
	if !ok {
		return objectPayload{}, false
	}
	typ := def.Type()
	if typ.code != TypeObject {
		return objectPayload{}, false
	}
	return typ.object, true
}
