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

import "testing"

func TestPathIsDescendantOf(t *testing.T) {
	ancestor := Root().Field(testField("spec"))
	descendant := Root().Field(testField("spec")).Field(testField("replicas"))

	requireEqual(t, descendant.IsDescendantOf(ancestor), true)
	requireEqual(t, ancestor.IsDescendantOf(ancestor), false)
	requireEqual(t, Root().Field(testField("status")).IsDescendantOf(ancestor), false)
}

func TestPathIsAncestorOf(t *testing.T) {
	ancestor := Root().Field(testField("spec"))
	descendant := Root().Field(testField("spec")).Field(testField("replicas"))

	requireEqual(t, ancestor.IsAncestorOf(descendant), true)
	requireEqual(t, ancestor.IsAncestorOf(ancestor), false)
	requireEqual(t, Root().Field(testField("status")).IsAncestorOf(descendant), false)
}

func TestPathOverlaps(t *testing.T) {
	left := Root().Field(testField("spec"))
	right := Root().Field(testField("spec")).Field(testField("replicas"))
	other := Root().Field(testField("status"))

	requireEqual(t, left.Overlaps(left), true)
	requireEqual(t, left.Overlaps(right), true)
	requireEqual(t, right.Overlaps(left), true)
	requireEqual(t, left.Overlaps(other), false)
}
