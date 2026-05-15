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

import (
	"time"

	"arcoris.dev/health"
)

// formatTimestamp converts a timestamp into the adapter's public representation.
func formatTimestamp(ts time.Time) string {
	if ts.IsZero() {
		return ""
	}

	return ts.UTC().Format(time.RFC3339Nano)
}

// durationMillis converts a duration into whole milliseconds for public output.
func durationMillis(d time.Duration) int64 {
	if d <= 0 {
		return 0
	}

	return d.Milliseconds()
}

// formatReason converts a core health reason into the public string form.
//
// Invalid reasons are collapsed to the stable "invalid" marker instead of
// echoing arbitrary raw strings back to the client.
func formatReason(reason health.Reason) string {
	if reason == health.ReasonNone {
		return ""
	}
	if !reason.IsValid() {
		return "invalid"
	}

	return string(reason)
}
