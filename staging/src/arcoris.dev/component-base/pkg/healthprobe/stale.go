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

package healthprobe

import "time"

// isStale reports whether a cached snapshot with age should be considered stale.
func isStale(age time.Duration, staleAfter time.Duration) bool {
	if staleAfter <= 0 {
		return false
	}
	if age < 0 {
		return false
	}

	return age > staleAfter
}

// validateStaleAfter validates the configured stale window.
func validateStaleAfter(staleAfter time.Duration) error {
	if staleAfter < 0 {
		return InvalidStaleAfterError{StaleAfter: staleAfter}
	}

	return nil
}
