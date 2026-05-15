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

// Package healthhttp adapts package health reports to safe HTTP health endpoints.
//
// # Package scope
//
// healthhttp is a transport adapter over package health. It turns evaluator
// reports into HTTP handlers, endpoint paths, status-code mappings, and safe
// text or JSON responses.
//
// The package intentionally does not define health checks, own a registry,
// execute checks on a schedule, choose evaluator timeouts, manage lifecycle
// state, handle process signals, run background loops, expose metrics, log
// diagnostics, attach tracing spans, or make admission, routing, scheduling, or
// restart decisions.
//
// # Safe exposure model
//
// Public responses are safe by construction. healthhttp never exposes
// health.Result.Cause, panic stacks, raw errors, context causes, credentials,
// connection strings, tenant identifiers, or internal network addresses. JSON
// responses use dedicated DTOs instead of embedding package health values
// directly so future additions to health.Report or health.Result cannot leak
// into the adapter surface accidentally.
//
// # Default endpoints
//
// InstallDefaults registers only target-specific endpoints:
//
//   - /startupz
//   - /livez
//   - /readyz
//
// DefaultHealthPath is provided for callers that want to install a general
// health endpoint explicitly. It is not registered by default because it does
// not carry stable target semantics across systems.
//
// # Non-goals
//
// healthhttp does not own health contracts, registry mutation, evaluator
// execution policy, middleware stacks, authentication, authorization, or
// router-specific pattern semantics. It remains a small adapter layer over the
// health core package.
package healthhttp
