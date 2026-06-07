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

package health

import "time"

// Result describes one health observation produced by a checker, gate, cached
// probe, or another health source.
//
// Result is intentionally transport-neutral. It does not define HTTP response
// codes, gRPC serving states, JSON rendering, metric labels, restart decisions,
// readiness outcomes, admission behavior, or scheduler policy. Those concerns
// belong to adapter packages and higher-level runtime owners.
//
// The zero value is an unnamed StatusUnknown result. It is safe to allocate in
// larger structs before the first health observation is available. Callers that
// publish, aggregate, or expose results SHOULD normalize them at the boundary
// where checker ownership and observation time are known.
//
// Name identifies the logical check that produced the result. A checker-owned
// result SHOULD set Name to the checker name. Aggregators MAY fill an empty Name
// from the checker that returned the result.
//
// Status is the primary operational health state.
//
// Reason is a stable, machine-readable explanation for Status when a reason is
// useful. A healthy result MAY leave Reason empty. Non-healthy results SHOULD
// provide a reason when the owner can classify the condition without leaking
// private implementation detail.
//
// Message is a safe, short, human-readable explanation. Message MUST NOT contain
// private operational data. Transport adapters may expose Message directly.
//
// Observed is the time at which the result was observed or normalized. A zero
// value means the result has not been timestamped yet.
//
// Duration is the amount of time spent producing the observation. A zero value is
// valid. Negative durations are invalid and SHOULD be normalized before
// aggregation or exposure.
//
// Cause preserves the internal lower-level failure cause. Cause is intentionally
// not a public diagnostic field. Transport adapters MUST NOT expose Cause by
// default. Logs, tests, and owner-controlled diagnostics may inspect Cause when
// they have permission to handle internal details.
type Result struct {
	// Name identifies the logical check that produced this result. Empty Name is
	// allowed while a result is still detached from checker ownership.
	Name string

	// Status is the primary health state observed by the checker.
	Status Status

	// Reason is the stable machine-readable classification for Status.
	Reason Reason

	// Message is the safe human-readable explanation for diagnostics and
	// transport adapters.
	Message string

	// Observed records when the result was produced or normalized. A zero value
	// means no observation timestamp has been attached yet.
	Observed time.Time

	// Duration records how long the observation took. Negative durations are
	// structurally invalid and should be normalized at ownership boundaries.
	Duration time.Duration

	// Cause preserves an internal lower-level cause for owner diagnostics. Public
	// adapters must not expose Cause by default.
	Cause error
}

// Healthy returns a healthy result for name.
//
// A healthy result intentionally has no reason or message by default. Callers
// that need a human-readable explanation for diagnostics may set Message with a
// follow-up value transformation, but most healthy results should remain compact.
func Healthy(name string) Result {
	return Result{
		Name:   name,
		Status: StatusHealthy,
	}
}

// Starting returns a starting result for name.
//
// Starting means the checked component or subsystem is still bootstrapping. It
// is not a terminal failure. Target policies decide whether starting is accepted
// for a specific health target.
func Starting(name string, reason Reason, message string) Result {
	return Result{
		Name:    name,
		Status:  StatusStarting,
		Reason:  reason,
		Message: message,
	}
}

// Degraded returns a degraded result for name.
//
// Degraded means the checked component or subsystem still has usable capability,
// but is operating with reduced capability, reduced confidence, or active
// protective behavior. Degraded MUST remain distinct from unhealthy so runtime
// owners can make target-specific decisions.
func Degraded(name string, reason Reason, message string) Result {
	return Result{
		Name:    name,
		Status:  StatusDegraded,
		Reason:  reason,
		Message: message,
	}
}

// Unhealthy returns an unhealthy result for name.
//
// Unhealthy is the strongest negative health observation. It describes the
// checked scope only; it does not by itself prescribe restart, traffic removal,
// admission closure, or scheduler exclusion.
func Unhealthy(name string, reason Reason, message string) Result {
	return Result{
		Name:    name,
		Status:  StatusUnhealthy,
		Reason:  reason,
		Message: message,
	}
}

// Unknown returns an unknown result for name.
//
// Unknown means the checker could not produce a reliable health observation. It
// is useful for timeouts, cancellations, uninitialized cached results, missing
// state, invalid caller-controlled input, or inconclusive checks.
func Unknown(name string, reason Reason, message string) Result {
	return Result{
		Name:    name,
		Status:  StatusUnknown,
		Reason:  reason,
		Message: message,
	}
}
