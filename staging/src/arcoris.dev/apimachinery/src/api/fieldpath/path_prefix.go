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

// Parent returns the immediate semantic parent path.
func (p Path) Parent() (Path, bool) {
	if len(p.elements) == 0 {
		return Path{}, false
	}

	return Path{elements: cloneElements(p.elements[:len(p.elements)-1])}, true
}

// HasPrefix reports whether prefix is a structural prefix of p.
func (p Path) HasPrefix(prefix Path) bool {
	if len(prefix.elements) > len(p.elements) {
		return false
	}

	for i := range prefix.elements {
		if !p.elements[i].Equal(prefix.elements[i]) {
			return false
		}
	}

	return true
}
