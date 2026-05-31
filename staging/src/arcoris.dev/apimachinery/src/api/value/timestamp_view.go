// Copyright 2026 The ARCORIS Authors.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package value

import "time"

// Timestamp returns the timestamp payload when v is KindTimestamp.
//
// For every other kind, Timestamp returns the zero time and ok=false. Timestamp
// values returned by constructors have already had monotonic clock data stripped.
func (v Value) Timestamp() (time.Time, bool) {
	if v.kind != KindTimestamp {
		return time.Time{}, false
	}

	return v.timestampValue, true
}
