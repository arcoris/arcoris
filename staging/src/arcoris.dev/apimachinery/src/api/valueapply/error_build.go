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

package valueapply

import (
	"fmt"

	"arcoris.dev/apimachinery/api/fieldpath"
	"arcoris.dev/apimachinery/api/internal/diagnostic"
)

// errorAt builds a valueapply diagnostic without a nested cause.
func errorAt(path fieldpath.Path, err error, reason ErrorReason, detail string) error {
	return &Error{
		Record: diagnostic.Record[ErrorReason]{
			Path:   path.String(),
			Err:    err,
			Reason: reason,
			Detail: detail,
		},
	}
}

// wrapAt builds a valueapply diagnostic that preserves a lower-level cause.
func wrapAt(path fieldpath.Path, err error, reason ErrorReason, detail string, cause error) error {
	return &Error{
		Record: diagnostic.Record[ErrorReason]{
			Path:   path.String(),
			Err:    err,
			Reason: reason,
			Detail: detail,
			Cause:  cause,
		},
	}
}

// errorfAt formats the diagnostic detail for callers with dynamic context.
func errorfAt(path fieldpath.Path, err error, reason ErrorReason, format string, args ...any) error {
	return errorAt(path, err, reason, fmt.Sprintf(format, args...))
}
