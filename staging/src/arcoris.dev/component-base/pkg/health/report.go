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

import "time"

// Report describes one target-level health evaluation.
//
// A Report is the Evaluator output for a concrete Target. It aggregates the most
// severe Status from the individual check Results while preserving every Result
// in deterministic registry order. Report intentionally remains a plain value so
// callers can store, copy, render, or adapt it without owning evaluator state.
//
// Report does not define transport behavior. HTTP status mapping, gRPC serving
// state, metrics, logging, restart policy, admission policy, routing, and
// scheduler decisions belong outside package health.
type Report struct {
	// Target is the health scope that was evaluated.
	Target Target

	// Status is the aggregate target status.
	//
	// Evaluator computes Status as the most severe Result status for the target.
	// A report with no checks uses StatusUnknown because no affirmative health
	// observation exists.
	Status Status

	// Observed is the time at which the report was produced.
	Observed time.Time

	// Duration is the evaluator-observed elapsed time for the target evaluation.
	//
	// Evaluator clamps negative durations to zero so mutable test clocks cannot
	// produce invalid runtime reports.
	Duration time.Duration

	// Checks contains normalized check Results in registry order.
	//
	// The slice is caller-owned after Evaluate returns. Evaluator does not retain
	// full reports or mutate returned Results.
	Checks []Result
}
