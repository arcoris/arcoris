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

// FloatFormatMode controls float value encoding.
type FloatFormatMode uint8

const (
	// FloatFormatDefault defers to the package default during Resolve.
	FloatFormatDefault FloatFormatMode = iota

	// FloatFormatReject rejects value.KindFloat during generic JSON encoding.
	FloatFormatReject

	// FloatFormatShortestRoundTrip writes finite floats with shortest round-trip text.
	FloatFormatShortestRoundTrip
)

// isKnownFloatFormatMode reports whether mode is part of the public enum.
func isKnownFloatFormatMode(mode FloatFormatMode) bool {
	switch mode {
	case FloatFormatDefault, FloatFormatReject, FloatFormatShortestRoundTrip:
		return true
	default:
		return false
	}
}
