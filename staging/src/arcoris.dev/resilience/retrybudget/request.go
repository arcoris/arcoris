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

package retrybudget

// Request is the typed admission request consumed by retry-budget adapters.
//
// TryAdmitRetry remains the direct retry-budget API. TryAdmit(Request) is the
// admission-compatible API for generic callers that work through
// admission.Admitter. Request exists so retry-budget admission does not use
// untyped, unrelated, or cross-domain request values in generic contracts.
//
// Request is intentionally empty because one admission attempt always represents
// one retry attempt. It carries no amount, context, priority, tenant, operation
// class, deadline, metadata, queueing policy, or retry classification. Those
// dimensions belong to higher-level policy and admission orchestration, not this
// spend-only retry-budget primitive.
type Request struct{}
