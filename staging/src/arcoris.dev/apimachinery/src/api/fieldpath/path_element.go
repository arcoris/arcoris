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
// Field and key elements both store string names, but they remain distinct
// kinds because descriptor-aware layers know whether a step addresses a fixed
// object field or a dynamic map key. Index and selector elements model ordered
// and associative list addressing respectively.
type Element struct {
	kind     ElementKind
	name     string
	index    int
	selector Selector
}

// Kind returns the semantic category stored in e.
func (e Element) Kind() ElementKind {
	return e.kind
}

// Name returns the field name or key string for name-bearing elements.
func (e Element) Name() string {
	if e.kind != ElementField && e.kind != ElementKey {
		return ""
	}

	return e.name
}

// Index returns the list index for index elements.
func (e Element) Index() int {
	if e.kind != ElementIndex {
		return 0
	}

	return e.index
}

// Selector returns a detached selector for selector elements.
func (e Element) Selector() Selector {
	if e.kind != ElementSelector {
		return Selector{}
	}

	return e.selector.clone()
}

// clone returns a detached element copy.
func (e Element) clone() Element {
	if e.kind == ElementSelector {
		e.selector = e.selector.clone()
	}

	return e
}
