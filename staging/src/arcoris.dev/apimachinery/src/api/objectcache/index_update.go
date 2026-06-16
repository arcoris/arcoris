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

import "arcoris.dev/apimachinery/api/objectstore"

func (idx indexes) add(item objectstore.ListItem) {
	key := item.Key
	objectName := key.Object
	addIndexKey(idx.byNamespace, objectName.Namespace, key)
	addIndexKey(idx.byName, objectName.Name, key)
	addIndexKey(idx.byObject, objectNameKey{
		namespace: objectName.Namespace,
		name:      objectName.Name,
	}, key)

	for labelKey, value := range item.State.Object.ObjectMeta.Labels {
		addIndexKey(idx.byLabelKey, labelKey, key)
		addIndexKey(idx.byLabelValue, labelValueKey{key: labelKey, value: value}, key)
	}

	for annotationKey, value := range item.State.Object.ObjectMeta.Annotations {
		addIndexKey(idx.byAnnotationKey, annotationKey, key)
		addIndexKey(idx.byAnnotationValue, annotationValueKey{
			key:   annotationKey,
			value: value,
		}, key)
	}
}

func (idx indexes) remove(item objectstore.ListItem) {
	key := item.Key
	objectName := key.Object
	removeIndexKey(idx.byNamespace, objectName.Namespace, key)
	removeIndexKey(idx.byName, objectName.Name, key)
	removeIndexKey(idx.byObject, objectNameKey{
		namespace: objectName.Namespace,
		name:      objectName.Name,
	}, key)

	for labelKey, value := range item.State.Object.ObjectMeta.Labels {
		removeIndexKey(idx.byLabelKey, labelKey, key)
		removeIndexKey(idx.byLabelValue, labelValueKey{key: labelKey, value: value}, key)
	}

	for annotationKey, value := range item.State.Object.ObjectMeta.Annotations {
		removeIndexKey(idx.byAnnotationKey, annotationKey, key)
		removeIndexKey(idx.byAnnotationValue, annotationValueKey{
			key:   annotationKey,
			value: value,
		}, key)
	}
}
