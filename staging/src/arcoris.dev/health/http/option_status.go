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

// WithStatusCodes configures the full HTTP status code mapping used by the
// handler.
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

// WithPassedStatus configures the HTTP status code used when a report passes
// the configured target policy.
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

// WithFailedStatus configures the HTTP status code used when evaluation
// succeeds but the report fails the configured target policy.
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

// WithErrorStatus configures the HTTP status code used for adapter-boundary
// failures such as invalid handler state.
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
