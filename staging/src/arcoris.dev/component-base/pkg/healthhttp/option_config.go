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

// config contains normalized handler configuration.
//
// The struct is package-private because configuration is meaningful only at the
// adapter boundary. Callers express intent through Options instead of mutating
// internal state directly.
type config struct {
	policy      health.TargetPolicy
	format      Format
	detailLevel DetailLevel
	statusCodes HTTPStatusCodes
}

// defaultConfig returns the safe-by-default handler configuration for target.
//
// Defaults intentionally keep the adapter conservative: text responses, no
// per-check details, package-health target policy, and standard HTTP status-code
// mapping.
func defaultConfig(target health.Target) config {
	return config{
		policy:      health.DefaultPolicy(target),
		format:      FormatText,
		detailLevel: DetailNone,
		statusCodes: DefaultStatusCodes(),
	}
}
