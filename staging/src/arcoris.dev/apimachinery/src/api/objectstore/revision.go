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
// commits. They are monotonic within one store and may have gaps. They are not
// API resource versions, object generations, wall-clock timestamps, distributed
// sequence numbers, or globally comparable values.
//
// Revision(0) is invalid for committed State and Change values. Higher
// list-watch layers may still use zero as an initial collection boundary when
// no committed mutation has been observed yet.
type Revision uint64

// IsZero reports whether r is the zero boundary or unset committed revision.
func (r Revision) IsZero() bool {
	return r == 0
}

// IsValid reports whether r can identify committed store state or a committed change.
func (r Revision) IsValid() bool {
	return !r.IsZero()
}
