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

import "arcoris.dev/health"

// Option configures a health HTTP handler at construction time.
//
// Options are adapter-scoped. They may change representation, safe detail
// exposure, target-policy interpretation, and HTTP status code mapping, but they
// must not alter registry contents, evaluator execution policy, or health-core
// contracts.
type Option func(*config) error

// applyOptions applies options to config in order.
//
// Nil options are rejected explicitly so handler construction cannot silently
// drop part of the caller's intended configuration.
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
// This changes only how an already-evaluated report is interpreted for HTTP
// pass/fail purposes. It does not alter how checks run or how report statuses
// are produced by package health.
func WithPolicy(policy health.TargetPolicy) Option {
	return func(config *config) error {
		config.policy = policy
		return nil
	}
}

// WithFormat configures the response format used by the handler.
//
// Format changes representation only. It does not enable content negotiation or
// widen the set of diagnostics exposed by the adapter.
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
// Even DetailAll remains subject to the package's safe DTO boundary and never
// exposes Result.Cause, panic stacks, raw errors, or context causes.
func WithDetailLevel(level DetailLevel) Option {
	return func(config *config) error {
		if err := validateDetailLevel(level); err != nil {
			return err
		}

		config.detailLevel = level
		return nil
	}
}
