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

// TimeView exposes read-only TypeTime payload data.
type TimeView struct {
	// payload is a detached copy of the time-of-day descriptor payload.
	payload timePayload
}

// Time returns a time-of-day view when t is TypeTime.
func (t Type) Time() (TimeView, bool) {
	return TimeView{payload: cloneTimePayload(t.timeOfDay)}, t.code == TypeTime
}
