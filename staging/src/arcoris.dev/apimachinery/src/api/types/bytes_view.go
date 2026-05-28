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

// BytesView exposes read-only TypeBytes payload data.
type BytesView struct {
	// payload is a detached copy of the bytes descriptor payload.
	payload bytesPayload
}

// Bytes returns a bytes view when t is TypeBytes.
func (t Type) Bytes() (BytesView, bool) {
	return BytesView{payload: cloneBytesPayload(t.bytes)}, t.code == TypeBytes
}

// MinLen returns the bytes minimum length rule.
func (v BytesView) MinLen() (int, bool) {
	return v.payload.minLen.value, v.payload.minLen.set
}

// MaxLen returns the bytes maximum length rule.
func (v BytesView) MaxLen() (int, bool) {
	return v.payload.maxLen.value, v.payload.maxLen.set
}
