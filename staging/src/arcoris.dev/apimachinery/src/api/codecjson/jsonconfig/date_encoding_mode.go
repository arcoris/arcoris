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

// DateEncodingMode controls value.KindDate encoding.
type DateEncodingMode uint8

const (
	// DateEncodingDefault defers to the package default during Resolve.
	DateEncodingDefault DateEncodingMode = iota

	// DateEncodingReject rejects dates in descriptor-agnostic JSON.
	DateEncodingReject

	// DateEncodingISO8601 is reserved for descriptor-aware date profiles.
	DateEncodingISO8601
)

// isKnownDateEncodingMode reports whether mode is part of the public enum.
func isKnownDateEncodingMode(mode DateEncodingMode) bool {
	switch mode {
	case DateEncodingDefault, DateEncodingReject, DateEncodingISO8601:
		return true
	default:
		return false
	}
}
