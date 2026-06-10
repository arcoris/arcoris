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

// NewKeyElement constructs one dynamic map-key path element.
func NewKeyElement(key MapKey) Element {
	return Element{kind: ElementKey, key: key}
}

// KeyElementFromString validates key and constructs one dynamic map-key element.
func KeyElementFromString(key string) (Element, error) {
	mapKey, err := NewMapKey(key)
	if err != nil {
		return Element{}, err
	}

	return NewKeyElement(mapKey), nil
}

// MustKeyElement constructs one dynamic map-key element or panics.
func MustKeyElement(key string) Element {
	element, err := KeyElementFromString(key)
	if err != nil {
		panic(err)
	}

	return element
}
