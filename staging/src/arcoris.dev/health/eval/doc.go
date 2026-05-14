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

// Package eval provides synchronous health evaluation over package health
// contracts.
//
// The package owns evaluator execution mechanics: reading registered checks,
// executing them with cooperative contexts, normalizing panics, cancellations,
// timeouts, and invalid results, applying sequential or bounded-parallel
// execution policy, and returning a health.Report.
//
// The package does not define health statuses, targets, reasons, results,
// reports, registries, gates, transport adapters, periodic probes, metrics,
// logging, tracing, lifecycle transitions, restart policy, routing policy,
// admission policy, or scheduling decisions.
package eval
