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

// Package health provides transport-neutral health contracts for ARCORIS
// component internals.
//
// # Package scope
//
// The root package owns the shared health language and in-process contracts:
//
//   - Status, Reason, Target, Result, and Report values;
//   - Checker and CheckFunc contracts;
//   - Registry storage for target-scoped checks;
//   - Gate for owner-published mutable health state;
//   - TargetPolicy helpers for interpreting target status;
//   - shared identifier validation for stable health names and reasons.
//
// It does not publish metrics, expose transports, run periodic probes, perform
// service discovery, drive lifecycle transitions, or make restart, admission,
// routing, scheduling, logging, tracing, or alerting decisions.
//
// # Optional modules
//
// Optional behavior is provided by separate modules:
//
//   - arcoris.dev/healtheval synchronously evaluates registered checks;
//   - arcoris.dev/healthprobe periodically evaluates and caches observations;
//   - arcoris.dev/healthhttp adapts reports and evaluators to HTTP endpoints;
//   - arcoris.dev/healthgrpc adapts reports to the standard gRPC health service;
//   - arcoris.dev/healthtest provides health-domain test helpers.
//
// Root package health must remain independent from those modules. Adapters and
// runtime helpers depend inward on health, while health does not import
// evaluation, HTTP, gRPC, probe, test-helper, runtime lifecycle, metrics,
// logging, or tracing packages.
//
// # File ownership
//
//   - check.go owns Checker and check name validation;
//   - check_func.go owns function-backed checker adapters;
//   - identifier.go owns shared lower_snake_case identifier syntax;
//   - status.go owns Status values and status ordering;
//   - reason.go owns Reason values and reason classification;
//   - target.go owns Target values and target enumeration;
//   - target_error.go owns target validation errors;
//   - result.go owns single-check Result values;
//   - report.go owns target-level Report values;
//   - policy.go owns target status policy;
//   - registry.go owns target-scoped check registration;
//   - registry_error.go owns registry error sentinels and typed errors;
//   - registry_validate.go owns registration batch validation;
//   - gate.go owns mutable owner-published checker state;
//   - shutdown.go owns shutdown and drain check adapters.
package health
