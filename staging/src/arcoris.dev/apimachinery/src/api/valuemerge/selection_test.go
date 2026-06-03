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

package valuemerge

import "testing"

func TestSelectionExact(t *testing.T) {
	path := root().Field("spec")
	got := selectAt(pathSet(path), path)

	if !got.exact || !got.descendants.IsEmpty() {
		t.Fatalf("selection = %#v", got)
	}
}

func TestSelectionDescendant(t *testing.T) {
	path := root().Field("spec")
	selected := path.Field("replicas")
	got := selectAt(pathSet(selected), path)

	if got.exact || !got.descendants.Has(selected) {
		t.Fatalf("selection = %#v", got)
	}
}

func TestSelectionIrrelevant(t *testing.T) {
	path := root().Field("spec")
	got := selectAt(pathSet(root().Field("metadata")), path)

	if got.selected() {
		t.Fatalf("selection = %#v", got)
	}
}

func TestSelectionAncestorIsNotSelectedLocally(t *testing.T) {
	path := root().Field("spec").Field("replicas")
	got := selectAt(pathSet(root().Field("spec")), path)

	if got.selected() {
		t.Fatalf("selection = %#v", got)
	}
}

func TestSelectionMapKey(t *testing.T) {
	path := root().Field("labels")
	selected := path.Key("app")
	got := selectAt(pathSet(selected), path)

	if !got.descendants.Has(selected) {
		t.Fatalf("selection descendants = %s", got.descendants)
	}
}

func TestSelectionListIndex(t *testing.T) {
	path := root()
	selected := path.Index(1)
	got := selectAt(pathSet(selected), path)

	if !got.descendants.Has(selected) {
		t.Fatalf("selection descendants = %s", got.descendants)
	}
}
