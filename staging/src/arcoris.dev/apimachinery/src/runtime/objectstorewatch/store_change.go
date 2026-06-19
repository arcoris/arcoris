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

import "arcoris.dev/apimachinery/api/objectstore"

// createdChange builds the committed create transition captured by the wrapper.
func createdChange(key objectstore.Key, after objectstore.State) (objectstore.Change, error) {
	return objectstore.NewCreatedChange(key, after)
}

// updatedChange builds the committed update transition captured by the wrapper.
func updatedChange(
	key objectstore.Key,
	before objectstore.State,
	after objectstore.State,
) (objectstore.Change, error) {
	return objectstore.NewUpdatedChange(key, before, after)
}

// deletedChange builds the committed delete transition captured by the wrapper.
//
// The delete revision must come from DeleteResult.Revision, not from the
// deleted live state, because the latter is the previous live revision.
func deletedChange(key objectstore.Key, result objectstore.DeleteResult) (objectstore.Change, error) {
	return objectstore.NewDeletedChange(key, result.Deleted, result.Revision)
}
