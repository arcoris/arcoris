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

package objectwatch

import "arcoris.dev/apimachinery/api/objectstore"

// Start describes the requested beginning of a watch stream.
//
// Start is part of the continuity contract. StartAfterRevision requires
// historical delivery after a known committed revision. StartAtCurrent
// deliberately asks for no historical catch-up.
type Start struct {
	// Mode identifies whether the stream catches up from history or starts at
	// the source's current point.
	Mode StartMode
	// Revision is the source-local revision boundary for StartAfterRevision.
	// It must be zero for StartAtCurrent.
	Revision objectstore.Revision
}

// AfterRevision constructs a historical start after revision.
//
// The revision must be non-zero because zero is not a committed progress
// boundary in objectstore.
func AfterRevision(revision objectstore.Revision) (Start, error) {
	start := Start{Mode: StartAfterRevision, Revision: revision}
	if err := start.Validate(); err != nil {
		return Start{}, err
	}

	return start, nil
}

// AtCurrent constructs a start at the source's current progress point.
//
// It requests no historical changes and therefore carries a zero Revision.
func AtCurrent() Start {
	return Start{Mode: StartAtCurrent}
}

// IsZero reports whether s carries no start mode or revision.
func (s Start) IsZero() bool {
	return s.Mode == 0 && s.Revision.IsZero()
}

// IsValid reports whether s passes watch start validation.
func (s Start) IsValid() bool {
	return s.Validate() == nil
}
