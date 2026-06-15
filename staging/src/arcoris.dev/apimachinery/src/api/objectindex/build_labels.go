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

// indexLabels records item position under label-existence and label-value keys.
//
// The index trusts the already-loaded object metadata as provided by the
// caller. Syntax validation belongs to metadata/query construction layers.
func (idx *Index) indexLabels(pos int, item objectstore.ListItem) {
	for key, value := range item.State.Object.ObjectMeta.Labels {
		idx.byLabelKey[key] = append(idx.byLabelKey[key], pos)
		valueKey := labelValueKey{key: key, value: value}
		idx.byLabelValue[valueKey] = append(idx.byLabelValue[valueKey], pos)
	}
}
