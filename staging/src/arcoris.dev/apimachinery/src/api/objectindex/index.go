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

package objectindex

import (
	"arcoris.dev/apimachinery/api/meta/annotations"
	metaidentity "arcoris.dev/apimachinery/api/meta/identity"
	"arcoris.dev/apimachinery/api/meta/labels"
	"arcoris.dev/apimachinery/api/objectstore"
)

// Index is a static in-memory index over already-loaded object list items.
//
// Index values are immutable by convention after Build returns. The item slice
// is detached from caller slice mutations, but the ListItem values are shallow:
// objectstore.State and reference-bearing payloads are not cloned.
type Index struct {
	items []objectstore.ListItem

	byNamespace       map[metaidentity.Namespace][]int
	byName            map[metaidentity.Name][]int
	byObjectName      map[objectNameKey][]int
	byLabelKey        map[labels.Key][]int
	byLabelValue      map[labelValueKey][]int
	byAnnotationKey   map[annotations.Key][]int
	byAnnotationValue map[annotationValueKey][]int
}

// Len returns the number of items retained by idx.
func (idx Index) Len() int {
	return len(idx.items)
}

// IsZero reports whether idx contains no indexed items.
func (idx Index) IsZero() bool {
	return len(idx.items) == 0
}
