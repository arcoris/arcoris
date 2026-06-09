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

package objectlifecycle

import (
	metaidentity "arcoris.dev/apimachinery/api/meta/identity"
	"arcoris.dev/apimachinery/api/objectstore"
)

// keyFor constructs the objectstore key for a resolved resource and object name.
func keyFor(op Operation, resolved resolvedResource, name metaidentity.ObjectName) (objectstore.Key, error) {
	key, err := objectstore.NewKey(resolved.gvr, name)
	if err != nil {
		return objectstore.Key{}, errorFor(op, ErrorReasonInvalidRequest, objectstore.Key{}, ErrInvalidRequest, err)
	}

	return key, nil
}
