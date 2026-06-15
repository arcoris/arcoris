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

import "arcoris.dev/apimachinery/api/objectstore"

// indexIdentity records item position under each storage-identity lookup key.
//
// This indexes objectstore.ListItem.Key.Object, not object metadata fields.
// Duplicate object keys are retained by appending every position.
func (idx *Index) indexIdentity(pos int, item objectstore.ListItem) {
	objectName := item.Key.Object
	idx.byNamespace[objectName.Namespace] = append(idx.byNamespace[objectName.Namespace], pos)
	idx.byName[objectName.Name] = append(idx.byName[objectName.Name], pos)
	key := objectNameKey{
		namespace: objectName.Namespace,
		name:      objectName.Name,
	}
	idx.byObjectName[key] = append(idx.byObjectName[key], pos)
}
