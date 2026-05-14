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

package health

// ExecutionMode identifies how Evaluator schedules checks during one target
// evaluation.
//
// ExecutionMode is evaluator-owned execution policy. It does not change checker
// contracts, registry semantics, result normalization, report aggregation,
// target policy, transport mapping, periodic probing, logging, metrics, tracing,
// retries, restart policy, admission policy, or scheduler behavior.
//
// The zero value is ExecutionSequential. This keeps Evaluator conservative and
// predictable unless a component owner explicitly opts into bounded parallel
// execution.
type ExecutionMode uint8

const (
	// ExecutionSequential evaluates checks one by one in Registry registration
	// order.
	//
	// Sequential execution is the default because it has the simplest load
	// profile, the smallest concurrency surface, and preserves the historical
	// Evaluator behavior.
	ExecutionSequential ExecutionMode = iota

	// ExecutionParallel evaluates checks concurrently with a bounded maximum
	// concurrency.
	//
	// Parallel execution is useful for independent I/O-bound checks such as
	// database, queue, cache, and storage probes. It MUST be explicitly
	// configured by the component owner because it may increase instantaneous
	// load on dependencies.
	ExecutionParallel
)

// String returns the stable diagnostic name of mode.
//
// The returned value is intended for diagnostics, tests, logs, and error
// messages. It is not a versioned wire format.
func (mode ExecutionMode) String() string {
	switch mode {
	case ExecutionSequential:
		return "sequential"
	case ExecutionParallel:
		return "parallel"
	default:
		return "invalid"
	}
}

// IsValid reports whether mode is one of the execution modes defined by this
// package.
func (mode ExecutionMode) IsValid() bool {
	switch mode {
	case ExecutionSequential,
		ExecutionParallel:
		return true
	default:
		return false
	}
}
