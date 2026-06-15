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

// Build constructs a static shallow index over items.
//
// Build copies the input slice structure so later caller appends or element
// replacement cannot change the index. It does not clone object states,
// metadata maps, or payload values, and it does not validate object identity,
// metadata syntax, resource contracts, uniqueness, or state correctness.
func Build(items []objectstore.ListItem) Index {
	if len(items) == 0 {
		return Index{}
	}

	idx := Index{
		items:             append([]objectstore.ListItem(nil), items...),
		byNamespace:       make(map[metaidentity.Namespace][]int, len(items)),
		byName:            make(map[metaidentity.Name][]int, len(items)),
		byObjectName:      make(map[objectNameKey][]int, len(items)),
		byLabelKey:        make(map[labels.Key][]int),
		byLabelValue:      make(map[labelValueKey][]int),
		byAnnotationKey:   make(map[annotations.Key][]int),
		byAnnotationValue: make(map[annotationValueKey][]int),
	}

	for pos, item := range idx.items {
		idx.indexIdentity(pos, item)
		idx.indexLabels(pos, item)
		idx.indexAnnotations(pos, item)
	}

	return idx
}
