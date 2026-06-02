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

// Insert returns a detached set containing path.
//
// Insert preserves sorted duplicate-free representation and does not mutate the
// receiver. The inserted path must be valid; malformed paths panic because this
// method has no error return and Set cannot store invalid paths.
func (s Set) Insert(path Path) Set {
	if err := path.Validate(); err != nil {
		panic(err)
	}

	index, found := findSetPath(s.paths, path)
	if found {
		return Set{paths: clonePaths(s.paths)}
	}

	paths := make([]Path, 0, len(s.paths)+1)
	paths = appendSetPaths(paths, s.paths[:index])
	paths = appendSetPath(paths, path)
	paths = appendSetPaths(paths, s.paths[index:])

	return Set{paths: paths}
}

// Delete returns a detached set without path.
//
// Deleting a missing path is a no-op. The receiver is never mutated.
func (s Set) Delete(path Path) Set {
	index, found := findSetPath(s.paths, path)
	if !found {
		return Set{paths: clonePaths(s.paths)}
	}

	paths := make([]Path, 0, len(s.paths)-1)
	paths = appendSetPaths(paths, s.paths[:index])
	paths = appendSetPaths(paths, s.paths[index+1:])

	return Set{paths: paths}
}
