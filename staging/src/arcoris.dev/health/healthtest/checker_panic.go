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

	"arcoris.dev/health"
)

// PanicChecker is a health.Checker that panics from Check.
//
// Evaluator recovers checker panics and stores panic details in Result.Cause.
// Adapter and integration tests can use PanicChecker to verify that recovery
// behavior without defining local panic-only checker types.
type PanicChecker struct {
	// NameValue is returned by Name.
	NameValue string

	// Value is the panic value raised by Check.
	//
	// Value is intentionally any so tests can cover string, error, and structured
	// panic values exactly as production code might accidentally raise them.
	Value any
}

// NewPanicChecker returns a checker that panics with value.
//
// The name is not validated here; registry or evaluator setup remains the owner
// of checker-name validation in tests that need that boundary.
func NewPanicChecker(name string, val any) PanicChecker {
	return PanicChecker{NameValue: name, Value: val}
}

// Name returns the configured checker name.
func (c PanicChecker) Name() string {
	return c.NameValue
}

// Check panics with the configured value.
func (c PanicChecker) Check(context.Context) health.Result {
	panic(c.Value)
}
