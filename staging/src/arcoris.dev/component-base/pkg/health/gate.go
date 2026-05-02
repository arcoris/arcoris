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
	"context"
	"errors"
	"fmt"
	"sync"
)

var (
	// ErrInvalidGateResult identifies a Result that cannot be stored in a Gate.
	//
	// Gate results must have a valid Status and a non-negative Duration. The gate
	// may fill an empty result name with its own name, but it does not accept
	// structurally invalid health observations.
	ErrInvalidGateResult = errors.New("health: invalid gate result")

	// ErrMismatchedGateResult identifies a Result whose non-empty name does not
	// match the owning Gate name.
	//
	// A Gate is itself a Checker and therefore owns one stable check name. Stored
	// results must either leave Name empty so the gate can fill it, or use the
	// exact gate name.
	ErrMismatchedGateResult = errors.New("health: mismatched gate result")
)

// Gate is a concurrency-safe mutable Checker backed by the latest Result.
//
// Gate is intended for component-owned health state that changes because of
// lifecycle, admission, overload, drain, dependency, or runtime-control events.
// It is the right primitive when the health state is already known by the
// component owner and should be read cheaply by an Evaluator.
//
// Gate does not execute work, poll dependencies, start goroutines, create
// timers, expose endpoints, emit metrics, or decide restart/admission behavior.
// It only stores the latest owner-published Result and returns it through the
// Checker contract.
//
// Gate is useful for:
//
//   - startup gates;
//   - readiness gates;
//   - admission gates;
//   - overload gates;
//   - drain gates;
//   - fatal runtime gates;
//   - dependency state cached by another owner.
//
// The zero Gate is not usable because it has no stable check name. Use NewGate
// or NewUnknownGate to construct a gate with explicit ownership.
//
// A Gate must not be copied after first use. Copying a live gate can split one
// logical health state into independent copies and produce inconsistent reports.
type Gate struct {
	mu sync.RWMutex

	name   string
	result Result
}

// NewGate returns a Gate with name and initial result.
//
// name MUST be a valid check name. initial MUST be structurally valid. If
// initial.Name is empty, NewGate fills it with name. If initial.Name is non-empty
// and differs from name, construction fails with ErrMismatchedGateResult.
//
// NewGate does not set observation time or duration. Time ownership belongs to
// the component publishing the result or to Evaluator when it normalizes a
// checker result.
func NewGate(name string, initial Result) (*Gate, error) {
	if err := ValidateCheckName(name); err != nil {
		return nil, err
	}

	result, err := normalizeGateResult(name, initial)
	if err != nil {
		return nil, err
	}

	return &Gate{
		name:   name,
		result: result,
	}, nil
}

