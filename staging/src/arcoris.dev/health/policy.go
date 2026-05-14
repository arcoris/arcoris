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

// TargetPolicy defines how a Status is interpreted for one health target.
//
// TargetPolicy is intentionally transport-neutral. It does not define HTTP
// status codes, gRPC serving states, restart behavior, traffic routing,
// admission decisions, scheduler decisions, or alerting behavior. It only
// answers whether a health status passes the target-specific health contract.
//
// StatusHealthy always passes. StatusUnhealthy always fails. Invalid status
// values always fail. StatusUnknown fails by default and cannot be allowed by
// this policy because unknown health is not an affirmative operational signal.
//
// StatusStarting and StatusDegraded are target-sensitive. A liveness target may
// allow them because a starting or degraded component can still be alive and
// making progress. Startup and readiness targets normally reject them unless a
// component owner explicitly chooses a more permissive policy.
//
// The zero value is conservative: it allows only StatusHealthy.
type TargetPolicy struct {
	// AllowStarting controls whether StatusStarting passes the target.
	//
	// Starting means startup or initialization is still in progress. It may be
	// acceptable for liveness, but it should normally fail startup and readiness
	// targets.
	AllowStarting bool

	// AllowDegraded controls whether StatusDegraded passes the target.
	//
	// Degraded means the component still has usable capability, but operates with
	// reduced confidence, partial capacity, or protective behavior. It may be
	// acceptable for liveness and selected readiness policies, but it should be an
	// explicit target-owner decision.
	AllowDegraded bool
}

// StartupPolicy returns the default policy for TargetStartup.
//
// Startup requires a fully healthy startup result. Starting, degraded, unknown,
// unhealthy, and invalid statuses fail. This is deliberately strict because the
// startup target answers whether bootstrap has completed, not whether bootstrap
// is still making progress.
func StartupPolicy() TargetPolicy {
	return TargetPolicy{}
}

// LivePolicy returns the default policy for TargetLive.
//
// Liveness is intentionally tolerant of StatusStarting and StatusDegraded. A
// component can be starting or degraded while still making progress and therefore
// should not be restarted solely because of those states.
//
// StatusUnknown and StatusUnhealthy still fail. Unknown health is not an
// affirmative progress signal, and unhealthy represents a direct negative
// liveness observation.
func LivePolicy() TargetPolicy {
	return TargetPolicy{
		AllowStarting: true,
		AllowDegraded: true,
	}
}

// ReadyPolicy returns the default policy for TargetReady.
//
// Readiness is conservative by default. Only StatusHealthy passes. Starting,
// degraded, unknown, unhealthy, and invalid statuses fail so a component is not
// treated as ready for new work unless it has produced an affirmative ready
// signal.
//
// Components that can safely receive selected work while degraded may provide a
// custom policy with AllowDegraded set to true.
func ReadyPolicy() TargetPolicy {
	return TargetPolicy{}
}

// DefaultPolicy returns the built-in policy for target.
//
// TargetUnknown and invalid target values return the zero-value conservative
// policy. Callers that need to reject non-concrete targets should validate the
// target separately. This function intentionally avoids returning an error so it
// can be used in low-level aggregation paths without allocating or branching on
// error values.
func DefaultPolicy(target Target) TargetPolicy {
	switch target {
	case TargetStartup:
		return StartupPolicy()
	case TargetLive:
		return LivePolicy()
	case TargetReady:
		return ReadyPolicy()
	default:
		return TargetPolicy{}
	}
}

// Passes reports whether status passes this target policy.
//
// Passes is a target-policy decision, not a transport or orchestration decision.
// A passing status may later be mapped to an HTTP 2xx response by an HTTP
// adapter, to SERVING by a gRPC adapter, or to another owner-defined outcome, but
// those mappings are outside the health core package.
func (p TargetPolicy) Passes(status Status) bool {
	switch status {
	case StatusHealthy:
		return true
	case StatusStarting:
		return p.AllowStarting
	case StatusDegraded:
		return p.AllowDegraded
	case StatusUnknown,
		StatusUnhealthy:
		return false
	default:
		return false
	}
}

// Fails reports whether status fails this target policy.
//
// Fails is the logical inverse of Passes. It is provided for readability in
// aggregation, tests, and report helpers.
func (p TargetPolicy) Fails(status Status) bool {
	return !p.Passes(status)
}

// Allows reports whether status is one of the non-healthy statuses explicitly
// allowed by this policy.
//
// Allows is narrower than Passes: StatusHealthy passes every valid policy, but it
// is not considered explicitly allowed. This helper is useful in tests and
// diagnostics that need to verify which exceptional states a target policy
// tolerates.
func (p TargetPolicy) Allows(status Status) bool {
	switch status {
	case StatusStarting:
		return p.AllowStarting
	case StatusDegraded:
		return p.AllowDegraded
	default:
		return false
	}
}

// WithStarting returns a copy of p with StatusStarting allowed or rejected.
//
// The method supports fluent policy construction while keeping TargetPolicy a
// plain value type.
func (p TargetPolicy) WithStarting(allow bool) TargetPolicy {
	p.AllowStarting = allow
	return p
}

// WithDegraded returns a copy of p with StatusDegraded allowed or rejected.
//
// The method supports fluent policy construction while keeping TargetPolicy a
// plain value type.
func (p TargetPolicy) WithDegraded(allow bool) TargetPolicy {
	p.AllowDegraded = allow
	return p
}
