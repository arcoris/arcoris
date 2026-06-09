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

// BytesView exposes read-only DescriptorBytes payload data.
type BytesView struct {
	// payload is a detached copy of the bytes descriptor payload.
	payload bytesPayload
}

// AsBytes returns a bytes view when desc is DescriptorBytes.
func (desc Descriptor) AsBytes() (BytesView, bool) {
	if desc.code != DescriptorBytes {
		return BytesView{}, false
	}

	return BytesView{payload: cloneBytesPayload(desc.bytes)}, true
}

// MinBytes returns the bytes minimum length rule.
func (v BytesView) MinBytes() (int, bool) {
	return v.payload.minBytes.value, v.payload.minBytes.set
}

// MaxBytes returns the bytes maximum length rule.
func (v BytesView) MaxBytes() (int, bool) {
	return v.payload.maxBytes.value, v.payload.maxBytes.set
}
