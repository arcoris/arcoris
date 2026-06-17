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

package objectquery

import (
	"arcoris.dev/apimachinery/api/fieldpath"
	"arcoris.dev/apimachinery/api/objectstore"
	"arcoris.dev/apimachinery/api/objectsurface"
	"arcoris.dev/apimachinery/api/value"
)

// fieldLookup is the internal result of resolving a selectable field path.
// present distinguishes missing from present null.
type fieldLookup struct {
	// value is meaningful only when present is true.
	value value.Value
	// present reports that every path segment existed.
	present bool
}

// lookupFieldValue resolves a registered surface-relative field path against a
// list item. It intentionally supports only explicitly queryable surfaces.
func lookupFieldValue(item objectstore.ListItem, ref FieldRef) fieldLookup {
	var root value.Value
	kinds := objectsurface.Kinds()
	switch ref.Surface {
	case kinds.Desired():
		root = item.State.Object.Desired
	case kinds.Observed():
		if item.State.Object.Observed == nil {
			return fieldLookup{}
		}
		root = *item.State.Object.Observed
	default:
		return fieldLookup{}
	}
	if root.IsZero() {
		return fieldLookup{}
	}

	current := root
	ok := true
	ref.Path.ForEach(func(_ int, element fieldpath.Element) bool {
		current, ok = lookupFieldElement(current, element)
		return ok
	})
	if !ok {
		return fieldLookup{}
	}

	return fieldLookup{value: current, present: true}
}

// lookupFieldElement applies one semantic fieldpath element to current.
func lookupFieldElement(current value.Value, element fieldpath.Element) (value.Value, bool) {
	switch element.Kind() {
	case fieldpath.ElementField:
		field, _ := element.AsField()
		return lookupRecordMember(current, field.String())
	case fieldpath.ElementKey:
		key, _ := element.AsKey()
		return lookupRecordMember(current, key.String())
	case fieldpath.ElementIndex:
		index, _ := element.AsIndex()
		list, ok := current.AsList()
		if !ok {
			return value.Value{}, false
		}
		return list.At(index)
	case fieldpath.ElementSelector:
		selector, _ := element.AsSelector()
		return lookupSelectedListItem(current, selector)
	default:
		return value.Value{}, false
	}
}

// lookupRecordMember treats field and key path elements as exact record member
// lookups. Descriptor-aware layers decide whether a record represents an object
// or dynamic map.
func lookupRecordMember(current value.Value, name string) (value.Value, bool) {
	record, ok := current.AsRecord()
	if !ok {
		return value.Value{}, false
	}

	return record.Get(value.MemberName(name))
}

// lookupSelectedListItem resolves an associative-list selector by scanning
// concrete list items in order and matching every selector entry exactly.
func lookupSelectedListItem(current value.Value, selector fieldpath.Selector) (value.Value, bool) {
	list, ok := current.AsList()
	if !ok {
		return value.Value{}, false
	}

	for i := 0; i < list.Len(); i++ {
		item, _ := list.At(i)
		if selectorMatchesValue(selector, item) {
			return item, true
		}
	}

	return value.Value{}, false
}

// selectorMatchesValue reports whether item is a record with every selector
// entry present and equal to the selector literal.
func selectorMatchesValue(selector fieldpath.Selector, item value.Value) bool {
	record, ok := item.AsRecord()
	if !ok {
		return false
	}

	matched := true
	selector.ForEach(func(_ int, entry fieldpath.SelectorEntry) bool {
		actual, ok := record.Get(value.MemberName(entry.Field().String()))
		if !ok || !literalMatchesValue(entry.Value(), actual) {
			matched = false
			return false
		}
		return true
	})

	return matched
}

// literalMatchesValue compares a fieldpath selector literal with the concrete
// Value kinds available to selectors.
func literalMatchesValue(literal fieldpath.Literal, actual value.Value) bool {
	switch literal.Kind() {
	case fieldpath.LiteralBool:
		expected, _ := literal.AsBool()
		got, ok := actual.AsBool()
		return ok && got == expected
	case fieldpath.LiteralString:
		expected, _ := literal.AsString()
		got, ok := actual.AsString()
		return ok && got == expected
	case fieldpath.LiteralInteger:
		got, ok := actual.AsInteger()
		if !ok {
			return false
		}
		if expected, ok := literal.AsInt64(); ok {
			return got.Equal(value.NewIntegerFromInt64(expected))
		}
		expected, ok := literal.AsUint64()
		return ok && got.Equal(value.NewIntegerFromUint64(expected))
	default:
		return false
	}
}
