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
	"arcoris.dev/apimachinery/api/objectstore"
	"arcoris.dev/apimachinery/runtime/objectreconciler"
)

// FromSnapshot resolves request.Key in snapshot and returns its current state.
//
// found is false when the snapshot can serve the key but the object is absent.
// When found is true, Current.Revision is the object's State.Revision, not
// snapshot.Revision.
func FromSnapshot(
	request objectreconciler.Request,
	snapshot objectreconciler.Snapshot,
) (Current, bool, error) {
	if err := request.Validate(); err != nil {
		return Current{}, false, errorWith(ErrInvalidRequest, err)
	}

	result, err := snapshot.View.Get(request.Key)
	if err != nil {
		return Current{}, false, errorWith(ErrInvalidSnapshot, err)
	}
	if !result.Found {
		return Current{}, false, nil
	}

	current, err := FromItem(objectstore.ListItem{
		Key:   result.Key,
		State: result.State,
	})
	if err != nil {
		return Current{}, false, err
	}

	return current, true, nil
}
