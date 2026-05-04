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
// The package is intentionally an adapter layer. It does not define health check
// contracts, own a health registry, execute checks directly, run periodic probes,
// publish metrics, log diagnostics, manage lifecycle state, handle operating
// system signals, or make restart, routing, admission, scheduling, or alerting
// decisions.
//
// The owned responsibilities are default health endpoint paths, HTTP method
// policy, handler configuration, HTTP status code mapping, safe text and JSON
// rendering, and mux installation helpers.
//
// Public HTTP responses never expose health.Result.Cause, panic stacks, raw
// errors, context causes, connection strings, credentials, internal addresses,
// tenant identifiers, or other private diagnostic data.
//
// InstallDefaults registers only target-specific endpoints:
//
//   - /startupz;
//   - /livez;
//   - /readyz.
//
// Compatibility paths such as /healthz and /health are provided as constants but
// are not installed by default because they do not have universal target
// semantics across systems.
package healthhttp
