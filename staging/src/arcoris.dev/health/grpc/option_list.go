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

package healthgrpc

// WithMaxListServices configures the maximum number of services List will
// render before returning RESOURCE_EXHAUSTED.
//
// The limit is an adapter safety guardrail. It prevents unexpectedly large
// service maps from being generated on a single List call without changing the
// configured service mappings or Source evaluation behavior.
func WithMaxListServices(max int) Option {
	return func(cfg *config) error {
		if err := validateMaxListServices(max); err != nil {
			return err
		}

		cfg.maxListServices = max
		return nil
	}
}
