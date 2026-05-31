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

// Bool returns the boolean payload when v is KindBool.
//
// For every other kind, Bool returns false and ok=false. The false value is
// only meaningful when ok is true.
func (v Value) Bool() (bool, bool) {
	if v.kind != KindBool {
		return false, false
	}

	return v.boolValue, true
}
