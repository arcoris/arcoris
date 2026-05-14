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

// Package healthtest provides reusable health-domain test fixtures for packages
// built on top of package health.
//
// # Package scope
//
// healthtest owns deterministic test helpers that speak the package-health
// model directly: simple checkers, controlled checkers, panic checkers, report
// fixtures, registry builders, evaluator builders, fake report sources, and
// health-specific assertions.
//
// The package is intended for adapter, integration, and external health tests.
// It is not production code and it does not define new runtime health behavior.
// Helpers intentionally return package health values and interfaces rather than
// introducing another health model.
//
// # Non-goals
//
// healthtest is not a generic testutil package. It does not provide HTTP, gRPC,
// Kubernetes, metrics, logging, tracing, lifecycle, signal, generic clock,
// generic wait, generic equality, generic panic, or generic error assertions.
// Transport-specific helpers such as HTTP mux recorders and gRPC fake streams
// belong to their transport packages.
//
// # Relationship to health
//
// Package health remains the owner of health contracts, validation, evaluator
// behavior, status aggregation, and target policy semantics. healthtest imports
// health and builds test fixtures around its public API.
//
// Internal tests inside package health that need unexported implementation
// details should keep package-local helpers. Importing healthtest from those
// package-internal tests would make health import a package that imports health,
// which risks an import cycle and blurs the ownership boundary.
//
// # File ownership
//
//   - checker.go owns simple function-backed and static checkers;
//   - checker_controlled.go owns blocking and sequence checkers;
//   - checker_panic.go owns panic-producing checkers;
//   - source.go owns health report sources compatible with adapter tests;
//   - registry.go owns target-group registry builders;
//   - evaluator.go owns evaluator builders with test-safe defaults;
//   - result.go owns canonical result fixtures;
//   - report.go owns canonical report fixtures;
//   - assert.go owns health-domain assertions.
package healthtest
