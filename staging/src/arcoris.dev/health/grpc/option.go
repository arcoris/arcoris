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

// Option configures a Server at construction time.
//
// Options are adapter-scoped. They configure service-name mappings, polling
// interval, list-size guardrails, and clock ownership. They must not mutate
// health registries, evaluator execution policy, or package-health contracts.
type Option func(*config) error

// applyOptions applies options in caller order to normalized configuration.
//
// Option order is part of the package contract: later options may intentionally
// replace earlier configuration, such as the default service mapping or Watch
// interval. Nil options are rejected as configuration errors instead of being
// skipped silently.
func applyOptions(cfg *config, opts ...Option) error {
	for _, opt := range opts {
		if opt == nil {
			return ErrNilOption
		}
		if err := opt(cfg); err != nil {
			return err
		}
	}

	return nil
}
