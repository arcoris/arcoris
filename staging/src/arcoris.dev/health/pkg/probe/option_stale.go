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

import "time"

// WithStaleAfter configures the cache staleness window.
//
// A zero staleAfter disables stale detection. Positive values enable stale
// detection. Negative values are rejected. Changing the probe schedule does not
// automatically change staleAfter; callers that need a specific freshness ratio
// should configure both explicitly.
func WithStaleAfter(staleAfter time.Duration) Option {
	return func(cfg *config) error {
		if err := validateStaleAfter(staleAfter); err != nil {
			return err
		}

		cfg.staleAfter = staleAfter
		return nil
	}
}
