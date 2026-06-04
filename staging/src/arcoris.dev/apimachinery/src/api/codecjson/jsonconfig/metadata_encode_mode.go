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

// MetadataEncodeMode controls object envelope metadata emission.
type MetadataEncodeMode uint8

const (
	// MetadataDefault defers to the package default during Resolve.
	MetadataDefault MetadataEncodeMode = iota

	// MetadataOmitZero omits zero metadata.
	MetadataOmitZero

	// MetadataEmitEmpty emits metadata even when it is zero.
	MetadataEmitEmpty
)

// isKnownMetadataEncodeMode reports whether mode is part of the public enum.
func isKnownMetadataEncodeMode(mode MetadataEncodeMode) bool {
	switch mode {
	case MetadataDefault, MetadataOmitZero, MetadataEmitEmpty:
		return true
	default:
		return false
	}
}
