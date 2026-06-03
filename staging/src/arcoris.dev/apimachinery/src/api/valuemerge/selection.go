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

import "arcoris.dev/apimachinery/api/fieldpath"

// mergeSelection classifies selected paths relative to one merge node.
type mergeSelection struct {
	// exact reports that the current path itself is selected.
	exact bool

	// descendants contains selected paths strictly below the current path.
	descendants fieldpath.Set
}

// selected reports whether the current node must do merge work.
func (s mergeSelection) selected() bool {
	return s.exact || !s.descendants.IsEmpty()
}

// selectAt classifies the selected field set for one semantic path.
func selectAt(fields fieldpath.Set, path fieldpath.Path) mergeSelection {
	selection := mergeSelection{}

	for _, selected := range fields.Paths() {
		switch {
		case selected.Equal(path):
			selection.exact = true
		case selected.IsDescendantOf(path):
			selection.descendants = selection.descendants.Insert(selected)
		}
	}

	return selection
}

// hasSelectedChild reports whether fields selects path or any descendant.
func hasSelectedChild(fields fieldpath.Set, path fieldpath.Path) bool {
	return selectAt(fields, path).selected()
}
