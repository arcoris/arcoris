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

// Build returns an immutable registry from b.
//
// Build copies registered check slices. Later builder mutations do not affect
// previously built registries. A nil builder builds an empty registry.
func (b *Builder) Build() (*Registry, error) {
	if b == nil || len(b.checks) == 0 {
		return &Registry{}, nil
	}

	sets := make(map[health.Target]health.CheckSet, len(b.checks))
	for _, target := range health.ConcreteTargets() {
		checks := b.checks[target]
		if len(checks) == 0 {
			continue
		}

		set, err := health.NewCheckSet(target, checks...)
		if err != nil {
			return nil, err
		}
		sets[target] = set
	}

	return newRegistry(sets), nil
}
