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

package objectstorewatch

import (
	"arcoris.dev/apimachinery/api/objectstore"
	"arcoris.dev/apimachinery/api/objectwatch"
)

// Store is a composite contract for observable object stores.
//
// Store is not a concrete implementation and does not add reflector behavior.
// Writers and lifecycle layers should normally depend on objectstore.Store.
// Reflectors should normally depend on ListerWatcher because they do not need
// write methods. Components that require both mutation and list-to-watch
// continuity may depend on objectstorewatch.Store.
type Store interface {
	objectstore.Store
	CollectionLister
	objectwatch.Source
}

// CapableStore is a Store that can describe its watch capabilities.
type CapableStore interface {
	Store
	objectwatch.CapabilityReporter
}
