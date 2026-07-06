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

package objectreconcilewrite

import "arcoris.dev/apimachinery/api/objectlifecycle"

// Delete builds a lifecycle request that deletes the current object.
//
// Expected is Current.Revision. Delete only constructs the request; not-found,
// conflict, and tombstone behavior remain objectlifecycle responsibilities.
func (c Current) Delete() (objectlifecycle.DeleteRequest, error) {
	if err := c.validate(); err != nil {
		return objectlifecycle.DeleteRequest{}, err
	}

	return objectlifecycle.DeleteRequest{
		Resource: c.key.Resource,
		Object:   c.key.Object,
		Expected: c.Revision(),
	}, nil
}
