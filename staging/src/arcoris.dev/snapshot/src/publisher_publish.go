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

// Publish publishes next and returns the resulting lightweight snapshot.
//
// Publish does not clone next. Callers must not mutate next after publication.
func (p *Publisher[T]) Publish(next T) Snapshot[T] {
	return p.PublishStamped(next).Snapshot()
}

// PublishStamped publishes next and returns the resulting stamped snapshot.
//
// PublishStamped assigns a fresh source-local revision, records the local
// publication time using the configured PassiveClock, stores a new immutable
// record, and returns the published stamped snapshot. Concurrent publish calls
// are serialized so readers cannot observe revision rollback. If revision
// overflow is detected, PublishStamped panics before storing a new record.
func (p *Publisher[T]) PublishStamped(next T) Stamped[T] {
	p.mu.Lock()
	defer p.mu.Unlock()

	rev := p.nextRevision.Next()
	updated := p.passiveClock().Now()
	rec := &record[T]{
		revision: rev,
		updated:  updated,
		value:    next,
	}

	p.nextRevision = rev
	p.ptr.Store(rec)

	return Stamped[T]{
		Revision: rev,
		Updated:  updated,
		Value:    next,
	}
}
