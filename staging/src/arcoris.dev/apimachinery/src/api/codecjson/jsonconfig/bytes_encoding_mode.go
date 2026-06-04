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

package jsonconfig

// BytesEncodingMode controls value.KindBytes encoding.
type BytesEncodingMode uint8

const (
	// BytesEncodingDefault defers to the package default during Resolve.
	BytesEncodingDefault BytesEncodingMode = iota

	// BytesEncodingReject rejects bytes in descriptor-agnostic JSON.
	BytesEncodingReject

	// BytesEncodingBase64Std is reserved for descriptor-aware bytes profiles.
	BytesEncodingBase64Std

	// BytesEncodingBase64URL is reserved for descriptor-aware bytes profiles.
	BytesEncodingBase64URL
)

// isKnownBytesEncodingMode reports whether mode is part of the public enum.
func isKnownBytesEncodingMode(mode BytesEncodingMode) bool {
	switch mode {
	case BytesEncodingDefault, BytesEncodingReject, BytesEncodingBase64Std, BytesEncodingBase64URL:
		return true
	default:
		return false
	}
}
