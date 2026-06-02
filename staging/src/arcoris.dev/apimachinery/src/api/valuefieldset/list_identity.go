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

package valuefieldset

import (
	"arcoris.dev/apimachinery/api/fieldpath"
	"arcoris.dev/apimachinery/api/types"
	"arcoris.dev/apimachinery/api/value"
)

// TODO: move ListMap selector extraction into a narrow shared internal package
// when valuecompare/apply need the same identity behavior.

// listMapSelector extracts the stable semantic identity for one associative-list item.
func (e *extractor) listMapSelector(
	indexPath fieldpath.Path,
	item value.Value,
	element types.Type,
	keys []types.FieldName,
	depth int,
) (fieldpath.Selector, error) {
	if item.Kind() != value.KindObject {
		return fieldpath.Selector{}, errorAt(
			indexPath,
			ErrInvalidListKey,
			ErrorReasonInvalidListKey,
			"list map item is not an object",
		)
	}

	objectDescriptor, err := e.listMapObjectDescriptor(indexPath, element, depth)
	if err != nil {
		return fieldpath.Selector{}, err
	}

	itemObject, _ := item.Object()
	entries := make([]fieldpath.SelectorEntry, 0, len(keys))

	for _, key := range keys {
		entry, err := e.listMapSelectorEntry(
			indexPath,
			itemObject,
			objectDescriptor,
			key,
			depth,
		)
		if err != nil {
			return fieldpath.Selector{}, err
		}

		entries = append(entries, entry)
	}

	selector, err := fieldpath.NewSelector(entries...)
	if err != nil {
		return fieldpath.Selector{}, wrapAt(
			indexPath,
			ErrInvalidListKey,
			ErrorReasonInvalidListKey,
			"list map selector is invalid",
			err,
		)
	}

	return selector, nil
}

// listMapSelectorEntry extracts one field/value pair for a ListMap selector.
func (e *extractor) listMapSelectorEntry(
	indexPath fieldpath.Path,
	itemObject value.ObjectView,
	objectDescriptor types.ObjectView,
	key types.FieldName,
	depth int,
) (fieldpath.SelectorEntry, error) {
	keyName := string(key)
	keyPath := indexPath.Field(keyName)

	keyType, ok := listMapKeyDescriptor(objectDescriptor, key)
	if !ok {
		return fieldpath.SelectorEntry{}, errorfAt(
			keyPath,
			ErrInvalidDescriptor,
			ErrorReasonInvalidDescriptor,
			"list map key field %q is not declared by the element descriptor",
			keyName,
		)
	}

	keyValue, ok := itemObject.Get(keyName)
	if !ok {
		return fieldpath.SelectorEntry{}, errorfAt(
			keyPath,
			ErrInvalidListKey,
			ErrorReasonMissingListKey,
			"list map key field %q is missing",
			keyName,
		)
	}
	if keyValue.IsNull() {
		return fieldpath.SelectorEntry{}, errorfAt(
			keyPath,
			ErrInvalidListKey,
			ErrorReasonInvalidListKey,
			"list map key field %q is null",
			keyName,
		)
	}

	literal, err := e.listMapLiteral(keyPath, keyValue, keyType, depth)
	if err != nil {
		return fieldpath.SelectorEntry{}, err
	}

	return fieldpath.NewSelectorEntry(keyName, literal), nil
}

// listMapObjectDescriptor returns the object descriptor behind a ListMap element.
func (e *extractor) listMapObjectDescriptor(
	path fieldpath.Path,
	element types.Type,
	depth int,
) (types.ObjectView, error) {
	switch element.Code() {
	case types.TypeObject:
		objectView, ok := element.Object()
		if !ok {
			return types.ObjectView{}, errorAt(
				path,
				ErrInvalidDescriptor,
				ErrorReasonInvalidDescriptor,
				"list map element descriptor is not an object",
			)
		}

		return objectView, nil
	case types.TypeRef:
		resolved, err := e.resolveRefDescriptor(path, element, depth)
		if err != nil {
			return types.ObjectView{}, err
		}
		if resolved.Code() != types.TypeObject {
			return types.ObjectView{}, errorAt(
				path,
				ErrInvalidDescriptor,
				ErrorReasonInvalidDescriptor,
				"resolved list map element descriptor is not an object",
			)
		}

		objectView, ok := resolved.Object()
		if !ok {
			return types.ObjectView{}, errorAt(
				path,
				ErrInvalidDescriptor,
				ErrorReasonInvalidDescriptor,
				"resolved list map element descriptor is not an object",
			)
		}

		return objectView, nil
	default:
		return types.ObjectView{}, errorAt(
			path,
			ErrInvalidDescriptor,
			ErrorReasonInvalidDescriptor,
			"list map element descriptor is not an object or reference",
		)
	}
}

