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

package value

import "fmt"

// errorf builds a structured value construction error with formatted detail.
//
// It is used for direct validation failures where there is no nested cause. The
// caller chooses the stable path, sentinel, and reason; only Detail is formatted
// for humans.
func errorf(path string, err error, reason ErrorReason, format string, args ...any) error {
	return &Error{
		Path:   path,
		Err:    err,
		Reason: reason,
		Detail: fmt.Sprintf(format, args...),
	}
}

// nested wraps a lower-level construction failure without hiding its identity.
//
// The nested cause remains available through errors.Is/errors.As while the outer
// error gives the current value path and construction classification.
func nested(path string, err error, reason ErrorReason, detail string, cause error) error {
	return &Error{
		Path:   path,
		Err:    err,
		Reason: reason,
		Detail: detail,
		Cause:  cause,
	}
}
