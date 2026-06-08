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

package lifecycle

import (
	"math"
	"strconv"
)

const (
	// ZeroRevision is the initial lifecycle revision before any transition commits.
	ZeroRevision Revision = 0

	// errRevisionOverflow is the stable panic value for impossible revision wrap.
	errRevisionOverflow = "lifecycle: revision overflow"
)

// Revision is the monotonic commit sequence number for a lifecycle controller.
//
// Revision zero means no transition has been committed. Committed transitions
// use one-based revisions. Revision is local to one Controller instance and is
// not a distributed version, storage revision, wall-clock timestamp, or snapshot
// package revision.
type Revision uint64

// IsZero reports whether r is the initial no-transition revision.
func (r Revision) IsZero() bool {
	return r == ZeroRevision
}

// Next returns the next monotonic lifecycle revision.
//
// Next panics on overflow instead of wrapping to ZeroRevision because revision
// wrap would make old snapshots indistinguishable from a fresh controller.
func (r Revision) Next() Revision {
	if r == Revision(math.MaxUint64) {
		panic(errRevisionOverflow)
	}

	return r + 1
}

// String returns the base-10 diagnostic form of r.
func (r Revision) String() string {
	return strconv.FormatUint(uint64(r), 10)
}
