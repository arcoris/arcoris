/*
  Copyright 2026 The ARCORIS Authors

  Licensed under the Apache License, Version 2.0 (the "License");
  you may not use this file except in compliance with the License.
  You may obtain a copy of the License at

      http://www.apache.org/licenses/LICENSE-2.0

  Unless required by applicable law or agreed to in writing, software
  distributed under the License is distributed on an "AS IS" BASIS,
  WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
  See the License for the specific language governing permissions and
  limitations under the License.
*/

package snapshot

// Revision is a local logical version of one snapshot source.
//
// A Revision identifies a committed value published by one Store, Publisher, or
// domain-specific source. Revisions are monotonic only within that source. Code
// must not compare revisions from unrelated sources as if they formed a global
// ordering.
//
// ZeroRevision is reserved for the absence of a committed publication. Store
// starts at revision 1 because it has an initial value immediately after
// construction. A zero-value Publisher reports ZeroRevision before its first
// Publish.
type Revision uint64

// ZeroRevision is the zero revision value.
//
// ZeroRevision means that no committed value has been observed or published for a
// source. It is not a valid committed Store revision because Store always has an
// initial value.
const ZeroRevision Revision = 0

// IsZero reports whether r is ZeroRevision.
func (r Revision) IsZero() bool {
	return r == ZeroRevision
}

// Next returns the next revision after r.
//
// Next panics on uint64 overflow. Overflow would break the package's monotonic
// revision contract and is treated as a programmer error rather than silently
// wrapping to ZeroRevision.
func (r Revision) Next() Revision {
	if r == ^Revision(0) {
		panic("snapshot: revision overflow")
	}

	return r + 1
}

// ChangedSince reports whether r differs from prev.
//
// The method is intentionally equality-based. It does not imply that revisions
// from different sources can be ordered globally.
func (r Revision) ChangedSince(prev Revision) bool {
	return r != prev
}
