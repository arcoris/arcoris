/*
  Copyright 2026 The ARCORIS Authors

  Licensed under the Apache License, Version 2.0 (the "License");
  you may not use this file except in compliance with the License.
  You may obtain a copy of the License at

      http://www.apache.org/licenses/LICENSE-2.0

  Unless required by applicable law or agreed to in writing, software
  distributed under the License is distributed on an "AS IS" BASIS,
  WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
  See the License for the specific language governing permissions and
  limitations under the License.
*/

package health

import (
	"errors"
	"sync"
)

// Registry owns the set of checks registered for health targets.
//
// Registry is a transport-neutral setup primitive. It does not execute checks,
// aggregate reports, apply target policy, expose endpoints, collect metrics, log
// diagnostics, or manage periodic probes. Those responsibilities belong to
// evaluator, report, adapter, observability, and probe-runner layers.
//
// Checks are stored in registration order. Deterministic ordering is part of the
// registry contract because reports, tests, and diagnostic output should not
// depend on Go map iteration order.
//
// The zero value is usable and represents an empty registry. NewRegistry is
// still preferred in component setup code because it makes registry ownership
// explicit.
//
// Registry is safe for concurrent use. In normal component construction, callers
// SHOULD register checks during setup and treat the registry as effectively
// immutable after startup. Dynamic registration is supported by synchronization,
// but owners that mutate a live registry must define their own lifecycle and
// diagnostic expectations for those changes.
type Registry struct {
	mu sync.RWMutex

	// checks stores checkers in registration order for each concrete target.
	checks map[Target][]Checker

	// names stores check names for duplicate detection within each target.
	names map[Target]map[string]struct{}
}

// NewRegistry returns an empty registry.
//
// NewRegistry has the same behavior as a zero-value Registry, but explicit
// construction makes registry ownership visible in component setup.
func NewRegistry() *Registry {
	return &Registry{}
}

// Register adds checks to target in registration order.
//
// target MUST be concrete. TargetUnknown and invalid target values are rejected
// with an error classified as ErrInvalidTarget.
//
// Each checker MUST be non-nil and MUST expose a valid stable name. Duplicate
// names are rejected within the same target and reported as ErrDuplicateCheck.
// The same name MAY appear under different targets.
//
// Register is atomic for a single call: if any supplied checker is invalid or
// duplicated, none of the supplied checkers are added to the registry. Batch
// validation and existing-registry conflicts are aggregated with errors.Join so
// callers can classify every child error with errors.Is and errors.As.
func (r *Registry) Register(target Target, checks ...Checker) error {
	if !target.IsConcrete() {
		return InvalidTargetError{Target: target}
	}

	prepared, err := prepareChecks(target, checks)
	if err != nil {
		return err
	}

	r.mu.Lock()
	defer r.mu.Unlock()

	r.initLocked()

	existingNames := r.names[target]
	if existingNames == nil {
		existingNames = make(map[string]struct{}, len(prepared))
	}

	var conflicts []error
	for _, check := range prepared {
		if _, exists := existingNames[check.Name]; exists {
			conflicts = append(conflicts, DuplicateCheckError{
				Target:        target,
				Name:          check.Name,
				Index:         check.Index,
				PreviousIndex: -1,
			})
		}
	}
	if len(conflicts) > 0 {
		return errors.Join(conflicts...)
	}

	for _, check := range prepared {
		r.checks[target] = append(r.checks[target], check.Checker)
		existingNames[check.Name] = struct{}{}
	}

	r.names[target] = existingNames

	return nil
}

// Checks returns checks registered for target in registration order.
//
// The returned slice is a defensive copy. Callers may sort, append to, truncate,
// or otherwise modify it without mutating the registry.
//
// Checks returns nil for non-concrete targets and for concrete targets with no
// registered checks.
func (r *Registry) Checks(target Target) []Checker {
	if !target.IsConcrete() {
		return nil
	}

	r.mu.RLock()
	defer r.mu.RUnlock()

	checks := r.checks[target]
	if len(checks) == 0 {
		return nil
	}

	out := make([]Checker, len(checks))
	copy(out, checks)

	return out
}

// Has reports whether name is registered under target.
//
// Has is a read-side convenience helper. It returns false for non-concrete
// targets and invalid check names instead of returning validation errors.
func (r *Registry) Has(target Target, name string) bool {
	if !target.IsConcrete() || !ValidCheckName(name) {
		return false
	}

	r.mu.RLock()
	defer r.mu.RUnlock()

	names := r.names[target]
	if len(names) == 0 {
		return false
	}

	_, exists := names[name]
	return exists
}

// Len returns the number of checks registered under target.
//
// Len returns zero for non-concrete targets.
func (r *Registry) Len(target Target) int {
	if !target.IsConcrete() {
		return 0
	}

	r.mu.RLock()
	defer r.mu.RUnlock()

	return len(r.checks[target])
}

// Empty reports whether no concrete target has registered checks.
func (r *Registry) Empty() bool {
	r.mu.RLock()
	defer r.mu.RUnlock()

	for _, target := range ConcreteTargets() {
		if len(r.checks[target]) > 0 {
			return false
		}
	}

	return true
}

// Targets returns concrete targets that currently have at least one registered
// check.
//
// Targets are returned in ConcreteTargets order. TargetUnknown is never returned
// because it is a sentinel, not an evaluable health target.
func (r *Registry) Targets() []Target {
	r.mu.RLock()
	defer r.mu.RUnlock()

	targets := make([]Target, 0, len(ConcreteTargets()))
	for _, target := range ConcreteTargets() {
		if len(r.checks[target]) > 0 {
			targets = append(targets, target)
		}
	}

	return targets
}

// initLocked initializes internal maps.
//
// The caller MUST hold r.mu for writing.
func (r *Registry) initLocked() {
	if r.checks == nil {
		r.checks = make(map[Target][]Checker, len(ConcreteTargets()))
	}
	if r.names == nil {
		r.names = make(map[Target]map[string]struct{}, len(ConcreteTargets()))
	}
}
