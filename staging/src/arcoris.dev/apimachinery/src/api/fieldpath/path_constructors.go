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

// Root returns the semantic root path.
//
// The root path contains no elements and formats as "$".
func Root() Path {
	return Path{}
}

// NewPath constructs a validated path from semantic elements.
//
// The resulting path owns its internal element slice. Later mutations to the
// caller's input slice do not affect the returned path.
func NewPath(elements ...Element) (Path, error) {
	path := Path{
		elements: cloneElements(elements),
	}

	if err := path.ValidateStructure(); err != nil {
		return Path{}, err
	}

	return path, nil
}

// MustPath constructs a validated path or panics when elements are malformed.
//
// It is intended for tests and static semantic-path fixtures.
func MustPath(elements ...Element) Path {
	path, err := NewPath(elements...)
	if err != nil {
		panic(err)
	}

	return path
}
