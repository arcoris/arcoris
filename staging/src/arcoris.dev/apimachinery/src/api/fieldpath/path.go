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

// Path is an ordered semantic payload location.
//
// A Path is immutable by API convention. Callers build larger paths from
// smaller ones through append-style helpers. Slices returned by the API are
// always detached so comparison, caching, and sharing stay predictable.
type Path struct {
	elements []Element
}

// Len returns the number of path elements below root.
func (p Path) Len() int {
	return len(p.elements)
}

// IsRoot reports whether p is the semantic root path.
func (p Path) IsRoot() bool {
	return len(p.elements) == 0
}

// Elements returns detached path elements in order.
func (p Path) Elements() []Element {
	return cloneElements(p.elements)
}

// Element returns one path element by index.
func (p Path) Element(index int) (Element, bool) {
	if index < 0 || index >= len(p.elements) {
		return Element{}, false
	}

	return p.elements[index], true
}

// ForEach visits path elements in order until fn returns false.
func (p Path) ForEach(fn func(index int, element Element) bool) {
	for i, element := range p.elements {
		if !fn(i, element) {
			return
		}
	}
}
