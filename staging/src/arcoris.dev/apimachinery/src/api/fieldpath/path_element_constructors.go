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

// FieldElement constructs one fixed-field path element.
//
// Construction is intentionally lightweight. Validation is deferred so callers
// can build whole paths first and then receive one structured path-level error.
func FieldElement(name string) Element {
	return Element{
		kind: ElementField,
		name: name,
	}
}

// KeyElement constructs one dynamic map-key path element.
func KeyElement(key string) Element {
	return Element{
		kind: ElementKey,
		name: key,
	}
}

// IndexElement constructs one positional list-index path element.
func IndexElement(index int) Element {
	return Element{
		kind:  ElementIndex,
		index: index,
	}
}

// SelectorElement constructs one associative-list selector path element.
//
// The selector payload is cloned on the boundary so later caller mutations do
// not affect the stored semantic path element.
func SelectorElement(selector Selector) Element {
	return Element{
		kind:     ElementSelector,
		selector: selector.clone(),
	}
}
