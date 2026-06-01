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

// IsDescendantOf reports whether p is a strict descendant of ancestor.
//
// The receiver must be longer than ancestor and share ancestor as a structural
// prefix.
func (p Path) IsDescendantOf(ancestor Path) bool {
	return p.Len() > ancestor.Len() && p.HasPrefix(ancestor)
}

// IntersectsSubtree reports whether p and other lie in the same subtree.
//
// Two paths intersect structurally when either one is a prefix of the other.
// This makes the operation symmetric and useful for subtree filtering.
func (p Path) IntersectsSubtree(other Path) bool {
	return p.HasPrefix(other) || other.HasPrefix(p)
}
