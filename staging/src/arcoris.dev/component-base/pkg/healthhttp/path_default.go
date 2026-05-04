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

// DefaultPath returns the default HTTP path for target.
//
// DefaultPath only maps concrete health targets that have unambiguous default
// probe paths:
//
//   - health.TargetStartup -> DefaultStartupPath;
//   - health.TargetLive -> DefaultLivePath;
//   - health.TargetReady -> DefaultReadyPath.
//
// The function returns false for health.TargetUnknown and invalid target values.
// It also intentionally does not return DefaultHealthPath or
// DefaultHealthPlainPath because those compatibility paths do not have a
// universal target meaning.
func DefaultPath(target health.Target) (string, bool) {
	switch target {
	case health.TargetStartup:
		return DefaultStartupPath, true
	case health.TargetLive:
		return DefaultLivePath, true
	case health.TargetReady:
		return DefaultReadyPath, true
	default:
		return "", false
	}
}