// listMapKeyDescriptor returns the descriptor for one declared selector field.
func listMapKeyDescriptor(
	objectView types.ObjectView,
	key types.FieldName,
) (types.Type, bool) {
	for _, fieldDescriptor := range objectView.Fields() {
		if fieldDescriptor.Name() == key {
			return fieldDescriptor.Type(), true
		}
	}

	return types.Type{}, false
}

// listMapLiteral converts one concrete key value to a selector literal.
func (e *extractor) listMapLiteral(
	path fieldpath.Path,
	val value.Value,
	descriptor types.Type,
	depth int,
) (fieldpath.Literal, error) {
	switch descriptor.Code() {
	case types.TypeBool:
		return boolListMapLiteral(path, val)
	case types.TypeString:
		return stringListMapLiteral(path, val)
	case types.TypeInt8,
		types.TypeInt16,
		types.TypeInt32,
		types.TypeInt64:
		return signedListMapLiteral(path, val)
	case types.TypeUint8,
		types.TypeUint16,
		types.TypeUint32,
		types.TypeUint64:
		return unsignedListMapLiteral(path, val)
	case types.TypeRef:
		resolved, err := e.resolveRefDescriptor(path, descriptor, depth)
		if err != nil {
			return fieldpath.Literal{}, err
		}

		return e.listMapLiteral(path, val, resolved, depth+1)
	default:
		return fieldpath.Literal{}, errorfAt(
			path,
			ErrInvalidListKey,
			ErrorReasonInvalidListKey,
			"list map key value kind %s cannot satisfy descriptor %s",
			val.Kind(),
			descriptor.Code(),
		)
	}
}

// boolListMapLiteral converts a boolean key to a selector literal.
func boolListMapLiteral(path fieldpath.Path, val value.Value) (fieldpath.Literal, error) {
	if val.Kind() != value.KindBool {
		return fieldpath.Literal{}, listMapLiteralKindError(path, val.Kind(), value.KindBool)
	}

	booleanValue, _ := val.Bool()
	return fieldpath.BoolLiteral(booleanValue), nil
}

// stringListMapLiteral converts a string key to a selector literal.
func stringListMapLiteral(path fieldpath.Path, val value.Value) (fieldpath.Literal, error) {
	if val.Kind() != value.KindString {
		return fieldpath.Literal{}, listMapLiteralKindError(path, val.Kind(), value.KindString)
	}

	text, _ := val.String()
	return fieldpath.StringLiteral(text), nil
}

// signedListMapLiteral converts a signed integer key to a selector literal.
func signedListMapLiteral(path fieldpath.Path, val value.Value) (fieldpath.Literal, error) {
	if val.Kind() != value.KindInteger {
		return fieldpath.Literal{}, listMapLiteralKindError(path, val.Kind(), value.KindInteger)
	}

	integerValue, _ := val.Integer()
	signedValue, ok := integerValue.Int64()
	if !ok {
		return fieldpath.Literal{}, errorAt(
			path,
			ErrInvalidListKey,
			ErrorReasonInvalidListKey,
			"list map key integer does not fit signed selector type",
		)
	}

	return fieldpath.Int64Literal(signedValue), nil
}

// unsignedListMapLiteral converts an unsigned integer key to a selector literal.
func unsignedListMapLiteral(path fieldpath.Path, val value.Value) (fieldpath.Literal, error) {
	if val.Kind() != value.KindInteger {
		return fieldpath.Literal{}, listMapLiteralKindError(path, val.Kind(), value.KindInteger)
	}

	integerValue, _ := val.Integer()
	unsignedValue, ok := integerValue.Uint64()
	if !ok {
		return fieldpath.Literal{}, errorAt(
			path,
			ErrInvalidListKey,
			ErrorReasonInvalidListKey,
			"list map key integer does not fit unsigned selector type",
		)
	}

	return fieldpath.Uint64Literal(unsignedValue), nil
}

// listMapLiteralKindError reports an unconvertible concrete key kind.
func listMapLiteralKindError(
	path fieldpath.Path,
	actual value.Kind,
	expected value.Kind,
) error {
	return errorfAt(
		path,
		ErrInvalidListKey,
		ErrorReasonInvalidListKey,
		"list map key value kind %s cannot become selector literal %s",
		actual,
		expected,
	)
}
