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

// Integer returns the integer payload when v is KindInteger.
//
// For every other kind, Integer returns the zero Integer and ok=false. The zero
// Integer is valid payload data, so callers must check ok.
func (v Value) Integer() (Integer, bool) {
	if v.kind != KindInteger {
		return Integer{}, false
	}

	return v.integerValue, true
}
