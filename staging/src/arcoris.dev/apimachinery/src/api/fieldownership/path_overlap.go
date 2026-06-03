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

package fieldownership

import "arcoris.dev/apimachinery/api/fieldpath"

// pathsOverlap reports whether two semantic paths describe overlapping structure.
//
// Exact matches overlap. Ancestor/descendant pairs overlap in either direction.
// Sibling paths do not overlap.
func pathsOverlap(a fieldpath.Path, b fieldpath.Path) bool {
	return a.HasPrefix(b) || b.HasPrefix(a)
}
