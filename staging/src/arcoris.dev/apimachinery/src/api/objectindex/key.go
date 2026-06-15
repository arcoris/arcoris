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
)

// objectNameKey is the comparable identity key used for namespace/name lookups.
//
// The index intentionally uses objectstore.Key.Object identity instead of
// object payload metadata. Storage identity is the list-item source of truth,
// while object metadata may evolve independently.
type objectNameKey struct {
	namespace metaidentity.Namespace
	name      metaidentity.Name
}

// labelValueKey is the comparable key for exact label key/value matches.
type labelValueKey struct {
	key   labels.Key
	value labels.Value
}

// annotationValueKey is the comparable key for exact annotation key/value matches.
type annotationValueKey struct {
	key   annotations.Key
	value annotations.Value
}
