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

// listMapSelector extracts the semantic identity selector for one ListMap item.
//
// api/types validates that map key fields are required and have stable scalar
// descriptor types. This function remains defensive because concrete payloads
// can still be missing keys, null, or carry the wrong value kind.
func (v *validator) listMapSelector(
	indexPath fieldpath.Path,
	item value.Value,
	keys []types.FieldName,
) (fieldpath.Selector, bool) {
	if !v.requireKind(indexPath, item, value.KindObject, types.TypeObject) {
		return fieldpath.Selector{}, false
	}

	itemObject, _ := item.Object()
	entries := make([]fieldpath.SelectorEntry, 0, len(keys))

	for _, key := range keys {
		keyName := string(key)
		keyPath := indexPath.Field(keyName)

		keyValue, ok := itemObject.Get(keyName)
		if !ok {
			v.addf(
				keyPath,
				ErrInvalidListKey,
				ErrorReasonMissingListKey,
				"list map key field %q is missing",
				keyName,
			)
			return fieldpath.Selector{}, false
		}
		if keyValue.IsNull() {
			v.addf(
				keyPath,
				ErrInvalidListKey,
				ErrorReasonInvalidListKey,
				"list map key field %q is null",
				keyName,
			)
			return fieldpath.Selector{}, false
		}

		literal, ok := v.listMapLiteral(keyPath, keyValue)
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

// listMapLiteral converts one concrete key value to a selector literal.
func (v *validator) listMapLiteral(
	path fieldpath.Path,
	val value.Value,
) (fieldpath.Literal, bool) {
	switch val.Kind() {
	case value.KindBool:
		b, _ := val.Bool()
		return fieldpath.BoolLiteral(b), true
	case value.KindString:
		text, _ := val.String()
		return fieldpath.StringLiteral(text), true
	case value.KindInteger:
		integer, _ := val.Integer()
		if signed, ok := integer.Int64(); ok {
			return fieldpath.Int64Literal(signed), true
		}
		unsigned, ok := integer.Uint64()
		if ok {
			return fieldpath.Uint64Literal(unsigned), true
		}
	}

	v.addf(
		path,
		ErrInvalidListKey,
		ErrorReasonInvalidListKey,
		"value kind %s cannot be used as a list map key",
		val.Kind(),
	)

	return fieldpath.Literal{}, false
}
