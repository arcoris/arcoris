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

package valuevalidation

import (
	"arcoris.dev/apimachinery/api/fieldpath"
	"arcoris.dev/apimachinery/api/types"
	"arcoris.dev/apimachinery/api/value"
)

// tryListMapSelector extracts the semantic identity selector for one ListMap item.
//
// The extraction path is intentionally quiet for ordinary payload failures.
// Missing keys, null keys, wrong key kinds, and non-object items are reported by
// recursive validation at the physical index path when selector extraction
// fails. This keeps key diagnostics centralized in the same descriptor-aware
// traversal used for every other object field.
func (v *validator) tryListMapSelector(
	indexPath fieldpath.Path,
	item value.Value,
	element types.Type,
	keys []types.FieldName,
) (fieldpath.Selector, bool) {
	if item.Kind() != value.KindObject {
		return fieldpath.Selector{}, false
	}

	objectView, ok := v.listMapObjectDescriptor(element)
	if !ok {
		return fieldpath.Selector{}, false
	}

	itemObject, _ := item.Object()
	entries := make([]fieldpath.SelectorEntry, 0, len(keys))

	for _, key := range keys {
		keyName := string(key)
		keyType, ok := listMapKeyDescriptor(objectView, key)
		if !ok {
			return fieldpath.Selector{}, false
		}

		keyValue, ok := itemObject.Get(keyName)
		if !ok {
			return fieldpath.Selector{}, false
		}
		if keyValue.IsNull() {
			return fieldpath.Selector{}, false
		}

		literal, ok := v.listMapLiteral(keyValue, keyType)
		if !ok {
			return fieldpath.Selector{}, false
		}

		entries = append(entries, fieldpath.NewSelectorEntry(keyName, literal))
	}

	selector, err := fieldpath.NewSelector(entries...)
	if err != nil {
		v.wrap(indexPath, ErrInvalidListKey, ErrorReasonInvalidListKey, "list map selector is invalid", err)
		return fieldpath.Selector{}, false
	}

	return selector, true
}

// listMapObjectDescriptor returns the object descriptor behind a ListMap element.
func (v *validator) listMapObjectDescriptor(element types.Type) (types.ObjectView, bool) {
	switch element.Code() {
	case types.TypeObject:
		return element.Object()
	case types.TypeRef:
		resolved, ok := v.resolveListMapRef(element, make(map[types.TypeName]bool))
		if !ok {
			return types.ObjectView{}, false
		}

		return resolved.Object()
	default:
		return types.ObjectView{}, false
	}
}

// listMapKeyDescriptor returns the descriptor for one declared selector field.
func listMapKeyDescriptor(objectView types.ObjectView, key types.FieldName) (types.Type, bool) {
	for _, fieldDescriptor := range objectView.Fields() {
		if fieldDescriptor.Name() == key {
			return fieldDescriptor.Type(), true
		}
	}

	return types.Type{}, false
}

// listMapLiteral converts one concrete key value to a selector literal.
func (v *validator) listMapLiteral(val value.Value, descriptor types.Type) (fieldpath.Literal, bool) {
	switch descriptor.Code() {
	case types.TypeBool:
		if val.Kind() != value.KindBool {
			return fieldpath.Literal{}, false
		}

		b, _ := val.Bool()
		return fieldpath.BoolLiteral(b), true
	case types.TypeString:
		if val.Kind() != value.KindString {
			return fieldpath.Literal{}, false
		}

		text, _ := val.String()
		return fieldpath.StringLiteral(text), true
	case types.TypeInt8,
		types.TypeInt16,
		types.TypeInt32,
		types.TypeInt64:
		return signedListMapLiteral(val)
	case types.TypeUint8,
		types.TypeUint16,
		types.TypeUint32,
		types.TypeUint64:
		return unsignedListMapLiteral(val)
	case types.TypeRef:
		resolved, ok := v.resolveListMapRef(descriptor, make(map[types.TypeName]bool))
		if !ok {
			return fieldpath.Literal{}, false
		}

		return v.listMapLiteral(val, resolved)
	}

	return fieldpath.Literal{}, false
}

// signedListMapLiteral converts a signed integer key to a selector literal.
func signedListMapLiteral(val value.Value) (fieldpath.Literal, bool) {
	if val.Kind() != value.KindInteger {
		return fieldpath.Literal{}, false
	}

	integer, _ := val.Integer()
	signed, ok := integer.Int64()
	if !ok {
		return fieldpath.Literal{}, false
	}

	return fieldpath.Int64Literal(signed), true
}

// unsignedListMapLiteral converts an unsigned integer key to a selector literal.
func unsignedListMapLiteral(val value.Value) (fieldpath.Literal, bool) {
	if val.Kind() != value.KindInteger {
		return fieldpath.Literal{}, false
	}

	integer, _ := val.Integer()
	unsigned, ok := integer.Uint64()
	if !ok {
		return fieldpath.Literal{}, false
	}

	return fieldpath.Uint64Literal(unsigned), true
}

// resolveListMapRef quietly resolves a descriptor reference for selector extraction.
func (v *validator) resolveListMapRef(
	descriptor types.Type,
	resolving map[types.TypeName]bool,
) (types.Type, bool) {
	view, ok := descriptor.Ref()
	if !ok || v.resolver == nil {
		return types.Type{}, false
	}

	name := view.Name()
	if resolving[name] {
		return types.Type{}, false
	}

	definition, ok := v.resolver.ResolveType(name)
	if !ok {
		return types.Type{}, false
	}

	resolving[name] = true
	resolved := definition.Type()
	if resolved.Code() == types.TypeRef {
		target, ok := v.resolveListMapRef(resolved, resolving)
		delete(resolving, name)
		return target, ok
	}

	delete(resolving, name)
	return resolved, true
}
