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

package healthhttp

import "arcoris.dev/component-base/pkg/health"

// Option configures a health HTTP handler at construction time.
//
// Options are applied to a private config value before a handler is created.
// They do not mutate an already constructed handler and are not retained after
// construction.
//
// Options must stay limited to HTTP adapter configuration:
//
//   - target policy;
//   - response format;
//   - response detail level;
//   - HTTP status code mapping.
//
// Options must not configure health checks, registries, evaluator execution,
// lifecycle transitions, signals, logging, metrics, tracing, authentication,
// authorization, routing decisions, scheduling decisions, admission decisions,
// or periodic probe execution.
type Option func(*config) error

// config contains normalized health HTTP handler configuration.
//
// The type is intentionally package-private. Public callers configure handlers
// through Option constructors, while NewHandler receives a fully normalized
// config.
//
// config must remain a small adapter-level object. If future features need
// independent domains such as query-parameter detail escalation, custom
// renderers, or middleware integration, they should be introduced as separate
// focused options rather than turning config into a global server configuration
// object.
type config struct {
	// policy decides whether the evaluated health report passes the handler's
	// target.
	//
	// The default policy is derived from the target:
	//
	//   - startup uses health.StartupPolicy();
	//   - live uses health.LivePolicy();
	//   - ready uses health.ReadyPolicy();
	//
	// Callers may override the policy with WithPolicy.
	policy health.TargetPolicy

	// format controls response representation.
	//
	// The zero/default format is FormatText because probe clients primarily
	// consume status codes and text is the smallest safe response body.
	format Format

	// detailLevel controls how much check-level information renderers may expose.
	//
	// The zero/default level is DetailNone because health endpoints may be
	// reachable by load balancers, orchestrators, and infrastructure probes.
	detailLevel DetailLevel

	// statusCodes maps handler outcomes to HTTP status codes.
	//
	// The default mapping is:
	//
	//   - passed report -> 200 OK;
	//   - failed health report -> 503 Service Unavailable;
	//   - adapter/evaluator boundary error -> 500 Internal Server Error.
	statusCodes HTTPStatusCodes
}

// defaultConfig returns the default handler configuration for target.
//
// The function assumes target has already been validated by the caller. Invalid
// and unknown targets receive the conservative zero-value target policy, which
// only passes healthy status. NewHandler should reject invalid targets before a
// handler is constructed.
func defaultConfig(target health.Target) config {
	return config{
		policy:      defaultTargetPolicy(target),
		format:      FormatText,
		detailLevel: DetailNone,
		statusCodes: DefaultStatusCodes(),
	}
}

// defaultTargetPolicy returns the default health policy for target.
//
// The mapping mirrors package health semantics while keeping healthhttp free
// from hidden target-specific behavior:
//
//   - startup requires healthy status;
//   - live allows starting and degraded status;
//   - ready requires healthy status.
func defaultTargetPolicy(target health.Target) health.TargetPolicy {
	switch target {
	case health.TargetStartup:
		return health.StartupPolicy()
	case health.TargetLive:
		return health.LivePolicy()
	case health.TargetReady:
		return health.ReadyPolicy()
	default:
		return health.TargetPolicy{}
	}
}

// applyOptions applies options to config in order.
//
// Later options win. Nil options are rejected with ErrNilOption so conditional
// option construction bugs are visible at handler construction time instead of
// being silently ignored.
func applyOptions(config *config, options ...Option) error {
	for _, option := range options {
		if option == nil {
			return ErrNilOption
		}
		if err := option(config); err != nil {
			return err
		}
	}

	return nil
}

// WithPolicy configures the target policy used by the handler.
//
// The policy decides whether a health.Report passes after Evaluator finishes
// successfully. It does not change the evaluated target and does not affect how
// checks are executed.
//
// health.TargetPolicy has no invalid state: the zero value is a strict policy
// that only passes healthy status. Therefore WithPolicy does not perform
// validation.
func WithPolicy(policy health.TargetPolicy) Option {
	return func(config *config) error {
		config.policy = policy
		return nil
	}
}

// WithFormat configures the response format used by the handler.
//
// Invalid formats are rejected at construction time so ServeHTTP never has to
// guess how to render a response.
func WithFormat(format Format) Option {
	return func(config *config) error {
		if err := validateFormat(format); err != nil {
			return err
		}

		config.format = format
		return nil
	}
}

// WithDetailLevel configures the amount of safe check-level detail exposed by
// the handler renderer.
//
// Detail level affects only response body detail. It does not affect evaluation,
// target policy, HTTP status code mapping, logging, metrics, or authorization.
func WithDetailLevel(level DetailLevel) Option {
	return func(config *config) error {
		if err := validateDetailLevel(level); err != nil {
			return err
		}

		config.detailLevel = level
		return nil
	}
}

// WithStatusCodes configures the HTTP status code mapping used by the handler.
//
// Zero fields are replaced by defaults before validation. This lets callers
// override only the codes they need:
//
//	WithStatusCodes(HTTPStatusCodes{Failed: http.StatusTooManyRequests})
//
// The normalized mapping must satisfy HTTPStatusCodes.Validate.
func WithStatusCodes(codes HTTPStatusCodes) Option {
	return func(config *config) error {
		codes = codes.Normalize()
		if err := codes.Validate(); err != nil {
			return err
		}

		config.statusCodes = codes
		return nil
	}
}

// WithPassedStatus configures the HTTP status code used when a report passes the
// handler target policy.
//
// The code must be a 2xx status. Other fields keep their current configured
// values.
func WithPassedStatus(code int) Option {
	return func(config *config) error {
		codes := config.statusCodes
		codes.Passed = code
		codes = codes.Normalize()

		if err := codes.Validate(); err != nil {
			return err
		}

		config.statusCodes = codes
		return nil
	}
}

// WithFailedStatus configures the HTTP status code used when health evaluation
// succeeds but the report fails the handler target policy.
//
// The code must be a 4xx or 5xx status. Other fields keep their current
// configured values.
func WithFailedStatus(code int) Option {
	return func(config *config) error {
		codes := config.statusCodes
		codes.Failed = code
		codes = codes.Normalize()

		if err := codes.Validate(); err != nil {
			return err
		}

		config.statusCodes = codes
		return nil
	}
}

// WithErrorStatus configures the HTTP status code used when the HTTP adapter
// cannot produce a reliable health response because of a handler, evaluator, or
// configuration boundary error.
//
// The code must be a 5xx status. Other fields keep their current configured
// values.
func WithErrorStatus(code int) Option {
	return func(config *config) error {
		codes := config.statusCodes
		codes.Error = code
		codes = codes.Normalize()

		if err := codes.Validate(); err != nil {
			return err
		}

		config.statusCodes = codes
		return nil
	}
}
