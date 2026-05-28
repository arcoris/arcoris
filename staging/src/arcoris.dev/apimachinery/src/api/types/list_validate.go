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

// validateList checks element type, length limits, semantics, and map keys.
func validateList(t Type, resolver Resolver, path string, resolving map[TypeName]bool) error {
	if t.list.elem == nil {
		return typeError(path+".elem", ErrInvalidType)
	}
	if err := validateType(*t.list.elem, resolver, path+".elem", resolving); err != nil {
		return err
	}
	if err := validateLengthLimits(t.list.minLen, t.list.maxLen, path+".len"); err != nil {
		return err
	}
	if !t.list.semantics.IsValid() {
		return typeError(path+".semantics", ErrInvalidType)
	}
	if t.list.semantics != ListMap {
		return nil
	}
	if len(t.list.mapKeys) == 0 {
		return typeError(path+".mapKeys", ErrInvalidField)
	}
	object, ok := listMapObject(*t.list.elem, resolver)
	if !ok {
		return typeError(path+".elem", ErrInvalidType)
	}
	fields := make(map[FieldName]FieldDescriptor, len(object.fields))
	for _, field := range object.fields {
		fields[field.name] = field
	}
	for _, key := range t.list.mapKeys {
		if !key.IsValid() {
			return typeError(path+".mapKeys", ErrInvalidField)
		}
		field, ok := fields[key]
		if !ok || !field.IsRequired() {
			return typeError(path+".mapKeys", ErrInvalidField)
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
