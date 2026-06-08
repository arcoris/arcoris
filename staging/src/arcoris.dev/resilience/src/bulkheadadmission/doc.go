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

// Package bulkheadadmission maps bulkhead observations to admission results.
//
// The package is an adapter. It does not own capacity accounting, queues,
// waiters, retry behavior, fairness, metrics, logging, tracing, health checks,
// or worker scheduling. Core accounting and lease ownership remain in package
// bulkhead. This package only gives the direct non-blocking acquisition result
// the generic admission.Result shape and maps refusals to standard admission
// reasons.
//
// Denied admission uses a high-level capacity-exhausted admission reason while
// preserving the precise bulkhead.Observation in metadata. Consumers that need
// to distinguish insufficient availability from debt should inspect that
// metadata instead of parsing the generic admission reason.
package bulkheadadmission
