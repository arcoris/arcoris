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

package fieldpath

// newError builds a structured field-path error without a nested cause.
func newError(err error, reason ErrorReason, detail string) error {
	return &Error{
		Err:    err,
		Reason: reason,
		Detail: detail,
	}
}

// nested wraps a lower-level failure without hiding its identity.
func nested(err error, reason ErrorReason, detail string, cause error) error {
	return &Error{
		Err:    err,
		Reason: reason,
		Detail: detail,
		Cause:  cause,
	}
}
