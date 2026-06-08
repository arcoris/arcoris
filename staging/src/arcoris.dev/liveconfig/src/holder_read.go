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

package liveconfig

import "arcoris.dev/snapshot"

// Snapshot returns the current lightweight live configuration snapshot.
//
// Snapshot delegates to the internal snapshot.Publisher. It does not take the
// Holder write mutex, does not clone the value, and does not update LastError.
// The returned Value must be treated as immutable.
func (h *Holder[T]) Snapshot() snapshot.Snapshot[T] {
	requireHolder(h)
	return h.pub.Snapshot()
}

// Stamped returns the current stamped live configuration snapshot.
//
// Stamped includes the local publication time assigned when the current value
// was accepted. It has the same immutability and read-side behavior as Snapshot.
func (h *Holder[T]) Stamped() snapshot.Stamped[T] {
	requireHolder(h)
	return h.pub.Stamped()
}

// Revision returns the current source-local configuration revision.
//
// Revision is a cheap read-side change check for consumers that do not need the
// value itself. The revision is local to this holder.
func (h *Holder[T]) Revision() snapshot.Revision {
	requireHolder(h)
	return h.pub.Revision()
}

// LastError returns the most recent rejected Apply error.
//
// LastError is diagnostic state for reload loops and operators. It is set only
// when Apply returns a normalization or validation error. Clone, normalization,
// validation, or equality panics propagate to the caller and do not replace the
// previous LastError. Any successful Apply clears LastError, including an
// EqualFunc no-op, because the most recent candidate was accepted.
func (h *Holder[T]) LastError() error {
	requireHolder(h)

	h.mu.Lock()
	defer h.mu.Unlock()

	return h.lastErr
}
