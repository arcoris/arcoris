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

// Event is one item in an object watch stream.
//
// Exactly one payload family is meaningful for each Kind. EventChanged carries
// Change, EventProgress carries only Revision, and EventRestartRequired carries
// Restart plus an optional Revision progress boundary.
type Event struct {
	// Kind identifies the event semantics and therefore which payload fields
	// are valid.
	Kind EventKind
	// Revision is the event progress boundary. For EventChanged it must equal
	// Change.Revision. For EventProgress it is the progress boundary. For
	// EventRestartRequired zero means the source cannot name a reliable
	// boundary, while non-zero names the last boundary it can report.
	Revision objectstore.Revision
	// Change is populated only for EventChanged and is validated as a committed
	// objectstore.Change.
	Change objectstore.Change
	// Restart is populated only for EventRestartRequired.
	Restart RestartReason
}

// IsZero reports whether e has no event fields set.
func (e Event) IsZero() bool {
	return e.Kind == 0 &&
		e.Revision.IsZero() &&
		e.Change.IsZero() &&
		e.Restart == 0
}

// IsValid reports whether e passes event validation.
func (e Event) IsValid() bool {
	return e.Validate() == nil
}

// Clone returns a detached copy of e.
func (e Event) Clone() Event {
	return Event{
		Kind:     e.Kind,
		Revision: e.Revision,
		Change:   e.Change.Clone(),
		Restart:  e.Restart,
	}
}
