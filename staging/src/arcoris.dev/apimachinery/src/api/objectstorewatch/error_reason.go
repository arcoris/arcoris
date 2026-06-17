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

package objectstorewatch

// ErrorReason identifies a stable objectstorewatch diagnostic class.
type ErrorReason string

const (
	// ErrorReasonInvalidSnapshot identifies a malformed list-to-watch snapshot.
	ErrorReasonInvalidSnapshot ErrorReason = "invalid_snapshot"
	// ErrorReasonInvalidBoundary identifies a malformed watch boundary.
	ErrorReasonInvalidBoundary ErrorReason = "invalid_boundary"
)

// IsValid reports whether r is a known objectstorewatch error reason.
func (r ErrorReason) IsValid() bool {
	return r == ErrorReasonInvalidSnapshot || r == ErrorReasonInvalidBoundary
}

// String returns stable diagnostic text for r.
func (r ErrorReason) String() string {
	if !r.IsValid() {
		return "unknown"
	}

	return string(r)
}
