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

// DurationEncodingMode controls value.KindDuration encoding.
type DurationEncodingMode uint8

const (
	// DurationEncodingDefault defers to the package default during Resolve.
	DurationEncodingDefault DurationEncodingMode = iota

	// DurationEncodingReject rejects durations in descriptor-agnostic JSON.
	DurationEncodingReject

	// DurationEncodingISO8601 is reserved for descriptor-aware duration profiles.
	DurationEncodingISO8601
)

// isKnownDurationEncodingMode reports whether mode is part of the public enum.
func isKnownDurationEncodingMode(mode DurationEncodingMode) bool {
	switch mode {
	case DurationEncodingDefault, DurationEncodingReject, DurationEncodingISO8601:
		return true
	default:
		return false
	}
}
