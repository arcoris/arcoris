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

// NewSelectorElement validates selector and constructs one selector element.
//
// The selector payload is cloned on the boundary so later caller mutations do
// not affect the stored semantic path element.
func NewSelectorElement(selector Selector) (Element, error) {
	element := Element{kind: ElementSelector, selector: selector.clone()}
	if err := element.ValidateStructure(); err != nil {
		return Element{}, err
	}

	return element, nil
}

// MustSelectorElement constructs one selector element or panics.
func MustSelectorElement(selector Selector) Element {
	element, err := NewSelectorElement(selector)
	if err != nil {
		panic(err)
	}

	return element
}
