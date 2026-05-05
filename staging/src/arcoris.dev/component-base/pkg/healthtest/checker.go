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

package healthtest

import (
	"context"

	"arcoris.dev/component-base/pkg/health"
)

// Checker is a lightweight health.Checker test fixture.
//
// Constructors do not validate names. Some tests need intentionally invalid
// checkers to verify registry and evaluator boundaries.
//
// Checker does not call health.NewCheck because that constructor owns production
// checker validation. healthtest intentionally lets tests choose whether they
// need a valid checker, an invalid checker, a static result, or a custom Check
// function.
type Checker struct {
	// NameValue is returned by Name.
	NameValue string

	// ResultValue is returned by Check when Func is nil.
	//
	// The value is returned as-is so tests can exercise evaluator normalization,
	// result-name mismatches, invalid statuses, or missing observations when
	// those behaviors are under test.
	ResultValue health.Result

	// Func is called by Check when non-nil.
	//
	// Func lets tests observe context, count calls externally, or construct
	// results dynamically without defining another local checker type.
	Func func(context.Context) health.Result
}

// StaticChecker returns a Checker that always returns result.
//
// The result name is not rewritten to match name. That preserves mismatch cases
// for evaluator tests and adapter tests that verify normalization boundaries.
func StaticChecker(name string, result health.Result) Checker {
	return Checker{NameValue: name, ResultValue: result}
}

// FuncChecker returns a Checker backed by fn.
//
// A nil fn is allowed and results in the zero ResultValue path. This keeps the
// type simple and lets registry tests decide which validation boundary they are
// exercising.
func FuncChecker(name string, fn func(context.Context) health.Result) Checker {
	return Checker{NameValue: name, Func: fn}
}

// HealthyChecker returns a Checker that reports healthy.
func HealthyChecker(name string) Checker {
	return StaticChecker(name, HealthyResult(name))
}

// StartingChecker returns a Checker that reports starting.
func StartingChecker(name string) Checker {
	return StaticChecker(name, StartingResult(name))
}

// DegradedChecker returns a Checker that reports degraded.
func DegradedChecker(name string, reason health.Reason) Checker {
	return StaticChecker(name, DegradedResult(name, reason))
}

// UnhealthyChecker returns a Checker that reports unhealthy.
func UnhealthyChecker(name string, reason health.Reason) Checker {
	return StaticChecker(name, UnhealthyResult(name, reason))
}

// UnknownChecker returns a Checker that reports unknown.
func UnknownChecker(name string, reason health.Reason) Checker {
	return StaticChecker(name, UnknownResult(name, reason))
}

// Name returns the configured checker name.
func (c Checker) Name() string {
	return c.NameValue
}

// Check returns the configured result or calls Func.
func (c Checker) Check(ctx context.Context) health.Result {
	if c.Func != nil {
		return c.Func(ctx)
	}

	return c.ResultValue
}
