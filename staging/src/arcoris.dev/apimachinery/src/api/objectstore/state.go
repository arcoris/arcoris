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

package objectstore

import (
	"arcoris.dev/apimachinery/api/object"
	"arcoris.dev/apimachinery/api/objectownership"
	"arcoris.dev/apimachinery/api/value"
)

// State is the committed live state value stored for one object key.
//
// Object contains the authoritative live object envelope. Ownership contains
// the canonical object ownership state associated with that live object.
// Revision is the store-local commit revision assigned by the store.
type State struct {
	// Object is the committed value-backed API object envelope.
	Object object.Object[value.Value, value.Value]
	// Ownership is the committed ownership state for modeled object surfaces.
	Ownership objectownership.State
	// Revision is the store-local revision assigned at commit time.
	Revision Revision
}