// NewUnknownGate returns a Gate initialized with StatusUnknown.
//
// This constructor is appropriate when a component wants to publish a gate during
// setup before the first concrete health observation is available.
func NewUnknownGate(name string) (*Gate, error) {
	return NewGate(
		name,
		Unknown(
			name,
			ReasonNotObserved,
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

// Check returns the latest Result stored in g.
//
// Check ignores ctx because Gate is an in-memory state holder and does not block,
// perform I/O, wait on another goroutine, or acquire external resources. The ctx
// parameter is still part of the Checker contract so Gate can be used anywhere a
// Checker is expected.
//
// A nil Gate returns an unnamed unknown Result with ErrNilChecker as Cause. Nil
// gates should normally be rejected by Registry, but this defensive behavior
// prevents accidental panics in direct use.
func (g *Gate) Check(ctx context.Context) Result {
	if g == nil {
		return Unknown(
			"",
			ReasonNotObserved,
			"health gate is nil",
		).WithCause(ErrNilChecker)
	}

	g.mu.RLock()
	defer g.mu.RUnlock()

	return g.result
}

// Set replaces the latest Result stored in g.
//
// result must either have an empty Name or the same Name as the gate. Empty names
// are normalized to the gate name. Structurally invalid results are rejected and
// leave the gate unchanged.
//
// Set does not modify Observed, Duration, Reason, Message, or Cause except for
// filling an empty Name. The component owner remains responsible for publishing
// meaningful observation metadata when it has that information.
func (g *Gate) Set(result Result) error {
	if g == nil {
		return ErrNilChecker
	}

	normalized, err := normalizeGateResult(g.name, result)
	if err != nil {
		return err
	}

	g.mu.Lock()
	defer g.mu.Unlock()

	g.result = normalized

	return nil
}

// Unknown stores a StatusUnknown result in g.
//
// Unknown is appropriate when the owner cannot currently produce a reliable
// health observation, for example before initialization, after losing a cached
// dependency observation, or after an inconclusive owner-controlled check.
func (g *Gate) Unknown(reason Reason, message string) error {
	if g == nil {
		return ErrNilChecker
	}

	return g.Set(Unknown(g.name, reason, message))
}

// Starting stores a StatusStarting result in g.
//
// Starting is appropriate while the owner is bootstrapping and making progress
// but has not completed initialization for the relevant health scope.
func (g *Gate) Starting(reason Reason, message string) error {
	if g == nil {
		return ErrNilChecker
	}

	return g.Set(Starting(g.name, reason, message))
}

// Healthy stores a StatusHealthy result in g.
//
// Healthy publishes an affirmative health observation for the gate's scope. It
// does not by itself grant readiness, admission, scheduling, or routing.
func (g *Gate) Healthy() error {
	if g == nil {
		return ErrNilChecker
	}

	return g.Set(Healthy(g.name))
}

// Degraded stores a StatusDegraded result in g.
//
// Degraded is appropriate when the owner still has usable capability but is
// operating with reduced confidence, partial capacity, or active protective
// behavior.
func (g *Gate) Degraded(reason Reason, message string) error {
	if g == nil {
		return ErrNilChecker
	}

	return g.Set(Degraded(g.name, reason, message))
}

// Unhealthy stores a StatusUnhealthy result in g.
//
// Unhealthy publishes a strong negative health observation for the gate's scope.
// It does not by itself prescribe restart, traffic removal, admission closure, or
// scheduler exclusion.
func (g *Gate) Unhealthy(reason Reason, message string) error {
	if g == nil {
		return ErrNilChecker
	}

	return g.Set(Unhealthy(g.name, reason, message))
}

// normalizeGateResult validates result for storage in a gate and fills an empty
// result name with the gate name.
func normalizeGateResult(name string, result Result) (Result, error) {
	if result.Name == "" {
		result.Name = name
	} else if result.Name != name {
		return Result{}, MismatchedGateResultError{
			GateName:   name,
			ResultName: result.Name,
		}
	}

	if !result.IsValid() {
		return Result{}, InvalidGateResultError{
			GateName: name,
			Result:   result,
		}
	}

	return result, nil
}

// InvalidGateResultError describes a structurally invalid Result rejected by a
// Gate.
//
// InvalidGateResultError is classified as ErrInvalidGateResult. Callers should
// use errors.Is for classification and inspect GateName or Result only for
// diagnostics.
type InvalidGateResultError struct {
	GateName string
	Result   Result
}

// Error returns the invalid gate result message.
func (e InvalidGateResultError) Error() string {
	return fmt.Sprintf(
		"%v: gate=%q status=%s duration=%s",
		ErrInvalidGateResult,
		e.GateName,
		e.Result.Status.String(),
		e.Result.Duration,
	)
}

// Is reports whether target matches the invalid gate result classification.
func (e InvalidGateResultError) Is(target error) bool {
	return target == ErrInvalidGateResult
}

// MismatchedGateResultError describes a Result whose name does not match its
// owning Gate.
//
// MismatchedGateResultError is classified as ErrMismatchedGateResult. Callers
// should use errors.Is for classification and inspect GateName or ResultName only
// for diagnostics.
type MismatchedGateResultError struct {
	GateName   string
	ResultName string
}

// Error returns the mismatched gate result message.
func (e MismatchedGateResultError) Error() string {
	return fmt.Sprintf(
		"%v: gate=%q result=%q",
		ErrMismatchedGateResult,
		e.GateName,
		e.ResultName,
	)
}

// Is reports whether target matches the mismatched gate result classification.
func (e MismatchedGateResultError) Is(target error) bool {
	return target == ErrMismatchedGateResult
}
