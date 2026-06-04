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

// TimestampEncodingMode controls value.KindTimestamp encoding.
type TimestampEncodingMode uint8

const (
	// TimestampEncodingDefault defers to the package default during Resolve.
	TimestampEncodingDefault TimestampEncodingMode = iota

	// TimestampEncodingReject rejects timestamps in descriptor-agnostic JSON.
	TimestampEncodingReject

	// TimestampEncodingRFC3339 is reserved for descriptor-aware timestamp profiles.
	TimestampEncodingRFC3339
)

// isKnownTimestampEncodingMode reports whether mode is part of the public enum.
func isKnownTimestampEncodingMode(mode TimestampEncodingMode) bool {
	switch mode {
	case TimestampEncodingDefault, TimestampEncodingReject, TimestampEncodingRFC3339:
		return true
	default:
		return false
	}
}
