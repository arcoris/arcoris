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

// Package healthgrpc adapts package health reports to the standard gRPC health
// checking service.
//
// # Package scope
//
// healthgrpc implements grpc.health.v1.Health using the standard protobuf and
// service definitions from google.golang.org/grpc/health/grpc_health_v1. It maps
// gRPC service names to health.Target values, evaluates those targets through a
// Source, and converts health.Status plus health.TargetPolicy into gRPC serving
// statuses.
//
// The package is a transport adapter only. It does not define health checks, own
// a registry, store independent mutable serving state as the source of truth,
// manage lifecycle transitions, perform service discovery, emit metrics, log
// diagnostics, attach tracing spans, or make restart, admission, routing,
// scheduling, or overload-control decisions.
//
// # Relationship to health
//
// Package health owns the transport-neutral health model: Target, Status,
// TargetPolicy, Report, and Registry. healthgrpc depends only on a Source
// interface, which is satisfied by *eval.Evaluator from package health/eval.
// The Source remains the owner of synchronous evaluation behavior, check
// execution policy, panic normalization, status aggregation, and report
// construction.
//
// # Service names
//
// gRPC service names are string transport identities. They are intentionally not
// Kubernetes GroupVersionKind, GroupVersionResource, or any other API machinery
// type. API-aware layers may pass schema identifiers as strings when that is
// their transport contract, but this package must not import API machinery or
// reinterpret service names as resource identifiers.
//
// The empty service name is the standard whole-server gRPC health service. By
// default it maps to health.TargetReady with health.ReadyPolicy. Optional target
// service mappings expose "startup", "live", and "ready" as transport names for
// the built-in concrete targets.
//
// # Safe exposure model
//
// Public gRPC responses expose only grpc.health.v1 serving status values. They
// never expose health.Result.Cause, panic stacks, raw errors, context causes,
// credentials, connection strings, tenant identifiers, internal network
// addresses, or source error text. Evaluation failures are converted to generic
// gRPC errors for Check and to UNKNOWN status values for List and Watch.
//
// # Watch model
//
// Watch is implemented as a per-stream polling adapter over Source. Each stream
// owns its ticker and sends an initial status followed by changes only. The
// package does not create a global background runner, does not retain stream
// references after Watch returns, and does not depend on arcoris.dev/health/probe.
//
// # Concurrency
//
// Server configuration is immutable after NewServer returns. Check, List, Watch,
// Services, HasService, and Target may be called concurrently as long as the
// supplied Source is safe for the same evaluation pattern. healthgrpc does not
// add mutable serving overrides or cache state, so it does not need package-wide
// locks on request paths.
//
// # File ownership
//
//   - source.go owns the Source boundary;
//   - source_nil.go owns typed-nil Source detection;
//   - server.go owns the Server type and immutable fields;
//   - server_constructor.go owns Server construction and service indexing;
//   - server_services.go owns read-only service inspection;
//   - register.go owns standard gRPC registration;
//   - service.go owns the service mapping model;
//   - service_default.go owns built-in service mappings;
//   - service_validate.go owns service mapping validation and normalization;
//   - service_error.go owns service mapping errors;
//   - status.go owns health-to-gRPC status conversion;
//   - rpc_error.go owns stable public gRPC error messages;
//   - option.go owns the Option contract and application order;
//   - config.go owns normalized adapter configuration;
//   - option_clock.go owns clock configuration;
//   - option_services.go owns service mapping options;
//   - option_watch.go owns Watch cadence options;
//   - option_list.go owns List guardrail options;
//   - option_error.go owns option and construction errors;
//   - check.go owns the Health.Check method;
//   - list.go owns the Health.List method;
//   - list_evaluate.go owns per-target List evaluation normalization;
//   - watch.go owns the Health.Watch method;
//   - watch_status.go owns Watch status evaluation and stream normalization.
//
// # Dependency policy
//
// Production code depends on the Go standard library, arcoris.dev/health,
// arcoris.dev/chrono/clock, and the standard gRPC health packages. It must not
// import arcoris.dev/apimachinery, arcoris.dev/health/probe, custom protobufs,
// logging, metrics, or tracing dependencies.
package healthgrpc
