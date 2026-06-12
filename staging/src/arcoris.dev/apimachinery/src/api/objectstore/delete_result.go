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

// DeleteResult is the successful result of a committed delete.
type DeleteResult struct {
	// Deleted is the detached live state that was tombstoned. Its Revision is
	// the previous live revision that matched the caller's expected revision.
	Deleted State
	// Revision is the store-local tombstone commit revision.
	Revision Revision
}

// IsZero reports whether r contains no deleted state or delete revision.
func (r DeleteResult) IsZero() bool {
	return r.Deleted.Revision.IsZero() && r.Revision.IsZero()
}
