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

package objectwatch

// RestartReason explains why a watch stream can no longer prove continuity.
type RestartReason uint8

const (
	// RestartHistoryUnavailable means the requested revision is no longer
	// available or cannot be served.
	RestartHistoryUnavailable RestartReason = iota + 1
	// RestartContinuityLost means an established stream can no longer prove
	// that no committed changes were skipped.
	RestartContinuityLost
	// RestartSourceReset means the source reset invalidated stream continuity.
	RestartSourceReset
)

// IsValid reports whether r is a supported restart reason.
func (r RestartReason) IsValid() bool {
	return r == RestartHistoryUnavailable ||
		r == RestartContinuityLost ||
		r == RestartSourceReset
}
