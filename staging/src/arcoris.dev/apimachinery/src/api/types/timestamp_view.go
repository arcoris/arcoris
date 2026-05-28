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

package types

// TimestampView exposes read-only TypeTimestamp payload data.
type TimestampView struct {
	// payload is a detached copy of the timestamp descriptor payload.
	payload timestampPayload
}

// Timestamp returns a timestamp view when t is TypeTimestamp.
func (t Type) Timestamp() (TimestampView, bool) {
	if t.code != TypeTimestamp {
		return TimestampView{}, false
	}
	return TimestampView{payload: cloneTimestampPayload(t.timestamp)}, true
}
