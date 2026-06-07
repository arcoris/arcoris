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

package eval

import (
	"time"

	"arcoris.dev/chrono/clock"
	"arcoris.dev/health"
)

// Evaluator executes resolved health checks and returns target reports.
//
// Evaluator is transport-neutral. It does not expose HTTP handlers, map gRPC
// serving states, log diagnostics, emit metrics, perform retries, run periodic
// probes, or decide restart, admission, routing, or scheduling behavior. It only
// owns the synchronous evaluation boundary for checks returned by a
// health.CheckResolver.
//
// Evaluation is deterministic with respect to resolver order. The default
// execution policy is sequential. Component owners may configure bounded
// parallel execution globally or per target. Parallel execution preserves
// health.Report.Checks order by resolver order even when checks finish out
// of order.
//
// Evaluator applies a cooperative context to every check and, when a timeout is
// configured, also enforces a caller-visible result boundary. A checker that
// ignores its context may continue running after the evaluator has returned a
// timeout result. health.Checker implementations SHOULD observe ctx whenever
// they can block, perform I/O, wait on another goroutine, or acquire external
// resources.
//
// Evaluator recovers checker panics and converts them into unhealthy results
// with health.ReasonPanic. Panic details are preserved only in
// health.Result.Cause and MUST NOT be exposed by public adapters by default.
type Evaluator struct {
	resolver health.CheckResolver

	clock          clock.PassiveClock
	defaultTimeout time.Duration
	targetTimeouts map[health.Target]time.Duration

	executionPolicy         ExecutionPolicy
	targetExecutionPolicies map[health.Target]ExecutionPolicy
}
