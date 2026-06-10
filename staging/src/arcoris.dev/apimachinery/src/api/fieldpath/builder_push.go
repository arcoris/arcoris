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

// Push appends one element.
func (b *Builder) Push(element Element) {
	b.elements = append(b.elements, element.clone())
}

// PushField appends one fixed-field element.
func (b *Builder) PushField(name FieldName) {
	b.Push(NewFieldElement(name))
}

// PushKey appends one dynamic map-key element.
func (b *Builder) PushKey(key MapKey) {
	b.Push(NewKeyElement(key))
}

// PushIndex appends one positional list-index element.
func (b *Builder) PushIndex(index int) error {
	element, err := NewIndexElement(index)
	if err != nil {
		return err
	}

	b.Push(element)
	return nil
}

// PushSelector appends one associative-list selector element.
func (b *Builder) PushSelector(selector Selector) error {
	element, err := NewSelectorElement(selector)
	if err != nil {
		return err
	}

	b.Push(element)
	return nil
}

// Pop removes the last element.
func (b *Builder) Pop() bool {
	if b == nil || len(b.elements) == 0 {
		return false
	}

	b.elements = b.elements[:len(b.elements)-1]
	return true
}
