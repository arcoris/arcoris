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
	"context"

	"arcoris.dev/apimachinery/api/objectstore"
)

// CollectionLister produces boundary-safe collection reads for structural collections.
//
// ListCollection is stronger than raw objectstore.Store.List because it returns
// a validated CollectionRead intended for watch continuation on the same
// ListerWatcher or Store boundary. It does not itself prove future history
// availability; that proof happens when the same source serves Watch.
type CollectionLister interface {
	// ListCollection reads collection and returns a validated collection read.
	ListCollection(context.Context, objectstore.ListRequest) (CollectionRead, error)
}
