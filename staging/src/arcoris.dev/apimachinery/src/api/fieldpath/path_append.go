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

// Append returns a detached path with elements appended in order.
//
// The receiver is never mutated. The returned path owns a fresh element slice.
func (p Path) Append(elements ...Element) Path {
	appended := cloneElements(p.elements)
	appended = append(appended, cloneElements(elements)...)

	return Path{
		elements: appended,
	}
}

// Field appends one fixed-field semantic element.
//
// This helper keeps path-building call sites readable while the underlying
// constructor remains explicit as FieldElement.
func (p Path) Field(name string) Path {
	return p.Append(FieldElement(name))
}

// Key appends one dynamic map-key semantic element.
func (p Path) Key(key string) Path {
	return p.Append(KeyElement(key))
}

// Index appends one positional list-index semantic element.
func (p Path) Index(index int) Path {
	return p.Append(IndexElement(index))
}

// Select appends one associative-list selector semantic element.
//
// Selector addressing is intended for associative lists whose identity comes
// from one or more key fields rather than from unstable positional indexes.
func (p Path) Select(selector Selector) Path {
	return p.Append(SelectorElement(selector))
}
