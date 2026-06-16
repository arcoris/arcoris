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
)

// objectNameKey is the comparable index key for exact namespace/name identity.
type objectNameKey struct {
	// namespace is the storage identity namespace. The zero namespace is a valid
	// bucket for global objects.
	namespace metaidentity.Namespace

	// name is the storage identity name.
	name metaidentity.Name
}

// labelValueKey is the comparable index key for one label key/value pair.
type labelValueKey struct {
	// key is the canonical label key.
	key labels.Key

	// value is the canonical label value stored on an item.
	value labels.Value
}

// annotationValueKey is the comparable index key for one annotation key/value
// pair.
type annotationValueKey struct {
	// key is the canonical annotation key.
	key annotations.Key

	// value is the canonical annotation value stored on an item.
	value annotations.Value
}
