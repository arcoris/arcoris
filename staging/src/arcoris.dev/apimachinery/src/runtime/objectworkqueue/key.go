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

package objectworkqueue

import "arcoris.dev/apimachinery/api/objectstore"

// keyID is the queue's private comparable identity for an item.
//
// objectstore.Key is a comparable value type, so the queue can use the
// authoritative API object identity directly without string encoding,
// reflection, marshaling, or lossy normalization.
type keyID = objectstore.Key

// keyForItem returns the stable deduplication identity for item.
func keyForItem(item Item) keyID {
	return item.Key
}
