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

// Change describes one committed object store transition.
//
// Change is a value-level API machinery contract. It is not a watch stream,
// runtime event bus, API server response, admission hook, authorization hook,
// audit log, cache record, or background-processing primitive. Future
// watch/cache/audit/controller layers may consume Change values, but this
// package does not implement those layers.
//
// Change only proves that the value is structurally shaped like a committed
// transition. It does not prove the transition occurred in a particular Store
// instance.
type Change struct {
	// Kind identifies which committed transition shape this value represents.
	Kind ChangeKind

	// Key is the object store key affected by the transition.
	Key Key

	// Revision is the store-local commit revision for this transition.
	Revision Revision

	// Before is the committed live state before an update or deletion.
	//
	// Created changes require Before to be zero because no live state existed
	// in the represented transition.
	Before State

	// After is the committed live state after a create or update.
	//
	// Deleted changes require After to be zero because tombstones are modeled by
	// the transition revision, not by a synthetic committed State.
	After State
}

// IsZero reports whether c has no transition kind, key, revision, or states.
func (c Change) IsZero() bool {
	return c.Kind == 0 &&
		c.Key.Equal(Key{}) &&
		c.Revision.IsZero() &&
		isZeroState(c.Before) &&
		isZeroState(c.After)
}

// Clone returns a detached copy of c.
func (c Change) Clone() Change {
	return Change{
		Kind:     c.Kind,
		Key:      c.Key,
		Revision: c.Revision,
		Before:   c.Before.Clone(),
		After:    c.After.Clone(),
	}
}

// IsValid reports whether c is structurally shaped like a committed transition.
func (c Change) IsValid() bool {
	return ValidateChange(c) == nil
}

// Validate checks that c is structurally shaped like its committed transition.
func (c Change) Validate() error {
	return ValidateChange(c)
}
