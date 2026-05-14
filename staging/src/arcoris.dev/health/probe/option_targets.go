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

import "arcoris.dev/health"

// WithTargets configures the health targets probed by Runner.
//
// Targets are evaluated sequentially in the order supplied here. The same order
// is used by Snapshots when returning cached observations. At least one concrete
// target is required, and duplicates are rejected.
func WithTargets(targets ...health.Target) Option {
	return func(config *config) error {
		normalized, err := normalizeTargets(targets)
		if err != nil {
			return err
		}

		config.targets = normalized
		return nil
	}
}
