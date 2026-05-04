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
type Option func(*config) error

// config contains normalized health HTTP handler configuration.
type config struct {
	policy      health.TargetPolicy
	format      Format
	detailLevel DetailLevel
	statusCodes HTTPStatusCodes
}

// defaultConfig returns the default handler configuration for target.
func defaultConfig(target health.Target) config {
	return config{
		policy:      health.DefaultPolicy(target),
		format:      FormatText,
		detailLevel: DetailNone,
		statusCodes: DefaultStatusCodes(),
	}
}

// applyOptions applies options to config in order.
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
func WithPolicy(policy health.TargetPolicy) Option {
	return func(config *config) error {
		config.policy = policy
		return nil
	}
}

// WithFormat configures the response format used by the handler.
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

// WithPassedStatus configures the HTTP status code used when a report passes.
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

// WithFailedStatus configures the HTTP status code used when a report fails.
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

// WithErrorStatus configures the HTTP status code used for adapter errors.
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
