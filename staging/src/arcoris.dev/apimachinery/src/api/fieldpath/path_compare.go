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

import "cmp"

// Equal reports whether e and other are the same semantic path element.
func (e Element) Equal(other Element) bool {
	return e.Compare(other) == 0
}

// Compare imposes deterministic total ordering on path elements.
func (e Element) Compare(other Element) int {
	switch {
	case e.kind < other.kind:
		return -1
	case e.kind > other.kind:
		return 1
	}

	switch e.kind {
	case ElementField:
		return cmp.Compare(e.field, other.field)
	case ElementKey:
		return cmp.Compare(e.key, other.key)
	case ElementIndex:
		return cmp.Compare(e.index, other.index)
	case ElementSelector:
		return e.selector.Compare(other.selector)
	default:
		return 0
	}
}

// Equal reports whether p and other identify the same semantic payload location.
func (p Path) Equal(other Path) bool {
	return p.Compare(other) == 0
}

// Compare imposes deterministic total ordering on semantic field paths.
func (p Path) Compare(other Path) int {
	for i := 0; i < len(p.elements) && i < len(other.elements); i++ {
		if cmp := p.elements[i].Compare(other.elements[i]); cmp != 0 {
			return cmp
		}
	}

	switch {
	case len(p.elements) < len(other.elements):
		return -1
	case len(p.elements) > len(other.elements):
		return 1
	default:
		return 0
	}
}
