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

	// hasDescendant reports that at least one selected path is below this path.
	hasDescendant bool
}

// selected reports whether the current node must do merge work.
func (s mergeSelection) selected() bool {
	return s.exact || s.hasDescendant
}

// selectAt classifies the selected field set for one semantic path.
func selectAt(fields fieldpath.Set, path fieldpath.Path) mergeSelection {
	selection := mergeSelection{}

	fields.ForEach(func(_ int, selected fieldpath.Path) bool {
		switch {
		case selected.Equal(path):
			selection.exact = true
		case selected.IsDescendantOf(path):
			selection.hasDescendant = true
		}
		return true
	})

	return selection
}

// hasSelectedChild reports whether fields selects path or any descendant.
func hasSelectedChild(fields fieldpath.Set, path fieldpath.Path) bool {
	return selectAt(fields, path).selected()
}
