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

package fieldpath

// Element is one semantic step in a structured field path.
//
// Field and key elements both store named tokens, but they remain distinct
// kinds because descriptor-aware layers know whether a step addresses a fixed
// object field or a dynamic map key. Index and selector elements model ordered
// and associative list addressing respectively.
type Element struct {
	kind     ElementKind
	field    FieldName
	key      MapKey
	index    int
	selector Selector
}

// Kind returns the semantic category stored in e.
func (e Element) Kind() ElementKind {
	return e.kind
}

// AsField returns the fixed field name for field elements.
func (e Element) AsField() (FieldName, bool) {
	if e.kind != ElementField {
		return "", false
	}

	return e.field, true
}

// AsKey returns the dynamic map key for key elements.
func (e Element) AsKey() (MapKey, bool) {
	if e.kind != ElementKey {
		return "", false
	}

	return e.key, true
}

// AsIndex returns the list index for index elements.
func (e Element) AsIndex() (int, bool) {
	if e.kind != ElementIndex {
		return 0, false
	}

	return e.index, true
}

// AsSelector returns a detached selector for selector elements.
func (e Element) AsSelector() (Selector, bool) {
	if e.kind != ElementSelector {
		return Selector{}, false
	}

	return e.selector.clone(), true
}
