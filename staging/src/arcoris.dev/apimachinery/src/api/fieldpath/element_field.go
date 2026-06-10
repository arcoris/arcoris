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

// NewFieldElement constructs one fixed-field path element.
func NewFieldElement(name FieldName) Element {
	return Element{kind: ElementField, field: name}
}

// FieldElementFromString validates name and constructs one fixed-field element.
func FieldElementFromString(name string) (Element, error) {
	fieldName, err := NewFieldName(name)
	if err != nil {
		return Element{}, err
	}

	return NewFieldElement(fieldName), nil
}

// MustFieldElement constructs one fixed-field element or panics.
func MustFieldElement(name string) Element {
	element, err := FieldElementFromString(name)
	if err != nil {
		panic(err)
	}

	return element
}
