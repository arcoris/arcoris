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

import (
	"arcoris.dev/apimachinery/api/fieldownership"
	"arcoris.dev/apimachinery/api/objectlifecycle"
	"arcoris.dev/apimachinery/api/value"
)

// UpdateObserved builds a lifecycle request that replaces the current object's
// Observed surface.
//
// Expected is Current.Revision. The observed payload is cloned so later caller
// mutation of composite values cannot change the constructed request.
func (c Current) UpdateObserved(
	observed value.Value,
	owner fieldownership.Owner,
) (objectlifecycle.UpdateObservedRequest, error) {
	if err := c.validate(); err != nil {
		return objectlifecycle.UpdateObservedRequest{}, err
	}

	return objectlifecycle.UpdateObservedRequest{
		Resource: c.key.Resource,
		Object:   c.key.Object,
		Observed: observed.Clone(),
		Owner:    owner,
		Expected: c.Revision(),
	}, nil
}
