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

package healthgate

import (
	"context"
	"sync"

	"arcoris.dev/health"
)

// Gate is a concurrency-safe mutable checker backed by the latest result.
//
// Gate is intended for component-owned health state that changes because of
// lifecycle, admission, overload, drain, dependency, or runtime-control events.
// It is useful when the health state is already known by the component owner and
// should be read cheaply by a health evaluator.
//
// Gate does not execute work, poll dependencies, start goroutines, create
// timers, expose endpoints, emit metrics, or decide restart/admission behavior.
// It only stores the latest owner-published health.Result and returns it through
// the health.Checker contract.
//
// The zero Gate is not usable because it has no stable check name. Use New or
// NewUnknown to construct a gate with explicit ownership. A Gate must not be
// copied after first use.
type Gate struct {
	// mu protects the mutable result field.
	mu sync.RWMutex

	// name is the stable health check name owned by the gate. It is immutable
	// after construction, so Name reads it without locking.
	name string

	// result is the latest owner-published health observation.
	result health.Result
}

var _ health.Checker = (*Gate)(nil)

// New returns a Gate with name and initial result.
//
// name must be a valid health check name. initial must be structurally valid. If
// initial.Name is empty, New fills it with name. If initial.Name is non-empty and
// differs from name, construction fails with ErrMismatchedGateResult.
func New(name string, initial health.Result) (*Gate, error) {
	if err := health.ValidateCheckName(name); err != nil {
		return nil, err
	}

	res, err := normalizeGateResult(name, initial)
	if err != nil {
		return nil, err
	}

	return &Gate{
		name:   name,
		result: res,
	}, nil
}

// NewUnknown returns a Gate initialized with StatusUnknown.
//
// This constructor is appropriate when a component wants to publish a gate
// during setup before the first concrete health observation is available.
func NewUnknown(name string) (*Gate, error) {
	return New(
		name,
		health.Unknown(
			name,
			health.ReasonNotObserved,
			"health gate has not observed a state yet",
		),
	)
}

// Name returns the stable check name owned by g.
func (g *Gate) Name() string {
	if g == nil {
		return ""
	}

	return g.name
}

// Check returns the latest result stored in g.
//
// Check ignores ctx because Gate is an in-memory state holder. A nil Gate
// returns an unnamed unknown result with health.ErrNilChecker as Cause. Nil gates
// should normally be rejected by resolver construction, but this defensive
// behavior prevents accidental panics in direct use.
func (g *Gate) Check(ctx context.Context) health.Result {
	return g.Snapshot()
}

// Snapshot returns the latest result stored in g without using the Checker
// context parameter.
//
// Snapshot exists for owner code that wants to inspect the gate directly. It has
// the same nil-gate defensive behavior as Check.
func (g *Gate) Snapshot() health.Result {
	if g == nil {
		return health.Unknown(
			"",
			health.ReasonNotObserved,
			"health gate is nil",
		).WithCause(health.ErrNilChecker)
	}

	g.mu.RLock()
	defer g.mu.RUnlock()

	return g.result
}
