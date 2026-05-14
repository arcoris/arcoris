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

package probe

// Option configures a Runner at construction time.
//
// Options are applied by NewRunner to a private config value before the Runner is
// created. They do not mutate an already constructed Runner and are not retained
// as option values after construction.
//
// Options must remain limited to probe-owned mechanics: clock, schedule,
// snapshot staleness window, initial probe behavior, and target list.
// They must not configure health checks, registries, evaluator execution policy,
// HTTP/gRPC adapters, metrics, logging, tracing, lifecycle transitions, restart
// policy, admission, routing, scheduling, retries, or retry delay growth.
type Option func(*config) error

// applyOptions applies options to normalized configuration in order.
//
// Later options win for single-value domains. WithTargets replaces the previous
// target list. Nil options are rejected explicitly so conditional option
// composition cannot silently drop configuration.
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
