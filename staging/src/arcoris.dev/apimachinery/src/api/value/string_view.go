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

// String returns the string payload when v is KindString.
//
// For every other kind, String returns the empty string and ok=false. The empty
// string is a valid string payload, so callers must check ok.
func (v Value) String() (string, bool) {
	if v.kind != KindString {
		return "", false
	}

	return v.stringValue, true
}
