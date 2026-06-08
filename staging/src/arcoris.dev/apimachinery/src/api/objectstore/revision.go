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

// Revision is a store-local monotonic commit revision.
//
// Revisions are assigned by a concrete store when Create, Update, or Delete
// commits. They are not API resource versions, object generations, wall-clock
// timestamps, distributed sequence numbers, or globally comparable values.
type Revision uint64

// IsZero reports whether r is the invalid/unset revision.
func (r Revision) IsZero() bool {
	return r == 0
}

// IsValid reports whether r can identify committed store state.
func (r Revision) IsValid() bool {
	return !r.IsZero()
}
