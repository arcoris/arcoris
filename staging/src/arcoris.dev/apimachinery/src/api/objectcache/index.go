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

package objectcache

import (
	"arcoris.dev/apimachinery/api/meta/annotations"
	metaidentity "arcoris.dev/apimachinery/api/meta/identity"
	"arcoris.dev/apimachinery/api/meta/labels"
	"arcoris.dev/apimachinery/api/objectstore"
)

// indexes stores private query acceleration buckets for cache-owned items.
//
// The canonical items map remains the source of live objects. Indexes narrow
// candidates only; objectquery.Predicate.Match still owns final semantics.
type indexes struct {
	// byNamespace narrows identity namespace equality.
	byNamespace map[metaidentity.Namespace]keySet

	// byName narrows identity name equality across namespaces.
	byName map[metaidentity.Name]keySet

	// byObject narrows the combined namespace/name identity requirement.
	byObject map[objectNameKey]keySet

	// byLabelKey narrows label exists requirements.
	byLabelKey map[labels.Key]keySet
	// byLabelValue narrows label equals/in requirements.
	byLabelValue map[labelValueKey]keySet

	// byAnnotationKey narrows annotation exists requirements.
	byAnnotationKey map[annotations.Key]keySet
	// byAnnotationValue narrows annotation equals/in requirements.
	byAnnotationValue map[annotationValueKey]keySet
}

// newIndexes returns empty but ready-to-use index buckets.
func newIndexes() indexes {
	return indexes{
		byNamespace:       map[metaidentity.Namespace]keySet{},
		byName:            map[metaidentity.Name]keySet{},
		byObject:          map[objectNameKey]keySet{},
		byLabelKey:        map[labels.Key]keySet{},
		byLabelValue:      map[labelValueKey]keySet{},
		byAnnotationKey:   map[annotations.Key]keySet{},
		byAnnotationValue: map[annotationValueKey]keySet{},
	}
}

// addIndexKey records one storage key in a single comparable bucket.
func addIndexKey[K comparable](buckets map[K]keySet, bucket K, key objectstore.Key) {
	set := buckets[bucket]
	if set == nil {
		set = keySet{}
		buckets[bucket] = set
	}
	set.add(key)
}

// removeIndexKey removes one storage key and deletes the bucket when it becomes
// empty, keeping later candidate planning simple.
func removeIndexKey[K comparable](buckets map[K]keySet, bucket K, key objectstore.Key) {
	set := buckets[bucket]
	if set == nil {
		return
	}
	set.remove(key)
	if len(set) == 0 {
		delete(buckets, bucket)
	}
}
