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

// DurationView exposes read-only TypeDuration payload data.
type DurationView struct {
	// payload is a detached copy of the duration descriptor payload.
	payload durationPayload
}

// Duration returns a duration view when t is TypeDuration.
func (t Type) Duration() (DurationView, bool) {
	if t.code != TypeDuration {
		return DurationView{}, false
	}

	return DurationView{payload: cloneDurationPayload(t.duration)}, true
}
