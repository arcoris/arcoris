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

package snapshot

// Snapshot returns the latest lightweight snapshot.
//
// Snapshot is lock-free. If no value has been published, Snapshot returns the
// zero Snapshot[T]. The returned value is not cloned; it is the immutable value
// stored in the latest published record.
func (p *Publisher[T]) Snapshot() Snapshot[T] {
	stamped := p.Stamped()
	return stamped.Snapshot()
}

// Stamped returns the latest stamped snapshot.
//
// Stamped is lock-free. If no value has been published, Stamped returns the zero
// Stamped[T]. The returned value is not cloned.
func (p *Publisher[T]) Stamped() Stamped[T] {
	rec := p.ptr.Load()
	if rec == nil {
		return Stamped[T]{}
	}

	return Stamped[T]{
		Revision: rec.revision,
		Updated:  rec.updated,
		Value:    rec.value,
	}
}

// Revision returns the revision of the latest visible published record.
//
// Revision returns ZeroRevision before the first Publish.
func (p *Publisher[T]) Revision() Revision {
	rec := p.ptr.Load()
	if rec == nil {
		return ZeroRevision
	}

	return rec.revision
}
