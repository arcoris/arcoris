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

// Delete delegates to the backend and publishes the committed delete change.
//
// The backend's DeleteResult supplies both the deleted live state and tombstone
// revision, which are the exact fields required to construct a valid
// objectstore.ChangeDeleted value.
func (s *Store) Delete(
	ctx context.Context,
	key objectstore.Key,
	expected objectstore.Revision,
) (objectstore.DeleteResult, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	result, err := s.backend.Delete(ctx, key, expected)
	if err != nil {
		return objectstore.DeleteResult{}, err
	}

	change := objectstore.Change{
		Kind:     objectstore.ChangeDeleted,
		Key:      key,
		Revision: result.Revision,
		Before:   result.Deleted,
	}
	if err := objectstore.ValidateChange(change); err != nil {
		s.closeAllWithContinuityLoss(err)
		return objectstore.DeleteResult{}, continuityError(err)
	}

	s.publishValidated(change)
	return result, nil
}
