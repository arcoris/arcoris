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

// AsBytes returns a detached bytes payload when v is KindBytes.
//
// The returned slice has its own backing array. Mutating it cannot affect the
// stored Value. For non-byte values, AsBytes returns nil and ok=false.
func (v Value) AsBytes() ([]byte, bool) {
	if v.kind != KindBytes {
		return nil, false
	}

	return cloneBytes(v.bytesValue), true
}
