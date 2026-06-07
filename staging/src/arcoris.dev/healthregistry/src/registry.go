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

// Registry is an immutable in-process health check resolver.
//
// Registry implements health.CheckResolver. It stores validated target-scoped
// check sets built by Builder. Registry is safe for concurrent reads and does
// not require a mutex because its state is immutable after construction.
type Registry struct {
	// sets contains immutable target-owned check sets.
	sets map[health.Target]health.CheckSet

	// targets stores concrete targets with at least one registered check in
	// health.ConcreteTargets order.
	targets []health.Target
}

var _ health.CheckResolver = (*Registry)(nil)

// newRegistry indexes immutable check sets for read-only resolution.
func newRegistry(sets map[health.Target]health.CheckSet) *Registry {
	targets := make([]health.Target, 0, len(sets))
	for _, target := range health.ConcreteTargets() {
		set, ok := sets[target]
		if !ok || set.Empty() {
			continue
		}
		targets = append(targets, target)
	}

	return &Registry{
		sets:    sets,
		targets: targets,
	}
}

// ResolveChecks returns checks registered for target.
//
// target must be concrete. A concrete target with no registered checks resolves
// to an empty CheckSet.
func (r *Registry) ResolveChecks(target health.Target) (health.CheckSet, error) {
	if !target.IsConcrete() {
		return health.CheckSet{}, health.InvalidTargetError{Target: target}
	}
	if r == nil {
		return health.NewCheckSet(target)
	}

	set, ok := r.sets[target]
	if !ok {
		return health.NewCheckSet(target)
	}

	return set, nil
}

// Targets returns concrete targets that have at least one registered check.
func (r *Registry) Targets() []health.Target {
	if r == nil || len(r.targets) == 0 {
		return nil
	}

	targets := make([]health.Target, len(r.targets))
	copy(targets, r.targets)

	return targets
}

// Len returns the number of checks registered for target.
func (r *Registry) Len(target health.Target) int {
	if r == nil || !target.IsConcrete() {
		return 0
	}

	return r.sets[target].Len()
}

// Empty reports whether r has no registered checks.
func (r *Registry) Empty() bool {
	return r == nil || len(r.targets) == 0
}

// Has reports whether name is registered under target.
func (r *Registry) Has(target health.Target, name string) bool {
	if r == nil || !target.IsConcrete() {
		return false
	}

	return r.sets[target].Has(name)
}
