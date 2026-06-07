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

package healthregistry

import "arcoris.dev/health"

// Builder is a mutable setup object for health check registration.
//
// Builder validates target-scoped check ownership before publishing an immutable
// Registry. The zero value is usable. Builder is intended for single-owner
// setup code and is not safe for concurrent mutation.
type Builder struct {
	// checks stores mutable target-scoped checker registration in caller order.
	checks map[health.Target][]health.Checker

	// names indexes registered check names by target for duplicate detection.
	// Values are indexes into checks[target].
	names map[health.Target]map[string]int
}

// NewBuilder returns an empty mutable registry builder.
func NewBuilder() *Builder {
	return &Builder{}
}

// init lazily prepares the builder maps.
//
// The builder zero value is intentionally usable, so mutating operations call
// init before touching storage. Every concrete target receives a name index to
// keep later registration code simple and nil-map free.
func (b *Builder) init() {
	if b.checks == nil {
		b.checks = make(map[health.Target][]health.Checker, len(health.ConcreteTargets()))
	}
	if b.names == nil {
		b.names = make(map[health.Target]map[string]int, len(health.ConcreteTargets()))
	}
	for _, target := range health.ConcreteTargets() {
		if b.names[target] == nil {
			b.names[target] = make(map[string]int)
		}
	}
}
