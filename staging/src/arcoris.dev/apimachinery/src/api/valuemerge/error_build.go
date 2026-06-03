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

package valuemerge

import (
	"fmt"

	"arcoris.dev/apimachinery/api/fieldpath"
	"arcoris.dev/apimachinery/api/internal/diagnostic"
)

// errorAt creates a merge diagnostic at a semantic payload path.
func errorAt(path fieldpath.Path, err error, reason ErrorReason, detail string) error {
	return &Error{
		Record: diagnostic.NewRecord(path.String(), err, reason, detail),
	}
}

// errorfAt creates a merge diagnostic with formatted detail text.
func errorfAt(path fieldpath.Path, err error, reason ErrorReason, format string, args ...any) error {
	return errorAt(path, err, reason, fmt.Sprintf(format, args...))
}

// wrapAt attaches valuemerge classification while preserving a lower-level cause.
func wrapAt(
	path fieldpath.Path,
	err error,
	reason ErrorReason,
	detail string,
	cause error,
) error {
	return &Error{
		Record: diagnostic.WrapRecord(path.String(), err, reason, detail, cause),
	}
}
