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
	"reflect"

	"arcoris.dev/apimachinery/api/objectstore"
)

// sameState compares committed cache state for update/delete preconditions.
//
// objectstore.State does not currently expose a semantic equality helper. The
// cache therefore compares the full value model exactly, including revision,
// object envelope, desired/observed values, and ownership state.
func sameState(left objectstore.State, right objectstore.State) bool {
	return reflect.DeepEqual(left, right)
}
