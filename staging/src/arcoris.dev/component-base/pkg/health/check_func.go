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
)

// ErrNilCheckFunc identifies a nil function passed to a health checker adapter.
//
// A function-backed checker cannot produce a health observation without an
// implementation function. Constructors return ErrNilCheckFunc instead of
// creating a checker that panics later during evaluation.
var ErrNilCheckFunc = errors.New("health: nil check function")

// CheckFunc evaluates a health condition and returns a Result.
//
// CheckFunc is the function form of Checker.Check. It is useful when a package
// needs a lightweight checker without defining a dedicated type. The same
// execution contract as Checker.Check applies: implementations SHOULD observe
// ctx when evaluation can block and SHOULD be safe to call repeatedly.
type CheckFunc func(ctx context.Context) Result

// ErrorCheckFunc evaluates a health condition and returns an error.
//
// ErrorCheckFunc is a compatibility convenience for simple checks that only need
// binary healthy/unhealthy semantics. A nil error maps to StatusHealthy. A
// non-nil error maps to StatusUnhealthy with ReasonFatal and the original error
// preserved as Result.Cause.
//
// More expressive checks SHOULD use CheckFunc so they can return StatusStarting,
// StatusDegraded, StatusUnknown, and precise Reason values.
type ErrorCheckFunc func(ctx context.Context) error

// NewCheck returns a named Checker backed by fn.
//
// NewCheck validates name and fn at construction time. The returned checker
// preserves Result values produced by fn, except that an empty result name is
// normalized to the checker name. Full defensive normalization, observation
// timestamps, duration measurement, timeout handling, and panic recovery belong
// to Evaluator.
func NewCheck(name string, fn CheckFunc) (Checker, error) {
	if err := ValidateCheckName(name); err != nil {
		return nil, err
	}
	if fn == nil {
		return nil, ErrNilCheckFunc
	}

	return checkFunc{
		name: name,
		fn:   fn,
	}, nil
}

// MustCheck returns a named Checker backed by fn and panics if construction
// fails.
//
// MustCheck is intended for package-level declarations and tests where invalid
// checker definitions are programmer errors. Runtime configuration paths SHOULD
// use NewCheck and return the error to the owner.
func MustCheck(name string, fn CheckFunc) Checker {
	checker, err := NewCheck(name, fn)
	if err != nil {
		panic(err)
	}

	return checker
}

// NewErrorCheck returns a named Checker backed by an error-returning function.
//
// NewErrorCheck is intentionally conservative. A returned error means the check
// is unhealthy and classified with ReasonFatal. Use NewCheck when the caller
// needs to distinguish timeout, cancellation, degradation, startup, draining,
// overload, dependency outage, or other domain-specific reasons.
func NewErrorCheck(name string, fn ErrorCheckFunc) (Checker, error) {
	if fn == nil {
		return nil, ErrNilCheckFunc
	}

	return NewCheck(name, func(ctx context.Context) Result {
		if err := fn(ctx); err != nil {
			return Unhealthy(
				name,
				ReasonFatal,
				"health check failed",
			).WithCause(err)
		}

		return Healthy(name)
	})
}

// MustErrorCheck returns a named Checker backed by an error-returning function
// and panics if construction fails.
//
// MustErrorCheck is intended for package-level declarations and tests where
// invalid checker definitions are programmer errors. Runtime configuration paths
// SHOULD use NewErrorCheck and return the error to the owner.
func MustErrorCheck(name string, fn ErrorCheckFunc) Checker {
	checker, err := NewErrorCheck(name, fn)
	if err != nil {
		panic(err)
	}

	return checker
}

// checkFunc adapts a function into Checker.
//
// checkFunc is intentionally private so the package can preserve construction
// invariants. Callers should create function-backed checkers through NewCheck,
// MustCheck, NewErrorCheck, or MustErrorCheck.
type checkFunc struct {
	name string
	fn   CheckFunc
}

// Name returns the stable check name.
func (c checkFunc) Name() string {
	return c.name
}

// Check evaluates the underlying function and returns its result.
//
// Check fills an empty result name with the checker name. It does not set
// observation time, measure duration, recover panics, apply timeout, or validate
// status. Those responsibilities belong to Evaluator so all checker
// implementations receive consistent boundary behavior.
func (c checkFunc) Check(ctx context.Context) Result {
	result := c.fn(ctx)
	if result.Name == "" {
		result.Name = c.name
	}

	return result
}
